package model

import (
	"io/fs"
	"time"
)

type Status string

const (
	StatusQueued  Status = "queued"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusFailed  Status = "failed"
)

type Job struct {
	SourcePath string
	TargetPath string
	SourceFmt   string
	TargetFmt   string
	Metadata    fs.FileInfo
}

type Result struct {
	Job         Job
	Status      Status
	Err         error
	InputBytes  int64
	OutputBytes int64
	Duration    time.Duration
	StartedAt   time.Time
	EndedAt     time.Time
}

type Summary struct {
	Discovered     int
	Processed      int
	Successes      int
	Failures       int
	InputBytes     int64
	OutputBytes    int64
	BytesPerSecond float64
	JobsPerSecond   float64
}

func (s Summary) CompressionRatio() float64 {
	if s.InputBytes == 0 {
		return 0
	}
	return float64(s.OutputBytes) / float64(s.InputBytes)
}
