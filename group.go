package goldapps

// Represents a email group.
// Email is the id and main email for the group.
// Members is a lost of email addresses that are members of this group.
// Aliases are alternative email addresses for the group.
type Group struct {
	Email   string
	Members []string
}

func (group Group) equals(other Group) (bool) {
	if group.Email != other.Email {return false}
	if len(group.Members) != len(other.Members) {return false}

	for _,member := range group.Members  {
		contains := false
		for _,otherMember := range other.Members  {
			if member == otherMember {
				contains = true
				break
			}
			if !contains {return false}
		}
	}

	return true
}