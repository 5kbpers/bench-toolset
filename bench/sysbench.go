package bench

import (
	"fmt"

	"github.com/5kbpers/bench-toolset/metrics"
	"github.com/5kbpers/bench-toolset/workload"
)

func EvalSysbenchRecords(records []*workload.Record, intervalSecs int, warmupSecs int, cutTailSecs int, kNumber int, percent float64) []*Result {
	recordsMap := groupRecords(records)
	if intervalSecs > 0 {
		for t, rs := range recordsMap {
			recordsMap[t] = splitRecordChunks(rs[warmupSecs:len(rs)-cutTailSecs], intervalSecs)
			fmt.Printf("Aggregate records with interval %d, got %d records.\n", intervalSecs, len(recordsMap[t]))
		}
	}
	results := make([]*Result, 0, 6*len(recordsMap))
	for _, rs := range recordsMap {
		counts := make(metrics.TaggedValueSlice, len(rs))
		avgLats := make(metrics.TaggedValueSlice, len(rs))
		p95Lats := make(metrics.TaggedValueSlice, 0, len(rs))
		p99Lats := make(metrics.TaggedValueSlice, 0, len(rs))
		for i, r := range rs {
			counts[i] = metrics.WithTag(r.Count, r.Tag)
			avgLats[i] = metrics.WithTag(r.AvgLatInMs, r.Tag)
			if r.P95LatInMs > 0 {
				p95Lats = append(p95Lats, metrics.WithTag(r.P95LatInMs, r.Tag))
			}
			if r.P99LatInMs > 0 {
				p99Lats = append(p99Lats, metrics.WithTag(r.P99LatInMs, r.Tag))
			}
		}
		results = append(results, calculateResults("", "tps", counts, kNumber, percent, "")...)
		results = append(results, calculateResults("", "avg-lat", avgLats, kNumber, percent, "ms")...)
		if len(p95Lats) > 0 {
			results = append(results, calculateResults("", "p95-lat", p95Lats, kNumber, percent, "ms")...)
		}
		if len(p99Lats) > 0 {
			results = append(results, calculateResults("", "p99-lat", p99Lats, kNumber, percent, "ms")...)
		}
	}
	return results
}
