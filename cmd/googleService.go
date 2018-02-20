package main

import (

	"github.com/cthit/goldapps"
	"github.com/cthit/goldapps/admin"
	"github.com/cthit/goldapps/ioutil"
	"github.com/cthit/goldapps/google"
	"golang.org/x/net/context"
)

func getGoogleService(keyPath string, adminMail string) (goldapps.GroupUpdateService, error) {
	
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
