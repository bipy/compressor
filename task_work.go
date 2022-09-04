package main

import (
	"bytes"
	"compressor/utils"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
)

func doTask(t *Task) error {
	file, err := os.Open(t.Input)
	if err != nil {
		return err
	}

	filename, err := utils.Touch(t.Output)
	if err != nil {
		return err
	}
	t.Output = filename

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, config.JpegQuality)
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
			t.Data = []byte(err.Error())
			failCh <- t
			continue
		}
		outCh <- t
	}
}
