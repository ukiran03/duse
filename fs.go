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
			size, err := concurrentDirSize(path)
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

// concurrentDirSize returns the total size of all files in the directory tree
// rooted at path. It skips directories it cannot read (e.g., permission denied)
// and returns the sum even if some errors occurred.
func concurrentDirSize(path string) (int64, error) {
	var total atomic.Int64
	var wg sync.WaitGroup
	// Limit concurrent directory reads (I/O bound)
	sema := make(chan struct{}, 2*runtime.NumCPU())

	var walker func(string) error
	walker = func(p string) error {
		sema <- struct{}{}        // acquire
		defer func() { <-sema }() // defer release

		entries, err := os.ReadDir(p)
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				fmt.Fprintf(os.Stderr, "Permission denied, skipping [%s]\n", p)
				return nil // continue with siblings
			}
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", p, err)
			return err
		}

		for _, entry := range entries {
			fullpath := filepath.Join(p, entry.Name())
			if entry.IsDir() {
				wg.Go(func() {
					// ignore error for now (or propagate via channel)
					_ = walker(fullpath)
				})
			} else {
				info, err := entry.Info()
				if err == nil {
					total.Add(info.Size())
				}
			}
		}
		return nil
	}

	wg.Go(func() {
		_ = walker(path)
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
