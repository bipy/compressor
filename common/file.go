package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Touch file
func Touch(filename string, dirMutex *sync.Mutex, maxRetry int) (name string, err error) {
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)
	dirMutex.Lock()
	defer dirMutex.Unlock()
	for i := 1; i < maxRetry; i++ {
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
