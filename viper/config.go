package main

import (
	"fmt"

	"github.com/spf13/viper"
	"io/ioutil"
)

func loadConfig() {

	viper.SetConfigName("config")         // name of config file (without extension)
	viper.AddConfigPath("/etc/google-ldap-sync/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.google-ldap-sync/") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error while reading config file: %s", err))
	}

}

func getGoogleServiceKey() (*[]byte, error){

	content, err := ioutil.ReadFile(viper.GetString("gapps.servicekeyfile"))
	if err != nil {
		panic(err)
	}

	return &content, nil
}
