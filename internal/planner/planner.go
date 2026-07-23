package planner

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"compressor/internal/model"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type Planner struct {
	requestedTarget string
	targetFmt       string
	targetExt       string
	sourceRoot      string
	accepted        map[string]struct{}
	probe           ProbeFunc
	targetOK        bool
	mu              sync.Mutex
	reserved        map[string]struct{}
}

type ProbeFunc func(path string) (format string, ok bool, err error)

type DestinationKind int

const (
	DestinationDirectory DestinationKind = iota
	DestinationFile
)

type Destination struct {
	Path string
	Kind DestinationKind
}

func Dir(path string) Destination {
	return Destination{Path: path, Kind: DestinationDirectory}
}

func File(path string) Destination {
	return Destination{Path: path, Kind: DestinationFile}
}

func New(outputFormat string, accepted []string, sourceRoot string) *Planner {
	targetFmt, targetExt := normalizeOutputFormat(outputFormat)

	acceptedSet := make(map[string]struct{}, len(accepted))
	if len(accepted) == 0 {
		for format := range supportedInputFormats {
			acceptedSet[format] = struct{}{}
		}
	} else {
		for _, format := range accepted {
			if canonical, _, ok := normalizeFormat(format); ok && supportedInputFormats[canonical] {
				acceptedSet[canonical] = struct{}{}
			}
		}
	}

	return &Planner{
		requestedTarget: outputFormat,
		targetFmt:       targetFmt,
		targetExt:       targetExt,
		sourceRoot:      normalizePath(sourceRoot),
		accepted:        acceptedSet,
		probe:           defaultProbe,
		targetOK:        supportedOutputFormats[targetFmt],
		reserved:        make(map[string]struct{}),
	}
}

func (p *Planner) Plan(srcPath string, dest Destination) (model.Job, bool, error) {
	if !p.targetOK {
		return model.Job{}, false, fmt.Errorf("unsupported output format %q", p.requestedTarget)
	}

	sourceFmt, ok, err := p.detectSourceFormat(srcPath)
	if err != nil {
		return model.Job{}, false, err
	}
	if !ok {
		return model.Job{}, false, nil
	}

	metadata, err := os.Stat(srcPath)
	if err != nil {
		return model.Job{}, false, err
	}

	targetPath, err := p.allocateTargetPath(dest, srcPath)
	if err != nil {
		return model.Job{}, false, err
	}

	return model.Job{
		SourcePath: srcPath,
		TargetPath: targetPath,
		SourceFmt:  sourceFmt,
		TargetFmt:  p.targetFmt,
		Metadata:   metadata,
	}, true, nil
}

func (p *Planner) detectSourceFormat(srcPath string) (string, bool, error) {
	ext := filepath.Ext(srcPath)
	format, _, ok := normalizeFormat(ext)
	if !ok || !supportedInputFormats[format] {
		return p.probeSourceFormat(srcPath)
	}
	if _, ok := p.accepted[format]; ok {
		probed, ok, err := p.probeSourceFormat(srcPath)
		if err != nil {
			return "", false, err
		}
		if !ok || probed != format {
			return "", false, nil
		}
		return format, true, nil
	}
	return "", false, nil
}

func normalizeOutputFormat(format string) (canonical string, display string) {
	canonical, display, ok := normalizeFormat(format)
	if !ok || !supportedOutputFormats[canonical] {
		return "", ""
	}
	return canonical, display
}

func normalizeFormat(format string) (canonical string, display string, ok bool) {
	trimmed := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(format), "."))
	if trimmed == "" {
		return "", "", false
	}

	switch trimmed {
	case "jpg", "jpeg":
		return "jpeg", trimmed, true
	case "tif", "tiff":
		return "tiff", trimmed, true
	case "heif", "heic":
		return "heic", trimmed, true
	default:
		return trimmed, trimmed, true
	}
}

func (p *Planner) probeSourceFormat(srcPath string) (string, bool, error) {
	format, ok, err := p.probe(srcPath)
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	canonical, _, ok := normalizeFormat(format)
	if !ok || !supportedInputFormats[canonical] {
		return "", false, nil
	}
	if _, ok := p.accepted[canonical]; !ok {
		return "", false, nil
	}
	return canonical, true, nil
}

func (p *Planner) allocateTargetPath(dest Destination, srcPath string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	candidate := p.targetCandidate(dest, srcPath)
	if !p.isReservedOrExistsLocked(candidate) {
		p.reserved[candidate] = struct{}{}
		return candidate, nil
	}

	dir := filepath.Dir(candidate)
	ext := filepath.Ext(candidate)
	name := strings.TrimSuffix(filepath.Base(candidate), ext)
	for i := 1; ; i++ {
		next := filepath.Join(dir, fmt.Sprintf("%s-%d%s", name, i, ext))
		if !p.isReservedOrExistsLocked(next) {
			p.reserved[next] = struct{}{}
			return next, nil
		}
	}
}

func (p *Planner) targetCandidate(dest Destination, srcPath string) string {
	switch dest.Kind {
	case DestinationDirectory:
		rel := p.relativeSourcePath(srcPath)
		candidate := filepath.Join(dest.Path, rel)
		return replaceExtension(candidate, p.targetExt)
	case DestinationFile:
		return replaceExtension(dest.Path, p.targetExt)
	default:
		return replaceExtension(dest.Path, p.targetExt)
	}
}

var supportedInputFormats = map[string]bool{
	"jpeg": true,
	"png":  true,
	"gif":  true,
	"bmp":  true,
	"tiff": true,
	"webp": true,
	"heic": true,
}

var supportedOutputFormats = map[string]bool{
	"jpeg": true,
	"png":  true,
	"webp": true,
}

func replaceExtension(path, ext string) string {
	base := strings.TrimSuffix(path, filepath.Ext(path))
	return base + "." + ext
}

func (p *Planner) relativeSourcePath(srcPath string) string {
	normalizedSrc := normalizePath(srcPath)
	if p.sourceRoot == "" {
		if filepath.IsAbs(normalizedSrc) {
			return filepath.Base(normalizedSrc)
		}
		return filepath.Clean(normalizedSrc)
	}

	rel, err := filepath.Rel(p.sourceRoot, normalizedSrc)
	if err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return rel
	}
	if filepath.IsAbs(normalizedSrc) {
		return filepath.Base(normalizedSrc)
	}
	return filepath.Clean(normalizedSrc)
}

func normalizePath(path string) string {
	if path == "" {
		return ""
	}
	clean := filepath.Clean(path)
	if abs, err := filepath.Abs(clean); err == nil {
		return abs
	}
	return clean
}

func defaultProbe(path string) (string, bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", false, err
	}
	defer f.Close()

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", false, err
	}
	data, err := io.ReadAll(io.LimitReader(f, 64<<10))
	if err != nil {
		return "", false, err
	}
	if format, ok := sniffHEICFormat(data); ok {
		return format, true, nil
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || cfg.Width <= 0 || cfg.Height <= 0 {
		return "", false, err
	}
	return strings.ToLower(format), true, nil
}

func (p *Planner) isReservedOrExistsLocked(path string) bool {
	if _, ok := p.reserved[path]; ok {
		return true
	}
	_, err := os.Stat(path)
	return err == nil
}

func sniffHEICFormat(data []byte) (string, bool) {
	if len(data) < 12 {
		return "", false
	}

	fourcc := string(data[4:8])
	if fourcc != "ftyp" {
		return "", false
	}

	switch string(data[8:12]) {
	case "heic", "heif":
		return "heic", true
	default:
		return "", false
	}
}
