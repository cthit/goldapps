package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func init() {
	err := loadConfig()
	if err != nil {
		panic(err)
	}
}

func main() {

	service, err := getGoogleService(viper.GetString("gapps.servicekeyfile"), viper.GetString("gapps.adminaccount"))
	if err != nil {
		fmt.Println(viper.GetString("gapps.servicekeyfile"))
		panic(err)
	}

	groups, err := service.Groups()
	if err != nil {
		panic(err)
	}

	for _, group := range *groups {
		fmt.Printf("%v, %v \n", group.Name, group.Email)
	}

	//getLDAPGroups()
}
