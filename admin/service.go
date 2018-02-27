package admin

import (
	"google.golang.org/api/admin/directory/v1" // Imports as admin

	"io/ioutil"

	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
	"github.com/cthit/goldapps"
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

func (s GoogleService) DeleteGroup(group goldapps.Group) (error) {
	return s.deleteGroup(group.Email)
}

func (s GoogleService) UpdateGroup(groupUpdate goldapps.GroupUpdate) (error) {
	new := admin.Group{
		Email: groupUpdate.Before.Email,
	}

	// Add all new members
	for _,member := range groupUpdate.After.Members {
		exists := false
		for _,existingMember := range groupUpdate.Before.Members {
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
	for _,existingMember := range groupUpdate.Before.Members {
		keep := false
		for _,member := range groupUpdate.After.Members {
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
	for _,alias := range groupUpdate.After.Aliases {
		exists := false
		for _,existingAlias := range groupUpdate.Before.Aliases {
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
	for _,existingAlias := range groupUpdate.Before.Aliases {
		keep := false
		for _,alias := range groupUpdate.After.Aliases {
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

func (s GoogleService) AddGroup(group goldapps.Group) (error) {
	new := admin.Group{
		Email:   group.Email,
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
	panic("Not implemented!")
}

func (s GoogleService) getGroup(email string) (admin.Group, error)  {
	group, err := s.service.Groups.Get(email).Do()

	return *group, err
}

func (s GoogleService) addGroup(group admin.Group) (error) {
	_, err := s.service.Groups.Insert(&group).Do()
	return err
}

func (s GoogleService) updateGroup(group admin.Group) (error) {
	_, err := s.service.Groups.Update(group.Email, &group).Do()
	return err
}

func (s GoogleService) deleteGroup(email string) error {
	err := s.service.Groups.Delete(email).Do()
	return err
}

func (s GoogleService) deleteMember(groupEmail string, member string, ) error {
	return s.service.Members.Delete(groupEmail, member).Do()
}

func (s GoogleService) addMember(groupEmail string, memberEmail string) error {
	_, err := s.service.Members.Insert(groupEmail, &admin.Member{Email:memberEmail}).Do()
	return err
}

func (s GoogleService) deleteAlias(groupEmail string, alias string, ) error {
	return s.service.Groups.Aliases.Delete(groupEmail, alias).Do()
}

func (s GoogleService) addAlias(groupEmail string, alias string) error {
	_, err := s.service.Groups.Aliases.Insert(groupEmail, &admin.Alias{Alias:alias}).Do()
	return err
}



/*
	==============
	Old code below
	==============

 */


 /*
func (s *GoogleService) updateMembers(g admin.Group, members []string) error {
	err := s.cleanupMembers(g, members)
	if err != nil {
		return err
	}

	return s.pushMembers(g, members)
}


func (s *GoogleService) cleanupMembers(g *admin.Group, members []string) error {
	current, err := s.members(g)
	if err != nil {
		return err
	}

	for _, cMember := range current {
		if !memberContains(members, cMember) {
			err = s.deleteMember(cMember, g)
			if err != nil {
				return err
			}
		}
	}
	return nil
}



func (s *GoogleService) pushMembers(g admin.Group, members []string) error {
	for _, member := range members {
		mem, err := s.Service.Members.Get(g.Email, member).Do()
		if err != nil {
			new := &admin.Member{
				Email: member,
			}
			mem, err = s.Service.Members.Insert(g.Email, new).Do()
			if err != nil {
				return err
			}
		} else {
			mem.Email = member
			mem, err = s.Service.Members.Update(g.Email, mem.Email, mem).Do()
		}
	}
	return nil
}
*/


/*
func (s GoogleService) Groups() ([]goldapps.Group, error) {
	groups, err := s.Service.Groups.List().Customer("my_customer").Do()
	if err != nil {
		return nil, err
	}

	for groups.NextPageToken != "" {
		newGroups, err := s.Service.Groups.List().Customer("my_customer").PageToken(groups.NextPageToken).Do()
		if err != nil {
			return nil, err
		}

		groups.Groups = append(groups.Groups, newGroups.Groups...)
		groups.NextPageToken = newGroups.NextPageToken
	}

	fGroups := make([]futureMembers, len(groups.Groups))

	for key, group := range groups.Groups {
		fGroups[key] = futureMembers{}
		fGroups[key].Start(&s, group)
	}

	uGroups := make([]goldapps.Group, len(groups.Groups))

	for key, group := range groups.Groups {

		members, err := fGroups[key].Members()
		if err != nil {
			return nil, err
		}

		new := goldapps.Group{
			Members: *members,
			Email:   group.Email,
			Aliases: group.Aliases,
		}

		uGroups[key] = new
	}

	return uGroups, nil

}

func (s *GoogleService) asyncMembers(g *admin.Group, ret chan memberResponse) {
	m, err := s.members(g)
	r := memberResponse{
		Members: m,
		Error:   err,
	}

	ret <- r
}

func (s *GoogleService) delayedMemberRequest(group *admin.Group, pageToken string, delay float64) (*admin.Members, error) {

	if delay > 1 {
		time.Sleep(time.Duration(delay))
	}else{
		delay = 2
	}

	var members *admin.Members = nil
	var err error = nil
	if pageToken == "" {
		members, err = s.Service.Members.List(group.Email).Do()
	} else {
		members, err = s.Service.Members.List(group.Email).PageToken(pageToken).Do()
	}
	if err != nil {
		if err.Error() == "googleapi: Error 403: Request rate higher than configured., quotaExceeded" {
			members, err = s.delayedMemberRequest(group, pageToken, math.Pow(delay, 2))
		}
	}
	return members, err
}

func (s *GoogleService) memberRequest(group *admin.Group, pageToken string) (*admin.Members, error) {

	return s.delayedMemberRequest(group, pageToken, 1)
}

func (s *GoogleService) members(group *admin.Group) (*[]string, error) {

	members, err := s.memberRequest(group, "")
	if err != nil {
		return nil, err
	}

	for members.NextPageToken != "" {
		newMembers, err := s.memberRequest(group, members.NextPageToken)
		if err != nil {
			return nil, err
		}

		members.Members = append(members.Members, newMembers.Members...)
		members.NextPageToken = newMembers.NextPageToken
	}

	uMembers := make([]string, len(members.Members))

	for key, member := range members.Members {
		new := member.Email

		uMembers[key] = new
	}

	return &uMembers, nil

}

*/