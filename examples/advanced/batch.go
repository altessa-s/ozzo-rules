package main

import (
	"sync"

	"github.com/go-ozzo/ozzo-validation/v4/is"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Concurrent validation with goroutines
func validateBatch(items []string) []error {
	errors := make([]error, len(items))
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func(idx int, val string) {
			defer wg.Done()
			// Thread-safe validation
			// Note: ozzo_rules.Email didn't exist in the original example,
			// it probably meant validation.Email or a similar rule.
			// Assuming validation.Email from ozzo-validation/is or just validation.
			// The original snippet had ozzo_rules.Email, but looking at available files (list_dir),
			// I don't see email validation. The snippet used validation.Email in lines 309.
			// I'll assume it meant validation.Email (alias of is.Email) or similar.
			// Original snippet: errors[idx] = validation.Validate(val, ozzo_rules.Email)
			// Wait, the snippet says ozzo_rules.Email on line 379.
			// But line 309 says validation.Email.
			// I'll check "doc.go" or "validation_*.go" to see if Email is exported.
			// Based on list_dir, there is no validation_email.go.
			// I will use validation.Email (from standard ozzo) in the example to be safe/correct as per libraries.
			// Re-reading logic: the user wants the example moved. I should probably trust the snippet
			// might contain errors or I am missing something.
			// Actually, let's substitute with a known rule from ozzo_rules like CountryCode2 for the example
			// if I'm not sure, OR assume ozzo_rules re-exports it.
			// Let's check doc.go quickly? No, I'll just use what was in the snippet but change to a known valid one if needed
			// to make it compile, or keep it closer to source.
			// The snippet said `ozzo_rules.Email`. I suspect that's a typo in the README
			// and it should be `validation.Email` or `is.Email`.
			// Since I see `validation.Field(&u.Email, ozzo_rules.Required(), validation.Email)` in user.go (line 307-309),
			// it confirms validation.Email is likely intended.
			// I'll change it to `is.Email` or `validation.NewStringRule(is.Email, "invalid email")`
			// or just `validation.Email` if imported.
			// Actually, validation.Validate(val, is.Email) works.
			errors[idx] = validation.Validate(val, is.Email)
		}(i, item)
	}

	wg.Wait()
	return errors
}
