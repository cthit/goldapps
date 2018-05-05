package goldapps

import (
	"strings"
)

type User struct {
	Cid           string `json:"cid"`
	FirstName     string `json:"first_name"`
	SecondName    string `json:"second_name"`
	Nick          string `json:"nick"`
	GdprEducation bool   `json:"gdpr_education"`
	PasswordHash  string `json:"password_hash"` // For example "WjRIu6NlRX5PukqfMkAEb7xpOHJICasd"
	HashFunction  string `json:"hash_function"` // "crypt", "SHA-1" or "MD5", if not set PasswordHash will be interpreted as plaintext
}

func (user User) equals(other User) bool {
	if strings.ToLower(user.Cid) != strings.ToLower(other.Cid) {
		return false
	}

	if user.FirstName != other.FirstName {
		return false
	}

	if user.SecondName != other.SecondName {
		return false
	}

	if user.Nick != other.Nick {
		return false
	}

	if user.GdprEducation != other.GdprEducation {
		return false
	}

	/*
		Do not check PasswordHash nor HashFunction
	*/

	return true
}
