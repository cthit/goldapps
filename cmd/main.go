package main

import (
	"fmt"
	"github.com/cthit/goldapps"
	"github.com/cthit/goldapps/admin"
	"github.com/cthit/goldapps/json"
	"github.com/spf13/viper"
	"regexp"
)

func init() {
	loadFlags()

	err := loadConfig()
	if err != nil {
		fmt.Println("Failed to load config.")
		panic(err)
	}
	fmt.Println("Loaded config.")
}

func main() {

	fmt.Println("Setting up provider")
	provider := getProvider()

	fmt.Println("Setting up consumer")
	consumer := getConsumer()

	fmt.Println("Collecting groups from the provider...")
	providerGroups, err := provider.GetGroups()
	if err != nil {
		fmt.Println("Failed to collect groups from provider")
		panic(err)
	}
	fmt.Printf("%d groups collected.\n", len(providerGroups))

	fmt.Println("Collecting users from the provider...")
	providerUsers, err := provider.GetUsers()
	if err != nil {
		fmt.Println("Failed to collect users from provider")
		panic(err)
	}
	fmt.Printf("%d users collected.\n", len(providerUsers))

	fmt.Println("Collecting groups from the consumer...")
	consumerGroups, err := consumer.GetGroups()
	if err != nil {
		fmt.Println("Failed to collect groups from consumer")
		panic(err)
	}
	fmt.Printf("%d groups collected.\n", len(consumerGroups))

	fmt.Println("Collecting users from the consumer...")
	consumerUsers, err := consumer.GetUsers()
	if err != nil {
		fmt.Println("Failed to collect users from consumer")
		panic(err)
	}
	fmt.Printf("%d users collected.\n", len(consumerUsers))

	fmt.Println("Colculating difference between the consumer and provider groups.")
	proposedGroupChanges := goldapps.GroupActionsRequired(consumerGroups, providerGroups)
	groupChanges := getGroupChanges(proposedGroupChanges)

	fmt.Println("Colculating difference between the consumer and provider users.")
	proposedUserChanges := goldapps.UserActionsRequired(consumerUsers, providerUsers)
	userChanges := getUserChanges(proposedUserChanges)

	if flags.interactive {
		proceed := askBool(
			fmt.Sprintf(
				"Are you sure you want to commit these (groups + users) additions(%d + %d), deletions(%d + %d) and updates(%d + %d)?",
				len(groupChanges.Additions),
				len(userChanges.Additions),
				len(groupChanges.Deletions),
				len(userChanges.Deletions),
				len(groupChanges.Updates),
				len(userChanges.Updates),
			),
			true,
		)
		if !proceed {
			fmt.Println("Done! (No changes made) Stopping application...")
			return
		}
	}
	if flags.dryRun {
		fmt.Println("Done! (No changes made, dryrun) Stopping application...")
		return
	}

	groupChangesPerformed, err1 := groupChanges.Commit(consumer)
	userChangesPerformed, err2 := userChanges.Commit(consumer)
	if err1 == nil && err2 == nil {
		fmt.Println("All actions performed!")
		return
	} else {
		fmt.Println("All actions could not be performed...")
		fmt.Printf("\t For groups:\n")
		fmt.Printf("\t\t Performed %d out of %d Additions\n", len(groupChangesPerformed.Additions), len(groupChanges.Additions))
		fmt.Printf("\t\t Performed %d out of %d Deletions\n", len(groupChangesPerformed.Deletions), len(groupChanges.Deletions))
		fmt.Printf("\t\t Performed %d out of %d Updates\n", len(groupChangesPerformed.Updates), len(groupChanges.Updates))
		fmt.Printf("\t\t Error: %s", err1.Error())
		fmt.Printf("\t For users:\n")
		fmt.Printf("\t\t Performed %d out of %d Additions\n", len(userChangesPerformed.Additions), len(userChanges.Additions))
		fmt.Printf("\t\t Performed %d out of %d Deletions\n", len(userChangesPerformed.Deletions), len(userChanges.Deletions))
		fmt.Printf("\t\t Performed %d out of %d Updates\n", len(userChangesPerformed.Updates), len(userChanges.Updates))
		fmt.Printf("\t\t Error: %s", err2.Error())
	}
}

func getGroupChanges(proposedChanges goldapps.GroupActions) goldapps.GroupActions {
	if !flags.interactive && flags.noInteraction {
		fmt.Printf(
			"(Groups) Automaticly accepting %d addition, %d deletions and %d updates\n",
			len(proposedChanges.Additions),
			len(proposedChanges.Deletions),
			len(proposedChanges.Updates),
		)
	} else {
		// Handle additions
		fmt.Printf("(Groups) Additions (%d):\n", len(proposedChanges.Additions))
		if len(proposedChanges.Additions) > 0 {
			for _, group := range proposedChanges.Additions {
				fmt.Printf("\t%v\n", group)
			}
			add := askBool(
				fmt.Sprintf("(Groups) Do you want to commit those %d additions?", len(proposedChanges.Additions)),
				true,
			)
			if !add {
				proposedChanges.Additions = nil
			}
		}

		// Handle Deletions
		fmt.Printf("(Groups) Deletions (%d):\n", len(proposedChanges.Deletions))
		if len(proposedChanges.Deletions) > 0 {
			for _, group := range proposedChanges.Deletions {
				fmt.Printf("\t%v\n", group)
			}
			add := askBool(
				fmt.Sprintf("(Groups) Do you want to commit those %d deletions?", len(proposedChanges.Deletions)),
				true,
			)
			if !add {
				proposedChanges.Deletions = nil
			}
		}

		// Handle changes
		fmt.Printf("(Groups) Changes (%d):\n", len(proposedChanges.Updates))
		if len(proposedChanges.Updates) > 0 {
			for _, update := range proposedChanges.Updates {
				fmt.Printf("\tUpdate:\n")
				fmt.Printf("\t\tFrom:\n")
				fmt.Printf("\t\t\t%v\n", update.Before)
				fmt.Printf("\t\tTo:\n")
				fmt.Printf("\t\t\t%v\n", update.After)
			}
			add := askBool(
				fmt.Sprintf("(Groups) Do you want to commit those %d updates?", len(proposedChanges.Updates)),
				true,
			)
			if !add {
				proposedChanges.Updates = nil
			}
		}
	}
	return proposedChanges
}

func getUserChanges(proposedChanges goldapps.UserActions) goldapps.UserActions {
	if !flags.interactive && flags.noInteraction {
		fmt.Printf(
			"(Users) Automaticly accepting %d addition, %d deletions and %d updates\n",
			len(proposedChanges.Additions),
			len(proposedChanges.Deletions),
			len(proposedChanges.Updates),
		)
	} else {
		// Handle additions
		fmt.Printf("(Users) Additions (%d):\n", len(proposedChanges.Additions))
		if len(proposedChanges.Additions) > 0 {
			for _, user := range proposedChanges.Additions {
				fmt.Printf("\t%v\n", user)
			}
			add := askBool(
				fmt.Sprintf("(Users) Do you want to commit those %d additions?", len(proposedChanges.Additions)),
				true,
			)
			if !add {
				proposedChanges.Additions = nil
			}
		}

		// Handle Deletions
		fmt.Printf("(Users) Deletions (%d):\n", len(proposedChanges.Deletions))
		if len(proposedChanges.Deletions) > 0 {
			for _, user := range proposedChanges.Deletions {
				fmt.Printf("\t%v\n", user)
			}
			add := askBool(
				fmt.Sprintf("(Users) Do you want to commit those %d deletions?", len(proposedChanges.Deletions)),
				true,
			)
			if !add {
				proposedChanges.Deletions = nil
			}
		}

		// Handle changes
		fmt.Printf("(Users) Changes (%d):\n", len(proposedChanges.Updates))
		if len(proposedChanges.Updates) > 0 {
			for _, update := range proposedChanges.Updates {
				fmt.Printf("\tUpdate:\n")
				fmt.Printf("\t\tFrom:\n")
				fmt.Printf("\t\t\t%v\n", update.Before)
				fmt.Printf("\t\tTo:\n")
				fmt.Printf("\t\t\t%v\n", update.After)
			}
			add := askBool(
				fmt.Sprintf("(Users) Do you want to commit those %d updates?", len(proposedChanges.Updates)),
				true,
			)
			if !add {
				proposedChanges.Updates = nil
			}
		}
	}
	return proposedChanges
}

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
		provider, err := NewLdapService()
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
