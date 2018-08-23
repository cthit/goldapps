package goldapps

import (
	"fmt"
	"bytes"
)

// Set of action to be performed on a set of users
type UserActions struct {
	Updates   []UserUpdate
	Additions []User
	Deletions []User
}
func (actions UserActions) Amount() int {
	return len(actions.Additions) + len(actions.Deletions) + len(actions.Updates)
}

// Set of actions that could not be performed with accompanying errors
type UserActionErrors struct {
	Updates   []UserUpdateError
	Additions []UserAddOrDelError
	Deletions []UserAddOrDelError
}
type UserUpdateError struct {
	Action UserUpdate
	Error  error
}
type UserAddOrDelError struct {
	Action User
	Error  error
}
func (actions UserActionErrors) Amount() int {
	return len(actions.Additions) + len(actions.Deletions) + len(actions.Updates)
}
func (actions UserActionErrors) String() string {
	builder := bytes.Buffer{}
	for _,deletion := range actions.Deletions  {
		builder.WriteString(fmt.Sprintf("Deletion of user \"%s\" failed with error %s\n", deletion.Action.Cid, deletion.Error.Error()))
	}
	for _,update := range actions.Updates  {
		builder.WriteString(fmt.Sprintf("Update of user \"%s\" failed with error %s\n", update.Action.After.Cid, update.Error.Error()))
	}
	for _,addition := range actions.Additions  {
		builder.WriteString(fmt.Sprintf("Addition of user \"%s\" failed with error %s\n", addition.Action.Cid, addition.Error.Error()))
	}
	return builder.String()
}

// Data struct representing how a user should look before and after an update
// Allows for efficient updates as application doesn't have to re-upload whole user
type UserUpdate struct {
	Before User
	After  User
}

// Commits a set of actions to a service.
// Returns all actions performed and a error if not all actions could be performed for some reason.
func (actions UserActions) Commit(service UpdateService) UserActionErrors {

	errors := UserActionErrors{}

	if len(actions.Deletions) > 0 {
		fmt.Println("(Users) Performing deletions")
		printProgress(0, len(actions.Deletions), 0)
		for deletionsIndex, user := range actions.Deletions {
			err := service.DeleteUser(user)
			if err != nil {
				// Save error
				errors.Deletions = append(errors.Deletions, UserAddOrDelError{Action: user, Error: err})
			}
			printProgress(deletionsIndex+1, len(actions.Deletions), len(errors.Deletions))
		}
	}

	if len(actions.Updates) > 0 {
		fmt.Println("(USers) Performing updates")
		printProgress(0, len(actions.Updates), 0)
		for updatesIndex, update := range actions.Updates {
			err := service.UpdateUser(update)
			if err != nil {
				// Save error
				errors.Updates = append(errors.Updates, UserUpdateError{Action: update, Error: err})
			}
			printProgress(updatesIndex+1, len(actions.Updates), len(errors.Updates))
		}
	}

	if len(actions.Additions) > 0 {
		fmt.Println("(Groups) Performing additions")
		printProgress(0, len(actions.Additions), 0)
		for additionsIndex, user := range actions.Additions {
			err := service.AddUser(user)
			if err != nil {
				// Save error
				errors.Additions = append(errors.Additions, UserAddOrDelError{Action: user, Error: err})
			}
			printProgress(additionsIndex+1, len(actions.Additions), len(errors.Additions))
		}
	}

	return errors
}

// Determines actions required to make the "old" user list look as the "new" user list.
// Returns a list with those actions.
func UserActionsRequired(old []User, new []User) UserActions {
	requiredActions := UserActions{}

	for _, newUser := range new {
		exists := false
		for _, oldUser := range old {
			if newUser.Same(oldUser) {
				// User exists
				exists = true
				// check if user has to be updates
				if !newUser.Equals(oldUser) {
					// Add User update
					requiredActions.Updates = append(requiredActions.Updates, UserUpdate{
						Before: oldUser,
						After:  newUser,
					})
				}
				break
			}
		}

		// Add user creation action if user doesn't exist
		if !exists {
			requiredActions.Additions = append(requiredActions.Additions, newUser)
		}
	}

	for _, oldUser := range old {
		// check if user should be removed
		removed := true
		for _, newUser := range new {
			if oldUser.Same(newUser) {
				removed = false
				break
			}
		}

		if removed {
			// Add user deletion action
			requiredActions.Deletions = append(requiredActions.Deletions, oldUser)
		}
	}

	return requiredActions
}
