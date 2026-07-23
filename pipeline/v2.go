package pipeline

import (
	"context"
	"strings"
	"time"

	"compressor/internal/discover"
	"compressor/internal/executor"
	"compressor/internal/model"
	"compressor/internal/planner"
	"compressor/internal/reporter"
	"compressor/internal/tui"
	"compressor/loader"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func RunV2(cfg *loader.Config) error {
	start := time.Now()

	accepted := strings.Fields(cfg.AcceptedInputFormat)
	dest := planner.Dir(cfg.OutputPath)
	if cfg.SingleFileMode {
		dest = planner.File(cfg.OutputPath)
	}

	p := planner.New(cfg.OutputType, accepted, cfg.InputPath)

	exec := executor.New(executor.Config{
		Workers:   cfg.ThreadCount,
		Quality:   cfg.Quality,
		MaxWidth:  cfg.MaxWidth,
		MaxHeight: cfg.MaxHeight,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := reporter.New()
	jobs := make(chan model.Job, cfg.ThreadCount*2)
	results := exec.Execute(ctx, jobs)

	// Start TUI or plain text reporter
	done := make(chan struct{})
	var plainReporter *reporter.PlainTextReporter
	var pProgram *tea.Program

	if tui.IsTTY() {
		model := tui.NewModel(r, done)
		pProgram = tea.NewProgram(model)
		go func() {
			if _, err := pProgram.Run(); err != nil {
				log.Errorf("[TUI] error: %v", err)
			}
		}()
	} else {
		plainReporter = reporter.NewPlainText(r)
		plainReporter.Start()
	}

	// Discovery goroutine
	go func() {
		defer close(jobs)
		for d := range discover.Discover(cfg.InputPath) {
			if d.Err != nil {
				log.Warnf("[Discover] error: %v", d.Err)
				continue
			}
			r.AddDiscovered(1)
			job, ok, err := p.Plan(d.Path, dest)
			if err != nil {
				log.Warnf("[Planner] error: %v path=%s", err, d.Path)
				continue
			}
			if !ok {
				continue
			}
			select {
			case jobs <- job:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Collect results
	for result := range results {
		r.Report(result)
		if plainReporter != nil {
			plainReporter.ReportJob(result)
		}
	}

	elapsed := time.Since(start)

	// Stop reporters
	if plainReporter != nil {
		plainReporter.Stop()
	} else {
		close(done)
		time.Sleep(100 * time.Millisecond) // Let TUI render final state
	}

	snap := r.Snapshot()
	log.Infof("[Summary] total=%d success=%d fail=%d", snap.Processed, snap.Successes, snap.Failures)
	log.Infof("[Summary] input=%.1fMB output=%.1fMB ratio=%.2f",
		float64(snap.InputBytes)/(1024*1024),
		float64(snap.OutputBytes)/(1024*1024),
		snap.CompressionRatio)
	log.Infof("[Summary] elapsed=%v throughput=%.1fMB/s",
		elapsed.Round(time.Millisecond),
		snap.BytesPerSecond/(1024*1024))

	return nil
}
