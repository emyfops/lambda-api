package util

import (
	"os"
	"path/filepath"
)

// WalkFiles walks through the files under the specified root directory and calls the provided walk function for each file.
// It skips directories and only processes regular files.
func WalkFiles(root string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		return walkFn(path, info, err)
	})
}
