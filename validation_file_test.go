// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestFileRule(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	subFile := filepath.Join(subDir, "subfile.txt")
	if err := os.WriteFile(subFile, []byte("sub content"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		value      any
		searchPath []string
		expectErr  bool
	}{
		// Valid files - absolute paths
		{"valid absolute path", testFile, nil, false},
		{"valid absolute path in subdir", subFile, nil, false},

		// Invalid files - absolute paths
		{"non-existent absolute path", filepath.Join(tempDir, "missing.txt"), nil, true},
		{"directory instead of file", tempDir, nil, false}, // os.Stat succeeds for directories

		// Valid files - with search paths
		{"valid relative with search path", "test.txt", []string{tempDir}, false},
		{"valid relative in subdir", "subfile.txt", []string{subDir}, false},
		{"valid relative multiple search paths", "subfile.txt", []string{tempDir, subDir}, false},
		{"valid absolute ignores search path", testFile, []string{"/wrong/path"}, false},

		// Invalid files - with search paths
		{"invalid relative no search path", "test.txt", nil, true},
		{"invalid relative wrong search path", "test.txt", []string{"/wrong/path"}, true},
		{"missing file with search path", "missing.txt", []string{tempDir}, true},

		// Edge cases
		{"empty string", "", nil, true},
		{"nil value", nil, nil, true},
		{"non-string value", 123, nil, true},
		{"pointer to valid file", fileStrPtr(testFile), nil, false},
		{"pointer to invalid file", fileStrPtr("missing.txt"), nil, true},
		{"empty search paths", "test.txt", []string{}, true},

		// Directory handling with search paths
		{"directory with search path", "subdir", []string{tempDir}, true}, // directories are rejected in findFile
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := File(tt.searchPath...).Validate(tt.value)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestFileRule_When(t *testing.T) {
	err := File().When(false).Validate("nonexistent.txt")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = File().When(true).Validate("nonexistent.txt")
	if err == nil {
		t.Error("expected error when condition is true and file doesn't exist")
	}
}

func TestFileRule_Error(t *testing.T) {
	customMsg := "custom file error"
	err := File().Error(customMsg).Validate("nonexistent.txt")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestFileRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := File().ErrorObject(customErr).Validate("nonexistent.txt")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestFileRule_ComplexSearchPaths(t *testing.T) {
	// Create complex directory structure
	tempDir := t.TempDir()

	dirs := []string{
		filepath.Join(tempDir, "config"),
		filepath.Join(tempDir, "data"),
		filepath.Join(tempDir, "templates"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}

	// Create files in different directories
	files := map[string]string{
		"config.yml":    filepath.Join(dirs[0], "config.yml"),
		"database.json": filepath.Join(dirs[1], "database.json"),
		"email.tmpl":    filepath.Join(dirs[2], "email.tmpl"),
	}

	for _, path := range files {
		if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Test finding files in different search paths
	tests := []struct {
		name       string
		filename   string
		searchDirs []string
		shouldFind bool
	}{
		{"find config in config dir", "config.yml", []string{dirs[0]}, true},
		{"find database in data dir", "database.json", []string{dirs[1]}, true},
		{"find template in templates dir", "email.tmpl", []string{dirs[2]}, true},
		{"find with multiple search paths", "config.yml", dirs, true},
		{"not found in wrong dir", "config.yml", []string{dirs[1]}, false},
		{"find in second search path", "database.json", []string{dirs[0], dirs[1]}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := File(tt.searchDirs...).Validate(tt.filename)
			if tt.shouldFind {
				if err != nil {
					t.Errorf("expected to find file, got error: %v", err)
				}
			} else {
				if err == nil {
					t.Error("expected not to find file, but validation passed")
				}
			}
		})
	}
}

func TestFindFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		filePath    string
		searchPaths []string
		expected    string
	}{
		{"absolute path exists", testFile, []string{"/wrong"}, testFile},
		{"absolute path not exists", "/nonexistent/file.txt", []string{tempDir}, ""},
		{"relative path found", "test.txt", []string{tempDir}, testFile},
		{"relative path not found", "test.txt", []string{"/wrong"}, ""},
		{"empty search paths", "test.txt", []string{}, ""},
		{"multiple search paths", "test.txt", []string{"/wrong", tempDir}, testFile},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findFile(tt.filePath, tt.searchPaths)
			if tt.expected == "" {
				if result != "" {
					t.Errorf("expected empty string, got %q", result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			}
		})
	}
}

func BenchmarkFileValidation(b *testing.B) {
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "bench.txt")
	if err := os.WriteFile(testFile, []byte("benchmark"), 0644); err != nil {
		b.Fatal(err)
	}

	b.Run("ExistingFile", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = File().Validate(testFile)
		}
	})

	b.Run("NonExistingFile", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = File().Validate("/nonexistent/file.txt")
		}
	})

	b.Run("WithSearchPath", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = File(tempDir).Validate("bench.txt")
		}
	})
}

// Helper function
func fileStrPtr(s string) *string {
	return &s
}
