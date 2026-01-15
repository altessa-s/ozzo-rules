// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"

	lru "github.com/hashicorp/golang-lru/v2"
)

const (
	// defaultRegexCacheCapacity is the default capacity for the regex cache
	defaultRegexCacheCapacity = 1000
)

// globalRegexCache is a thread-safe LRU cache for compiled regular expressions
var globalRegexCache = mustNewRegexCache(defaultRegexCacheCapacity)

// mustNewRegexCache creates a new LRU cache for regex patterns, panics on error
func mustNewRegexCache(capacity int) *lru.Cache[string, *regexp.Regexp] {
	cache, err := lru.New[string, *regexp.Regexp](capacity)
	if err != nil {
		// This should never happen with positive capacity
		panic(err)
	}
	return cache
}

// compileRegex returns a compiled regex from cache or compiles and caches it
func compileRegex(pattern string) (*regexp.Regexp, error) {
	// Try to get from cache first
	if re, exists := globalRegexCache.Get(pattern); exists {
		return re, nil
	}

	// Not in cache, compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// Cache the compiled regex
	globalRegexCache.Add(pattern, re)
	return re, nil
}
