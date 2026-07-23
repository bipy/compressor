package reporter

import (
	"sync"
	"sync/atomic"
	"time"

	"compressor/internal/model"
)

type Reporter struct {
	mu          sync.RWMutex
	startedAt   time.Time
	discovered  atomic.Int64
	processed   atomic.Int64
	successes   atomic.Int64
	failures    atomic.Int64
	inputBytes  atomic.Int64
	outputBytes atomic.Int64
	recentJobs  []model.Result
	maxRecent   int
	failures_detail []model.Result
}

func New() *Reporter {
	return &Reporter{
		startedAt:  time.Now(),
		maxRecent:  20,
		recentJobs: make([]model.Result, 0, 20),
	}
}

func (r *Reporter) AddDiscovered(n int) {
	r.discovered.Add(int64(n))
}

func (r *Reporter) Report(result model.Result) {
	r.processed.Add(1)
	r.inputBytes.Add(result.InputBytes)
	r.outputBytes.Add(result.OutputBytes)

	switch result.Status {
	case model.StatusDone:
		r.successes.Add(1)
	case model.StatusFailed:
		r.failures.Add(1)
		r.mu.Lock()
		r.failures_detail = append(r.failures_detail, result)
		r.mu.Unlock()
	}

	r.mu.Lock()
	if len(r.recentJobs) >= r.maxRecent {
		r.recentJobs = r.recentJobs[1:]
	}
	r.recentJobs = append(r.recentJobs, result)
	r.mu.Unlock()
}

func (r *Reporter) Snapshot() Snapshot {
	r.mu.RLock()
	recent := make([]model.Result, len(r.recentJobs))
	copy(recent, r.recentJobs)
	r.mu.RUnlock()

	elapsed := time.Since(r.startedAt)
	processed := r.processed.Load()
	inputB := r.inputBytes.Load()
	outputB := r.outputBytes.Load()

	var bps, jps float64
	if elapsed.Seconds() > 0 {
		bps = float64(inputB) / elapsed.Seconds()
		jps = float64(processed) / elapsed.Seconds()
	}

	var ratio float64
	if inputB > 0 {
		ratio = float64(outputB) / float64(inputB)
	}

	return Snapshot{
		Discovered:     r.discovered.Load(),
		Processed:      processed,
		Successes:      r.successes.Load(),
		Failures:       r.failures.Load(),
		InputBytes:     inputB,
		OutputBytes:    outputB,
		Elapsed:        elapsed,
		BytesPerSecond: bps,
		JobsPerSecond:  jps,
		CompressionRatio: ratio,
		RecentJobs:     recent,
	}
}

func (r *Reporter) Failures() []model.Result {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.Result, len(r.failures_detail))
	copy(out, r.failures_detail)
	return out
}

type Snapshot struct {
	Discovered      int64
	Processed       int64
	Successes       int64
	Failures        int64
	InputBytes      int64
	OutputBytes     int64
	Elapsed         time.Duration
	BytesPerSecond  float64
	JobsPerSecond   float64
	CompressionRatio float64
	RecentJobs      []model.Result
}
