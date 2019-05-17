package services

import (
	"github.com/cthit/goldapps/internal/pkg/model"
)

type CollectionService interface {
	GetGroups() ([]model.Group, error)
	GetUsers() ([]model.User, error)
}

type UpdateService interface {
	DeleteGroup(model.Group) error
	UpdateGroup(model.GroupUpdate) error
	AddGroup(model.Group) error
	DeleteUser(model.User) error
	UpdateUser(model.UserUpdate) error
	AddUser(model.User) error
	CollectionService
}
