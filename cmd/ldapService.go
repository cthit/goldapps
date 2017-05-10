package main

import (
	glsync "github.com/hulthe/google-ldap-sync"
	"github.com/hulthe/google-ldap-sync/ldap"
)

func getLDAPService(url string, serverName string, userName string, password string) (glsync.GroupService, error) {

	service, err := ldap.NewLDAPService(url, serverName, userName, password)
	if err != nil {
		return nil, err
	}

	return service, nil

}
