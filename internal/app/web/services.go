package web

import (
	"fmt"
	"regexp"

	"github.com/cthit/goldapps/internal/pkg"
	"github.com/cthit/goldapps/internal/pkg/actions"
	"github.com/cthit/goldapps/internal/pkg/duplicates"
	"github.com/cthit/goldapps/internal/pkg/model"
	"github.com/cthit/goldapps/internal/pkg/services"
	"github.com/cthit/goldapps/internal/pkg/services/admin"
	"github.com/cthit/goldapps/internal/pkg/services/gamma"
	"github.com/cthit/goldapps/internal/pkg/services/json"
	"github.com/spf13/viper"
)

func getConsumer(toJson string) (services.UpdateService, error) {
	var consumer services.UpdateService
	var err error

	isJson, _ := regexp.MatchString(`.+\.json$`, toJson)
	if !isJson {
		consumer, err = admin.NewGoogleService(
			viper.GetString("gapps.consumer.servicekeyfile"),
			viper.GetString("gapps.consumer.adminaccount"))
		return consumer, nil
	} else {
		consumer, err = json.NewJsonService(toJson)
	}
	if err != nil {
		fmt.Println("Failed to get consumer")
		return nil, err
	}
	return consumer, nil
}

func getProvider(fromJson string) (services.CollectionService, error) {
	var provider services.CollectionService
	var err error

	isJson, _ := regexp.MatchString(`.+\.json$`, fromJson)
	if !isJson {
		provider, err = gamma.CreateGammaService(
			viper.GetString("gamma.provider.apiKey"),
			viper.GetString("gamma.provider.url"))
	} else {
		provider, err = json.NewJsonService(fromJson)
	}
	if err != nil {
		fmt.Println("Failed to get provider")
		return nil, err
	}
	return provider, nil
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

func getChangeSuggestions(fromJson string, toJson string) (actions.UserActions, actions.GroupActions, error) {
	fmt.Println("Setting up providers")
	provider, err := getProvider(fromJson)
	if err != nil {
		return actions.UserActions{}, actions.GroupActions{}, err
	}

	fmt.Println("Setting up services")
	consumer, err := getConsumer(toJson)
	if err != nil {
		return actions.UserActions{}, actions.GroupActions{}, err
	}

	// Collect users and groups
	var providerUsers model.Users
	var consumerUsers model.Users
	var providerGroups model.Groups
	var consumerGroups model.Groups

	fmt.Println("Collecting groups from the providers...")
	providerGroups = collectGroups(provider)
	fmt.Println("Collecting groups from the services...")
	consumerGroups = collectGroups(consumer)

	fmt.Println("Collecting users from the providers...")
	providerUsers = collectUsers(provider)
	fmt.Println("Collecting users from the services...")
	consumerUsers = collectUsers(consumer)

	// Get and process additions
	providerGroups, providerUsers = pkg.AddAdditions(providerGroups, providerUsers, "additions.json")

	// Check for and handle duplicates
	providerUsers, providerGroups = duplicates.RemoveDuplicates(providerUsers, providerGroups)

	// Get changes to make
	fmt.Println("Colculating difference between the services and providers groups.")
	proposedGroupChanges := actions.GroupActionsRequired(consumerGroups, providerGroups)

	fmt.Println("Colculating difference between the services and providers users.")
	proposedUserChanges := actions.UserActionsRequired(consumerUsers, providerUsers)

	return proposedUserChanges, proposedGroupChanges, nil
}
