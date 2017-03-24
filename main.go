package main

import (
	"google.golang.org/api/admin/directory/v1"
	"fmt"
)

func init()  {
	loadConfig()
}

func main() {
	jsonKey, err := getGoogleServiceKey()
	if err != nil {
		panic(err)
	}

	service, err := getGoogleService(jsonKey, admin.AdminDirectoryGroupScope)
	if err != nil {
		panic(err)
	}

	groups, err := getGroups(service)
	if err != nil {
		panic(err)
	}

	for _, group := range groups.Groups {
		fmt.Printf("%v, %v \n", group.Name, group.Email)
	}

	//getLDAPGroups()
}
