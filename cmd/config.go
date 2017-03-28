package main

import (
	"github.com/spf13/viper"
)

func loadConfig() error {

	viper.SetConfigName("config")                   // name of config file (without extension)
	viper.AddConfigPath("/etc/google-ldap-sync/")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.google-ldap-sync/") // call multiple times to add many search paths
	viper.AddConfigPath(".")                        // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	return err

}

