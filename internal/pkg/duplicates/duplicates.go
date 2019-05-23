package duplicates

import (
	"github.com/cthit/goldapps/internal/pkg/model"
	"strings"
)

func RemoveDuplicates(users model.Users, groups model.Groups) (model.Users, model.Groups) {

	// Compare Users with Groups
	for i, user := range users {
		for k := 0; k < len(groups); k++ {
			// Check if any cid conflicts with any group name
			if model.CompareEmails(user.Cid, extractIdentifier(groups[k].Email)) {
				if groups[k].Expendable {
					groups = removeArrayGroup(groups, k)
					k-- // don't breaking the loop
				} else {
					// No good strategy exists, simply panic and let admins handle the situation
					// This would probably also cause tremendous problems in other applications
					panic(user.Cid + "==" + extractIdentifier(groups[k].Email))
				}
			}
			// Check if any user nick conflicts with any group name
			if model.CompareEmails(user.Nick, extractIdentifier(groups[k].Email)) {
				if groups[k].Expendable {
					groups = removeArrayGroup(groups, k)
					k-- // don't breaking the loop
				} else {
					// Nicks are not that important
					users[i].Nick = ""
				}
			}

			for aliasIndex, alias := range groups[k].Aliases {
				// Check if any cid conflicts with any group alias
				if model.CompareEmails(user.Cid, extractIdentifier(alias)) {
					if groups[k].Expendable {
						groups[k] = removeAlias(groups[k], aliasIndex)
					} else {
						// No good strategy exists, simply panic and let admins handle the situation
						// This would probably also cause tremendous problems in other applications
						panic(user.Cid + "== (alias)" + extractIdentifier(groups[k].Email))
					}
				}
				// Check if any Nick conflicts with any group alias
				if model.CompareEmails(user.Nick, extractIdentifier(alias)) {
					if groups[k].Expendable {
						groups[k] = removeAlias(groups[k], aliasIndex)
					} else {
						// Nicks are not that important
						users[i].Nick = ""
					}
				}
			}
		}
	}

	// Compare Users with Users
	for i, user := range users {
		for j, otherUser := range users {
			// Don't check with itself
			if i != j {
				// Compare cids
				if model.CompareEmails(user.Cid, otherUser.Cid) {
					// Should not be able to happen
					panic("two users with cid: " + user.Cid)
				}
				// Compare Nicks
				if model.CompareEmails(user.Nick, otherUser.Nick) {
					// Nicks are not that important
					users[i].Nick = ""
					users[j].Nick = ""
				}
				// Compare cids with nicks
				if model.CompareEmails(user.Cid, otherUser.Nick) {
					// Nicks are not that important
					users[j].Nick = ""
				}
			}
		}
	}

	// Compare Groups with Groups
	for i, group := range groups {
		for j, otherGroup := range groups {
			// Don't check with itself
			if i != j {
				// Compare Emails
				if model.CompareEmails(group.Email, otherGroup.Email) {
					// Something is set up wrong
					panic("two groups with email: " + group.Email)
				}
				for _, alias := range group.Aliases {
					// Compare emails with aliases
					if model.CompareEmails(alias, otherGroup.Email) {
						// Something is set up wrong
						panic("two groups with alias/email: " + group.Email + ", " + otherGroup.Email)
					}
					for _, otherAlias := range otherGroup.Aliases {
						// Compare aliases with aliases
						if model.CompareEmails(alias, otherAlias) {
							// Something is set up wrong
							panic("two groups with alias: " + alias)
						}
					}
				}
			}
		}
	}
	return users, groups
}

func removeArrayGroup(s model.Groups, i int) model.Groups {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func removeArrayString(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func removeAlias(group model.Group, aliasIndex int) model.Group {
	group.Aliases = removeArrayString(group.Aliases, aliasIndex)
	return group
}

func extractIdentifier(email string) string {
	return strings.Split(email, "@")[0]
}
