package main

import (
	
	glsync "github.com/hulthe/google-ldap-sync"
	"github.com/hulthe/google-ldap-sync/admin"
	"github.com/hulthe/google-ldap-sync/ioutil"
	"github.com/hulthe/google-ldap-sync/google"
	"golang.org/x/net/context"
)

func getGoogleService(keyPath string, adminMail string) (glsync.GroupUpdateService, error) {
	
	jsonKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	
	// Parse jsonKey
	config, err := google.ConfigFromJSON(jsonKey, admin.Scopes()...)
	if err != nil {
		return nil, err
	}
	
	// Why do I need this??
	config.Subject = adminMail
	
	// Create a http client
	client := config.Client(context.Background())
	
	return admin.NewGoogleService(client)
}
