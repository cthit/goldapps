package cli

import (
	"bytes"
	"fmt"
)

func askBool(question string, preferred bool) bool {
	if preferred {
		fmt.Printf("%s (Y/n): ", question)
	} else {
		fmt.Printf("%s (y/N): ", question)
	}

	var input string
	fmt.Scanln(&input)

	if len(input) == 0 {
		return preferred
	} else if input == "Y" || input == "y" {
		return true
	} else if input == "N" || input == "n" {
		return false
	} else {
		fmt.Printf("'%s' is not a valid answer:\n", input)
		return askBool(question, preferred)
	}
}

func askString(question string, preferred string) string {
	fmt.Printf("%s (%s): ", question, preferred)

	var input string
	fmt.Scanln(&input)

	if len(input) == 0 {
		return preferred
	} else {
		return input
	}
}

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
