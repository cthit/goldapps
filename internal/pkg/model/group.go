package model

// Represents a email group.
// Email is the id and main email for the group.
// Members is a lost of email addresses that are members of this group.
// Aliases are alternative email addresses for the group.
type Group struct {
	Email      string   `json:"email"`
	Type       string   `json:"type"`
	Members    []string `json:"members"`
	Aliases    []string `json:"aliases"`
	Expendable bool     `json:"expendable"` // Not used in comparision
}

type Groups []Group

// Data struct representing how a group looks not and how it should look after an update
// Allows for efficient updates as application doesn't have to re-upload whole group
type GroupUpdate struct {
	Before Group
	After  Group
}

// Search for groupname(email) in list of groups
func (groups Groups) Contains(email string) bool {
	for _, group := range groups {
		if group.Email == email {
			return true
		}
	}
	return false
}

func (group Group) Same(other Group) bool {
	return CompareEmails(group.Email, other.Email)
}

func (group Group) Equals(other Group) bool {
	if !group.Same(other) {
		return false
	}

	if len(group.Members) != len(other.Members) {
		return false
	}
	for _, member := range group.Members {
		contains := false
		for _, otherMember := range other.Members {
			if CompareEmails(member, otherMember) {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}

	if len(group.Aliases) != len(other.Aliases) {
		return false
	}
	for _, alias := range group.Aliases {
		contains := false
		for _, otherAlias := range other.Aliases {
			if CompareEmails(alias, otherAlias) {
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
