package packer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuilder_Pack(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	os.MkdirAll(sourceDir, 0755)

	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("test content"), 0644)
	os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(sourceDir, "subdir", "nested.txt"), []byte("nested"), 0644)

	builder := NewBuilder([]string{".git", "node_modules", "*.bak"})
	outputPath := filepath.Join(tmpDir, "output.tar.gz")

	result, err := builder.Pack(sourceDir, outputPath)
	if err != nil {
		t.Fatalf("Pack failed: %v", err)
	}

	if result.FileCount != 2 {
		t.Errorf("expected 2 files, got %d", result.FileCount)
	}

	if result.SizeBytes == 0 {
		t.Error("size should be greater than 0")
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("output archive should exist")
	}
}

func TestBuilder_ShouldExclude(t *testing.T) {
	builder := NewBuilder([]string{".git", "node_modules", "*.log", "vendor"})

	tests := []struct {
		path     string
		isDir    bool
		expected bool
	}{
		{".git/config", false, true},
		{".git", true, true},
		{"node_modules/package/index.js", false, true},
		{"vendor/autoload.php", false, true},
		{"debug.log", false, true},
		{"app/Providers/AppServiceProvider.php", false, false},
		{"routes/web.php", false, false},
	}

	for _, tt := range tests {
		info, _ := os.Stat(tt.path)
		if info == nil {
			info = os.FileInfo(nil)
		}

		got := builder.shouldExclude(tt.path, info)
		if got != tt.expected {
			t.Errorf("shouldExclude(%q) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}

func TestBuilder_SetExcludes(t *testing.T) {
	builder := NewBuilder([]string{".git"})

	if len(builder.excludes) != 1 {
		t.Error("expected 1 exclude pattern")
	}

	builder.SetExcludes([]string{".git", "node_modules"})

	if len(builder.excludes) != 2 {
		t.Error("expected 2 exclude patterns after SetExcludes")
	}
}