package admin

import (
	"google.golang.org/api/admin/directory/v1" // Imports as admin

	"io/ioutil"

	"bytes"
	"fmt"
	"github.com/cthit/goldapps"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"time"
	"strings"
	"google.golang.org/api/googleapi"
)

type googleService struct {
	service *admin.Service
}

func NewGoogleService(keyPath string, adminMail string) (goldapps.UpdateService, error) {

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

	gs := googleService{
		service: service,
	}

	return gs, nil
}

func (s googleService) DeleteGroup(group goldapps.Group) error {
	return s.deleteGroup(group.Email)
}

func (s googleService) UpdateGroup(groupUpdate goldapps.GroupUpdate) error {
	newGroup := admin.Group{
		Email: groupUpdate.Before.Email,
	}

	// Add all new members
	for _, member := range groupUpdate.After.Members {
		exists := false
		for _, existingMember := range groupUpdate.Before.Members {
			if strings.ToLower(member) == strings.ToLower(existingMember) {
				exists = true
				break
			}
		}
		if !exists {
			err := s.addMember(groupUpdate.Before.Email, member)
			if err != nil {
				fmt.Printf("Failed to add menber %s\n",member)
				return err
			}
		}
	}

	// Remove all old members
	for _, existingMember := range groupUpdate.Before.Members {
		keep := false
		for _, member := range groupUpdate.After.Members {
			if strings.ToLower(existingMember) == strings.ToLower(member) {
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
			if strings.ToLower(alias) == strings.ToLower(existingAlias) {
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
			if strings.ToLower(existingAlias) == strings.ToLower(alias) {
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

	return s.updateGroup(newGroup)
}

func (s googleService) AddGroup(group goldapps.Group) error {
	newGroup := admin.Group{
		Email: group.Email,
	}

	err := s.addGroup(newGroup)
	if err != nil {
		return err
	}

	time.Sleep(time.Second*10)

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

func (s googleService) GetGroups() ([]goldapps.Group, error) {

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
	fmt.Printf("\rDone\n")

	return groups, nil

}

func (s googleService) getGroups(customer string) ([]admin.Group, error) {
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

func (s googleService) getMembers(email string) ([]string, error) {
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

func (s googleService) getGroup(email string) (admin.Group, error) {
	group, err := s.service.Groups.Get(email).Do()

	return *group, err
}

func (s googleService) addGroup(group admin.Group) error {
	_, err := s.service.Groups.Insert(&group).Do()
	return err
}

func (s googleService) updateGroup(group admin.Group) error {
	_, err := s.service.Groups.Update(group.Email, &group).Do()
	return err
}

func (s googleService) deleteGroup(email string) error {
	err := s.service.Groups.Delete(email).Do()
	return err
}

func (s googleService) deleteMember(groupEmail string, member string) error {
	return s.service.Members.Delete(groupEmail, member).Do()
}

func (s googleService) addMember(groupEmail string, memberEmail string) error {
	_, err := s.service.Members.Insert(groupEmail, &admin.Member{Email: memberEmail}).Do()
	return err
}

func (s googleService) deleteAlias(groupEmail string, alias string) error {
	return s.service.Groups.Aliases.Delete(groupEmail, alias).Do()
}

func (s googleService) addAlias(groupEmail string, alias string) error {
	_, err := s.service.Groups.Aliases.Insert(groupEmail, &admin.Alias{Alias: alias}).Do()
	return err
}

func (s googleService) AddUser(user goldapps.User) error {
	_,err := s.service.Users.Insert(&admin.User{
		Addresses:                  nil,
		AgreedToTerms:              false,
		Aliases:                    nil,
		ChangePasswordAtNextLogin:  false,
		CreationTime:               "",
		CustomSchemas:              nil,
		CustomerId:                 "",
		DeletionTime:               "",
		Emails:                     nil,
		Etag:                       "",
		ExternalIds:                nil,
		Gender:                     nil,
		HashFunction:               "",
		Id:                         "",
		Ims:                        nil,
		IncludeInGlobalAddressList: false,
		IpWhitelisted:              false,
		IsAdmin:                    false,
		IsDelegatedAdmin:           false,
		IsEnforcedIn2Sv:            false,
		IsEnrolledIn2Sv:            false,
		IsMailboxSetup:             false,
		Keywords:                   nil,
		Kind:                       "",
		Languages:                  nil,
		LastLoginTime:              "",
		Locations:                  nil,
		Name: &admin.UserName{
			FamilyName:      "",
			FullName:        "",
			GivenName:       "",
			ForceSendFields: nil,
			NullFields:      nil,
		},
		NonEditableAliases: nil,
		Notes:              nil,
		OrgUnitPath:        "",
		Organizations:      nil,
		Password:           "",
		Phones:             nil,
		PosixAccounts:      nil,
		PrimaryEmail:       "",
		Relations:          nil,
		SshPublicKeys:      nil,
		Suspended:          false,
		SuspensionReason:   "",
		ThumbnailPhotoEtag: "",
		ThumbnailPhotoUrl:  "",
		Websites:           nil,
		ServerResponse: googleapi.ServerResponse{
			HTTPStatusCode: 0,
			Header:         nil,
		},
		ForceSendFields: nil,
		NullFields:      nil,
	}).Do()
	return err
}



func (s googleService) DeleteUser(goldapps.User) error {
	panic("implement me")
}

func (s googleService) UpdateUser(goldapps.UserUpdate) error {
	panic("implement me")
}

func (s googleService) GetUsers() ([]goldapps.User, error) {
	users, err := s.getUsers("my_customer")
	if err != nil {
		return nil, err
	}

	fmt.Println(users[1].Name.GivenName)
	fmt.Println(users[1].Name.FamilyName)
	fmt.Println(users[1].PrimaryEmail)
	fmt.Println(users[1].Password) // does not work
	fmt.Println(users[1].HashFunction) // does not work
	externalId := users[1].ExternalIds.([]interface{})

	for index ,id := range externalId  {
		for key, value := range id.(map[string]interface{}) {
			fmt.Printf("(%d) %s: %s\n", index, key, value.(string))
		}
	}

	return nil, err
}

func (s googleService) getUsers(customer string) ([]admin.User, error) {
	users, err := s.service.Users.List().Customer(customer).Do()
	if err != nil {
		return nil, err
	}

	for users.NextPageToken != "" {
		newUsers, err := s.service.Users.List().Customer(customer).PageToken(users.NextPageToken).Do()
		if err != nil {
			return nil, err
		}

		users.Users = append(users.Users, newUsers.Users...)
		users.NextPageToken = newUsers.NextPageToken
	}

	result := make([]admin.User, len(users.Users))
	for i, user := range users.Users {
		result[i] = *user
	}

	return result, nil
}