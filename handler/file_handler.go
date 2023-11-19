package handler

import (
	"compressor/constant"
	"compressor/utils"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/samber/lo"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileHandler struct {
	Mu                  *sync.Mutex
	BasePath            string
	OutputPath          string
	AcceptedInputFormat []string
	OutputFormat        string
	Count               int
}

func (h *FileHandler) Travel(singleFile bool) (fileList []string, dstList []string) {
	if singleFile {
		return []string{h.BasePath}, []string{strings.TrimSuffix(h.BasePath, filepath.Ext(h.BasePath)) + h.OutputFormat}
	}

	// find all images
	err := filepath.WalkDir(h.BasePath, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			log.Errorf("[FileHandler] walk error. path=%s err=%v", path, e)
			return e
		}
		if !d.IsDir() {
			if ext := strings.ToLower(filepath.Ext(d.Name()))[1:]; lo.Contains(h.AcceptedInputFormat, ext) {
				fileList = append(fileList, path)
			}
		}
		return nil
	})
	if err != nil {
		log.Errorf("[FileHandler] travel failed. err=%v", err)
		return
	}

	dstList = lo.Map[string, string](fileList, func(item string, index int) string {
		newPath := filepath.Join(h.OutputPath, strings.TrimPrefix(item, h.BasePath))
		return strings.TrimSuffix(newPath, filepath.Ext(newPath)) + h.OutputFormat
	})
	return
}

func (h *FileHandler) Write(path string, data []byte) bool {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		log.Errorf("[FileHandler] cannot write data to file. path=%s err=%v", path, err)
		return false
	}
	return true
}

func (h *FileHandler) Touch(path string) (newPath string, ok bool) {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	ext := filepath.Ext(path)
	baseName := strings.TrimSuffix(path, ext)
	newPath = path

	// generate new file path
	for i := 0; i < constant.TouchMaxRetryTimes; i++ {
		if i > 0 {
			newPath = fmt.Sprintf("%s-%d%s", baseName, i, ext)
		}
		// check if path exists
		if utils.IsFileNotExist(newPath) {
			err := utils.CreateFile(newPath)
			if err != nil {
				log.Errorf("[FileHandler] cannot create new file. path=%s err=%v", newPath, err)
				return "", false
			}
			log.Debugf("[FileHandler] create new file. path=%s", newPath)
			return newPath, true
		}
	}
	log.Errorf("[FileHandler] create new file failed: exceed max retry limit. path=%s limit=%d",
		path, constant.TouchMaxRetryTimes)
	return "", false
}
