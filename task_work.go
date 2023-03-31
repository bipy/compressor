package main

import (
	"bytes"
	"compressor/utils"
	"errors"
	"github.com/disintegration/imaging"
	"github.com/samber/lo"
	"image"
)

func doTask(t *Task) error {
	img, err := imaging.Open(t.Input)
	if err != nil {
		return err
	}

	filename, err := utils.Touch(t.Output, touchMutex)
	if err != nil {
		return err
	}
	t.Output = filename

	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	if w == 0 || h == 0 {
		return errors.New("image size is 0")
	}

	var dst *image.NRGBA
	if w > config.MaxWidth && h > config.MaxHeight {
		if float64(config.MaxWidth)/float64(w) < float64(config.MaxHeight)/float64(h) {
			dst = imaging.Resize(img, config.MaxWidth, 0, imaging.Lanczos)
		} else {
			dst = imaging.Resize(img, 0, config.MaxHeight, imaging.Lanczos)
		}
	} else if w > config.MaxWidth {
		dst = imaging.Resize(img, config.MaxWidth, 0, imaging.Lanczos)
	} else if h > config.MaxHeight {
		dst = imaging.Resize(img, 0, config.MaxHeight, imaging.Lanczos)
	}

	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, lo.Ternary[image.Image](dst == nil, img, dst), imaging.JPEG, imaging.JPEGQuality(config.Quality))

	if err != nil {
		return err
	}
	t.Data = buf.Bytes()
	return nil
}

// compress job, multiple goroutine
func compress() {
	defer wg.Done()
	// get job from channel,
	// channel inCh will be closed by sender
	for t := range inCh {
		// check if success
		if err := doTask(&t); err != nil {
			// if failed, push to fail channel (multi-sender)
			t.Err = &err
			failCh <- t
			continue
		}
		outCh <- t
	}
}
