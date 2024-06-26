package web

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/cthit/goldapps/internal/pkg"
	"github.com/cthit/goldapps/internal/pkg/actions"
	"github.com/cthit/goldapps/internal/pkg/duplicates"
	"github.com/cthit/goldapps/internal/pkg/model"
	"github.com/cthit/goldapps/internal/pkg/services"
	"github.com/cthit/goldapps/internal/pkg/services/admin"
	"github.com/cthit/goldapps/internal/pkg/services/auth"
	"github.com/cthit/goldapps/internal/pkg/services/gamma"
	"github.com/cthit/goldapps/internal/pkg/services/json"
	"github.com/spf13/viper"
)

func getConsumer(toJson string) (services.UpdateService, error) {
	var consumer services.UpdateService
	var err error

	isJson, _ := regexp.MatchString(`.+\.json$`, toJson)
	if toJson == "gapps" {
		consumer, err = admin.NewGoogleService(
			viper.GetString("gapps.consumer.servicekeyfile"),
			viper.GetString("gapps.consumer.adminaccount"))
		return consumer, nil
	} else if isJson {
		consumer, err = json.NewJsonService(toJson)
	} else {
		fmt.Printf("Consumer '%s' was not found\n", toJson)
		return nil, errors.New("Consumer not found")
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
	if fromJson == "gamma" {
		provider, err = gamma.CreateGammaService(
			viper.GetString("gamma.provider.apiKey"),
			viper.GetString("gamma.provider.url"))
	} else if fromJson == "auth" {
		provider, _ = auth.CreateAuthService(
			viper.GetString("auth.provider.apiKey"),
			viper.GetString("auth.provider.url"))
	} else if isJson {
		provider, err = json.NewJsonService(fromJson)
	} else {
		fmt.Printf("Provider '%s' was not found\n", fromJson)
		return nil, errors.New("Provider not found")
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
	providerGroups, providerUsers = pkg.AddAdditions(providerGroups, providerUsers, viper.GetString("additions.file"))

	// Check for and handle duplicates
	providerUsers, providerGroups = duplicates.RemoveDuplicates(providerUsers, providerGroups)

	// Get changes to make
	fmt.Println("Colculating difference between the services and providers groups.")
	proposedGroupChanges := actions.GroupActionsRequired(consumerGroups, providerGroups)

	fmt.Println("Colculating difference between the services and providers users.")
	proposedUserChanges := actions.UserActionsRequired(consumerUsers, providerUsers)

	return proposedUserChanges, proposedGroupChanges, nil
}

func commitChanges(userChanges actions.UserActions, groupChanges actions.GroupActions, toJson string) bool {
	fmt.Println("Setting up services")
	consumer, err := getConsumer(toJson)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Commit changes
	userErrors := userChanges.Commit(consumer)
	groupErrors := groupChanges.Commit(consumer)

	// Print result
	if groupErrors.Amount() == 0 {
		fmt.Println("All groups actions performed!")
	} else {
		fmt.Printf("%d out of %d group actions performed\n", groupChanges.Amount()-groupErrors.Amount(), groupChanges.Amount())
		fmt.Print(groupErrors.String())
	}
	if userErrors.Amount() == 0 {
		fmt.Println("All users actions performed!")
	} else {
		fmt.Printf("%d out of %d group actions performed\n", userChanges.Amount()-userErrors.Amount(), groupChanges.Amount())
		fmt.Print(userErrors.String())
	}

	return userErrors.Amount() == 0 && groupErrors.Amount() == 0
}
