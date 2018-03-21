package json

import (
	"github.com/cthit/goldapps"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

type jsonService struct {
	path string
}

func (s jsonService) DeleteUser(goldapps.User) error {
	panic("implement me")
}

func (s jsonService) UpdateUser(goldapps.UserUpdate) error {
	panic("implement me")
}

func (s jsonService) AddUser(goldapps.User) error {
	panic("implement me")
}

func (s jsonService) GetUsers() ([]goldapps.User, error) {
	panic("implement me")
}

func NewJsonService(path string) (jsonService, error) {
	return jsonService{
		path: path,
	}, nil
}

func (s jsonService) save(groups []goldapps.Group) error {
	data, _ := json.Marshal(groups)

	err := ioutil.WriteFile(s.path, data, 0666)
	return err
}

func (s jsonService) DeleteGroup(group goldapps.Group) error {
	groups, err := s.GetGroups()
	if err != nil {
		return err
	}

	for i,g := range groups {
		if g.Email == group.Email{
			err = s.save(append(groups[:i], groups[i+1:]...))
			return err
		}
	}
	return fmt.Errorf("group not found %v", group)
}

func (s jsonService) UpdateGroup(groupUpdate goldapps.GroupUpdate) error {
	groups, err := s.GetGroups()
	if err != nil {
		return err
	}

	for i,g := range groups {
		if g.Email == groupUpdate.Before.Email{
			err = s.save(append(append(groups[:i], groupUpdate.After), groups[i+1:]...))
			return err
		}
	}
	return fmt.Errorf("group not found %v", groupUpdate.Before)
}

func (s jsonService) AddGroup(group goldapps.Group) error {
	groups, err := s.GetGroups()
	if err != nil {
		return err
	}

	groups = append(groups, group)

	err = s.save(groups)
	return err
}

func (s jsonService) GetGroups() ([]goldapps.Group, error) {

	bytes, err := ioutil.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	var data []goldapps.Group
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
