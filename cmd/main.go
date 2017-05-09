package main

import (
/*"fmt"
"github.com/spf13/viper"*/
)

func init() {
	err := loadConfig()
	if err != nil {
		panic(err)
	}
}

func main() {

	provider, err := getLDAPService()
	if err != nil {
		panic(err)
	}

	/*consumer, err := getGoogleService()
	if err != nil {

	}*/

	g, err := provider.Groups()
	if err != nil {
		panic(err)
	}

	if g != nil {

	}
	/*
		err = consumer.UpdateGroups(g)
		if err != nil {

		}
	*/
}
