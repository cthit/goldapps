package goldapps

import "strings"

// Represents a email group.
// Email is the id and main email for the group.
// Members is a lost of email addresses that are members of this group.
// Aliases are alternative email addresses for the group.
type Group struct {
	Email   string   `json:"email"`
	Type    string   `json:"type"`
	Members []string `json:"members"`
	Aliases []string `json:"aliases"`
}

type Groups []Group

func (groups Groups) Contains(email string) bool {
	for _, group := range groups {
		if group.Email == email {
			return true
		}
	}
	return false
}

func (group Group) equals(other Group) bool {
	if strings.ToLower(group.Email) != strings.ToLower(other.Email) {
		return false
	}
	if group.Type != other.Type {
		return false
	}
	if len(group.Members) != len(other.Members) {
		return false
	}
	if len(group.Aliases) != len(other.Aliases) {
		return false
	}

	for _, member := range group.Members {
		contains := false
		for _, otherMember := range other.Members {
			if strings.ToLower(member) == strings.ToLower(otherMember) {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}

	for _, alias := range group.Aliases {
		contains := false
		for _, otheralias := range other.Aliases {
			if strings.ToLower(alias) == strings.ToLower(otheralias) {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}

	return true
}
