package discover

import (
	"io/fs"
	"os"
	"path/filepath"
)

// Result is a streamed discovery item.
// When Err is non-nil, discovery stopped with an error.
type Result struct {
	Path string
	Err  error
}

// Discover streams regular file paths found under root.
func Discover(root string) <-chan Result {
	out := make(chan Result)

	go func() {
		defer close(out)

		info, err := os.Lstat(root)
		if err != nil {
			out <- Result{Err: err}
			return
		}

		if info.Mode().IsRegular() {
			out <- Result{Path: root}
			return
		}

		if !info.IsDir() {
			return
		}

		if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			switch typ := d.Type(); {
			case typ.IsRegular():
				out <- Result{Path: path}
				return nil
			case typ != 0:
				return nil
			default:
				info, err := d.Info()
				if err != nil {
					return nil
				}
				if !info.Mode().IsRegular() {
					return nil
				}

				out <- Result{Path: path}
				return nil
			}
		}); err != nil {
			out <- Result{Err: err}
		}
	}()

	return out
}
