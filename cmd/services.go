package main

import (
	"../ldap"
	"../json"
	"../admin"
	"github.com/spf13/viper"
	"../../goldapps"
	"fmt"
	"regexp"
)


func getConsumer() goldapps.UpdateService {
	var to string
	if flags.interactive {
		to = askString("Which consumer would you like to use, 'gapps' or '*.json?", "gapps")
	} else {
		to = flags.to
	}

	switch to {
	case "gapps":
		consumer, err := admin.NewGoogleService(
			viper.GetString("gapps.consumer.servicekeyfile"),
			viper.GetString("gapps.consumer.adminaccount"))
		if err != nil {
			fmt.Println("Failed to create gapps connection.")
			panic(err)
		}
		return consumer
	default:
		isJson, _ := regexp.MatchString(`.+\.json$`, to)
		if isJson {
			consumer, _ := json.NewJsonService(to)
			return consumer
		} else {
			fmt.Println("You must specify 'gapps' or '*.json' as consumer.")
			previous := flags.interactive
			flags.interactive = true
			defer func() {
				flags.interactive = previous
			}()
			return getConsumer()
		}
	}
}

func getProvider() goldapps.CollectionService {
	var from string
	if flags.interactive {
		from = askString("which provider would you like to use, 'ldap', 'gapps' or '*.json'?", "ldap")
	} else {
		from = flags.from
	}

	switch from {
	case "ldap":
		provider, err := newLdapService()
		if err != nil {
			fmt.Println("Failed to create LDAP connection.")
			panic(err)
		}
		return provider
	case "gapps":
		provider, err := admin.NewGoogleService(viper.GetString("gapps.provider.servicekeyfile"), viper.GetString("gapps.provider.adminaccount"))
		if err != nil {
			fmt.Println("Failed to create gapps connection, make sure you have setup gappsProvider in the config file.")
			panic(err)
		}
		return provider
	default:
		isJson, _ := regexp.MatchString(`.+\.json$`, from)
		if isJson {
			provider, _ := json.NewJsonService(from)
			return provider
		} else {
			fmt.Println("You must specify 'gapps', 'ldap' or '*.json' as provider.")
			previous := flags.interactive
			flags.interactive = true
			defer func() {
				flags.interactive = previous
			}()
			return getProvider()
		}
	}
}

func collectGroups(service goldapps.CollectionService) (goldapps.Groups) {
	groups, err := service.GetGroups()
	if err != nil {
		fmt.Println("Failed to collect groups")
		panic(err)
	}
	fmt.Printf("%d groups collected.\n", len(groups))
	return groups
}

func collectUsers(service goldapps.CollectionService) (goldapps.Users) {
	users, err := service.GetUsers()
	if err != nil {
		fmt.Println("Failed to collect users")
		panic(err)
	}
	fmt.Printf("%d users collected.\n", len(users))
	return users
}

func newLdapService() (*ldap.ServiceLDAP, error) {
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
