// Copyright 2017 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package control

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/coreos/etcd/pkg/report"
	"github.com/gyuho/dataframe"
	"golang.org/x/net/context"
)

type values struct {
	bytes      [][]byte
	strings    []string
	sampleSize int
}

func newValues(cfg Config) (v values, rerr error) {
	v.bytes = [][]byte{randBytes(cfg.Step2.ValueSize)}
	v.strings = []string{string(v.bytes[0])}
	v.sampleSize = 1
	return
}

type benchmark struct {
	bar        *pb.ProgressBar
	report     report.Report
	reportDone <-chan report.Stats
	stats      report.Stats

	reqHandlers []ReqHandler
	reqGen      func(chan<- request)
	reqDone     func()
	wg          sync.WaitGroup

	mu           sync.RWMutex
	inflightReqs chan request
}

// pass totalN in case that 'cfg' is manipulated
func newBenchmark(totalN int, clientsN int, reqHandlers []ReqHandler, reqDone func(), reqGen func(chan<- request)) (b *benchmark) {
	b = &benchmark{
		bar:         pb.New(totalN),
		reqHandlers: reqHandlers,
		reqGen:      reqGen,
		reqDone:     reqDone,
		wg:          sync.WaitGroup{},
	}
	b.inflightReqs = make(chan request, clientsN)

	b.bar.Format("Bom !")
	b.bar.Start()
	b.report = report.NewReportSample("%4.4f")
	return
}

// only useful when multiple ranges of requests are run with one report
func (b *benchmark) reset(clientsN int, reqHandlers []ReqHandler, reqDone func(), reqGen func(chan<- request)) {
	if len(reqHandlers) == 0 {
		panic(fmt.Errorf("got 0 reqHandlers"))
	}
	b.reqHandlers = reqHandlers
	b.reqDone = reqDone
	b.reqGen = reqGen

	// inflight requests will be dropped!
	b.mu.Lock()
	b.inflightReqs = make(chan request, clientsN)
	b.mu.Unlock()
}

func (b *benchmark) getInflightsReqs() (ch chan request) {
	b.mu.RLock()
	ch = b.inflightReqs
	b.mu.RUnlock()
	return
}

func (b *benchmark) startRequests() {
	for i := range b.reqHandlers {
		b.wg.Add(1)
		go func(rh ReqHandler) {
			defer b.wg.Done()
			for req := range b.getInflightsReqs() {
				if rh == nil {
					panic(fmt.Errorf("got nil rh"))
				}
				st := time.Now()
				err := rh(context.Background(), &req)
				b.report.Results() <- report.Result{Err: err, Start: st, End: time.Now()}
				b.bar.Increment()
			}
		}(b.reqHandlers[i])
	}
	go b.reqGen(b.getInflightsReqs())
	b.reportDone = b.report.Stats()
}

func (b *benchmark) waitRequestsEnd() {
	b.wg.Wait()
	if b.reqDone != nil {
		b.reqDone() // cancel connections
	}
}

func (b *benchmark) finishReports() {
	close(b.report.Results())
	b.bar.Finish()
	st := <-b.reportDone
	b.stats = st
}

func (b *benchmark) waitAll() {
	b.waitRequestsEnd()
	b.finishReports()
}

func printStats(st report.Stats) {
	// to be piped to cfg.Log via stdout when dbtester executed
	if len(st.Lats) > 0 {
		fmt.Printf("Total: %v\n", st.Total)
		fmt.Printf("Slowest: %f secs\n", st.Slowest)
		fmt.Printf("Fastest: %f secs\n", st.Fastest)
		fmt.Printf("Average: %f secs\n", st.Average)
		fmt.Printf("Requests/sec: %4.4f\n", st.RPS)
	}
	if len(st.ErrorDist) > 0 {
		for k, v := range st.ErrorDist {
			fmt.Printf("ERROR %q : %d\n", k, v)
		}
	} else {
		fmt.Println("ERRRO: 0")
	}
}

func saveDataLatencyDistributionSummary(cfg Config, st report.Stats) {
	fr := dataframe.New()

	c1 := dataframe.NewColumn("TOTAL-SECONDS")
	c1.PushBack(dataframe.NewStringValue(fmt.Sprintf("%4.4f", st.Total.Seconds())))
	if err := fr.AddColumn(c1); err != nil {
		plog.Fatal(err)
	}

	c2 := dataframe.NewColumn("REQUESTS-PER-SECOND")
	c2.PushBack(dataframe.NewStringValue(fmt.Sprintf("%4.4f", st.RPS)))
	if err := fr.AddColumn(c2); err != nil {
		plog.Fatal(err)
	}

	c3 := dataframe.NewColumn("SLOWEST-LATENCY-MS")
	c3.PushBack(dataframe.NewStringValue(fmt.Sprintf("%4.4f", 1000*st.Slowest)))
	if err := fr.AddColumn(c3); err != nil {
		plog.Fatal(err)
	}

	c4 := dataframe.NewColumn("FASTEST-LATENCY-MS")
	c4.PushBack(dataframe.NewStringValue(fmt.Sprintf("%4.4f", 1000*st.Fastest)))
	if err := fr.AddColumn(c4); err != nil {
		plog.Fatal(err)
	}

	c5 := dataframe.NewColumn("AVERAGE-LATENCY-MS")
	c5.PushBack(dataframe.NewStringValue(fmt.Sprintf("%4.4f", 1000*st.Average)))
	if err := fr.AddColumn(c5); err != nil {
		plog.Fatal(err)
	}

	c6 := dataframe.NewColumn("STDDEV-LATENCY-MS")
	c6.PushBack(dataframe.NewStringValue(fmt.Sprintf("%4.4f", 1000*st.Stddev)))
	if err := fr.AddColumn(c6); err != nil {
		plog.Fatal(err)
	}

	if len(st.ErrorDist) > 0 {
		for errName, errN := range st.ErrorDist {
			errcol := dataframe.NewColumn(fmt.Sprintf("ERROR: %q", errName))
			errcol.PushBack(dataframe.NewStringValue(errN))
			if err := fr.AddColumn(errcol); err != nil {
				plog.Fatal(err)
			}
		}
	} else {
		errcol := dataframe.NewColumn("ERROR")
		errcol.PushBack(dataframe.NewStringValue("0"))
		if err := fr.AddColumn(errcol); err != nil {
			plog.Fatal(err)
		}
	}

	if err := fr.CSVHorizontal(cfg.DataLatencyDistributionSummary); err != nil {
		plog.Fatal(err)
	}
}

func saveDataLatencyDistributionPercentile(cfg Config, st report.Stats) {
	pctls, seconds := report.Percentiles(st.Lats)
	c1 := dataframe.NewColumn("LATENCY-PERCENTILE")
	c2 := dataframe.NewColumn("LATENCY-MS")
	for i := range pctls {
		pct := fmt.Sprintf("p%.1f", pctls[i])
		if strings.HasSuffix(pct, ".0") {
			pct = strings.Replace(pct, ".0", "", -1)
		}

		c1.PushBack(dataframe.NewStringValue(pct))
		c2.PushBack(dataframe.NewStringValue(fmt.Sprintf("%f", 1000*seconds[i])))
	}

	fr := dataframe.New()
	if err := fr.AddColumn(c1); err != nil {
		plog.Fatal(err)
	}
	if err := fr.AddColumn(c2); err != nil {
		plog.Fatal(err)
	}
	if err := fr.CSV(cfg.DataLatencyDistributionPercentile); err != nil {
		plog.Fatal(err)
	}
}

func saveDataLatencyDistributionAll(cfg Config, st report.Stats) {
	min := int64(math.MaxInt64)
	max := int64(-100000)
	rm := make(map[int64]int64)
	for _, lt := range st.Lats {
		// convert second(float64) to millisecond
		ms := lt * 1000

		// truncate all digits below 10ms
		// (e.g. 125.11ms becomes 120ms)
		v := int64(math.Trunc(ms/10) * 10)
		if _, ok := rm[v]; !ok {
			rm[v] = 1
		} else {
			rm[v]++
		}

		if min > v {
			min = v
		}
		if max < v {
			max = v
		}
	}

	c1 := dataframe.NewColumn("LATENCY-MS")
	c2 := dataframe.NewColumn("COUNT")
	cur := min
	for {
		c1.PushBack(dataframe.NewStringValue(fmt.Sprintf("%d", int64(cur))))
		v, ok := rm[cur]
		if ok {
			c2.PushBack(dataframe.NewStringValue(fmt.Sprintf("%d", v)))
		} else {
			c2.PushBack(dataframe.NewStringValue("0"))
		}
		cur += 10
		if cur-10 == max { // was last point
			break
		}
	}
	fr := dataframe.New()
	if err := fr.AddColumn(c1); err != nil {
		plog.Fatal(err)
	}
	if err := fr.AddColumn(c2); err != nil {
		plog.Fatal(err)
	}
	if err := fr.CSV(cfg.DataLatencyDistributionAll); err != nil {
		plog.Fatal(err)
	}
}

func saveDataLatencyThroughputTimeseries(cfg Config, st report.Stats, tsToClientN map[int64]int) {
	// TODO: UNIX-TS from pkg/report data is time.Time.Unix
	// UNIX-TS from psn.CSV data is time.Time.UnixNano
	// we need some kind of way to combine those with matching timestamps
	c1 := dataframe.NewColumn("UNIX-TS")
	c2 := dataframe.NewColumn("CONTROL-CLIENT-NUM")
	c3 := dataframe.NewColumn("AVG-LATENCY-MS")
	c4 := dataframe.NewColumn("AVG-THROUGHPUT")
	for i := range st.TimeSeries {
		// this Timestamp is unix seconds
		c1.PushBack(dataframe.NewStringValue(fmt.Sprintf("%d", st.TimeSeries[i].Timestamp)))
		if len(tsToClientN) == 0 {
			c2.PushBack(dataframe.NewStringValue(fmt.Sprintf("%d", cfg.Step2.Clients)))
		} else {
			c2.PushBack(dataframe.NewStringValue(fmt.Sprintf("%d", tsToClientN[st.TimeSeries[i].Timestamp])))
		}
		c3.PushBack(dataframe.NewStringValue(fmt.Sprintf("%f", toMillisecond(st.TimeSeries[i].AvgLatency))))
		c4.PushBack(dataframe.NewStringValue(fmt.Sprintf("%d", st.TimeSeries[i].ThroughPut)))
	}

	fr := dataframe.New()
	if err := fr.AddColumn(c1); err != nil {
		plog.Fatal(err)
	}
	if err := fr.AddColumn(c2); err != nil {
		plog.Fatal(err)
	}
	if err := fr.AddColumn(c3); err != nil {
		plog.Fatal(err)
	}
	if err := fr.AddColumn(c4); err != nil {
		plog.Fatal(err)
	}
	if err := fr.CSV(cfg.DataLatencyThroughputTimeseries); err != nil {
		plog.Fatal(err)
	}

	// aggregate latency by the number of keys
	tss := processTimeSeries(st.TimeSeries, 1000, cfg.Step2.TotalRequests)
	ctt1 := dataframe.NewColumn("KEYS")
	ctt2 := dataframe.NewColumn("AVG-LATENCY-MS")
	for i := range tss {
		ctt1.PushBack(dataframe.NewStringValue(tss[i].keyNum))
		ctt2.PushBack(dataframe.NewStringValue(fmt.Sprintf("%f", toMillisecond(tss[i].avgLat))))
	}
	frr := dataframe.New()
	if err := frr.AddColumn(ctt1); err != nil {
		plog.Fatal(err)
	}
	if err := frr.AddColumn(ctt2); err != nil {
		plog.Fatal(err)
	}
	if err := frr.CSV(cfg.DataLatencyByKeyNumber); err != nil {
		plog.Fatal(err)
	}
}

func generateReport(cfg Config, h []ReqHandler, reqDone func(), reqGen func(chan<- request)) {
	b := newBenchmark(cfg.Step2.TotalRequests, cfg.Step2.Clients, h, reqDone, reqGen)
	b.startRequests()
	b.waitAll()

	printStats(b.stats)
	saveAllStats(cfg, b.stats, nil)
}

func saveAllStats(cfg Config, stats report.Stats, tsToClientN map[int64]int) {
	// cfg.DataLatencyDistributionSummary
	saveDataLatencyDistributionSummary(cfg, stats)

	// cfg.DataLatencyDistributionPercentile
	saveDataLatencyDistributionPercentile(cfg, stats)

	// cfg.DataLatencyDistributionAll
	saveDataLatencyDistributionAll(cfg, stats)

	// cfg.DataLatencyThroughputTimeseries
	saveDataLatencyThroughputTimeseries(cfg, stats, tsToClientN)
}

// processTimeSeries sorts all data points by its timestamp.
// And then aggregate by the cumulative throughput,
// in order to map the number of keys to the average latency.
//
//	type DataPoint struct {
//		Timestamp  int64
//		AvgLatency time.Duration
//		ThroughPut int64
//	}
//
// If unis is 1000 and the average throughput per second is 30,000
// and its average latency is 10ms, it will have 30 data points with
// latency 10ms.
func processTimeSeries(tss report.TimeSeries, unit int64, totalRequests int) keyNumToAvgLatencys {
	sort.Sort(tss)

	cumulKeyN := int64(0)
	maxKey := int64(0)

	rm := make(map[int64]time.Duration)

	// this data is aggregated by second
	// and we want to map number of keys to latency
	// so the range is the key
	// and the value is the cumulative throughput
	for _, ts := range tss {
		cumulKeyN += ts.ThroughPut
		if cumulKeyN < unit {
			// not enough data points yet
			continue
		}

		lat := ts.AvgLatency

		// cumulKeyN >= unit
		for cumulKeyN > maxKey {
			maxKey += unit
			rm[maxKey] = lat
		}
	}

	// fill-in empty rows
	for i := maxKey; i < int64(totalRequests); i += unit {
		if _, ok := rm[i]; !ok {
			rm[i] = time.Duration(0)
		}
	}
	if _, ok := rm[int64(totalRequests)]; !ok {
		rm[int64(totalRequests)] = time.Duration(0)
	}

	kss := []keyNumToAvgLatency{}
	for k, v := range rm {
		kss = append(kss, keyNumToAvgLatency{keyNum: k, avgLat: v})
	}
	sort.Sort(keyNumToAvgLatencys(kss))

	return kss
}

type keyNumToAvgLatency struct {
	keyNum int64
	avgLat time.Duration
}

type keyNumToAvgLatencys []keyNumToAvgLatency

func (t keyNumToAvgLatencys) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t keyNumToAvgLatencys) Len() int           { return len(t) }
func (t keyNumToAvgLatencys) Less(i, j int) bool { return t[i].keyNum < t[j].keyNum }
