package admin

import (
	"google.golang.org/api/admin/directory/v1"

	"github.com/cthit/goldapps"

	"math"
	"net/http"
	"time"
)

type GoogleService struct {
	Service *admin.Service
}

type ServiceConfig interface {
	GoogleServiceJSONKeyPath() string
	GoogleServiceAdmin() string
}

func NewGoogleService(client *http.Client) (*GoogleService, error) {

	service, err := admin.New(client)
	if err != nil {
		return nil, err
	}

	gs := &GoogleService{
		Service: service,
	}

	return gs, nil
}

func (s GoogleService) UpdateGroups(groups []goldapps.Group) error {

	// Grab remote groups
	remoteGroups, err := s.Groups()
	if err != nil {
		return err
	}

	// Remove obselate remote groups
	for _, remoteGroup := range remoteGroups {
		if !groupContains(groups, remoteGroup) {
			err = s.deleteGroup(remoteGroup.Email)
			if err != nil {
				return err
			}
		}
	}

	// Update local groups
	for _, group := range groups {

	}

	return nil

}

func (s GoogleService) getGroup(email string) (admin.Group, error)  {
	group, err := s.Service.Groups.Get(email).Do()

	return *group, err
}

func (s GoogleService) insertGroup(group admin.Group) (error) {
	_, err := s.Service.Groups.Insert(&group).Do()
	return err
}

func (s GoogleService) updateSingleGroup(group admin.Group) (error) {
	_, err := s.Service.Groups.Update(group.Email, &group).Do()
	return err
}

func (s GoogleService) updateGroup(g goldapps.Group) error {
	// Retrieves the group with the specified email
	group, err := s.getGroup(g.Email)
	if err != nil { // TODO: Should probably check for more specific error
		// Assumes that the retrieval failed because the group doesn't exist.
		// Creates a new group.
		new := admin.Group{
			Name:    g.Name,
			Email:   g.Email,
			Aliases: *g.Alias,
		}

		// Inserts the new group
		err := s.insertGroup(new)
		if err != nil {
			return err
		}

	} else { // The group already exist so update it
		group.Email = g.Email
		group.Name = g.Name
		group.Aliases = *g.Alias

		err = s.updateSingleGroup(group)
		if err != nil {
			return err
		}
	}

	// The members all need to be updated
	err = s.updateMembers(&group, g.Members)
	if err != nil {
		return err
	}

	return nil

}

func (s *GoogleService) updateMembers(g *admin.Group, members *[]goldapps.Member) error {
	err := s.cleanupMembers(g, members)
	if err != nil {
		return err
	}

	return s.pushMembers(g, members)
}

func (s *GoogleService) cleanupMembers(g *admin.Group, members *[]goldapps.Member) error {
	current, err := s.members(g)
	if err != nil {
		return err
	}

	for _, cMember := range *current {
		if !memberContains(*members, cMember) {
			err = s.deleteMember(cMember, g)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *GoogleService) deleteMember(member goldapps.Member, g *admin.Group) error {
	return s.Service.Members.Delete(g.Email, member.Email).Do()
}

func (s *GoogleService) pushMembers(g *admin.Group, members *[]goldapps.Member) error {
	for _, member := range *members {
		mem, err := s.Service.Members.Get(g.Email, member.Email).Do()
		if err != nil {
			new := &admin.Member{
				Email: member.Email,
			}
			mem, err = s.Service.Members.Insert(g.Email, new).Do()
			if err != nil {
				return err
			}
		} else {
			mem.Email = member.Email
			mem, err = s.Service.Members.Update(g.Email, mem.Email, mem).Do()
		}
	}
	return nil
}

func (s GoogleService) deleteGroup(Email string) error {
	err := s.Service.Groups.Delete(Email).Do()

	return err
}

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
			Name:    group.Name,
			Members: members,
			Email:   group.Email,
			Alias:   &group.Aliases,
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

	time.Sleep(time.Duration(delay))

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

	var members *admin.Members = nil
	var err error = nil
	if pageToken == "" {
		members, err = s.Service.Members.List(group.Email).Do()
	} else {
		members, err = s.Service.Members.List(group.Email).PageToken(pageToken).Do()
	}
	if err != nil {
		if err.Error() == "googleapi: Error 403: Request rate higher than configured., quotaExceeded" {
			members, err = s.delayedMemberRequest(group, pageToken, 2)
		}
	}
	return members, err
}

func (s *GoogleService) members(group *admin.Group) (*[]goldapps.Member, error) {

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

	uMembers := make([]goldapps.Member, len(members.Members))

	for key, member := range members.Members {
		new := goldapps.Member{
			Email: member.Email,
		}

		uMembers[key] = new
	}

	return &uMembers, nil

}
