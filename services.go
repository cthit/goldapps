package goldapps

type UpdateService interface {
	DeleteGroup(Group) error
	UpdateGroup(GroupUpdate) error
	AddGroup(Group) error
	DeleteUser(User) error
	UpdateUser(UserUpdate) error
	AddUser(User) error
	CollectionService
}

type CollectionService interface {
	GetGroups() ([]Group, error)
	GetUsers() ([]User, error)
}
