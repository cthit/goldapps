package google

type GroupUpdateService interface {
	UpdateGroups(*[]Group) error
	Groups() (*[]Group, error)
}

type GroupService interface {
	Groups() (*[]Group, error)
}
