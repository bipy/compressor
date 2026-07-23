package executor

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"compressor/internal/codec"
	"compressor/internal/model"
	"compressor/internal/transform"
	"compressor/internal/writer"
)

type Executor struct {
	workers int
	codec   *codec.Factory
	writer  *writer.Writer
	transformOpts transform.Options
	quality int
}

type Config struct {
	Workers      int
	Quality      int
	MaxWidth     int
	MaxHeight    int
}

func New(cfg Config) *Executor {
	workers := cfg.Workers
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	return &Executor{
		workers: workers,
		codec:   codec.NewFactory(),
		writer:  writer.New(),
		transformOpts: transform.Options{
			MaxWidth:        cfg.MaxWidth,
			MaxHeight:       cfg.MaxHeight,
			AutoOrientation: true,
		},
		quality: cfg.Quality,
	}
}

func (e *Executor) Execute(ctx context.Context, jobs <-chan model.Job) <-chan model.Result {
	results := make(chan model.Result, e.workers)

	go func() {
		defer close(results)

		var wg sync.WaitGroup
		sem := make(chan struct{}, e.workers)

		for job := range jobs {
			select {
			case <-ctx.Done():
				return
			case sem <- struct{}{}:
			}

			wg.Add(1)
			go func(j model.Job) {
				defer func() {
					<-sem
					wg.Done()
				}()

				result := e.processJob(ctx, j)
				select {
				case results <- result:
				case <-ctx.Done():
				}
			}(job)
		}

		wg.Wait()
	}()

	return results
}

func (e *Executor) processJob(ctx context.Context, job model.Job) (result model.Result) {
	result.Job = job
	result.StartedAt = time.Now()
	defer func() {
		result.EndedAt = time.Now()
		result.Duration = result.EndedAt.Sub(result.StartedAt)
		if r := recover(); r != nil {
			result.Status = model.StatusFailed
			result.Err = fmt.Errorf("panic: %v", r)
		}
	}()

	select {
	case <-ctx.Done():
		result.Status = model.StatusFailed
		result.Err = ctx.Err()
		return result
	default:
	}

	raw, err := os.ReadFile(job.SourcePath)
	if err != nil {
		result.Status = model.StatusFailed
		result.Err = fmt.Errorf("read: %w", err)
		return result
	}
	result.InputBytes = int64(len(raw))

	img, _, err := e.codec.Decode(bytesReader(raw))
	if err != nil {
		result.Status = model.StatusFailed
		result.Err = fmt.Errorf("decode: %w", err)
		return result
	}

	img, err = transform.Apply(img, e.transformOpts)
	if err != nil {
		result.Status = model.StatusFailed
		result.Err = fmt.Errorf("transform: %w", err)
		return result
	}

	encoded, err := e.codec.EncodeToBytes(img, job.TargetFmt, codec.EncodeOptions{Quality: e.quality})
	if err != nil {
		result.Status = model.StatusFailed
		result.Err = fmt.Errorf("encode: %w", err)
		return result
	}

	if err := e.writer.Write(job.TargetPath, encoded); err != nil {
		result.Status = model.StatusFailed
		result.Err = fmt.Errorf("write: %w", err)
		return result
	}

	result.OutputBytes = int64(len(encoded))
	result.Status = model.StatusDone
	return result
}

type bytesReaderWrapper struct {
	data []byte
	pos  int
}

func bytesReader(data []byte) *bytesReaderWrapper {
	return &bytesReaderWrapper{data: data}
}

func (r *bytesReaderWrapper) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("EOF")
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
