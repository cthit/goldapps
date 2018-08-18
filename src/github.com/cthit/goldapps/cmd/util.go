package main

import "fmt"

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
