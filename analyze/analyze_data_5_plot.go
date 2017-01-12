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

func (all *allAggregatedData) draw(cfg PlotConfig, cols ...dataframe.Column) error {
	// frame now contains
	// AVG-LATENCY-MS-etcd-v3.1-go1.7.4, AVG-LATENCY-MS-zookeeper-r3.4.9-java8, AVG-LATENCY-MS-consul-v0.7.2-go1.7.4
	pl, err := plot.New()
	if err != nil {
		return err
	}
	pl.Title.Text = fmt.Sprintf("%s, %s", all.title, cfg.YAxis)
	pl.X.Label.Text = cfg.XAxis
	pl.Y.Label.Text = cfg.YAxis
	pl.Legend.Top = true

	var ps []plot.Plotter
	for i, col := range cols {
		pt, err := points(col)
		if err != nil {
			return err
		}

		l, err := plotter.NewLine(pt)
		if err != nil {
			return err
		}
		l.Color = getRGB(all.headerToLegend[col.Header()], i)
		l.Dashes = plotutil.Dashes(i)
		ps = append(ps, l)

		pl.Legend.Add(all.headerToLegend[col.Header()], l)
	}
	pl.Add(ps...)

	for _, outputPath := range cfg.OutputPathList {
		if err = pl.Save(plotWidth, plotHeight, outputPath); err != nil {
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
		n, _ := v.Number()
		pts[i].X = float64(i)
		pts[i].Y = n
	}
	return pts, nil
}

func getRGB(legend string, i int) color.Color {
	legend = strings.ToLower(strings.TrimSpace(legend))
	if strings.HasPrefix(legend, "etcd") {
		return color.RGBA{24, 90, 169, 255} // blue
	}
	if strings.HasPrefix(legend, "zookeeper") {
		return color.RGBA{38, 169, 24, 255} // green
	}
	if strings.HasPrefix(legend, "consul") {
		return color.RGBA{198, 53, 53, 255} // red
	}
	if strings.HasPrefix(legend, "zetcd") {
		return color.RGBA{251, 206, 0, 255} // yellow
	}
	if strings.HasPrefix(legend, "cetcd") {
		return color.RGBA{116, 24, 169, 255} // purple
	}
	return plotutil.Color(i)
}
