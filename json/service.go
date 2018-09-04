package json

import (
	"encoding/json"
	"fmt"
	"github.com/cthit/goldapps"
	"io/ioutil"
	"os"
)

type Service struct {
	path string
}

type dataObject struct {
	Groups []goldapps.Group `json:"groups"`
	Users  []goldapps.User  `json:"users"`
}

func (s Service) DeleteUser(user goldapps.User) error {
	groups, err := s.GetGroups()
	if err != nil {
		return err
	}
	users, err := s.GetUsers()
	if err != nil {
		return err
	}

	for i, u := range users {
		if u.Cid == user.Cid {
			err = s.save(dataObject{
				groups,
				append(users[:i], users[i+1:]...),
			})
			return err
		}
	}
	return fmt.Errorf("user not found %v", user)
}

func (s Service) UpdateUser(update goldapps.UserUpdate) error {
	groups, err := s.GetGroups()
	if err != nil {
		return err
	}
	users, err := s.GetUsers()
	if err != nil {
		return err
	}

	for i, u := range users {
		if u.Cid == update.Before.Cid {
			err = s.save(dataObject{
				groups,
				append(append(users[:i], update.After), users[i+1:]...),
			})
			return err
		}
	}
	return fmt.Errorf("user not found %v", update.Before)
}

func (s Service) AddUser(user goldapps.User) error {
	groups, err := s.GetGroups()
	if err != nil {
		return err
	}
	users, err := s.GetUsers()
	if err != nil {
		return err
	}

	users = append(users, user)

	err = s.save(dataObject{
		groups,
		users,
	})
	return err
}

func (s Service) GetUsers() ([]goldapps.User, error) {

	data, err := s.get()
	if err != nil {
		return nil, err
	}

	return data.Users, nil
}

func NewJsonService(path string) (Service, error) {

	// Check if file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Create file
		_, err := os.Create("path")
		if err != nil {
			return Service{}, err
		}
		// Write empty object to file
		err = Service{path: path}.save(dataObject{})
		if err != nil {
			return Service{}, err
		}
	}

	return Service{
		path: path,
	}, nil
}

func (s Service) save(data dataObject) error {
	json, _ := json.Marshal(data)

	err := ioutil.WriteFile(s.path, json, 0666)
	return err
}

func (s Service) get() (dataObject, error) {

	bytes, err := ioutil.ReadFile(s.path)
	if err != nil {
		return dataObject{}, err
	}

	var data dataObject
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return dataObject{}, err
	}

	return data, nil
}

func (s Service) DeleteGroup(group goldapps.Group) error {
	groups, err := s.GetGroups()
	if err != nil {
		return err
	}
	users, err := s.GetUsers()
	if err != nil {
		return err
	}

	for i, g := range groups {
		if g.Email == group.Email {
			err = s.save(dataObject{append(groups[:i], groups[i+1:]...),
				users,
			})
			return err
		}
	}
	return fmt.Errorf("group not found %v", group)
}

func (s Service) UpdateGroup(groupUpdate goldapps.GroupUpdate) error {
	groups, err := s.GetGroups()
	if err != nil {
		return err
	}
	users, err := s.GetUsers()
	if err != nil {
		return err
	}

	for i, g := range groups {
		if g.Email == groupUpdate.Before.Email {
			err = s.save(dataObject{
				append(append(groups[:i], groupUpdate.After), groups[i+1:]...),
				users,
			})
			return err
		}
	}
	return fmt.Errorf("group not found %v", groupUpdate.Before)
}

func (s Service) AddGroup(group goldapps.Group) error {
	groups, err := s.GetGroups()
	if err != nil {
		return err
	}
	users, err := s.GetUsers()
	if err != nil {
		return err
	}

	groups = append(groups, group)

	err = s.save(dataObject{
		groups,
		users,
	})
	return err
}

func (s Service) GetGroups() ([]goldapps.Group, error) {

	data, err := s.get()
	if err != nil {
		return nil, err
	}

	return data.Groups, nil
}
