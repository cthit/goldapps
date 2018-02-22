package goldapps

type GroupUpdateService interface {
	DeleteGroup(group Group) (error)
	UpdateGroup(groupUpdate GroupUpdate) (error)
	AddGroup(group Group) (error)
}

type GroupService interface {
	GetGroups() ([]Group, error)
}
