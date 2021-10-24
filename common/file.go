package common

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Touch file
func Touch(filename *string, dirMutex *sync.Mutex, maxRetry int) bool {
	touched := false
	ext := filepath.Ext(*filename)
	baseName := (*filename)[:len(*filename)-len(ext)]
	dirMutex.Lock()
	defer dirMutex.Unlock()
	for i := 1; i < maxRetry; i++ {
		_, err := os.Stat(*filename)
		// if file exist
		if err == nil {
			*filename = fmt.Sprintf("%s (%d)%s", baseName, i, ext)
		} else if os.IsNotExist(err) {
			// touch
			_, err := os.OpenFile(*filename, os.O_RDONLY|os.O_CREATE, 0644)
			if err == nil {
				touched = true
			}
			break
		}
	}
	return touched
}
