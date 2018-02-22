package goldapps

type GroupUpdate struct {
	Before		Group
	After		Group
}


type Actions struct {
	Updates		[]GroupUpdate
	Additions	[]Group
	Deletions	[]Group
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

func ActionsRequired(from []Group, to []Group) (Actions) {
	requiredActions := Actions{}

	

	for _,toGroup := range to {
		for ; ;  {
			
		}
	}


	return requiredActions
}