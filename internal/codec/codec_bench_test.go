package codec

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

func createBenchImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	for y := 0; y < 200; y++ {
		for x := 0; x < 200; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 128, A: 255})
		}
	}
	return img
}

func BenchmarkEncode_JPEG(b *testing.B) {
	f := NewFactory()
	img := createBenchImage()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		f.Encode(&buf, img, "jpeg", EncodeOptions{Quality: 80})
	}
}

func BenchmarkEncode_PNG(b *testing.B) {
	f := NewFactory()
	img := createBenchImage()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		f.Encode(&buf, img, "png", EncodeOptions{})
	}
}

func BenchmarkDecode_JPEG(b *testing.B) {
	f := NewFactory()
	img := createBenchImage()

	var encoded bytes.Buffer
	f.Encode(&encoded, img, "jpeg", EncodeOptions{Quality: 80})
	data := encoded.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Decode(bytes.NewReader(data))
	}
}
