package reporter

import (
	"fmt"
	"os"
	"time"

	"compressor/internal/model"
)

type PlainTextReporter struct {
	reporter *Reporter
	done     chan struct{}
}

func NewPlainText(r *Reporter) *PlainTextReporter {
	return &PlainTextReporter{
		reporter: r,
		done:     make(chan struct{}),
	}
}

func (p *PlainTextReporter) Start() {
	go p.loop()
}

func (p *PlainTextReporter) Stop() {
	close(p.done)
}

func (p *PlainTextReporter) loop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.done:
			p.printFinal()
			return
		case <-ticker.C:
			p.printProgress()
		}
	}
}

func (p *PlainTextReporter) printProgress() {
	snap := p.reporter.Snapshot()
	fmt.Fprintf(os.Stderr, "\r[%d/%d] ✓%d ✗%d | %.1fMB/s | %s",
		snap.Processed,
		snap.Discovered,
		snap.Successes,
		snap.Failures,
		snap.BytesPerSecond/(1024*1024),
		snap.Elapsed.Round(time.Second),
	)
}

func (p *PlainTextReporter) printFinal() {
	snap := p.reporter.Snapshot()
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Printf("[Summary] total=%d success=%d fail=%d\n", snap.Processed, snap.Successes, snap.Failures)
	fmt.Printf("[Summary] input=%.1fMB output=%.1fMB ratio=%.2f\n",
		float64(snap.InputBytes)/(1024*1024),
		float64(snap.OutputBytes)/(1024*1024),
		snap.CompressionRatio)
	fmt.Printf("[Summary] elapsed=%v throughput=%.1fMB/s\n",
		snap.Elapsed.Round(time.Millisecond),
		snap.BytesPerSecond/(1024*1024))

	if snap.Failures > 0 {
		fmt.Printf("\nFailures:\n")
		for _, f := range p.reporter.Failures() {
			fmt.Printf("  %s: %v\n", f.Job.SourcePath, f.Err)
		}
	}
}

func (p *PlainTextReporter) ReportJob(result model.Result) {
	snap := p.reporter.Snapshot()
	if result.Status == model.StatusDone {
		fmt.Fprintf(os.Stderr, "\r[%d] %s (%.1fKB -> %.1fKB, %v)\n",
			snap.Successes,
			result.Job.SourcePath,
			float64(result.InputBytes)/1024,
			float64(result.OutputBytes)/1024,
			result.Duration.Round(time.Millisecond))
	}
}

func (p *PlainTextReporter) Done() chan struct{} {
	return p.done
}
