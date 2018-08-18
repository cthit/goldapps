package main

import (
	"../ldap"
	"github.com/spf13/viper"
)

func NewLdapService() (*ldap.ServiceLDAP, error) {
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

	return ldap.NewLDAPService(dbConfig, loginConfig, usersConfig, groupsConfig, customEntryConfigs)
}
