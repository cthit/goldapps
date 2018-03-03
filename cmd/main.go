package main

import (
	"fmt"
	"github.com/cthit/goldapps"
	"github.com/cthit/goldapps/admin"
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

	fmt.Println("Collecting groups from the consumer...")
	consumerGroups, err := consumer.GetGroups()
	if err != nil {
		fmt.Println("Failed to collect groups from consumer")
		panic(err)
	}
	fmt.Printf("%d groups collected.\n", len(consumerGroups))

	fmt.Println("Colculating difference between the consumer and provider.")
	proposedChanges := goldapps.ActionsRequired(consumerGroups, providerGroups)
	changes := getChanges(proposedChanges)

	if flags.interactive {
		proceed := askBool(
			fmt.Sprintf(
				"Are you sure you want to commit these additions(%d), deletions(%d) and updates(%d)?",
				len(changes.Additions),
				len(changes.Deletions),
				len(changes.Updates),
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

	performed, err := changes.Commit(consumer)
	if err == nil {
		fmt.Println("All actions performed!")
		return
	} else {
		fmt.Println("All actions could not be performed...")
		fmt.Printf("/t Performed %d out of %d Additions/n", len(performed.Additions), len(changes.Additions))
		fmt.Printf("/t Performed %d out of %d Deletions/n", len(performed.Deletions), len(changes.Deletions))
		fmt.Printf("/t Performed %d out of %d Updates/n", len(performed.Updates), len(changes.Updates))
	}

	/*data, err := json.Marshal(consumerGroups)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("gapps_groups.json", data, 0777)
	if err != nil {
		panic(err)
	}*/
}

func getChanges(proposedChanges goldapps.Actions) goldapps.Actions {
	if !flags.interactive && flags.noInteraction {
		fmt.Printf(
			"Automaticly accepting %d addition, %d deletions and %d updates\n",
			len(proposedChanges.Additions),
			len(proposedChanges.Deletions),
			len(proposedChanges.Updates),
		)
	} else {
		// Handle additions
		fmt.Printf("Additions (%d):\n", len(proposedChanges.Additions))
		for _, group := range proposedChanges.Additions {
			fmt.Printf("\t%v\n", group)
		}
		add := askBool(
			fmt.Sprintf("Do you want to commit those %d additions?", len(proposedChanges.Additions)),
			true,
		)
		if !add {
			proposedChanges.Additions = nil
		}

		// Handle Deletions
		fmt.Printf("Deletions (%d):\n", len(proposedChanges.Deletions))
		for _, group := range proposedChanges.Deletions {
			fmt.Printf("\t%v\n", group)
		}
		add = askBool(
			fmt.Sprintf("Do you want to commit those %d deletions?", len(proposedChanges.Deletions)),
			true,
		)
		if !add {
			proposedChanges.Deletions = nil
		}

		// Handle changes
		fmt.Printf("Changes (%d):\n", len(proposedChanges.Updates))
		for _, update := range proposedChanges.Updates {
			fmt.Printf("\tUpdate:\n")
			fmt.Printf("\t\tFrom:\n")
			fmt.Printf("\t\t\t%v\n", update.Before)
			fmt.Printf("\t\tTo:\n")
			fmt.Printf("\t\t\t%v\n", update.After)
		}
		add = askBool(
			fmt.Sprintf("Do you want to commit those %d updates?", len(proposedChanges.Updates)),
			true,
		)
		if !add {
			proposedChanges.Updates = nil
		}
	}
	return proposedChanges
}

func getConsumer() goldapps.GroupUpdateService {
	var to string
	if flags.interactive {
		to = askString("Which consumer would you like to use, 'gapps' or '*.json?", "gapps")
	} else {
		to = flags.to
	}

	switch to {
	case "gapps":
		consumer, err := admin.NewGoogleService(
			viper.GetString("gapps.servicekeyfile"),
			viper.GetString("gapps.adminaccount"))
		if err != nil {
			fmt.Println("Failed to create gapps connection.")
			panic(err)
		}
		return consumer
	default:
		isJson, _ := regexp.MatchString(`.+\.json$`, to)
		if isJson {
			panic("Not implemented!")
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

func getProvider() goldapps.GroupService {
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
		provider, err := admin.NewGoogleService(viper.GetString("gappsProvider.servicekeyfile"), viper.GetString("gappsProvider.adminaccount"))
		if err != nil {
			fmt.Println("Failed to create gapps connection, make sure you have setup gappsProvider in the config file.")
			panic(err)
		}
		return provider
	default:
		isJson, _ := regexp.MatchString(`.+\.json$`, from)
		if isJson {
			panic("Not implemented!")
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
