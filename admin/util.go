package admin

import "github.com/cthit/goldapps"

func memberContains(s []goldapps.Member, e goldapps.Member) bool {
	for _, a := range s {
		if a.Email == e.Email {
			return true
		}
	}
	return false
}

func groupContains(s []goldapps.Group, e goldapps.Group) bool {
	for _, a := range s {
		if a.Email == e.Email {
			return true
		}
	}
	return false
}
