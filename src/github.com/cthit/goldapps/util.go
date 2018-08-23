package goldapps

import (
	"bytes"
	"fmt"
	"strings"
)

// Done does include failed too
func printProgress(done, total, failed int) {
	p := (done * 100) / total
	builder := bytes.Buffer{}
	for i := 0; i <= 100; i++ {
		if i < p {
			builder.WriteByte('=')
		} else if i == p {
			builder.WriteByte('>')
		} else {
			builder.WriteByte(' ')
		}
	}
	fmt.Printf("\rProgress: [%s] %d/%d", builder.String(), done, total)

	// Add failed counter if necessary
	if failed != 0 {
		fmt.Printf(" (Failed: %d)", failed)
	}

	// Replace progressbar with done text
	if done == total {
		if failed != 0 {
			fmt.Printf("Done! (Failed: %d)", failed)
		}
		fmt.Printf("\rDone\n")
	}
}

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
