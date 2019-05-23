package model

import (
	"regexp"
	"strings"
)

func CompareEmails(email, other string) bool {
	return SanitizeEmail(email) == SanitizeEmail(other)
}

// Only work on the part before the @
// You are only supposed to send in the part to the left of the @
func SanitizeEmail(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	replacelist := map[string]string{
		"π": "pi",
		"å": "a",
		"ä": "a",
		"ö": "o",
		"ø": "o",
		"æ": "ae",
		" ": "-",
	}
	allowed := regexp.MustCompile("[a-z]|[0-9]|-")
	parts := strings.Split(s, "")
	for i := range parts {
		if !allowed.MatchString(parts[i]) {
			parts[i] = replacelist[parts[i]]
		}
	}
	return strings.Join(parts, "")
}
