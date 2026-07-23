package planner

import (
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestPlannerBuildsJobsForSupportedFormats(t *testing.T) {
	sourceRoot := t.TempDir()
	src := filepath.Join(sourceRoot, "album", "nested", "example.png")
	mustWritePNG(t, src)
	outputRoot := t.TempDir()

	p := New("jpg", []string{"jpg", "png", "webp"}, sourceRoot)
	job, ok, err := p.Plan(src, Dir(outputRoot))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected supported file")
	}
	if job.SourceFmt != "png" {
		t.Fatalf("expected source format png, got %q", job.SourceFmt)
	}
	if job.Metadata == nil {
		t.Fatal("expected metadata to be populated")
	}
	if job.Metadata.Name() != "example.png" {
		t.Fatalf("unexpected metadata name: %q", job.Metadata.Name())
	}
	if job.Metadata.Mode().Type()&fs.ModeType != 0 {
		t.Fatalf("expected regular file metadata, got mode %v", job.Metadata.Mode())
	}
	if job.TargetFmt != "jpeg" {
		t.Fatalf("expected normalized target format jpeg, got %q", job.TargetFmt)
	}
	want := filepath.Join(outputRoot, "album", "nested", "example.jpg")
	if job.TargetPath != want {
		t.Fatalf("unexpected target path: %s", job.TargetPath)
	}
}

func TestPlannerPreservesRelativeDirectoryStructure(t *testing.T) {
	sourceRoot := t.TempDir()
	src := filepath.Join(sourceRoot, "album", "nested", "example.png")
	mustWritePNG(t, src)
	outputRoot := t.TempDir()

	p := New("jpg", []string{"jpg", "png", "webp"}, sourceRoot)
	job, ok, err := p.Plan(src, Dir(outputRoot))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected supported file")
	}
	want := filepath.Join(outputRoot, "album", "nested", "example.jpg")
	if job.TargetPath != want {
		t.Fatalf("unexpected preserved target path: %s", job.TargetPath)
	}
}

func TestPlannerSkipsUnsupportedExtensions(t *testing.T) {
	p := New("jpg", []string{"jpg", "png", "webp"}, t.TempDir())
	src := filepath.Join(t.TempDir(), "example.gif")
	mustWriteFile(t, src)

	_, ok, err := p.Plan(src, Dir(t.TempDir()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected unsupported file to be skipped")
	}
}

func TestPlannerUsesProbeFallbackWhenExtensionUnsupported(t *testing.T) {
	p := New("jpg", []string{"jpg", "png", "webp"}, t.TempDir())
	src := filepath.Join(t.TempDir(), "example.bin")
	mustWriteFile(t, src)
	p.probe = func(string) (string, bool, error) {
		return "png", true, nil
	}

	job, ok, err := p.Plan(src, Dir(t.TempDir()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected probe fallback to accept the file")
	}
	if job.SourceFmt != "png" {
		t.Fatalf("expected probed source format png, got %q", job.SourceFmt)
	}
}

func TestPlannerDefaultsAcceptedFormatsWhenListEmpty(t *testing.T) {
	p := New("jpg", nil, t.TempDir())
	src := filepath.Join(t.TempDir(), "example.gif")
	mustWriteGIF(t, src)

	job, ok, err := p.Plan(src, Dir(t.TempDir()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected empty accepted list to default to supported inputs")
	}
	if job.SourceFmt != "gif" {
		t.Fatalf("expected gif source format, got %q", job.SourceFmt)
	}
}

func TestPlannerNormalizesExplicitOutputFilePathExtension(t *testing.T) {
	sourceRoot := t.TempDir()
	src := filepath.Join(sourceRoot, "example.png")
	mustWritePNG(t, src)
	outputFile := filepath.Join(t.TempDir(), "out", "example.png")

	p := New("jpg", []string{"jpg", "png", "webp"}, sourceRoot)
	job, ok, err := p.Plan(src, File(outputFile))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected supported file")
	}
	if job.TargetPath != replaceExtension(outputFile, "jpg") {
		t.Fatalf("unexpected target path: %s", job.TargetPath)
	}
}

func TestPlannerAddsCollisionSuffixForExistingTarget(t *testing.T) {
	sourceRoot := t.TempDir()
	src := filepath.Join(sourceRoot, "album", "nested", "example.png")
	mustWritePNG(t, src)
	outputRoot := t.TempDir()
	existing := filepath.Join(outputRoot, "album", "nested", "example.jpg")
	mustWriteFile(t, existing)

	p := New("jpg", []string{"jpg", "png", "webp"}, sourceRoot)
	job, ok, err := p.Plan(src, Dir(outputRoot))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected supported file")
	}
	if job.TargetPath == existing {
		t.Fatal("expected planner to avoid the colliding target path")
	}
	if job.TargetPath != filepath.Join(outputRoot, "album", "nested", "example-1.jpg") {
		t.Fatalf("unexpected collision-safe target path: %s", job.TargetPath)
	}
}

func TestPlannerAvoidsSameBatchCollisionAcrossPlans(t *testing.T) {
	sourceRoot := t.TempDir()
	src := filepath.Join(sourceRoot, "album", "nested", "example.png")
	mustWritePNG(t, src)
	outputRoot := t.TempDir()

	p := New("jpg", []string{"jpg", "png", "webp"}, sourceRoot)

	first, ok, err := p.Plan(src, Dir(outputRoot))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected supported file")
	}

	second, ok, err := p.Plan(src, Dir(outputRoot))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected supported file")
	}

	if first.TargetPath == second.TargetPath {
		t.Fatalf("expected different planned paths, both were %s", first.TargetPath)
	}
	if second.TargetPath != filepath.Join(outputRoot, "album", "nested", "example-1.jpg") {
		t.Fatalf("unexpected second planned path: %s", second.TargetPath)
	}
}

func TestPlannerAllocatesUniqueTargetPathsConcurrently(t *testing.T) {
	sourceRoot := t.TempDir()
	src := filepath.Join(sourceRoot, "album", "nested", "example.png")
	mustWritePNG(t, src)
	outputRoot := t.TempDir()

	p := New("jpg", []string{"jpg", "png", "webp"}, sourceRoot)

	var wg sync.WaitGroup
	wg.Add(2)

	results := make(chan string, 2)
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			job, ok, err := p.Plan(src, Dir(outputRoot))
			if err != nil {
				errs <- err
				return
			}
			if !ok {
				errs <- nil
				return
			}
			results <- job.TargetPath
		}()
	}

	wg.Wait()
	close(results)
	close(errs)

	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := make(map[string]struct{})
	for path := range results {
		got[path] = struct{}{}
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 unique target paths, got %v", got)
	}
	if _, ok := got[filepath.Join(outputRoot, "album", "nested", "example.jpg")]; !ok {
		t.Fatalf("expected base target path in results, got %v", got)
	}
	if _, ok := got[filepath.Join(outputRoot, "album", "nested", "example-1.jpg")]; !ok {
		t.Fatalf("expected collision-resolved target path in results, got %v", got)
	}
}

func TestPlannerExplicitDestinationKindsAreNotHeuristic(t *testing.T) {
	sourceRoot := t.TempDir()
	src := filepath.Join(sourceRoot, "album", "nested", "example.png")
	mustWritePNG(t, src)
	outputRoot := filepath.Join(t.TempDir(), "out.v2")

	p := New("jpg", []string{"jpg", "png", "webp"}, sourceRoot)
	dirJob, ok, err := p.Plan(src, Dir(outputRoot))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected supported file")
	}
	if dirJob.TargetPath != filepath.Join(outputRoot, "album", "nested", "example.jpg") {
		t.Fatalf("unexpected directory target path: %s", dirJob.TargetPath)
	}

	fileOutput := filepath.Join(t.TempDir(), "out.txt")
	fileJob, ok, err := p.Plan(src, File(fileOutput))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected supported file")
	}
	if fileJob.TargetPath != replaceExtension(fileOutput, "jpg") {
		t.Fatalf("unexpected explicit file target: %s", fileJob.TargetPath)
	}
}

func TestPlannerRejectsUnknownTargetFormat(t *testing.T) {
	p := New("jpgg", []string{"jpg", "png", "webp"}, t.TempDir())
	src := filepath.Join(t.TempDir(), "example.png")
	mustWritePNG(t, src)

	_, ok, err := p.Plan(src, Dir(t.TempDir()))
	if err == nil {
		t.Fatal("expected unknown target format to return an error")
	}
	if ok {
		t.Fatal("expected unknown target format to prevent job planning")
	}
	if !strings.Contains(err.Error(), "jpgg") {
		t.Fatalf("expected error to include original requested format, got %q", err.Error())
	}
}

func TestPlannerRejectsMislabeledKnownExtension(t *testing.T) {
	sourceRoot := t.TempDir()
	src := filepath.Join(sourceRoot, "example.png")
	mustWriteGIF(t, src)

	p := New("jpg", []string{"jpg", "png", "webp"}, sourceRoot)
	_, ok, err := p.Plan(src, Dir(t.TempDir()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected mislabeled png to be rejected by probe verification")
	}
}

func TestPlannerSurfacesProbeFailureForSupportedExtension(t *testing.T) {
	sourceRoot := t.TempDir()
	src := filepath.Join(sourceRoot, "broken.png")
	mustWriteFile(t, src)

	p := New("jpg", []string{"jpg", "png", "webp"}, sourceRoot)
	_, ok, err := p.Plan(src, Dir(t.TempDir()))
	if err == nil {
		t.Fatal("expected probe failure to return an error")
	}
	if ok {
		t.Fatal("expected probe failure to prevent job planning")
	}
}

func TestPlannerPreservesRelativeStructureWithRelativeSourceRoot(t *testing.T) {
	root := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})

	sourceRoot := filepath.Clean(filepath.Join(".", "inputs", "normalized"))
	src := filepath.Join(sourceRoot, "album", "nested", "example.png")
	mustWritePNG(t, src)
	outputRoot := filepath.Join(".", "out")

	p := New("jpg", []string{"jpg", "png", "webp"}, sourceRoot)
	job, ok, err := p.Plan(src, Dir(outputRoot))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected supported file")
	}
	want := filepath.Join(outputRoot, "album", "nested", "example.jpg")
	if job.TargetPath != want {
		t.Fatalf("unexpected target path: %s", job.TargetPath)
	}
}

func TestPlannerDoesNotReviveDisallowedKnownExtensionWithProbe(t *testing.T) {
	p := New("jpg", []string{"jpg", "png", "webp"}, t.TempDir())
	src := filepath.Join(t.TempDir(), "example.gif")
	mustWriteGIF(t, src)
	calls := 0
	p.probe = func(string) (string, bool, error) {
		calls++
		return "png", true, nil
	}

	_, ok, err := p.Plan(src, Dir(t.TempDir()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected known-but-disallowed extension to stay skipped")
	}
	if calls != 0 {
		t.Fatalf("expected probe to be skipped for disallowed extension, called %d times", calls)
	}
}

func TestPlannerUsesDefaultProbeForExtensionlessDetection(t *testing.T) {
	p := New("jpg", []string{"jpg", "png", "webp"}, t.TempDir())

	src := filepath.Join(t.TempDir(), "fixture.bin")
	mustWritePNG(t, src)

	job, ok, err := p.Plan(src, Dir(t.TempDir()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected default probe to identify the PNG fixture")
	}
	if job.SourceFmt != "png" {
		t.Fatalf("expected default probe to report png, got %q", job.SourceFmt)
	}
}

func TestPlannerUsesDefaultProbeForHeicHeifDetection(t *testing.T) {
	p := New("jpg", []string{"jpg", "png", "webp", "heic"}, t.TempDir())

	for _, brand := range []string{"heic", "heif"} {
		src := filepath.Join(t.TempDir(), brand+".bin")
		mustWriteHEIC(t, src, brand)

		job, ok, err := p.Plan(src, Dir(t.TempDir()))
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", brand, err)
		}
		if !ok {
			t.Fatalf("expected default probe to identify the %s fixture", brand)
		}
		if job.SourceFmt != "heic" {
			t.Fatalf("expected default probe to report heic for %s, got %q", brand, job.SourceFmt)
		}
	}
}

func mustWriteFile(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent directories: %v", err)
	}
	if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}

func mustWritePNG(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent directories: %v", err)
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create png fixture: %v", err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode png fixture: %v", err)
	}
}

func mustWriteGIF(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent directories: %v", err)
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create gif fixture: %v", err)
	}
	defer f.Close()

	img := image.NewPaletted(image.Rect(0, 0, 1, 1), color.Palette{color.White, color.Black})
	img.SetColorIndex(0, 0, 1)
	if err := gif.Encode(f, img, nil); err != nil {
		t.Fatalf("encode gif fixture: %v", err)
	}
}

func mustWriteHEIC(t *testing.T, path string, brand string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent directories: %v", err)
	}
	if brand != "heic" && brand != "heif" {
		t.Fatalf("unsupported heic brand for test: %s", brand)
	}
	data := []byte{
		0x00, 0x00, 0x00, 0x18,
		'f', 't', 'y', 'p',
		brand[0], brand[1], brand[2], brand[3],
		0x00, 0x00, 0x00, 0x00,
		brand[0], brand[1], brand[2], brand[3],
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write heic fixture: %v", err)
	}
}
