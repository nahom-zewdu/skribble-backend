// internal/metrics/latency.go
// This file defines the LatencyTracker struct and related functions for tracking latency metrics in the Skribble backend.
// It maintains a fixed-size list of latency samples and provides functions to add new latency measurements and calculate statistics such as minimum, 50th percentile, 95th percentile, 99th percentile, and maximum latency.

package metrics

import (
	"sort"
	"sync"
)

type LatencyTracker struct {
	mu sync.Mutex

	samples []int64
	maxSize int
}

var Latency = &LatencyTracker{
	maxSize: 5000,
}

func AddLatency(ms int64) {
	Latency.mu.Lock()
	defer Latency.mu.Unlock()

	if len(Latency.samples) >= Latency.maxSize {
		Latency.samples =
			Latency.samples[1:]
	}

	Latency.samples =
		append(
			Latency.samples,
			ms,
		)
}

func LatencyStats() map[string]int64 {
	Latency.mu.Lock()
	defer Latency.mu.Unlock()

	if len(Latency.samples) == 0 {
		return map[string]int64{}
	}

	values :=
		append(
			[]int64{},
			Latency.samples...,
		)

	sort.Slice(values,
		func(i, j int) bool {
			return values[i] < values[j]
		},
	)

	n := len(values)

	return map[string]int64{
		"min": values[0],
		"p50": values[n*50/100],
		"p95": values[n*95/100],
		"p99": values[n*99/100],
		"max": values[n-1],
	}
}
