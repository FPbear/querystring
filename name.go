package querystring

import (
	"strings"
	"unicode"
)

type Name string

const (
	CamelCase  Name = "camel"
	PascalCase Name = "pascal"
	SnakeCase  Name = "snake"
)

func (n Name) IsEmpty() bool {
	return n == ""
}

// Convert : Converts the string name to the specified case.
// The function returns the converted string.
// For example:
// Name("camel").Convert("hello_world") returns "helloWorld".
// Name("pascal").Convert("hello_world") returns "HelloWorld".
// Name("snake").Convert("HelloWorld") returns "hello_world".
// Name("").Convert("hello_world") returns "hello_world".
func (n Name) Convert(name string) string {
	var tag string
	switch n {
	case CamelCase:
		tag = ToCamelCase(name)
	case PascalCase:
		tag = ToPascalCase(name)
	case SnakeCase:
		tag = ToSnakeCase(name)
	default:
		tag = name
	}
	return tag
}

// ToCamelCase : Converts the string s to camel case.
// The string s must start with a letter and contain only letters, numbers, and the characters '_', '-', and '.'.
// The first letter of the string s is converted to lowercase.
// If the string s is empty or does not start with a letter, the function returns an empty string.
// If the string s contains characters other than letters, numbers, '_', '-', and '.', the function returns an empty string.
// The function returns the converted string.
// For example:
// ToCamelCase("foo-bar") returns "fooBar".
// ToCamelCase("foo_bar") returns "fooBar".
// ToCamelCase("foo.bar") returns "fooBar".
// ToCamelCase("foo-bar-") returns "".
// ToCamelCase("foo-bar-1") returns "".
// ToCamelCase("1foo-bar") returns "".
// ToCamelCase("foobar") returns "foobar".
// ToCamelCase("fooBar") returns "fooBar".
func ToCamelCase(s string) string {
	var l = len(s)
	if l == 0 {
		return ""
	}
	if !unicode.IsLetter(rune(s[0])) ||
		(!unicode.IsLetter(rune(s[l-1])) && !unicode.IsNumber(rune(s[l-1]))) {
		return ""
	}

	var builder strings.Builder
	var isSeparator bool
	builder.WriteRune(unicode.ToLower(rune(s[0])))
	for i := 1; i < l; i++ {
		r := rune(s[i])
		if !checkCompliance(r) {
			return ""
		}
		if checkSpeChars(r) {
			isSeparator = true
			continue
		}
		if isSeparator {
			builder.WriteRune(unicode.ToUpper(r))
			isSeparator = false
			continue
		}
		builder.WriteRune(r)
	}
	return builder.String()
}

// ToPascalCase : Converts the string s to pascal case.
// The string s must start with a letter and contain only letters, numbers, and the characters '_', '-', and '.'.
// The first letter of the string s is converted to uppercase.
// If the string s is empty or does not start with a letter, the function returns an empty string.
// If the string s contains characters other than letters, numbers, '_', '-', and '.', the function returns an empty string.
// The function returns the converted string.
// For example:
// ToPascalCase("foo-bar") returns "FooBar".
// ToPascalCase("foo_bar") returns "FooBar".
// ToPascalCase("foo.bar") returns "FooBar".
// ToPascalCase("foo-bar-") returns "".
// ToPascalCase("foo-bar-1") returns "".
// ToPascalCase("1foo-bar") returns "".
// ToPascalCase("foobar") returns "Foobar".
// ToPascalCase("fooBar") returns "Foobar".
// ToPascalCase("fooBar1") returns "Foobar1".
// ToPascalCase("FooBar") returns "FooBar".
// ToPascalCase("foo_2-bar") returns "Foo2Bar".
func ToPascalCase(s string) string {
	var l = len(s)
	if l == 0 {
		return ""
	}
	if !unicode.IsLetter(rune(s[0])) ||
		(!unicode.IsLetter(rune(s[l-1])) && !unicode.IsNumber(rune(s[l-1]))) {
		return ""
	}
	var builder strings.Builder
	var isSeparator bool
	builder.WriteRune(unicode.ToUpper(rune(s[0])))
	for i := 1; i < len(s); i++ {
		r := rune(s[i])
		if !checkCompliance(r) {
			return ""
		}
		if isSeparator {
			builder.WriteRune(unicode.ToUpper(r))
			isSeparator = false
			continue
		}
		if checkSpeChars(r) {
			isSeparator = true
			continue
		}
		builder.WriteRune(r)
	}
	return builder.String()
}

// ToSnakeCase : Converts the string s to snake case.
// The string s must start with a letter and contain only letters, numbers, and the characters '_', '-', and '.'.
// The first letter of the string s is converted to lowercase.
// If the string s is empty or does not start with a letter, the function returns an empty string.
// If the string s contains characters other than letters, numbers, '_', '-', and '.', the function returns an empty string.
// The function returns the converted string.
// For example:
// ToSnakeCase("foo-bar") returns "foo_bar".
// ToSnakeCase("foo_bar") returns "foo_bar".
// ToSnakeCase("foo.bar") returns "foo_bar".
// ToSnakeCase("foo-bar-") returns "".
// ToSnakeCase("foo-bar-1") returns "foo_bar_1".
// ToSnakeCase("1foo-bar") returns "".
func ToSnakeCase(input string) string {
	if len(input) == 0 {
		return ""
	}
	var l = len(input)
	if !unicode.IsLetter(rune(input[0])) ||
		(!unicode.IsLetter(rune(input[l-1])) && !unicode.IsNumber(rune(input[l-1]))) {
		return ""
	}
	var result strings.Builder
	var isSeparator bool
	for i, char := range input {
		if i == 0 {
			result.WriteRune(unicode.ToLower(char))
			continue
		}
		if checkSpeChars(char) && !isSeparator {
			result.WriteRune('_')
			isSeparator = true
			continue
		}
		if unicode.IsLetter(char) {
			if (!isSeparator && unicode.IsUpper(char)) || unicode.IsNumber(rune(input[i-1])) {
				result.WriteRune('_')
				isSeparator = true
			}
		} else if unicode.IsNumber(char) {
			if !isSeparator {
				result.WriteRune('_')
				isSeparator = true
			}
		}
		result.WriteRune(unicode.ToLower(char))
		isSeparator = false
	}

	return result.String()
}

// checkSpeChars : Check if the character is a special character.
// The function returns true if the character is a special character( '_', '-', '.' ); otherwise, it returns false.
// For example:
// checkSpeChars('-') returns true.
// checkSpeChars('a') returns false.
func checkSpeChars(r rune) bool {
	if r == '_' || r == '-' || r == '.' {
		return true
	}
	return false
}

func checkCompliance(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsNumber(r) || checkSpeChars(r) {
		return true
	}
	return false
}
