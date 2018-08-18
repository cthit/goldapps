package json

import (
	"encoding/json"
	"fmt"
	"../../goldapps"
	"io/ioutil"
)

type Service struct {
	path string
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
			err = s.save(data{
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
			err = s.save(data{
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

	err = s.save(data{
		groups,
		users,
	})
	return err
}

func (s Service) GetUsers() ([]goldapps.User, error) {

	bytes, err := ioutil.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	var data data
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return data.Users, nil
}

func NewJsonService(path string) (Service, error) {
	return Service{
		path: path,
	}, nil
}

func (s Service) save(data data) error {
	json, _ := json.Marshal(data)

	err := ioutil.WriteFile(s.path, json, 0666)
	return err
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
			err = s.save(data{append(groups[:i], groups[i+1:]...),
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
			err = s.save(data{
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

	err = s.save(data{
		groups,
		users,
	})
	return err
}

func (s Service) GetGroups() ([]goldapps.Group, error) {

	bytes, err := ioutil.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	var data data
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return data.Groups, nil
}
