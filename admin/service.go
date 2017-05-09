package admin

import (
	"google.golang.org/api/admin/directory/v1"

	glsync "github.com/hulthe/google-ldap-sync"

	"math"
	"net/http"
	"time"
)

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

type GoogleService struct {
	Service *admin.Service
}

func (s GoogleService) UpdateGroups(g []glsync.Group) error {

	remote, err := s.Groups()
	if err != nil {
		return err
	}

	for _, group := range remote {
		if !groupContains(g, group) {
			err = s.deleteGroup(group.Email)
		} else {
			err = s.updateGroup(&group)
		}

		if err != nil {
			return err
		}
	}

	return nil

}

func (s GoogleService) updateGroup(g *glsync.Group) error {
	group, err := s.Service.Groups.Get(g.Email).Do()
	if err != nil {
		new := &admin.Group{
			Name:    g.Name,
			Email:   g.Email,
			Aliases: *g.Alias,
		}

		group, err = s.Service.Groups.Insert(new).Do()
		if err != nil {
			return err
		}

	} else {
		group.Email = g.Email
		group.Name = g.Name
		group.Aliases = *g.Alias

		group, err = s.Service.Groups.Update(group.Email, group).Do()
		if err != nil {
			return err
		}
	}

	err = s.updateMembers(group, g.Members)
	if err != nil {
		return err
	}

	return nil

}

func (s *GoogleService) updateMembers(g *admin.Group, members *[]glsync.Member) error {
	err := s.cleanupMembers(g, members)
	if err != nil {
		return err
	}

	return s.pushMembers(g, members)
}

func (s *GoogleService) cleanupMembers(g *admin.Group, members *[]glsync.Member) error {
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

func (s *GoogleService) deleteMember(member glsync.Member, g *admin.Group) error {
	return s.Service.Members.Delete(g.Email, member.Email).Do()
}

func (s *GoogleService) pushMembers(g *admin.Group, members *[]glsync.Member) error {
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

func (s GoogleService) Groups() ([]glsync.Group, error) {
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

	uGroups := make([]glsync.Group, len(groups.Groups))

	for key, group := range groups.Groups {

		members, err := fGroups[key].Members()
		if err != nil {
			return nil, err
		}

		new := glsync.Group{
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

func (s *GoogleService) members(group *admin.Group) (*[]glsync.Member, error) {

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

	uMembers := make([]glsync.Member, len(members.Members))

	for key, member := range members.Members {
		new := glsync.Member{
			Email: member.Email,
		}

		uMembers[key] = new
	}

	return &uMembers, nil

}
