package admin

import (
	"google.golang.org/api/admin/directory/v1" // Imports as admin

	"io/ioutil"

	"bytes"
	"fmt"
	"github.com/cthit/goldapps"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

type GoogleService struct {
	service *admin.Service
}

func NewGoogleService(keyPath string, adminMail string) (*GoogleService, error) {

	jsonKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	// Parse jsonKey
	config, err := google.JWTConfigFromJSON(jsonKey, Scopes()...)
	if err != nil {
		return nil, err
	}

	// Why do I need this??
	config.Subject = adminMail

	// Create a http client
	client := config.Client(context.Background())

	service, err := admin.New(client)
	if err != nil {
		return nil, err
	}

	gs := &GoogleService{
		service: service,
	}

	return gs, nil
}

func (s GoogleService) DeleteGroup(group goldapps.Group) error {
	return s.deleteGroup(group.Email)
}

func (s GoogleService) UpdateGroup(groupUpdate goldapps.GroupUpdate) error {
	new := admin.Group{
		Email: groupUpdate.Before.Email,
	}

	// Add all new members
	for _, member := range groupUpdate.After.Members {
		exists := false
		for _, existingMember := range groupUpdate.Before.Members {
			if member == existingMember {
				exists = true
				break
			}
		}
		if !exists {
			err := s.addMember(groupUpdate.Before.Email, member)
			if err != nil {
				return err
			}
		}
	}

	// Remove all old members
	for _, existingMember := range groupUpdate.Before.Members {
		keep := false
		for _, member := range groupUpdate.After.Members {
			if existingMember == member {
				keep = true
				break
			}
		}
		if !keep {
			err := s.deleteMember(groupUpdate.Before.Email, existingMember)
			if err != nil {
				return err
			}
		}
	}

	// Add all new aliases
	for _, alias := range groupUpdate.After.Aliases {
		exists := false
		for _, existingAlias := range groupUpdate.Before.Aliases {
			if alias == existingAlias {
				exists = true
				break
			}
		}
		if !exists {
			err := s.addAlias(groupUpdate.Before.Email, alias)
			if err != nil {
				return err
			}
		}
	}

	// Remove all old aliases
	for _, existingAlias := range groupUpdate.Before.Aliases {
		keep := false
		for _, alias := range groupUpdate.After.Aliases {
			if existingAlias == alias {
				keep = true
				break
			}
		}
		if !keep {
			err := s.deleteAlias(groupUpdate.Before.Email, existingAlias)
			if err != nil {
				return err
			}
		}
	}

	return s.updateGroup(new)
}

func (s GoogleService) AddGroup(group goldapps.Group) error {
	new := admin.Group{
		Email: group.Email,
	}

	err := s.addGroup(new)
	if err != nil {
		return err
	}

	// Add members
	for _, member := range group.Members {
		err = s.addMember(group.Email, member)
		if err != nil {
			return err
		}
	}

	// Add Aliases
	for _, alias := range group.Aliases {
		err = s.addAlias(group.Email, alias)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s GoogleService) GetGroups() ([]goldapps.Group, error) {

	adminGroups, err := s.getGroups("my_customer")
	if err != nil {
		return nil, err
	}

	groups := make([]goldapps.Group, len(adminGroups))
	for i, group := range adminGroups {

		p := (i * 100) / len(groups)

		builder := bytes.Buffer{}
		for i := 0; i < 100; i++ {
			if i < p {
				builder.WriteByte('=')
			} else if i == p {
				builder.WriteByte('>')
			} else {
				builder.WriteByte(' ')
			}

		}

		fmt.Printf("\rProgress: [%s] %d/%d", builder.String(), i+1, len(groups))

		members, err := s.getMembers(group.Email)
		if err != nil {
			return nil, err
		}

		groups[i] = goldapps.Group{
			Email:   group.Email,
			Members: members,
			Aliases: group.Aliases,
		}
	}
	fmt.Printf("\rDone!\n")

	return groups, nil

}

func (s GoogleService) getGroups(customer string) ([]admin.Group, error) {
	groups, err := s.service.Groups.List().Customer(customer).Do()
	if err != nil {
		return nil, err
	}

	for groups.NextPageToken != "" {
		newGroups, err := s.service.Groups.List().Customer(customer).PageToken(groups.NextPageToken).Do()
		if err != nil {
			return nil, err
		}

		groups.Groups = append(groups.Groups, newGroups.Groups...)
		groups.NextPageToken = newGroups.NextPageToken
	}

	result := make([]admin.Group, len(groups.Groups))
	for i, group := range groups.Groups {
		result[i] = *group
	}

	return result, nil
}

func (s GoogleService) getMembers(email string) ([]string, error) {
	members, err := s.service.Members.List(email).Do()
	if err != nil {
		return nil, err
	}

	result := make([]string, len(members.Members))
	for i, member := range members.Members {
		result[i] = member.Email
	}

	return result, nil
}

func (s GoogleService) getGroup(email string) (admin.Group, error) {
	group, err := s.service.Groups.Get(email).Do()

	return *group, err
}

func (s GoogleService) addGroup(group admin.Group) error {
	_, err := s.service.Groups.Insert(&group).Do()
	return err
}

func (s GoogleService) updateGroup(group admin.Group) error {
	_, err := s.service.Groups.Update(group.Email, &group).Do()
	return err
}

func (s GoogleService) deleteGroup(email string) error {
	err := s.service.Groups.Delete(email).Do()
	return err
}

func (s GoogleService) deleteMember(groupEmail string, member string) error {
	return s.service.Members.Delete(groupEmail, member).Do()
}

func (s GoogleService) addMember(groupEmail string, memberEmail string) error {
	_, err := s.service.Members.Insert(groupEmail, &admin.Member{Email: memberEmail}).Do()
	return err
}

func (s GoogleService) deleteAlias(groupEmail string, alias string) error {
	return s.service.Groups.Aliases.Delete(groupEmail, alias).Do()
}

func (s GoogleService) addAlias(groupEmail string, alias string) error {
	_, err := s.service.Groups.Aliases.Insert(groupEmail, &admin.Alias{Alias: alias}).Do()
	return err
}
