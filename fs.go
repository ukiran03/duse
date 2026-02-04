package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

type FileEntry struct {
	Name  string
	Size  int64
	IsDir bool
}

func traverseFs(root string, depth int) []*FileEntry {
	var entries []*FileEntry

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing: %v [%v]\n", path, err)
			return fs.SkipDir
		}
		// Avoid processing the root itself
		if path == root {
			return nil
		}
		relPath, _ := filepath.Rel(root, path)
		relDepth := splitPathLength(relPath)

		if relDepth > depth {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		entry := &FileEntry{
			Name:  filepath.ToSlash(relPath), // shows depth significance
			IsDir: d.IsDir(),
		}
		if d.IsDir() {
			size, err := concurrnetDirSize(path)
			if err != nil {
				fmt.Println("Error getting directory size:", err)
				return fs.SkipDir
			}
			entry.Name += "/"
			entry.Size = size
		} else {
			info, err := d.Info()
			if err != nil {
				return nil
			}
			entry.Size = info.Size()
		}
		entries = append(entries, entry)
		return nil
	})
	if err != nil {
		fmt.Printf("Walk finished with error: %v\n", err)
	}
	return entries
}

// concurrnetDirSize produces the size of the directory at the path by
// recursively visiting all sub-dirs and files once (hopefully).
func concurrnetDirSize(path string) (int64, error) {
	var total atomic.Int64
	var wg sync.WaitGroup
	sema := make(chan struct{}, 2*runtime.NumCPU())

	var walker func(string)
	walker = func(p string) {
		sema <- struct{}{}
		entries, err := os.ReadDir(p)
		<-sema

		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				// NOTE: Are we really skipping
				fmt.Fprintf(os.Stderr, "Permission denied, skipping [%v]\n", p)
			} else {
				fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", p, err)
			}
			return
		}

		for _, entry := range entries {
			if entry.IsDir() {
				wg.Go(func() {
					walker(filepath.Join(p, entry.Name()))
				})
			} else {
				info, err := entry.Info()
				if err == nil {
					total.Add(info.Size())
				}
			}
		}
	}
	wg.Go(func() {
		walker(path)
	})
	wg.Wait()
	return total.Load(), nil
}

// splitPathLength tells how deep are we in the path
func splitPathLength(p string) int {
	p = filepath.Clean(p)
	if p == "." || p == "/" || p == "" {
		return 0
	}
	return len(strings.Split(p, string(os.PathSeparator)))
}
