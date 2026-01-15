# Package Description
[![Test](https://github.com/altessa-s/ozzo-rules/actions/workflows/test.yml/badge.svg)](https://github.com/altessa-s/ozzo-rules/actions/workflows/test.yml)
[![Lint](https://github.com/altessa-s/ozzo-rules/actions/workflows/lint.yml/badge.svg)](https://github.com/altessa-s/ozzo-rules/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/altessa-s/ozzo-rules)](https://goreportcard.com/report/github.com/altessa-s/ozzo-rules)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/altessa-s/ozzo-rules)](https://github.com/altessa-s/ozzo-rules/releases)


Package `ozzo_rules` provides a comprehensive set of additional validation rules for use with the [ozzo-validation](https://github.com/go-ozzo/ozzo-validation) package. It extends the standard validation capabilities with specialized validators for geographic data, contact information, identity verification, network addresses, and region-specific formats, all optimized with LRU caching for enhanced performance.

> This package was originally developed as an internal project at Altessa Solutions Inc. and is now open-sourced under the MIT License.


> AI tools were used during the open-source transition for: creating GitHub workflows and configuration files, updating and consolidating documentation, generating usage examples, and reviewing code for consistency.

## API Reference

### Core Types

| Type | Description |
|------|-------------|
| `Rule` | Common interface for all validation rules |

### Geographic Validators

| Rule | Description |
|------|-------------|
| `CountryCode2() CountryCode2Rule` | Validates ISO 3166-1 alpha-2 country codes |
| `LangCode2() LangCode2Rule` | Validates ISO 639-1 language codes |
| `ZipCode() ZipCodeRule` | Validates international postal codes using the postcode library |
| `Timezone() TimezoneRule` | Validates IANA timezone database names |

### Contact Information Validators

| Rule | Description |
|------|-------------|
| `Phone(countryCode string) PhoneRule` | Validates international phone numbers using Google's libphonenumber |

### Identity Validators

| Rule | Description |
|------|-------------|
| `Username() UsernameRule` | Validates usernames with alphanumeric, underscore, and hyphen characters |
| `Password() PasswordRule` | Validates password strength with configurable requirements |

### Date/Time Validators

| Rule | Description |
|------|-------------|
| `Birthdate() BirthdateRule` | Validates birthdates in YYYY-MM-DD format |
| `Timestamp() TimestampRule` | Validates Unix timestamps (seconds since epoch) |
| `Duration() DurationRule` | Validates Go duration strings (e.g., "1h30m", "5s") |
| `DurationLimit(min, max time.Duration) DurationLimitRule` | Validates durations within specified bounds |

### Network Validators

| Rule | Description |
|------|-------------|
| `ListenAddress() ListenAddressRule` | Validates network listen addresses (host:port format) |
| `ServerAddress() ServerAddressRule` | Validates server addresses (host:port or just host) |
| `PortRange() PortRangeRule` | Validates port range strings (e.g., "8080" or "8000-9000") |
| `URI() URIRule` | Validates URIs with various schemes |

### Russian-specific Validators

| Rule | Description |
|------|-------------|
| `RuINN() RuINNRule` | Validates Russian Individual Taxpayer Number |
| `RuSNILS() RuSNILSRule` | Validates Russian social insurance number |
| `RuOGRN() RuOGRNRule` | Validates Russian Primary State Registration Number |
| `RuBIC() RuBICRule` | Validates Russian Bank Identification Code |
| `RuFIO() RuFIORrule` | Validates Russian names (Cyrillic characters) |

### Apple Developer Validators

| Rule | Description |
|------|-------------|
| `ApnKey() ApnKeyRule` | Validates Apple Push Notification authentication keys |
| `ApnKeyId() ApnKeyIdRule` | Validates Apple Push Notification Key IDs |
| `ApnTeamId() ApnTeamIdRule` | Validates Apple Developer Team IDs |
| `BundleID() BundleIDRule` | Validates iOS Bundle Identifiers |

### Service Discovery Validators

| Rule | Description |
|------|-------------|
| `ConsulServiceName() ConsulServiceNameRule` | Validates Consul service names |
| `ConsulServiceAddr() ConsulServiceAddrRule` | Validates Consul service addresses |

### Utility Validators

| Rule | Description |
|------|-------------|
| `OneOf(values ...interface{}) OneOfRule` | Validates value is in allowed set |
| `Required() RequiredRule` | Enhanced required field validation |
| `Regex(pattern string) RegexRule` | Validates against regular expression |
| `File() FileRule` | Validates file paths |
| `Slug() SlugRule` | Validates URL-friendly slugs |

### Pagination Validators

| Rule | Description |
|------|-------------|
| `ListLimit() ListLimitRule` | Validates pagination limit values |
| `ListOffset() ListOffsetRule` | Validates pagination offset values |

### Financial Validators

| Rule | Description |
|------|-------------|
| `RsPIB() RsPIBRule` | Validates Serbian Tax Identification Number |
| `RsMBR() RsMBRRule` | Validates Serbian Company Registration Number |

## Examples

You can find all validation examples in the [examples](./examples) directory.

- [Basic Usage](./examples/basic/main.go)
- [Struct Validation](./examples/advanced/user.go)
- [API Request Validation](./examples/advanced/request.go)
- [Network Configuration](./examples/advanced/server.go)
- [Concurrent Validation](./examples/advanced/batch.go)

## Performance Notes

- **Regex Compilation**: All regex patterns are compiled once and cached using an LRU cache (capacity: 1000)
- **Memory Usage**: Minimal allocations per validation, with most validators using stack allocation
- **Concurrency**: All validators are thread-safe and can be used concurrently without locks
- **Benchmark Results**: Regex validation improved from 52.8ms to 0.1ms (528x faster) with caching

## Limitations

- **Phone Validation**: Requires explicit country code; cannot auto-detect from national formats
- **ZipCode Validation**: Limited by postcode library's country support
- **Language Codes**: Only ISO 639-1 (2-letter) codes supported; ISO 639-2/3 not supported
- **Timezone Validation**: Requires exact IANA timezone names or known abbreviations

## See Also

- [ozzo-validation](https://github.com/go-ozzo/ozzo-validation) - The base validation framework
- [libphonenumber](https://github.com/nyaruka/phonenumbers) - Phone number validation library
- [postcode](https://github.com/adrg/postcode) - International postal code validation
- [govalidator](https://github.com/asaskevich/govalidator) - Additional validation utilities
