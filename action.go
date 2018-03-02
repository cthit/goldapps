package goldapps

type GroupUpdate struct {
	Before Group
	After  Group
}

type Actions struct {
	Updates   []GroupUpdate
	Additions []Group
	Deletions []Group
}

// Commits a set of actions to a service.
// Returns all actions performed and a error if not all actions could be performed for some reason.
func (actions Actions) Commit(service GroupUpdateService) (Actions, error) {

	performedActions := Actions{}

	for _, update := range actions.Updates {
		err := service.UpdateGroup(update)
		if err != nil {
			return performedActions, err
		}

		performedActions.Updates = append(performedActions.Updates, update)
	}

	for _, group := range actions.Additions {
		err := service.AddGroup(group)
		if err != nil {
			return performedActions, err
		}

		performedActions.Additions = append(performedActions.Additions, group)
	}

	for _, group := range actions.Deletions {
		err := service.DeleteGroup(group)
		if err != nil {
			return performedActions, err
		}

		performedActions.Deletions = append(performedActions.Deletions, group)
	}

	return performedActions, nil
}

// Determines actions required to make the "old" group list look as the "new" group list.
// Returns a list with those actions.
func ActionsRequired(old []Group, new []Group) Actions {
	requiredActions := Actions{}

	for _, newGroup := range new {

		exists := false
		for _, oldGroup := range old {
			if newGroup.Email == oldGroup.Email {
				exists = true
				if !newGroup.equals(oldGroup) { // Group exists but is modified
					requiredActions.Updates = append(requiredActions.Updates, GroupUpdate{
						Before: oldGroup,
						After:  newGroup,
					})
				}
				break
			}
		}

		if !exists { // Group does not exist in old list
			requiredActions.Additions = append(requiredActions.Additions, newGroup)
		}
	}

	for _, oldGroup := range old {

		exists := false
		for _, newGroup := range new {
			if oldGroup.Email == newGroup.Email {
				exists = true
				break
			}
		}

		if !exists { // Old list has group but the new list doesn't
			requiredActions.Deletions = append(requiredActions.Deletions, oldGroup)
		}

	}

	return requiredActions
}
