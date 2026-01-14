// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"os"
	"path/filepath"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrFileInvalid is the error that returns when a value is not a valid file path.
var ErrFileInvalid = validation.NewError("validation_file", "invalid file path")

// File is a validation rule that checks if a value is a valid file path.
func File(searchPath ...string) FileRule {
	return FileRule{condition: true, searchPath: searchPath, err: ErrFileInvalid}
}

// FileRule is a rule that checks if a value is a valid file path.
type FileRule struct {
	err        validation.Error
	condition  bool
	searchPath []string
}

// Validate checks if the given value is valid or not.
func (r FileRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	if v == nil {
		return r.err
	}

	// Handle string pointers
	if ptr, ok := v.(*string); ok {
		if ptr == nil {
			return r.err
		}
		v = *ptr
	}

	value, err := validation.EnsureString(v)
	if err != nil {
		return r.err
	}

	if len(r.searchPath) > 0 {
		file := findFile(value, r.searchPath)
		if file == "" {
			return r.err
		}

		return nil
	}

	if _, err := os.Stat(value); err != nil {
		if os.IsNotExist(err) {
			return r.err
		}
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r FileRule) When(condition bool) FileRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r FileRule) Error(message string) FileRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r FileRule) ErrorObject(err validation.Error) FileRule {
	r.err = err
	return r
}

func findFile(filePath string, searchPaths []string) string {
	if filepath.IsAbs(filePath) {
		_, err := os.Stat(filePath)
		if err == nil {
			return filePath
		}
		return ""
	}

	for _, parent := range searchPaths {
		found, err := filepath.Abs(filepath.Join(parent, filePath))
		if err != nil {
			continue
		}
		fileInfo, err := os.Stat(found)
		if err == nil && !fileInfo.IsDir() {
			return found
		}
	}
	return ""
}
