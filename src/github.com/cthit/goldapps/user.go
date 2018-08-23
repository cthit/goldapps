package goldapps

import (
	"strings"
)

type User struct {
	Cid        string `json:"cid"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Nick       string `json:"nick"`
	Mail       string `json:"mail"`
}

type Users []User

// Search for username(cid) in list of groups
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

	// Don't check email as its not saved in every consumer atm
	/*if user.Mail != other.Mail {
		return false
	}*/

	return true
}
