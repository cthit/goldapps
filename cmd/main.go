package main

import (
	"github.com/spf13/viper"
	"fmt"
	"github.com/cthit/goldapps/ldap"
)

func init() {
	err := loadConfig()
	if err != nil {
		panic(err)
	}
}

func main() {

	/*provider, err := getGoogleService(viper.GetString("gapps.servicekeyfile"), viper.GetString("gapps.adminaccount"))
	if err != nil {
		panic(err)
	}*/

	/*consumer, err := getGoogleService()
	if err != nil {

	}*/

	dbConfig := ldap.ServerConfig{
		Url:        viper.GetString("ldap.url"),
		ServerName: viper.GetString("ldap.servername"),
	}

	groupsConfig := ldap.EntryConfig{
		BaseDN:     viper.GetString("ldap.groups.basedn"),
		Filter:     viper.GetString("ldap.groups.filter"),
		Attributes: viper.GetStringSlice("ldap.groups.attributes"),
	}

	usersConfig := ldap.EntryConfig{
		BaseDN:     viper.GetString("ldap.users.basedn"),
		Filter:     viper.GetString("ldap.users.filter"),
		Attributes: viper.GetStringSlice("ldap.users.attributes"),
	}

	// Add custom entries
	customEntryNames := viper.GetStringSlice("ldap.custom")
	customEntryConfigs := make([]ldap.CustomEntryConfig, 0)
	for _, entry := range customEntryNames {
		customEntryConfigs = append(customEntryConfigs,
			ldap.CustomEntryConfig{
				BaseDN:       viper.GetString("ldap." + entry + ".basedn"),
				Filter:       viper.GetString("ldap." + entry + ".filter"),
				ParentFilter: viper.GetString("ldap." + entry + ".parent_filter"),
				Attributes:   viper.GetStringSlice("ldap." + entry + ".attributes"),
				Mail:         viper.GetString("ldap." + entry + ".mail"),
			},
		)
	}

	loginConfig := ldap.LoginConfig{
		UserName: viper.GetString("ldap.user"),
		Password: viper.GetString("ldap.password"),
	}

	provider, err := ldap.NewLDAPService(dbConfig, loginConfig, usersConfig, groupsConfig, customEntryConfigs)

	if err != nil {
		panic(err)
	}

	g, err := provider.GetGroups()
	if err != nil {
		panic(err)
	}

	if g != nil {
		for _, group := range g {
			fmt.Println(group)
		}
	}
	/*
		err = consumer.UpdateGroups(g)
		if err != nil {

		}
	*/
}
