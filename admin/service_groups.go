package admin

import (
	"google.golang.org/api/admin/directory/v1" // Imports as admin
	"../../goldapps"
	"fmt"
	"time"
	"bytes"
)

func (s googleService) DeleteGroup(group goldapps.Group) error {
	err := s.adminService.Groups.Delete(group.Email).Do()
	return err
}

func (s googleService) UpdateGroup(groupUpdate goldapps.GroupUpdate) error {
	newGroup := admin.Group{
		Email: groupUpdate.Before.Email,
	}

	// Add all new members
	for _, member := range groupUpdate.After.Members {
		exists := false
		for _, existingMember := range groupUpdate.Before.Members {
			if goldapps.CompareEmails(member, existingMember) {
				exists = true
				break
			}
		}
		if !exists {
			_, err := s.adminService.Members.Insert(groupUpdate.Before.Email, &admin.Member{Email: member}).Do()
			if err != nil {
				fmt.Printf("Failed to add menber %s\n", member)
				return err
			}
		}
	}

	// Remove all old members
	for _, existingMember := range groupUpdate.Before.Members {
		keep := false
		for _, member := range groupUpdate.After.Members {
			if goldapps.CompareEmails(existingMember, member) {
				keep = true
				break
			}
		}
		if !keep {
			err := s.adminService.Members.Delete(groupUpdate.Before.Email, existingMember).Do()
			if err != nil {
				return err
			}
		}
	}

	// Add all new aliases
	for _, alias := range groupUpdate.After.Aliases {
		exists := false
		for _, existingAlias := range groupUpdate.Before.Aliases {
			if goldapps.CompareEmails(alias, existingAlias) {
				exists = true
				break
			}
		}
		if !exists {
			_, err := s.adminService.Groups.Aliases.Insert(groupUpdate.Before.Email, &admin.Alias{Alias: alias}).Do()
			if err != nil {
				return err
			}
		}
	}

	// Remove all old aliases
	for _, existingAlias := range groupUpdate.Before.Aliases {
		keep := false
		for _, alias := range groupUpdate.After.Aliases {
			if goldapps.CompareEmails(existingAlias, alias) {
				keep = true
				break
			}
		}
		if !keep {
			err := s.adminService.Groups.Aliases.Delete(groupUpdate.Before.Email, existingAlias).Do()
			if err != nil {
				return err
			}
		}
	}

	_, err := s.adminService.Groups.Update(groupUpdate.Before.Email, &newGroup).Do()
	return err
}

func (s googleService) AddGroup(group goldapps.Group) error {
	newGroup := admin.Group{
		Email: group.Email,
	}

	_, err := s.adminService.Groups.Insert(&newGroup).Do()
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 10)

	// Add members
	for _, member := range group.Members {
		_, err := s.adminService.Members.Insert(group.Email, &admin.Member{Email: member}).Do()
		if err != nil {
			return err
		}
	}

	// Add Aliases
	for _, alias := range group.Aliases {
		_, err := s.adminService.Groups.Aliases.Insert(group.Email, &admin.Alias{Alias: alias}).Do()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s googleService) GetGroups() ([]goldapps.Group, error) {

	adminGroups, err := s.getGoogleGroups(googleCustomer)
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

		members, err := s.getGoogleGroupMembers(group.Email)
		if err != nil {
			return nil, err
		}

		groups[i] = goldapps.Group{
			Email:   group.Email,
			Members: members,
			Aliases: group.Aliases,
		}
	}
	fmt.Printf("\rDone\n")

	return groups, nil

}

func (s googleService) getGoogleGroups(customer string) ([]admin.Group, error) {
	groups, err := s.adminService.Groups.List().Customer(customer).Do()
	if err != nil {
		return nil, err
	}

	for groups.NextPageToken != "" {
		newGroups, err := s.adminService.Groups.List().Customer(customer).PageToken(groups.NextPageToken).Do()
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

func (s googleService) getGoogleGroupMembers(email string) ([]string, error) {
	members, err := s.adminService.Members.List(email).Do()
	if err != nil {
		return nil, err
	}

	result := make([]string, len(members.Members))
	for i, member := range members.Members {
		result[i] = member.Email
	}

	return result, nil
}
