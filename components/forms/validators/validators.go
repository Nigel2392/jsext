package validators

import "regexp"

// Generic validators form forms.

// Regex for validation.
const (
	PhonenumberRegex       = `^[\+]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3}[-\s\.]?[0-9]{4,6}$`
	EmailRegex             = `^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`
	DigitRegex             = `^[0-9]+$`
	AlphaRegex             = `^[a-zA-Z]+$`
	AlphaNumericRegex      = `^[a-zA-Z0-9]+$`
	AlphaNumericSpaceRegex = `^[a-zA-Z0-9 ]+$`
	SlugRegex              = `^[a-z0-9-]+$`
)

// Returns true if any of the strings are empty.
func EmptyCheck(str ...string) bool {
	for _, s := range str {
		if s == "" {
			return true
		}
	}
	return false
}

// Returns true if valid.
// Returns false if length is less than min or greater than max.
func IsValidLength(min, max int, str ...string) bool {
	for _, s := range str {
		if len(s) < min || len(s) > max {
			return false
		}
	}
	return true
}

// Returns true if valid.
// Returns false if regex does not match for any of the strings.
func IsValidRegex(regex string, strings ...string) bool {
	for _, str := range strings {
		if !regexp.MustCompile(regex).MatchString(str) {
			return false
		}
	}
	return true
}

// Returns true if input is emailaddress.
func IsEmail(strings ...string) bool {
	return IsValidRegex(EmailRegex, strings...)
}

// Returns true if input is phonenumber.
func IsPhonenumber(strings ...string) bool {
	return IsValidRegex(PhonenumberRegex, strings...)
}

// Returns true if input is digit.
func IsDigit(strings ...string) bool {
	return IsValidRegex(DigitRegex, strings...)
}

// Returns true if input is an alphabetic string.
func IsAlpha(strings ...string) bool {
	return IsValidRegex(AlphaRegex, strings...)
}

// Returns true if input is an alphanumeric string.
func IsAlphaNumeric(strings ...string) bool {
	return IsValidRegex(AlphaNumericRegex, strings...)
}

// Returns true if input is an alphanumeric string with spaces.
func IsAlphaNumericWSpaces(strings ...string) bool {
	return IsValidRegex(AlphaNumericSpaceRegex, strings...)
}

// Returns true if input is slug format.
func IsSlug(strings ...string) bool {
	return IsValidRegex(SlugRegex, strings...)
}
