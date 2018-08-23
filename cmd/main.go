package main

import (
	"fmt"
	"../../goldapps"
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

	// Collect users and groups
	var providerUsers goldapps.Users
	var consumerUsers goldapps.Users
	var providerGroups goldapps.Groups
	var consumerGroups goldapps.Groups
	if !flags.onlyUsers {
		fmt.Println("Collecting groups from the provider...")
		providerGroups = collectGroups(provider)
		fmt.Println("Collecting groups from the consumer...")
		consumerGroups = collectGroups(consumer)
	}
	if !flags.onlyGroups {
		fmt.Println("Collecting users from the provider...")
		providerUsers = collectUsers(provider)
		fmt.Println("Collecting users from the consumer...")
		consumerUsers = collectUsers(consumer)
	}

	// Get and process additions
	providerGroups, providerUsers = addAdditions(providerGroups, providerUsers)

	// Check for and handle duplicates
	providerUsers, providerGroups = goldapps.RemoveDuplicates(providerUsers, providerGroups)

	// Get changes to make
	groupChanges := goldapps.GroupActions{}
	if !flags.onlyUsers {
		fmt.Println("Colculating difference between the consumer and provider groups.")
		proposedGroupChanges := goldapps.GroupActionsRequired(consumerGroups, providerGroups)
		groupChanges = getGroupChanges(proposedGroupChanges)
	}
	userChanges := goldapps.UserActions{}
	if !flags.onlyGroups {
		fmt.Println("Colculating difference between the consumer and provider users.")
		proposedUserChanges := goldapps.UserActionsRequired(consumerUsers, providerUsers)
		userChanges = getUserChanges(proposedUserChanges)
	}

	// Ask for confirmation if we are in interactive mode
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

	// Stop application if dryrun
	if flags.dryRun {
		fmt.Println("Done! (No changes made, dryrun) Stopping application...")
		return
	}

	// Commit changes
	groupErrors := groupChanges.Commit(consumer)
	userErrors := userChanges.Commit(consumer)

	// Print result
	if groupErrors.Amount() == 0 {
		fmt.Println("All groups actions performed!")
	} else {
		fmt.Printf("&d out of %d group actions performed\n", groupChanges.Amount()-groupErrors.Amount(), )
		fmt.Print(groupErrors.String())
	}
	if userErrors.Amount() == 0 {
		fmt.Println("All groups actions performed!")
	} else {
		fmt.Printf("&d out of %d group actions performed\n", userChanges.Amount()-userErrors.Amount(), )
		fmt.Print(userErrors.String())
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

