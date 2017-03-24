package main

import (
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
	"golang.org/x/net/context"
	"github.com/spf13/viper"
)

func getGroups(service *admin.Service) (*admin.Groups, error) {

	groups, err := service.Groups.List().Customer("my_customer").MaxResults(1000000).Do()
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve groups, make sure you are using the right scopes: %v", err)
	}

	return groups, nil

}

func getGoogleService(jsonKey *[]byte,  scope ...string) (*admin.Service, error) {

	// Parse jsonKey
	config, err := google.JWTConfigFromJSON(*jsonKey, scope...)
	if err != nil {
		return nil, fmt.Errorf("Could not pase jsonKey: %v", err)
	}

	// Why do I need this??
	config.Subject = viper.GetString("gapps.adminaccount")

	// Create a client
	client := config.Client(context.Background())

	// Create a service
	service, err := admin.New(client)
	if err != nil {
		return nil, fmt.Errorf("Could not create service: %v", err)
	}

	return service, nil
}