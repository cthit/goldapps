package goldapps

import (
	"fmt"
)

type UserUpdate struct {
	Before User
	After  User
}

type UserActions struct {
	Updates   []UserUpdate
	Additions []User
	Deletions []User
}

// Commits a set of actions to a service.
// Returns all actions performed and a error if not all actions could be performed for some reason.
func (actions UserActions) Commit(service UpdateService) (UserActions, error) {

	performedActions := UserActions{}

	if len(actions.Updates) > 0 {
		fmt.Println("(Users) Performing updates")
	}
	for _, update := range actions.Updates {
		err := service.UpdateUser(update)
		if err != nil {
			fmt.Println()
			return performedActions, err
		}

		performedActions.Updates = append(performedActions.Updates, update)
		printProgress(len(performedActions.Updates), len(actions.Updates))
	}

	if len(actions.Additions) > 0 {
		fmt.Println("(Users) Performing additions")
	}
	for _, user := range actions.Additions {
		err := service.AddUser(user)
		if err != nil {
			fmt.Println()
			return performedActions, err
		}

		performedActions.Additions = append(performedActions.Additions, user)
		printProgress(len(performedActions.Additions), len(actions.Additions))
	}

	if len(actions.Deletions) > 0 {
		fmt.Println("(Users) Performing deletions")
	}
	for _, user := range actions.Deletions {
		err := service.DeleteUser(user)
		if err != nil {
			fmt.Println()
			return performedActions, err
		}

		performedActions.Deletions = append(performedActions.Deletions, user)
		printProgress(len(performedActions.Deletions), len(actions.Deletions))
	}

	return performedActions, nil
}

// Determines actions required to make the "old" user list look as the "new" user list.
// Returns a list with those actions.
func UserActionsRequired(old []User, new []User) UserActions {
	requiredActions := UserActions{}

	for _, newUser := range new {

		exists := false
		for _, oldUser := range old {
			if newUser.Cid == oldUser.Cid {
				exists = true
				if !newUser.equals(oldUser) { // User exists but is modified
					requiredActions.Updates = append(requiredActions.Updates, UserUpdate{
						Before: oldUser,
						After:  newUser,
					})
				}
				break
			}
		}

		if !exists { // User does not exist in old list
			requiredActions.Additions = append(requiredActions.Additions, newUser)
		}
	}

	for _, oldUser := range old {

		exists := false
		for _, newUser := range new {
			if oldUser.Cid == newUser.Cid {
				exists = true
				break
			}
		}

		if !exists { // Old list has user but the new list doesn't
			requiredActions.Deletions = append(requiredActions.Deletions, oldUser)
		}

	}

	return requiredActions
}
