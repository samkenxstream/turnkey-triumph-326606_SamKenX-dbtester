// Copyright 2016 CoreOS, Inc.
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

package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

type timeSeries struct {
	timestamp  int64
	avgLatency time.Duration
	throughPut int64
}

type TimeSeries []timeSeries

func (t TimeSeries) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TimeSeries) Len() int           { return len(t) }
func (t TimeSeries) Less(i, j int) bool { return t[i].timestamp < t[j].timestamp }

type secondPoint struct {
	totalLatency time.Duration
	count        int64
}

type secondPoints struct {
	mu sync.Mutex
	tm map[int64]secondPoint
}

func newSecondPoints() *secondPoints {
	return &secondPoints{tm: make(map[int64]secondPoint)}
}

func (sp *secondPoints) Add(ts time.Time, lat time.Duration) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	tk := ts.Unix()
	if v, ok := sp.tm[tk]; !ok {
		sp.tm[tk] = secondPoint{totalLatency: lat, count: 1}
	} else {
		v.totalLatency += lat
		v.count += 1
		sp.tm[tk] = v
	}
}

func (sp *secondPoints) getTimeSeries() TimeSeries {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	tslice := make(TimeSeries, len(sp.tm))
	i := 0
	log.Printf("getTimeSeries has started on %d results\n", len(sp.tm))
	for k, v := range sp.tm {
		tslice[i] = timeSeries{
			timestamp:  k,
			avgLatency: time.Duration(v.totalLatency) / time.Duration(v.count),
			throughPut: v.count,
		}
		i++
		if i%100 == 0 {
			log.Printf("processing timeseries at %d / %d", i, len(sp.tm))
		}
	}
	sort.Sort(tslice)
	return tslice
}

func (ts TimeSeries) String() string {
	buf := new(bytes.Buffer)
	wr := csv.NewWriter(buf)
	if err := wr.Write([]string{"unix_ts", "avg_latency_ms", "throughput"}); err != nil {
		log.Fatal(err)
	}
	rows := [][]string{}
	for i := range ts {
		row := []string{
			fmt.Sprintf("%d", ts[i].timestamp),
			fmt.Sprintf("%f", toMillisecond(ts[i].avgLatency)),
			fmt.Sprintf("%d", ts[i].throughPut),
		}
		rows = append(rows, row)
	}
	if err := wr.WriteAll(rows); err != nil {
		log.Fatal(err)
	}
	wr.Flush()
	if err := wr.Error(); err != nil {
		log.Fatal(err)
	}
	txt := buf.String()
	if err := toFile(txt, csvResultPath); err != nil {
		log.Println(err)
	} else {
		log.Println("time series saved... uploading to Google cloud storage...")
		kbts, err := ioutil.ReadFile(googleCloudStorageJSONKeyPath)
		if err != nil {
			log.Fatal(err)
		}
		conf, err := google.JWTConfigFromJSON(
			kbts,
			storage.ScopeFullControl,
		)
		if err != nil {
			log.Fatal(err)
		}
		ctx := context.Background()
		aclient, err := storage.NewAdminClient(ctx, googleCloudProjectName, cloud.WithTokenSource(conf.TokenSource(ctx)))
		if err != nil {
			log.Fatal(err)
		}
		defer aclient.Close()

		if err := aclient.CreateBucket(context.Background(), googleCloudStorageBucketName, nil); err != nil {
			if !strings.Contains(err.Error(), "You already own this bucket. Please select another name") {
				log.Fatal(err)
			}
		}

		sctx := context.Background()
		sclient, err := storage.NewClient(sctx, cloud.WithTokenSource(conf.TokenSource(sctx)))
		if err != nil {
			log.Fatal(err)
		}
		defer sclient.Close()

		log.Printf("Uploading %s\n", csvResultPath)
		wc := sclient.Bucket(googleCloudStorageBucketName).Object(csvResultPath).NewWriter(context.Background())
		wc.ContentType = "text/plain"
		if _, err := wc.Write([]byte(txt)); err != nil {
			log.Fatal(err)
		}
		if err := wc.Close(); err != nil {
			log.Fatal(err)
		}
	}
	return fmt.Sprintf("\nSample in one second (unix latency throughput):\n%s", txt)
}
