package swallow

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const retention = 2 * time.Hour

// prune removes logs older than the retention period and afterwards any
// emptied origin directory, never the swallow dir itself. It is best-effort:
// failures never abort the run.
func prune(dir string) {
	cutoff := time.Now().Add(-retention)

	var subdirs []string
	_ = filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.IsDir() {
			if path != dir {
				subdirs = append(subdirs, path)
			}
			return nil
		}
		if info, err := entry.Info(); err == nil && info.ModTime().Before(cutoff) {
			_ = os.Remove(path)
		}
		return nil
	})

	sort.Slice(subdirs, func(i, j int) bool { return len(subdirs[i]) > len(subdirs[j]) })
	for _, subdir := range subdirs {
		_ = os.Remove(subdir)
	}
}
