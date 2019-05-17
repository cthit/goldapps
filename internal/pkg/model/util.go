package model

import (
	"strings"
)

func CompareEmails(email, other string) bool {
	return SanitizeEmail(email) == SanitizeEmail(other)
}

func SanitizeEmail(s string) string {
	s = strings.ToLower(s)
	s = strings.Replace(s, "π", "pi", -1)
	s = strings.Replace(s, "å", "a", -1)
	s = strings.Replace(s, "ä", "a", -1)
	s = strings.Replace(s, "ö", "o", -1)
	s = strings.Replace(s, "ö", "o", -1)
	s = strings.Replace(s, "ø", "o", -1)
	s = strings.Replace(s, "æ", "ae", -1)
	s = strings.Replace(s, " ", "-", -1)
	s = strings.Replace(s, ".", "", -1)
	return s
}
