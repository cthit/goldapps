package goldapps

import (
	"strings"
)

type User struct {
	Cid           string `json:"cid"`
	FirstName     string `json:"first_name"`
	SecondName    string `json:"second_name"`
	Nick          string `json:"nick"`
	Mail          string `json:"mail"`
	GdprEducation bool   `json:"gdpr_education"`
}

type Users []User

func (users Users) Contains(cid string) bool {
	for _, user := range users {
		if user.Cid == cid {
			return true
		}
	}
	return false
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

	/*if user.Mail != other.Mail {
		return false
	}*/

	/*if user.GdprEducation != other.GdprEducation {
		return false
	}*/

	/*
		Do not check PasswordHash nor HashFunction
	*/

	return true
}