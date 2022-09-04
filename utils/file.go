package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Touch crate & rename file
func Touch(filename string) (name string, err error) {
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)
	for i := 1; i < 128; i++ {
		_, err = os.Stat(filename)
		// if file exist
		if err == nil {
			filename = fmt.Sprintf("%s (%d)%s", baseName, i, ext)
			continue
		}
		if os.IsNotExist(err) {
			// touch
			_, err = os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0644)
			if err == nil {
				return filename, nil
			}
		}
	}
	if err != nil {
		return
	}
	return "", errors.New("untouchable")
}
