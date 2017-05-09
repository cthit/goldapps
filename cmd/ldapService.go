package main

import (
	
	glsync "github.com/hulthe/google-ldap-sync"
	"github.com/hulthe/google-ldap-sync/ldap"
)

func getLDAPService() (glsync.GroupService, error) {
	
	service, err := ldap.NewLDAPService()
	if err != nil {
		return nil, err
	}
	
	return service, nil
	
}
