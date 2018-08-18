package goldapps

import (
	"strings"
)

func CheckDuplicates(users Users, groups Groups) (Users, Groups) {

	// User <-> Group
	for i, user := range users {
		for k := 0; k < len(groups); k++ {
			if strings.ToLower(user.Cid+"@chalmers.it") == strings.ToLower(groups[k].Email) { // check cid with group mail
				panic(user.Cid + "@chalmers.it" + "==" + groups[k].Email) // panic, this is bad
			}
			if strings.ToLower(user.Nick+"@chalmers.it") == strings.ToLower(groups[k].Email) { // check nick with group mail
				if len(groups[k].Members) == 1 && strings.SplitN(groups[k].Members[0], "@", 2)[1] != "chalmers.it" { // special case for digit pateter.
					
					groups = remove(groups, k)
					k-- // dont breaking the loop
				} else {
					users[i].Nick = "" // Simply remove the conflicting Nick
				}
			}

			for _, alias := range groups[k].Aliases {
				if strings.ToLower(user.Cid+"@chalmers.it") == strings.ToLower(alias) { // check cid with all aliases
					panic(user.Cid + "@chalmers.it" + "==" + groups[k].Email) // panic, this is bad
				}
				if strings.ToLower(user.Nick+"@chalmers.it") == strings.ToLower(alias) { // check nick with all aliases
					users[i].Nick = "" // Remove nick because it's stupid
				}
			}
		}
	}

	// User <-> User
	for i, user := range users {
		for j, other := range users {
			if i != j { // Don't check with itself
				if strings.ToLower(user.Cid) == strings.ToLower(other.Cid) { // cid vs cid
					panic("two users with cid: " + user.Cid) // this shouldn't happen
				}
				if strings.ToLower(user.Nick) == strings.ToLower(other.Nick) { // don't compete over nicks
					users[i].Nick = ""
					users[j].Nick = ""
				}
				if strings.ToLower(user.Cid) == strings.ToLower(other.Nick) { // cid vs nick
					users[j].Nick = "" // cid takes precedence
				}
			}
		}
	}

	// Group <-> group
	for i, group := range groups {
		for j, other := range groups {
			if i != j { // don't check with itself
				if strings.ToLower(group.Email) == strings.ToLower(other.Email) { // mail vs mail
					panic("two groups with email: " + group.Email) // panic, something is set up wrong
				}
				for _, alias := range group.Aliases {
					if strings.ToLower(alias) == strings.ToLower(other.Email) { //email vs alias
						panic("two groups with alias/email: " + group.Email + ", " + other.Email) // panic, something is set up wrong
					}
					for _, oalias := range other.Aliases {
						if strings.ToLower(alias) == strings.ToLower(oalias) { // alias vs alias
							panic("two groups with alias: " + alias) // panic, something is set up wrong
						}
					}
				}
			}
		}
	}
	return users, groups
}

func remove(s Groups, i int) Groups {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}