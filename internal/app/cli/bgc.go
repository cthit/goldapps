package cli

import (
	"fmt"

	"github.com/cthit/goldapps/internal/pkg/actions"
	"github.com/cthit/goldapps/internal/pkg/duplicates"
	"github.com/cthit/goldapps/internal/pkg/model"
)

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
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

func Run() {

	fmt.Println("Setting up providers")
	provider := getProvider()

	fmt.Println("Setting up services")
	consumer := getConsumer()

	// Collect users and groups
	var providerUsers model.Users
	var consumerUsers model.Users
	var providerGroups model.Groups
	var consumerGroups model.Groups
	if !flags.onlyUsers {
		fmt.Println("Collecting groups from the providers...")
		providerGroups = collectGroups(provider)
		fmt.Println("Collecting groups from the services...")
		consumerGroups = collectGroups(consumer)
	}
	if !flags.onlyGroups {
		fmt.Println("Collecting users from the providers...")
		providerUsers = collectUsers(provider)
		fmt.Println("Collecting users from the services...")
		consumerUsers = collectUsers(consumer)
	}

	// Get and process additions
	providerGroups, providerUsers = addAdditions(providerGroups, providerUsers)

	// Check for and handle duplicates
	providerUsers, providerGroups = duplicates.RemoveDuplicates(providerUsers, providerGroups)

	// Get changes to make
	groupChanges := actions.GroupActions{}
	if !flags.onlyUsers {
		fmt.Println("Colculating difference between the services and providers groups.")
		proposedGroupChanges := actions.GroupActionsRequired(consumerGroups, providerGroups)
		groupChanges = getGroupChanges(proposedGroupChanges)
	}
	userChanges := actions.UserActions{}
	if !flags.onlyGroups {
		fmt.Println("Colculating difference between the services and providers users.")
		proposedUserChanges := actions.UserActionsRequired(consumerUsers, providerUsers)
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
}

func has(member string, members []string) bool {
	for _, v := range members {
		if model.CompareEmails(v, member) {
			return true
		}
	}
	return false
}

func printGroupDiff(before model.Group, after model.Group) {
	fmt.Printf("\n\tUpdate: ")
	if before.Email != after.Email {
		fmt.Printf("\t%s -> %s\n", before.Email, after.Email)
	} else {
		fmt.Printf("\t\t%s\n", after.Email)
	}
	if before.Type != after.Email {
		fmt.Printf("\t\t\t%s -> %s\n", before.Type, after.Type)
	}

	added := []string{}
	deleted := []string{}

	for _, member := range before.Members {
		if !has(member, after.Members) {
			deleted = append(deleted, member)
		}
	}

	for _, member := range after.Members {
		if !has(member, before.Members) {
			added = append(added, member)
		}
	}

	for _, del := range deleted {
		fmt.Printf("\t\t\t- %s\n", del)
	}

	for _, add := range added {
		fmt.Printf("\t\t\t+ %s\n", add)
	}

	fmt.Printf("\t\t\tAliases: %v -> %v\n", before.Aliases, after.Aliases)
}

func getGroupChanges(proposedChanges actions.GroupActions) actions.GroupActions {
	if !flags.interactive && flags.noInteraction {
		fmt.Printf(
			"(Groups) Automaticly accepting %d addition, %d deletions and %d updates\n",
			len(proposedChanges.Additions),
			len(proposedChanges.Deletions),
			len(proposedChanges.Updates),
		)
	} else {
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
				printGroupDiff(update.Before, update.After)
			}
			add := askBool(
				fmt.Sprintf("(Groups) Do you want to commit those %d updates?", len(proposedChanges.Updates)),
				true,
			)
			if !add {
				proposedChanges.Updates = nil
			}
		}

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
	}
	return proposedChanges
}

func getUserChanges(proposedChanges actions.UserActions) actions.UserActions {
	if !flags.interactive && flags.noInteraction {
		fmt.Printf(
			"(Users) Automaticly accepting %d addition, %d deletions and %d updates\n",
			len(proposedChanges.Additions),
			len(proposedChanges.Deletions),
			len(proposedChanges.Updates),
		)
	} else {
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

	}
	return proposedChanges
}
