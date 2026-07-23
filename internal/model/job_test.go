package model

import (
	"testing"
	"time"
)

func TestCoreRuntimeTypes(t *testing.T) {
	_ = Job{
		SourcePath: "source.jpg",
		TargetPath: "target.png",
		SourceFmt:  "jpg",
		TargetFmt:  "png",
		Metadata:   nil,
	}

	_ = Result{
		Job:       Job{},
		Status:    StatusRunning,
		Duration:  2 * time.Second,
		StartedAt: time.Unix(0, 0),
		EndedAt:   time.Unix(2, 0),
	}

	s := Summary{InputBytes: 200, OutputBytes: 100}
	if got := s.CompressionRatio(); got != 0.5 {
		t.Fatalf("expected 0.5, got %v", got)
	}

	if got := (Summary{}).CompressionRatio(); got != 0 {
		t.Fatalf("expected 0 compression ratio when input bytes are zero, got %v", got)
	}

	_ = Summary{
		Discovered:  1,
		Processed:   1,
		Successes:   1,
		Failures:    0,
		InputBytes:  200,
		OutputBytes: 100,
		BytesPerSecond:  50,
		JobsPerSecond:   2,
	}
}
