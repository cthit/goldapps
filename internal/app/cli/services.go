package cli

import (
	"fmt"
	"regexp"

	"github.com/cthit/goldapps/internal/pkg/model"
	"github.com/cthit/goldapps/internal/pkg/services"
	"github.com/cthit/goldapps/internal/pkg/services/admin"
	"github.com/cthit/goldapps/internal/pkg/services/auth"
	"github.com/cthit/goldapps/internal/pkg/services/gamma"
	"github.com/cthit/goldapps/internal/pkg/services/json"
	"github.com/cthit/goldapps/internal/pkg/services/ldap"
	"github.com/spf13/viper"
)

func getConsumer() services.UpdateService {
	var to string
	if flags.interactive {
		to = askString("Which services would you like to use, 'gapps' or '*.json?", "gapps")
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
			fmt.Println("You must specify 'gapps' or '*.json' as services.")
			previous := flags.interactive
			flags.interactive = true
			defer func() {
				flags.interactive = previous
			}()
			return getConsumer()
		}
	}
}

func getProvider() services.CollectionService {
	var from string
	if flags.interactive {
		from = askString("which providers would you like to use, 'ldap', 'gapps', 'gamma' or '*.json'?", "ldap")
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
	case "gamma":
		provider, err := gamma.CreateGammaService(
			viper.GetString("gamma.provider.apiKey"),
			viper.GetString("gamma.provider.url"))
		if err != nil {
			fmt.Println("Failed to connect to Gamma")
			panic(err)
		}
		return provider
	case "auth":
		provider, _ := auth.CreateAuthService(
			viper.GetString("auth.provider.apiKey"),
			viper.GetString("auth.provider.url"),
		)
		return provider
	default:
		isJson, _ := regexp.MatchString(`.+\.json$`, from)
		if isJson {
			provider, _ := json.NewJsonService(from)
			return provider
		} else {
			fmt.Println("You must specify 'gapps', 'ldap' or '*.json' as providers.")
			previous := flags.interactive
			flags.interactive = true
			defer func() {
				flags.interactive = previous
			}()
			return getProvider()
		}
	}
}

func collectGroups(service services.CollectionService) model.Groups {
	groups, err := service.GetGroups()
	if err != nil {
		fmt.Println("Failed to collect groups")
		panic(err)
	}
	fmt.Printf("%d groups collected.\n", len(groups))
	return groups
}

func collectUsers(service services.CollectionService) model.Users {
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
