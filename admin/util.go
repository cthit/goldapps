package admin

import glsync "github.com/hulthe/google-ldap-sync"

func memberContains(s []glsync.Member, e glsync.Member) bool {
	for _, a := range s {
		if a.Email == e.Email {
			return true
		}
	}
	return false
}

func groupContains(s []glsync.Group, e glsync.Group) bool {
	for _, a := range s {
		if a.Email == e.Email {
			return true
		}
	}
	return false
}
