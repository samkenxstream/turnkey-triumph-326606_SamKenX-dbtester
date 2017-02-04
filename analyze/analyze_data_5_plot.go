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

package analyze

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/gyuho/dataframe"
)

var (
	plotWidth  = 12 * vg.Inch
	plotHeight = 8 * vg.Inch
)

func init() {
	plot.DefaultFont = "Helvetica"
	plotter.DefaultLineStyle.Width = vg.Points(1.5)
	plotter.DefaultGlyphStyle.Radius = vg.Points(2.0)
}

// PlotConfig defines what to plot.
type PlotConfig struct {
	Column         string   `yaml:"column"`
	XAxis          string   `yaml:"x_axis"`
	YAxis          string   `yaml:"y_axis"`
	OutputPathList []string `yaml:"output_path_list"`
}

type pair struct {
	x dataframe.Column
	y dataframe.Column
}

type triplet struct {
	x      dataframe.Column
	minCol dataframe.Column
	avgCol dataframe.Column
	maxCol dataframe.Column
}

func (all *allAggregatedData) draw(cfg PlotConfig, pairs ...pair) error {
	// frame now contains
	// AVG-LATENCY-MS-etcd-v3.1-go1.7.4, AVG-LATENCY-MS-zookeeper-r3.4.9-java8, AVG-LATENCY-MS-consul-v0.7.2-go1.7.4
	plt, err := plot.New()
	if err != nil {
		return err
	}
	plt.Title.Text = fmt.Sprintf("%s, %s", all.title, cfg.YAxis)
	plt.X.Label.Text = cfg.XAxis
	plt.Y.Label.Text = cfg.YAxis
	plt.Legend.Top = true

	var ps []plot.Plotter
	for i, p := range pairs {
		pt, err := points(p.y)
		if err != nil {
			return err
		}

		l, err := plotter.NewLine(pt)
		if err != nil {
			return err
		}
		l.Color = getRGB(all.headerToLegend[p.y.Header()], i)
		l.Dashes = plotutil.Dashes(i)
		ps = append(ps, l)

		plt.Legend.Add(all.headerToLegend[p.y.Header()], l)
	}
	plt.Add(ps...)

	for _, outputPath := range cfg.OutputPathList {
		if err = plt.Save(plotWidth, plotHeight, outputPath); err != nil {
			return err
		}
	}
	return nil
}

func (all *allAggregatedData) drawXY(cfg PlotConfig, pairs ...pair) error {
	// frame now contains
	// KEYS-DB-TAG-X, AVG-LATENCY-MS-DB-TAG-Y, ...
	plt, err := plot.New()
	if err != nil {
		return err
	}
	plt.Title.Text = fmt.Sprintf("%s, %s", all.title, cfg.YAxis)
	plt.X.Label.Text = cfg.XAxis
	plt.Y.Label.Text = cfg.YAxis
	plt.Legend.Top = true

	var ps []plot.Plotter
	for i, p := range pairs {
		pt, err := pointsXY(p.x, p.y)
		if err != nil {
			return err
		}

		l, err := plotter.NewLine(pt)
		if err != nil {
			return err
		}
		l.Color = getRGB(all.headerToLegend[p.y.Header()], i)
		l.Dashes = plotutil.Dashes(i)
		ps = append(ps, l)

		plt.Legend.Add(all.headerToLegend[p.y.Header()], l)
	}
	plt.Add(ps...)

	for _, outputPath := range cfg.OutputPathList {
		if err = plt.Save(plotWidth, plotHeight, outputPath); err != nil {
			return err
		}
	}
	return nil
}

func (all *allAggregatedData) drawXYWithErrorPoints(cfg PlotConfig, triplets ...triplet) error {
	// frame now contains
	// KEYS-DB-TAG-X, MIN-LATENCY-MS-DB-TAG-Y, AVG-LATENCY-MS-DB-TAG-Y, MAX-LATENCY-MS-DB-TAG-Y, ...
	plt, err := plot.New()
	if err != nil {
		return err
	}
	plt.Title.Text = fmt.Sprintf("%s, %s", all.title, cfg.YAxis)
	plt.X.Label.Text = cfg.XAxis
	plt.Y.Label.Text = cfg.YAxis
	plt.Legend.Top = true

	var ps []plot.Plotter
	for i, triplet := range triplets {
		{
			pt, err := pointsXY(triplet.x, triplet.minCol)
			if err != nil {
				return err
			}
			l, err := plotter.NewLine(pt)
			if err != nil {
				return err
			}
			l.Color = getRGBII(all.headerToLegend[triplet.avgCol.Header()], i)
			l.Dashes = plotutil.Dashes(i)
			ps = append(ps, l)
			plt.Legend.Add(all.headerToLegend[triplet.avgCol.Header()]+" MIN", l)
		}
		{
			pt, err := pointsXY(triplet.x, triplet.avgCol)
			if err != nil {
				return err
			}
			l, err := plotter.NewLine(pt)
			if err != nil {
				return err
			}
			l.Color = getRGB(all.headerToLegend[triplet.avgCol.Header()], i)
			l.Dashes = plotutil.Dashes(i)
			ps = append(ps, l)
			plt.Legend.Add(all.headerToLegend[triplet.avgCol.Header()], l)
		}
		{
			pt, err := pointsXY(triplet.x, triplet.maxCol)
			if err != nil {
				return err
			}
			l, err := plotter.NewLine(pt)
			if err != nil {
				return err
			}
			l.Color = getRGBIII(all.headerToLegend[triplet.avgCol.Header()], i)
			l.Dashes = plotutil.Dashes(i)
			ps = append(ps, l)
			plt.Legend.Add(all.headerToLegend[triplet.avgCol.Header()]+" MAX", l)
		}
	}
	plt.Add(ps...)

	for _, outputPath := range cfg.OutputPathList {
		if err = plt.Save(plotWidth, plotHeight, outputPath); err != nil {
			return err
		}
	}
	return nil
}

func points(col dataframe.Column) (plotter.XYs, error) {
	bv, ok := col.BackNonNil()
	if !ok {
		return nil, fmt.Errorf("BackNonNil not found")
	}
	rowN, ok := col.FindLast(bv)
	if !ok {
		return nil, fmt.Errorf("not found %v", bv)
	}
	pts := make(plotter.XYs, rowN)
	for i := range pts {
		v, err := col.Value(i)
		if err != nil {
			return nil, err
		}
		n, _ := v.Float64()
		pts[i].X = float64(i)
		pts[i].Y = n
	}
	return pts, nil
}

func pointsXY(colX, colY dataframe.Column) (plotter.XYs, error) {
	bv, ok := colX.BackNonNil()
	if !ok {
		return nil, fmt.Errorf("BackNonNil not found")
	}
	rowN, ok := colX.FindLast(bv)
	if !ok {
		return nil, fmt.Errorf("not found %v", bv)
	}
	pts := make(plotter.XYs, rowN)
	for i := range pts {
		vx, err := colX.Value(i)
		if err != nil {
			return nil, err
		}
		x, _ := vx.Float64()

		vy, err := colY.Value(i)
		if err != nil {
			return nil, err
		}
		y, _ := vy.Float64()

		pts[i].X = x
		pts[i].Y = y
	}
	return pts, nil
}

func getRGB(legend string, i int) color.Color {
	tag := makeTag(legend)
	if strings.HasPrefix(tag, "etcd") {
		return color.RGBA{24, 90, 169, 255} // blue
	}
	if strings.HasPrefix(tag, "zookeeper") {
		return color.RGBA{38, 169, 24, 255} // green
	}
	if strings.HasPrefix(tag, "consul") {
		return color.RGBA{198, 53, 53, 255} // red
	}
	if strings.HasPrefix(tag, "zetcd") {
		return color.RGBA{251, 206, 0, 255} // yellow
	}
	if strings.HasPrefix(tag, "cetcd") {
		return color.RGBA{116, 24, 169, 255} // purple
	}
	return plotutil.Color(i)
}

func getRGBII(legend string, i int) color.Color {
	tag := makeTag(legend)
	if strings.HasPrefix(tag, "etcd") {
		return color.RGBA{37, 29, 191, 255} // deep-blue
	}
	if strings.HasPrefix(tag, "zookeeper") {
		return color.RGBA{7, 64, 35, 255} // deep-green
	}
	if strings.HasPrefix(tag, "consul") {
		return color.RGBA{212, 8, 46, 255} // deep-red
	}
	if strings.HasPrefix(tag, "zetcd") {
		return color.RGBA{229, 255, 0, 255} // deep-yellow
	}
	if strings.HasPrefix(tag, "cetcd") {
		return color.RGBA{255, 0, 251, 255} // deep-purple
	}
	return plotutil.Color(i)
}

func getRGBIII(legend string, i int) color.Color {
	tag := makeTag(legend)
	if strings.HasPrefix(tag, "etcd") {
		return color.RGBA{129, 212, 247, 255} // light-blue
	}
	if strings.HasPrefix(tag, "zookeeper") {
		return color.RGBA{129, 247, 152, 255} // light-green
	}
	if strings.HasPrefix(tag, "consul") {
		return color.RGBA{247, 156, 156, 255} // light-red
	}
	if strings.HasPrefix(tag, "zetcd") {
		return color.RGBA{245, 247, 166, 255} // light-yellow
	}
	if strings.HasPrefix(tag, "cetcd") {
		return color.RGBA{247, 166, 238, 255} // light-purple
	}
	return plotutil.Color(i)
}
