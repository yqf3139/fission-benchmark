// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package requester

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

const (
	barChar = "∎"
)

type Report struct {
	AvgTotal float64
	Fastest  float64
	Slowest  float64
	Average  float64
	Rps      float64

	trace bool //if trace is set, the following fields will be filled

	AvgConn   float64
	AvgDns    float64
	AvgReq    float64
	AvgRes    float64
	AvgDelay  float64
	ConnLats  []float64
	DnsLats   []float64
	ReqLats   []float64
	ResLats   []float64
	DelayLats []float64

	results chan *result
	Total   time.Duration

	ErrorDist      map[string]int
	StatusCodeDist map[string]int
	Lats           []float64
	SizeTotal      int64

	output string

	w io.Writer
}

type ReportWrapper struct {
	Report *Report
	Error  error
}

func NewReport(w io.Writer, size int, results chan *result, output string, total time.Duration, trace bool) *Report {
	return &Report{
		output:         output,
		results:        results,
		Total:          total,
		trace:          trace,
		StatusCodeDist: make(map[string]int),
		ErrorDist:      make(map[string]int),
		w:              w,
	}
}

func (r *Report) Finalize() {
	for res := range r.results {
		if res.err != nil {
			r.ErrorDist[res.err.Error()]++
		} else {
			r.Lats = append(r.Lats, res.duration.Seconds())
			r.AvgTotal += res.duration.Seconds()
			if r.trace {
				r.AvgConn += res.connDuration.Seconds()
				r.AvgDelay += res.delayDuration.Seconds()
				r.AvgDns += res.dnsDuration.Seconds()
				r.AvgReq += res.reqDuration.Seconds()
				r.AvgRes += res.resDuration.Seconds()
				r.ConnLats = append(r.ConnLats, res.connDuration.Seconds())
				r.DnsLats = append(r.DnsLats, res.dnsDuration.Seconds())
				r.ReqLats = append(r.ReqLats, res.reqDuration.Seconds())
				r.DelayLats = append(r.DelayLats, res.delayDuration.Seconds())
				r.ResLats = append(r.ResLats, res.resDuration.Seconds())
			}
			r.StatusCodeDist[fmt.Sprint(res.statusCode)]++
			if res.contentLength > 0 {
				r.SizeTotal += res.contentLength
			}
		}
	}
	r.Rps = float64(len(r.Lats)) / r.Total.Seconds()
	r.Average = r.AvgTotal / float64(len(r.Lats))
	if r.trace {
		r.AvgConn = r.AvgConn / float64(len(r.Lats))
		r.AvgDelay = r.AvgDelay / float64(len(r.Lats))
		r.AvgDns = r.AvgDns / float64(len(r.Lats))
		r.AvgReq = r.AvgReq / float64(len(r.Lats))
		r.AvgRes = r.AvgRes / float64(len(r.Lats))
	}
}

func (r *Report) printCSV() {
	r.printf("response-time")
	if r.trace {
		r.printf(",DNS+dialup,DNS,Request-write,Response-delay,Response-read")
	}
	r.printf("\n")
	for i, val := range r.Lats {
		r.printf("%4.4f", val)
		if r.trace {
			r.printf(",%4.4f,%4.4f,%4.4f,%4.4f,%4.4f", r.ConnLats[i], r.DnsLats[i], r.ReqLats[i],
				r.DelayLats[i], r.ResLats[i])
		}
		r.printf("\n")
	}
}

func (r *Report) Print(enableHistogram, enableLatencies bool) {
	if r.output == "csv" {
		r.printCSV()
		return
	}

	if len(r.Lats) > 0 {
		sort.Float64s(r.Lats)
		r.Fastest = r.Lats[0]
		r.Slowest = r.Lats[len(r.Lats)-1]
		r.printf("\nSummary:\n")
		r.printf("  Total:\t%4.4f secs\n", r.Total.Seconds())
		r.printf("  Slowest:\t%4.4f secs\n", r.Slowest)
		r.printf("  Fastest:\t%4.4f secs\n", r.Fastest)
		r.printf("  Average:\t%4.4f secs\n", r.Average)
		r.printf("  Requests/sec:\t%4.4f\n", r.Rps)
		if r.SizeTotal > 0 {
			r.printf("  Total data:\t%d bytes\n", r.SizeTotal)
			r.printf("  Size/request:\t%d bytes\n", r.SizeTotal/int64(len(r.Lats)))
		}
		if r.trace {
			r.printf("\nDetailed Report:\n")
			r.printSection("DNS+dialup", r.AvgConn, r.ConnLats)
			r.printSection("DNS-lookup", r.AvgDns, r.DnsLats)
			r.printSection("Request Write", r.AvgReq, r.ReqLats)
			r.printSection("Response Wait", r.AvgDelay, r.DelayLats)
			r.printSection("Response Read", r.AvgRes, r.ResLats)
		}
		r.printStatusCodes()
		if enableHistogram {
			r.printHistogram()
		}
		if enableLatencies {
			r.printLatencies()
		}
	}

	if len(r.ErrorDist) > 0 {
		r.printErrors()
	}
}

// printSection prints details for http-trace fields
func (r *Report) printSection(tag string, avg float64, lats []float64) {
	sort.Float64s(lats)
	fastest, slowest := lats[0], lats[len(lats)-1]
	r.printf("\n\t%s:\n", tag)
	r.printf("  \t\tAverage:\t%4.4f secs\n", avg)
	r.printf("  \t\tFastest:\t%4.4f secs\n", fastest)
	r.printf("  \t\tSlowest:\t%4.4f secs\n", slowest)
}

// Prints percentile latencies.
func (r *Report) printLatencies() {
	pctls := []int{10, 25, 50, 75, 90, 95, 99}
	data := make([]float64, len(pctls))
	j := 0
	for i := 0; i < len(r.Lats) && j < len(pctls); i++ {
		current := i * 100 / len(r.Lats)
		if current >= pctls[j] {
			data[j] = r.Lats[i]
			j++
		}
	}
	r.printf("\nLatency distribution:\n")
	for i := 0; i < len(pctls); i++ {
		if data[i] > 0 {
			r.printf("  %v%% in %4.4f secs\n", pctls[i], data[i])
		}
	}
}

func (r *Report) printHistogram() {
	bc := 10
	buckets := make([]float64, bc+1)
	counts := make([]int, bc+1)
	bs := (r.Slowest - r.Fastest) / float64(bc)
	for i := 0; i < bc; i++ {
		buckets[i] = r.Fastest + bs*float64(i)
	}
	buckets[bc] = r.Slowest
	var bi int
	var max int
	for i := 0; i < len(r.Lats); {
		if r.Lats[i] <= buckets[bi] {
			i++
			counts[bi]++
			if max < counts[bi] {
				max = counts[bi]
			}
		} else if bi < len(buckets)-1 {
			bi++
		}
	}
	r.printf("\nResponse time histogram:\n")
	for i := 0; i < len(buckets); i++ {
		// Normalize bar lengths.
		var barLen int
		if max > 0 {
			barLen = (counts[i]*40 + max/2) / max
		}
		r.printf("  %4.3f [%v]\t|%v\n", buckets[i], counts[i], strings.Repeat(barChar, barLen))
	}
}

// Prints status code distribution.
func (r *Report) printStatusCodes() {
	r.printf("\nStatus code distribution:\n")
	for code, num := range r.StatusCodeDist {
		r.printf("  [%s]\t%d responses\n", code, num)
	}
}

func (r *Report) printErrors() {
	r.printf("\nError distribution:\n")
	for err, num := range r.ErrorDist {
		r.printf("  [%d]\t%s\n", num, err)
	}
}

func (r *Report) printf(s string, v ...interface{}) {
	fmt.Fprintf(r.w, s, v...)
}
