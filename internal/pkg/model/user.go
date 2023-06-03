package model

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

// Data struct representing how a user should look before and after an update
// Allows for efficient updates as application doesn't have to re-upload whole user
type UserUpdate struct {
	Before User `json:"before"`
	After  User `json:"after"`
}

// Search for username(cid) in list of groups
func (users Users) Contains(cid string) bool {
	for _, user := range users {
		if user.Cid == cid {
			return true
		}
	}
	return false
}

func (user User) Same(other User) bool {
	return strings.ToLower(user.Cid) == strings.ToLower(other.Cid)
}

func (user User) Equals(other User) bool {
	if !user.Same(other) {
		return false
	}

	if user.FirstName != other.FirstName {
		return false
	}

	if user.SecondName != other.SecondName {
		return false
	}

	if SanitizeEmail(user.Nick) != SanitizeEmail(other.Nick) { // Because google uses nick for mail
		return false
	}

	// Don't check email as its not saved in every services atm
	/*if user.Mail != other.Mail {
		return false
	}*/

	return true
}
