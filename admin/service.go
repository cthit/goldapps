package admin

import (
	"github.com/cthit/goldapps"

	"google.golang.org/api/admin/directory/v1" // Imports as admin

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type googleService struct {
	google *admin.Service
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
		google: service,
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
				fmt.Printf("Failed to add menber %s\n", member)
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

	time.Sleep(time.Second * 10)

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
	groups, err := s.google.Groups.List().Customer(customer).Do()
	if err != nil {
		return nil, err
	}

	for groups.NextPageToken != "" {
		newGroups, err := s.google.Groups.List().Customer(customer).PageToken(groups.NextPageToken).Do()
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
	members, err := s.google.Members.List(email).Do()
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
	group, err := s.google.Groups.Get(email).Do()

	return *group, err
}

func (s googleService) addGroup(group admin.Group) error {
	_, err := s.google.Groups.Insert(&group).Do()
	return err
}

func (s googleService) updateGroup(group admin.Group) error {
	_, err := s.google.Groups.Update(group.Email, &group).Do()
	return err
}

func (s googleService) deleteGroup(email string) error {
	err := s.google.Groups.Delete(email).Do()
	return err
}

func (s googleService) deleteMember(groupEmail string, member string) error {
	return s.google.Members.Delete(groupEmail, member).Do()
}

func (s googleService) addMember(groupEmail string, memberEmail string) error {
	_, err := s.google.Members.Insert(groupEmail, &admin.Member{Email: memberEmail}).Do()
	return err
}

func (s googleService) deleteAlias(groupEmail string, alias string) error {
	return s.google.Groups.Aliases.Delete(groupEmail, alias).Do()
}

func (s googleService) addAlias(groupEmail string, alias string) error {
	_, err := s.google.Groups.Aliases.Insert(groupEmail, &admin.Alias{Alias: alias}).Do()
	return err
}

func (s googleService) AddUser(user goldapps.User) error {
	_, err := s.google.Users.Insert(buildAdminUser(user)).Do()
	return err
}

func buildAdminUser(user goldapps.User) *admin.User {
	return &admin.User{
		Name: &admin.UserName{
			FamilyName: user.SecondName,
			GivenName:  fmt.Sprintf("%s / %s", user.Nick, user.FirstName),
		},
		IncludeInGlobalAddressList: true,
		PrimaryEmail:               fmt.Sprintf("%s@chalmers.it", user.Cid),
		Emails: &[]admin.UserEmail{
			{
				Address: fmt.Sprintf("%s@chalmers.it", user.Nick),
				Primary: false,
				Type:    "other",
			},
		},
		Password:                  "RandomPassword", // Todo: how to do with passwords?
		ChangePasswordAtNextLogin: true,
		Suspended:                 !user.GdprEducation,
		SuspensionReason:          "You have not attended the GDPR education!",
	}
}

func (s googleService) DeleteUser(user goldapps.User) error {
	err := s.google.Users.Delete(fmt.Sprintf("%s@chalmers.it", user.Cid)).Do()
	return err
}

func (s googleService) UpdateUser(update goldapps.UserUpdate) error {
	_, err := s.google.Users.Update(
		fmt.Sprintf("%s@chalmers.it", update.Before.Cid),
		buildAdminUser(update.After),
	).Do()
	return err
}

func (s googleService) GetUsers() ([]goldapps.User, error) {
	adminUsers, err := s.getUsers("my_customer")
	if err != nil {
		return nil, err
	}
	users := make([]goldapps.User, len(adminUsers))

	for i, adminUser := range adminUsers {
		// Separating nick and firstName from (Nick / FirstName)
		givenName := strings.Split(adminUser.Name.GivenName, " / ")
		nick := givenName[0]
		firstName := ""
		if len(givenName) >= 2 {
			firstName = givenName[1]
		}

		// Extracting cid form (cid@example.ex)
		cid := strings.Split(adminUser.PrimaryEmail, "@")[0]

		// Check suspension and suspension reason to determine GDPR status
		gdpr := !(adminUser.Suspended && adminUser.SuspensionReason == "You have not attended the GDPR education!")

		users[i] = goldapps.User{
			Cid:           cid,
			FirstName:     firstName,
			SecondName:    adminUser.Name.FamilyName,
			Nick:          nick,
			GdprEducation: gdpr,
		}
	}

	return users, err
}

func (s googleService) getUsers(customer string) ([]admin.User, error) {
	users, err := s.google.Users.List().Customer(customer).Do()
	if err != nil {
		return nil, err
	}

	for users.NextPageToken != "" {
		newUsers, err := s.google.Users.List().Customer(customer).PageToken(users.NextPageToken).Do()
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
