package main

import (
	"github.com/cthit/goldapps"
	"github.com/cthit/goldapps/ldap"
)

func getLDAPService(url string, serverName string, userName string, password string) (goldapps.GroupService, error) {

	service, err := ldap.NewLDAPService(url, serverName, userName, password)
	if err != nil {
		return nil, err
	}

	return service, nil

}
