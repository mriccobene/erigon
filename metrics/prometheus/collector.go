// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package prometheus

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/VictoriaMetrics/metrics"
)

var (
	typeGaugeTpl           = "# TYPE %s gauge\n"
	typeCounterTpl         = "# TYPE %s counter\n"
	typeSummaryTpl         = "# TYPE %s summary\n"
	keyValueTpl            = "%s %v\n\n"
	keyCounterTpl          = "%s %v\n"
	keyQuantileTagValueTpl = "%s {quantile=\"%s\"} %v\n"
)

// collector is a collection of byte buffers that aggregate Prometheus reports
// for different metric types.
type collector struct {
	buff *bytes.Buffer
}

// newCollector creates a new Prometheus metric aggregator.
func newCollector() *collector {
	return &collector{
		buff: &bytes.Buffer{},
	}
}

func (c *collector) addCounter(name string, m *metrics.Counter) {
	c.writeCounter(name, m.Get())
}

func (c *collector) addGauge(name string, m *metrics.Gauge) {
	c.writeGauge(name, m.Get())
}

func (c *collector) addFloatCounter(name string, m *metrics.FloatCounter) {
	c.writeGauge(name, m.Get())
}

func (c *collector) addHistogram(name string, m *metrics.Histogram) {
	c.buff.WriteString(fmt.Sprintf(typeSummaryTpl, name))

	c.writeSummarySum(name, fmt.Sprintf("%f", m.GetSum()))
	c.writeSummaryCounter(name, len(m.GetDecimalBuckets()))
	c.buff.WriteRune('\n')
}

func (c *collector) addTimer(name string, m *metrics.Summary) {
	pv := m.GetQuantiles()
	ps := m.GetQuantileValues()

	var sum float64 = 0
	c.buff.WriteString(fmt.Sprintf(typeSummaryTpl, name))
	for i := range pv {
		c.writeSummaryPercentile(name, strconv.FormatFloat(pv[i], 'f', -1, 64), ps[i])
		sum += ps[i]
	}

	c.writeSummaryTime(name, fmt.Sprintf("%f", m.GetTime().Seconds()))
	c.writeSummarySum(name, fmt.Sprintf("%f", sum))
	c.writeSummaryCounter(name, len(ps))
	c.buff.WriteRune('\n')
}

func (c *collector) writeGauge(name string, value interface{}) {
	//c.buff.WriteString(fmt.Sprintf(typeGaugeTpl, name))
	c.buff.WriteString(fmt.Sprintf(keyValueTpl, name, value))
}

func (c *collector) writeCounter(name string, value interface{}) {
	//c.buff.WriteString(fmt.Sprintf(typeCounterTpl, name))
	c.buff.WriteString(fmt.Sprintf(keyValueTpl, name, value))
}

func (c *collector) writeSummaryCounter(name string, value interface{}) {
	name = name + "_count"
	c.buff.WriteString(fmt.Sprintf(keyCounterTpl, name, value))
}

func (c *collector) writeSummaryPercentile(name, p string, value interface{}) {
	c.buff.WriteString(fmt.Sprintf(keyQuantileTagValueTpl, name, p, value))
}

func (c *collector) writeSummarySum(name string, value string) {
	name = name + "_sum"
	c.buff.WriteString(fmt.Sprintf(keyCounterTpl, name, value))
}

func (c *collector) writeSummaryTime(name string, value string) {
	name = name + "_time"
	c.buff.WriteString(fmt.Sprintf(keyCounterTpl, name, value))
}