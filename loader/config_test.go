package loader

import (
	"testing"
)

func TestCheckRejectsNonPositiveThreadCount(t *testing.T) {
	for _, tc := range []int{0, -1} {
		cfg := &Config{
			ThreadCount: tc,
			InputPath:   t.TempDir(),
			Quality:     90,
			MaxWidth:    1,
			MaxHeight:   1,
		}

		if cfg.Check() {
			t.Fatalf("expected thread count %d to be rejected", tc)
		}
	}
}
