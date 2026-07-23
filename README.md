<div align="center">

# 🗜️ Compressor

### Batch image compression that just works

<br />

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-4caf50?style=flat-square)](LICENSE)
[![Go Report](https://goreportcard.com/badge/github.com/bipy/compressor?style=flat-square)](https://goreportcard.com/report/github.com/bipy/compressor)

[English](README.md) · [中文](README_zh_CN.md)

<br />

```
compressor -i ~/Photos -t webp -q 80
```

Point it at a folder. It handles the rest.

</div>

---

## What it does

Compressor takes a pile of images and makes them smaller. It reads JPEG, PNG, GIF, BMP, TIFF, WebP, and HEIC files, then outputs JPEG, PNG, or WebP — your choice.

It runs in parallel, uses your CPU fully by default, and streams files through the pipeline so processing kicks off before the scan finishes. No waiting around.

## Install

```bash
go install github.com/bipy/compressor@latest
```

<details>
<summary>Prefer building from source?</summary>

```bash
git clone https://github.com/bipy/compressor.git
cd compressor
go build
```

</details>

## Try it

```bash
# Everything in a folder
compressor -i ~/Photos

# One file, convert to WebP
compressor -i photo.png -o photo.webp

# Resize to fit 1920px wide, quality 75
compressor -i ~/Photos -t webp -width 1920 -q 75

# Slower but gentler on your machine
compressor -i ~/Photos -j 2
```

## All the knobs

| Flag | Default | What it does |
|:-----|:--------|:-------------|
| `-i` | *(required)* | Input file or directory |
| `-o` | auto-generated | Where to put the output |
| `-t` | `jpg` | Output format — `jpg`, `png`, or `webp` |
| `-q` | `90` | Quality, from 1 to 100 |
| `-j` | all CPU cores | Number of worker threads |
| `-accept` | everything | Which input formats to process |
| `-width` | keep original | Max width in pixels |
| `-height` | keep original | Max height in pixels |

## Under the hood

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ Discover │───▶│  Plan    │───▶│ Execute  │───▶│  Write   │───▶│  Report  │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘
```

**Discover** walks the input directory and streams file paths into the pipeline — no upfront scan delay.

**Plan** detects each file's format, figures out where the output goes, and handles naming conflicts automatically.

**Execute** hands jobs to a worker pool. Each worker decodes, resizes if needed, encodes, and writes — completely independently.

**Write** uses atomic writes (write to temp, then rename). If anything fails, no partial files are left behind.

**Report** aggregates metrics and feeds them to the live TUI dashboard — or plain text if you're piping output.

## Terminal UI

When you run Compressor in a terminal, you get a live dashboard:

```
Compressor
Total: 142 | ✓ 138 | ✗ 4 | Elapsed: 12s
Input: 2.1GB | Output: 856MB | Ratio: 0.41 | Speed: 178.5MB/s
[████████████████████░░░░░░░░░░] 67.2%

Recent:
  ✓ IMG_2024.jpg (2.1MB)
  ✓ IMG_2025.jpg (1.8MB)
  ✓ IMG_2026.jpg (2.3MB)

q: quit │ p: pause │ v: failures │ ?: help
```

Piping to a file? Running in CI? It detects that and switches to plain text automatically.

## Project layout

```
internal/
├── codec/       encode and decode — JPEG, PNG, WebP, GIF, BMP, TIFF, HEIC
├── discover/    streaming file walker
├── executor/    worker pool with panic recovery
├── model/       Job, Result, Summary types
├── planner/     format detection, path resolution, collision handling
├── reporter/    metrics aggregation
├── transform/   resize and auto-orientation
├── tui/         terminal dashboard (bubbletea)
└── writer/      atomic file writes
```

## Performance

Benchmarks on Apple M4 (single-threaded, per operation):

| Operation | Time per image | Allocated memory | What it means |
|:----------|---------------:|-----------------:|:--------------|
| Encode JPEG | 271 μs | 7.2 KB | ~3,700 images/sec |
| Encode PNG | 539 μs | 852 KB | ~1,850 images/sec (PNG is heavier) |
| Decode JPEG | 152 μs | 83 KB | ~6,500 images/sec — decode is faster than encode |
| Scan 1K files | 33 μs | 7.6 KB | 1000 files scanned in 33ms — discovery is nearly free |

> **Latency** = time to process one image. Lower is better.
> **Allocated memory** = heap memory used per operation. Lower means less GC pressure.
> These are single-worker numbers. Actual throughput scales linearly with `-j` threads — 8 cores ≈ 8× the images/sec.

## License

[MIT](LICENSE) — use it however you like.
