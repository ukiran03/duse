package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
)

func splitPath(p string) []string {
	p = filepath.Clean(p)
	if p == "." || p == "/" {
		return nil
	}
	var parts []string
	for {
		dir, file := filepath.Split(p)
		if file != "" {
			parts = append(parts, file) // O(1) average case
		}
		// Stop if we reached the root or the end of a relative path
		if dir == p || dir == "" || dir == "." || dir == "/" {
			break
		}
		p = filepath.Clean(dir)
	}
	return parts
}

func ConcurrnetDirSize(path string) (int64, error) {
	var total atomic.Int64
	var wg sync.WaitGroup
	sema := make(chan struct{}, 2*runtime.NumCPU())

	var walker func(string)
	walker = func(p string) {
		sema <- struct{}{}
		entries, err := os.ReadDir(p)
		<-sema
		if err != nil {
			log.Printf("Error reading %s: %v\n", p, err)
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
		if path == root {
			return nil
		}
		relPath, _ := filepath.Rel(root, path)
		relDepth := len(splitPath(relPath))

		if relDepth > depth {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		entry := &FileEntry{}
		entry.Name = d.Name()
		entry.IsDir = d.IsDir()
		if d.IsDir() {
			size, err := ConcurrnetDirSize(path)
			if err != nil {
				fmt.Println("Error getting directory size:", err)
				return fs.SkipDir
			}
			entry.Size = size
			entries = append(entries, entry)

			// Stop WalkDir from going into this folder because we
			// already calculated its size!
			return fs.SkipDir
		} else {
			info, err := d.Info()
			if err != nil {
				return nil // Skip this specific file if info fails
			}
			entry.Size = info.Size()
			entries = append(entries, entry)
		}

		return nil
	})
	if err != nil {
		fmt.Println("traverseFs error:", err)
	}
	return entries
}
