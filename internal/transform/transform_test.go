package transform

import (
	"image"
	"testing"
)

func TestApplyLeavesImageUntouchedWhenNoResizeNeeded(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 100, 50))

	got, err := Apply(src, Options{MaxWidth: 200, MaxHeight: 200})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Bounds() != src.Bounds() {
		t.Fatalf("unexpected bounds: %v", got.Bounds())
	}
}

func TestApplyResizesImageToFitWithinBounds(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 100, 50))

	got, err := Apply(src, Options{MaxWidth: 60, MaxHeight: 60})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Bounds() != image.Rect(0, 0, 60, 30) {
		t.Fatalf("unexpected bounds: %v", got.Bounds())
	}
}

func TestApplyRejectsInvalidImageBounds(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 0, 10))

	if _, err := Apply(src, Options{MaxWidth: 60, MaxHeight: 60}); err == nil {
		t.Fatalf("expected invalid bounds error")
	}
}

func TestApplyAutoOrientsImageBeforeResize(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 40, 20))

	got, err := Apply(src, Options{AutoOrientation: true, Orientation: 6})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Bounds() != image.Rect(0, 0, 20, 40) {
		t.Fatalf("unexpected bounds: %v", got.Bounds())
	}
}
