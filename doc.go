// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

// Package ozzo_rules provides a comprehensive set of additional validation rules for use with
// the ozzo-validation package. It extends the standard validation capabilities with specialized
// validators for geographic data, contact information, identity verification, network addresses,
// and region-specific formats, all optimized with LRU caching for enhanced performance.
//
// The package includes validators for common use cases such as country codes, phone numbers,
// email addresses, URLs, and various identity documents. All validators implement the
// ozzo-validation Rule interface and support conditional validation, custom error messages,
// and thread-safe concurrent usage.
//
// Example:
//
//	// Validate a user registration form
//	err := validation.ValidateStruct(&user,
//	    validation.Field(&user.Username, ozzo_rules.Required(), ozzo_rules.Username()),
//	    validation.Field(&user.Password, ozzo_rules.Password().MinLength(12)),
//	    validation.Field(&user.Country, ozzo_rules.CountryCode2()),
//	    validation.Field(&user.Phone, ozzo_rules.Phone(user.Country)),
//	)
package ozzo_rules
