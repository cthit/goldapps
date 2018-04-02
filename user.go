package goldapps

import "strings"

type User struct {
	Cid           string `json:"cid"`
	FirstName     string `json:"first_name"`
	SecondName    string `json:"second_name"`
	Nick          string `json:"nick"` // Must be sanitized
	Mail          string `json:"mail"` // Backup email?  must be investigated
	GdprEducation bool   `json:"gdpr_education"`
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

	if strings.ToLower(user.Mail) != strings.ToLower(other.Mail) {
		return false
	}

	if user.GdprEducation != other.GdprEducation {
		return false
	}

	return true
}
