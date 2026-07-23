package reporter

import (
	"fmt"
	"testing"
	"time"

	"compressor/internal/model"
)

func TestReporter_Basic(t *testing.T) {
	r := New()

	r.AddDiscovered(10)
	r.Report(model.Result{
		Status:      model.StatusDone,
		InputBytes:  1024,
		OutputBytes: 512,
		Duration:    time.Second,
	})
	r.Report(model.Result{
		Status: model.StatusFailed,
		Err:    fmt.Errorf("test error"),
	})

	snap := r.Snapshot()

	if snap.Discovered != 10 {
		t.Fatalf("expected discovered=10, got %d", snap.Discovered)
	}
	if snap.Processed != 2 {
		t.Fatalf("expected processed=2, got %d", snap.Processed)
	}
	if snap.Successes != 1 {
		t.Fatalf("expected successes=1, got %d", snap.Successes)
	}
	if snap.Failures != 1 {
		t.Fatalf("expected failures=1, got %d", snap.Failures)
	}
	if snap.InputBytes != 1024 {
		t.Fatalf("expected inputBytes=1024, got %d", snap.InputBytes)
	}
	if snap.OutputBytes != 512 {
		t.Fatalf("expected outputBytes=512, got %d", snap.OutputBytes)
	}
}

func TestReporter_CompressionRatio(t *testing.T) {
	r := New()

	r.Report(model.Result{
		Status:      model.StatusDone,
		InputBytes:  1000,
		OutputBytes: 500,
	})

	snap := r.Snapshot()
	if snap.CompressionRatio != 0.5 {
		t.Fatalf("expected ratio=0.5, got %f", snap.CompressionRatio)
	}
}

func TestReporter_RecentJobs(t *testing.T) {
	r := New()
	r.maxRecent = 3

	for i := 0; i < 5; i++ {
		r.Report(model.Result{
			Status:     model.StatusDone,
			InputBytes: int64(i * 100),
		})
	}

	snap := r.Snapshot()
	if len(snap.RecentJobs) != 3 {
		t.Fatalf("expected 3 recent jobs, got %d", len(snap.RecentJobs))
	}
}

func TestReporter_Failures(t *testing.T) {
	r := New()

	r.Report(model.Result{Status: model.StatusDone})
	r.Report(model.Result{Status: model.StatusFailed, Err: fmt.Errorf("err1")})
	r.Report(model.Result{Status: model.StatusFailed, Err: fmt.Errorf("err2")})

	failures := r.Failures()
	if len(failures) != 2 {
		t.Fatalf("expected 2 failures, got %d", len(failures))
	}
}
