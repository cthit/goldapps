package pkg

import (
	"fmt"
	"regexp"

	"github.com/cthit/goldapps/internal/pkg/model"
	"github.com/cthit/goldapps/internal/pkg/services/json"
)

func AddAdditions(providerGroups model.Groups, providerUsers model.Users, from string) (model.Groups, model.Users) {
	fmt.Println("Collecting additions")
	additionUsers, additionGroups := GetAdditions(from)
	if additionUsers != nil && additionGroups != nil {
		fmt.Printf("%d usersAdditions and %d groupAdditions collected.\n", len(additionUsers), len(additionGroups))

		fmt.Print("Merging groups... ")
		providerGroups = mergeAdditionGroups(additionGroups, providerGroups)
		fmt.Println("Done!")

		fmt.Print("Merging users...  ")
		providerUsers = mergeAdditionalUsers(additionUsers, providerUsers)
		fmt.Println("Done!")
	} else {
		fmt.Println("Skipping additions")
	}
	return providerGroups, providerUsers
}

func GetAdditions(from string) ([]model.User, []model.Group) {
	if from == "" {
		return nil, nil
	}

	isJson, _ := regexp.MatchString(`.+\.json$`, from)
	if !isJson {
		return nil, nil
	}

	provider, _ := json.NewJsonService(from)
	groups, err := provider.GetGroups()
	if err != nil {
		panic(err)
	}
	users, err := provider.GetUsers()
	if err != nil {
		panic(err)
	}
	return users, groups
}

func mergeAdditionGroups(additionGroups model.Groups, providerGroups model.Groups) model.Groups {
	for _, group := range additionGroups {
		found := false
		for i, pgroup := range providerGroups {
			if pgroup.Email == group.Email {
				found = true

				// Add properties if found
				for _, alias := range group.Aliases {
					aliasFound := false
					for _, other := range pgroup.Aliases {
						if other == alias {
							aliasFound = true
						}
					}
					if !aliasFound {
						providerGroups[i].Aliases = append(pgroup.Aliases, alias)
					}
				}

				for _, member := range group.Members {
					memberFound := false
					for _, other := range pgroup.Members {
						if other == member {
							memberFound = true
						}
					}
					if !memberFound {
						providerGroups[i].Members = append(providerGroups[i].Members, member)
					}
				}
			}
		}

		// Otherwise simply append it
		if !found {
			providerGroups = append(providerGroups, group)
		}
	}
	return providerGroups
}
func mergeAdditionalUsers(additionUsers model.Users, providerUsers model.Users) model.Users {
	for _, user := range additionUsers {
		found := false
		for i, pUser := range providerUsers {

			// Add Properties if found, never replace tho
			if pUser.Cid == user.Cid {
				found = true
				if user.Nick != "" {
					providerUsers[i].Nick = user.Nick
				}
				if user.SecondName != "" {
					providerUsers[i].SecondName = user.SecondName
				}
				if user.FirstName != "" {
					providerUsers[i].FirstName = user.FirstName
				}
				if user.Mail != "" {
					providerUsers[i].Mail = user.Mail
				}
			}
		}

		// Just add user if it wasn't found
		if !found {
			providerUsers = append(providerUsers, user)
		}
	}
	return providerUsers
}
