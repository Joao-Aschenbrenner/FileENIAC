package packer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Builder struct {
	excludes []string
}

type Result struct {
	ArchivePath string
	FileCount   int
	SizeBytes   int64
}

func NewBuilder(excludes []string) *Builder {
	return &Builder{
		excludes: excludes,
	}
}

func (b *Builder) Pack(sourceDir, outputPath string) (*Result, error) {
	file, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create archive: %w", err)
	}
	defer file.Close()

	gzw := gzip.NewWriter(file)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	var fileCount int
	var totalSize int64

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		if b.shouldExclude(relPath, info) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Name = relPath
		if info.IsDir() {
			header.Name += "/"
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			n, err := io.Copy(tw, f)
			f.Close()
			if err != nil {
				return err
			}
			totalSize += n
			fileCount++
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to pack directory: %w", err)
	}

	return &Result{
		ArchivePath: outputPath,
		FileCount:   fileCount,
		SizeBytes:   totalSize,
	}, nil
}

func (b *Builder) shouldExclude(relPath string, info os.FileInfo) bool {
	relPath = filepath.ToSlash(relPath)

	for _, pattern := range b.excludes {
		pattern = filepath.ToSlash(pattern)

		if pattern == ".git" && (relPath == ".git" || strings.HasPrefix(relPath, ".git/")) {
			return true
		}

		if strings.HasPrefix(pattern, ".git") && strings.Contains(relPath, ".git") {
			return true
		}

		if strings.HasPrefix(pattern, "*.") {
			ext := pattern[1:]
			if strings.HasSuffix(relPath, ext) {
				return true
			}
		}

		if pattern == relPath || strings.HasPrefix(relPath, pattern+"/") {
			return true
		}

		parts := strings.Split(relPath, "/")
		for _, part := range parts {
			if part == pattern {
				return true
			}
		}
	}

	return false
}

func (b *Builder) SetExcludes(excludes []string) {
	b.excludes = excludes
}
