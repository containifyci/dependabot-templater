package search

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type UniqueStringSlice struct {
	elements []string
	unqiue   map[string]bool
}

func (u *UniqueStringSlice) Add(s string) {
	if u.unqiue[NormalizePath(s)] {
		return // Already in the map
	}
	u.elements = append(u.elements, NormalizePath(s))
	u.unqiue[NormalizePath(s)] = true
}

func NormalizePath(path string) string {
	pdir := path
	if strings.HasPrefix(pdir, "../") {
		for strings.HasPrefix(pdir, "../") {
			pdir = pdir[3:]
		}
		folders := strings.Split(pdir, "/")
		pdir = filepath.Join(folders[1:]...)
	}
	return strings.TrimPrefix(filepath.Dir(pdir), "./")
}

func SearchForString(dir string, target string) ([]string, error) {
	foundFiles := UniqueStringSlice{
		unqiue: make(map[string]bool),
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.Contains(path, ".terraform") {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".tf") {
			return nil
		}
		if contains(path, target) {
			foundFiles.Add(path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return foundFiles.elements, nil
}

func SearchForFolder(dir string, targets ...string) ([]string, error) {
	foundFiles := UniqueStringSlice{
		unqiue: make(map[string]bool),
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, target := range targets {
			if strings.HasSuffix(strings.ToLower(filepath.Dir(path)), strings.ToLower(target)) {
				foundFiles.Add(path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return foundFiles.elements, nil
}

func SearchForFiles(dir string, targets ...string) ([]string, error) {
	foundFiles := UniqueStringSlice{
		unqiue: make(map[string]bool),
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, target := range targets {
			if strings.EqualFold(filepath.Base(path), target) {
				foundFiles.Add(path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return foundFiles.elements, nil
}

func contains(file string, target string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Printf("failed to close file: %s", err)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil || n == 0 {
			break
		}
		if strings.Contains(string(buf[:n]), target) {
			return true
		}
	}
	return false
}
