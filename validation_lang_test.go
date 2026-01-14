// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestLangCode2(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid language codes - Major languages
		{"valid English", "en", false},
		{"valid Spanish", "es", false},
		{"valid French", "fr", false},
		{"valid German", "de", false},
		{"valid Chinese", "zh", false},
		{"valid Japanese", "ja", false},
		{"valid Korean", "ko", false},
		{"valid Arabic", "ar", false},
		{"valid Russian", "ru", false},
		{"valid Portuguese", "pt", false},
		{"valid Italian", "it", false},
		{"valid Dutch", "nl", false},
		{"valid Polish", "pl", false},
		{"valid Turkish", "tr", false},
		{"valid Hindi", "hi", false},

		// Valid language codes - European languages
		{"valid Swedish", "sv", false},
		{"valid Norwegian", "no", false},
		{"valid Danish", "da", false},
		{"valid Finnish", "fi", false},
		{"valid Greek", "el", false},
		{"valid Czech", "cs", false},
		{"valid Hungarian", "hu", false},
		{"valid Romanian", "ro", false},
		{"valid Bulgarian", "bg", false},
		{"valid Croatian", "hr", false},
		{"valid Serbian", "sr", false},
		{"valid Slovak", "sk", false},
		{"valid Slovenian", "sl", false},
		{"valid Lithuanian", "lt", false},
		{"valid Latvian", "lv", false},
		{"valid Estonian", "et", false},
		{"valid Albanian", "sq", false},
		{"valid Macedonian", "mk", false},
		{"valid Maltese", "mt", false},
		{"valid Icelandic", "is", false},
		{"valid Irish", "ga", false},
		{"valid Welsh", "cy", false},
		{"valid Basque", "eu", false},
		{"valid Catalan", "ca", false},
		{"valid Galician", "gl", false},

		// Valid language codes - Asian languages
		{"valid Thai", "th", false},
		{"valid Vietnamese", "vi", false},
		{"valid Indonesian", "id", false},
		{"valid Malay", "ms", false},
		{"valid Filipino", "tl", false},
		{"valid Bengali", "bn", false},
		{"valid Urdu", "ur", false},
		{"valid Tamil", "ta", false},
		{"valid Telugu", "te", false},
		{"valid Marathi", "mr", false},
		{"valid Gujarati", "gu", false},
		{"valid Kannada", "kn", false},
		{"valid Malayalam", "ml", false},
		{"valid Sinhala", "si", false},
		{"valid Burmese", "my", false},
		{"valid Khmer", "km", false},
		{"valid Lao", "lo", false},
		{"valid Mongolian", "mn", false},
		{"valid Nepali", "ne", false},
		{"valid Persian", "fa", false},
		{"valid Hebrew", "he", false},
		{"valid Yiddish", "yi", false},
		{"valid Armenian", "hy", false},
		{"valid Georgian", "ka", false},
		{"valid Kazakh", "kk", false},
		{"valid Uzbek", "uz", false},
		{"valid Azerbaijani", "az", false},
		{"valid Turkmen", "tk", false},
		{"valid Tajik", "tg", false},
		{"valid Kyrgyz", "ky", false},

		// Valid language codes - African languages
		{"valid Swahili", "sw", false},
		{"valid Zulu", "zu", false},
		{"valid Xhosa", "xh", false},
		{"valid Afrikaans", "af", false},
		{"valid Somali", "so", false},
		{"valid Hausa", "ha", false},
		{"valid Yoruba", "yo", false},
		{"valid Igbo", "ig", false},
		{"valid Amharic", "am", false},
		{"valid Tigrinya", "ti", false},
		{"valid Oromo", "om", false},

		// Valid language codes - Other regions
		{"valid Maori", "mi", false},
		{"valid Hawaiian", "haw", true}, // 3-letter code, should be invalid
		{"valid Samoan", "sm", false},
		{"valid Tongan", "to", false},
		{"valid Fijian", "fj", false},

		// Valid language codes - Constructed languages
		{"valid Esperanto", "eo", false},
		{"valid Interlingua", "ia", false},
		{"valid Volapük", "vo", false},

		// Valid language codes - Ancient/Classical languages
		{"valid Latin", "la", false},
		{"valid Sanskrit", "sa", false},
		{"valid Ancient Greek", "grc", true}, // 3-letter code, should be invalid

		// Valid language codes - Sign languages
		{"valid Sign Language", "sgn", true}, // 3-letter code, should be invalid

		// Invalid language codes
		{"invalid three letters", "eng", true},
		{"invalid four letters", "engl", true},
		{"invalid single letter", "e", true},
		{"invalid numeric", "12", true},
		{"invalid alphanumeric", "e1", true},
		{"invalid special chars", "e$", true},
		{"invalid with space", "e n", true},
		{"invalid with dash", "e-n", true},
		{"invalid empty", "", true},
		{"invalid spaces", "  ", true},
		{"invalid non-existent", "xx", true},
		{"invalid zz", "zz", true},
		{"valid Afar", "aa", false}, // Afar language

		// Case sensitivity (ISO 639 codes are lowercase)
		{"invalid uppercase EN", "EN", true},
		{"invalid uppercase ES", "ES", true},
		{"invalid mixed case En", "En", true},
		{"invalid mixed case eN", "eN", true},

		// Edge cases
		{"nil value", nil, true},
		{"integer value", 12, true},
		{"float value", 12.34, true},
		{"boolean value", true, true},
		{"slice value", []string{"en"}, true},

		// Special codes
		{"valid Norwegian Bokmål", "nb", false},
		{"valid Norwegian Nynorsk", "nn", false},
		{"valid Chinese Simplified", "zh", false}, // zh is valid, variants like zh-CN are different
		{"valid Serbian Latin", "sr", false},      // sr is valid, script variants are different

		// Deprecated codes - no longer valid in govalidator
		{"invalid Indonesian old", "in", true}, // deprecated, use id
		{"invalid Hebrew old", "iw", true},     // deprecated, use he
		{"invalid Yiddish old", "ji", true},    // deprecated, use yi
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.Validate(tt.value, LangCode2())
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

func TestLangCode2_CommonMistakes(t *testing.T) {
	// Test common mistakes people make with language codes
	mistakes := []struct {
		name        string
		input       string
		description string
	}{
		{"ENG instead of en", "eng", "Should use 2-letter code en"},
		{"EN uppercase", "EN", "Language codes are lowercase"},
		{"english", "english", "Should use code, not name"},
		{"en-US", "en-US", "Locale code, not language code"},
		{"en_US", "en_US", "Locale code with underscore"},
		{"chi instead of zh", "chi", "Chinese 3-letter code"},
		{"ger instead of de", "ger", "German 3-letter code"},
		{"fre instead of fr", "fre", "French 3-letter code"},
		{"spa instead of es", "spa", "Spanish 3-letter code"},
		{"1a", "1a", "Starts with number"},
		{"a1", "a1", "Contains number"},
	}

	for _, m := range mistakes {
		t.Run(m.name, func(t *testing.T) {
			err := validation.Validate(m.input, LangCode2())
			if err == nil {
				t.Errorf("%s: expected error but got nil", m.description)
			}
		})
	}
}

func TestLangCode2_LanguageFamilies(t *testing.T) {
	// Test language codes by language families
	families := map[string][]string{
		"Germanic": {
			"en", // English
			"de", // German
			"nl", // Dutch
			"sv", // Swedish
			"no", // Norwegian
			"da", // Danish
			"is", // Icelandic
			"af", // Afrikaans
			"fy", // Frisian
			"lb", // Luxembourgish
			"yi", // Yiddish
		},
		"Romance": {
			"es", // Spanish
			"fr", // French
			"it", // Italian
			"pt", // Portuguese
			"ro", // Romanian
			"ca", // Catalan
			"gl", // Galician
			"la", // Latin
			"co", // Corsican
			"sc", // Sardinian
		},
		"Slavic": {
			"ru", // Russian
			"pl", // Polish
			"uk", // Ukrainian
			"cs", // Czech
			"sk", // Slovak
			"bg", // Bulgarian
			"hr", // Croatian
			"sr", // Serbian
			"sl", // Slovenian
			"mk", // Macedonian
			"be", // Belarusian
		},
		"Sino-Tibetan": {
			"zh", // Chinese
			"my", // Burmese
			"bo", // Tibetan
			"dz", // Dzongkha
		},
		"Indo-Aryan": {
			"hi", // Hindi
			"bn", // Bengali
			"pa", // Punjabi
			"gu", // Gujarati
			"mr", // Marathi
			"ur", // Urdu
			"ne", // Nepali
			"si", // Sinhala
			"sa", // Sanskrit
		},
		"Semitic": {
			"ar", // Arabic
			"he", // Hebrew
			"am", // Amharic
			"ti", // Tigrinya
			"mt", // Maltese
		},
		"Turkic": {
			"tr", // Turkish
			"az", // Azerbaijani
			"kk", // Kazakh
			"ky", // Kyrgyz
			"tk", // Turkmen
			"uz", // Uzbek
			"tt", // Tatar
			"ba", // Bashkir
		},
		"Austronesian": {
			"id", // Indonesian
			"ms", // Malay
			"tl", // Filipino/Tagalog
			"jv", // Javanese
			"su", // Sundanese
			"mg", // Malagasy
			"mi", // Maori
			"sm", // Samoan
			"to", // Tongan
			"fj", // Fijian
			"ty", // Tahitian
		},
	}

	for family, codes := range families {
		t.Run(family, func(t *testing.T) {
			for _, code := range codes {
				err := validation.Validate(code, LangCode2())
				if err != nil {
					t.Errorf("%s family: language code %s should be valid, got error: %v", family, code, err)
				}
			}
		})
	}
}

func TestLangCode2_WritingSystems(t *testing.T) {
	// Test languages by writing system
	writingSystems := map[string][]string{
		"Latin": {
			"en", "es", "fr", "de", "it", "pt", "nl", "pl", "cs", "hu",
			"ro", "tr", "vi", "id", "ms", "tl", "sw", "yo", "zu", "xh",
		},
		"Cyrillic": {
			"ru", "uk", "bg", "sr", "mk", "be", "kk", "ky", "tg", "mn",
		},
		"Arabic": {
			"ar", "fa", "ur", "ps", "sd", "ug", "ku",
		},
		"Chinese": {
			"zh",
		},
		"Devanagari": {
			"hi", "ne", "mr", "sa",
		},
		"Other": {
			"ja", // Japanese (Hiragana, Katakana, Kanji)
			"ko", // Korean (Hangul)
			"th", // Thai
			"he", // Hebrew
			"el", // Greek
			"ka", // Georgian
			"hy", // Armenian
			"am", // Amharic (Ge'ez)
			"ta", // Tamil
			"te", // Telugu
			"kn", // Kannada
			"ml", // Malayalam
			"si", // Sinhala
			"my", // Burmese
			"km", // Khmer
			"lo", // Lao
			"bo", // Tibetan
		},
	}

	for system, codes := range writingSystems {
		t.Run(system, func(t *testing.T) {
			for _, code := range codes {
				err := validation.Validate(code, LangCode2())
				if err != nil {
					t.Errorf("%s writing system: language code %s should be valid, got error: %v", system, code, err)
				}
			}
		})
	}
}

func TestLangCode2_SpecialCases(t *testing.T) {
	// Test special cases and edge conditions
	specialCases := []struct {
		name      string
		code      string
		expectErr bool
		note      string
	}{
		// Macrolanguages
		{"Arabic macrolanguage", "ar", false, "Covers multiple Arabic varieties"},
		{"Chinese macrolanguage", "zh", false, "Covers Mandarin, Cantonese, etc."},
		{"Malay macrolanguage", "ms", false, "Covers Malaysian and Indonesian"},

		// Special purpose codes
		{"No linguistic content", "zxx", true, "Valid in ISO 639-2/3 but not 639-1"},
		{"Undetermined", "und", true, "Valid in ISO 639-2/3 but not 639-1"},
		{"Multiple languages", "mul", true, "Valid in ISO 639-2/3 but not 639-1"},

		// Deprecated codes that might still be encountered
		{"Indonesian deprecated", "in", true, "Deprecated by ISO, use id"},
		{"Hebrew deprecated", "iw", true, "Deprecated by ISO, use he"},
		{"Yiddish deprecated", "ji", true, "Deprecated by ISO, use yi"},

		// Invalid constructed codes
		{"Invalid xx", "xx", true, "Often used as placeholder"},
		{"Invalid zz", "zz", true, "Often used as test code"},
		{"Invalid qq", "qq", true, "Sometimes used for unknown"},
	}

	for _, sc := range specialCases {
		t.Run(sc.name, func(t *testing.T) {
			err := validation.Validate(sc.code, LangCode2())
			if sc.expectErr && err == nil {
				t.Errorf("%s (%s): expected error but got nil - %s", sc.name, sc.code, sc.note)
			} else if !sc.expectErr && err != nil {
				t.Errorf("%s (%s): expected no error but got: %v - %s", sc.name, sc.code, err, sc.note)
			}
		})
	}
}

func TestLangCode2Rule_When(t *testing.T) {
	err := LangCode2().When(false).Validate("xx")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = LangCode2().When(true).Validate("xx")
	if err == nil {
		t.Error("expected error when condition is true and language code is invalid")
	}
}

func TestLangCode2Rule_Error(t *testing.T) {
	customMsg := "custom language code error"
	err := LangCode2().Error(customMsg).Validate("xx")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestLangCode2Rule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := LangCode2().ErrorObject(customErr).Validate("xx")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func BenchmarkLangCode2Validation(b *testing.B) {
	codes := []any{
		"en",
		"es",
		"zh",
		"INVALID",
		"eng",
		"EN",
		"12",
		"",
		nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, code := range codes {
			_ = validation.Validate(code, LangCode2())
		}
	}
}
