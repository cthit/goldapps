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

const gdprSuspensionText = "You have not attended the GDPR education!"

const googleDuplicateEntryError = "googleapi: Error 409: Entity already exists., duplicate"

// my_customer seems to work...
const googleCustomer = "my_customer"

type googleService struct {
	google *admin.Service
	admin  string
	domain string
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

	// Extract account and mail
	s := strings.Split(adminMail, "@")
	admin := s[0]
	domain := s[1]

	gs := googleService{
		google: service,
		admin:  admin,
		domain: domain,
	}

	return gs, nil
}

func (s googleService) DeleteGroup(group goldapps.Group) error {
	err := s.google.Groups.Delete(group.Email).Do()
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
			if strings.ToLower(member) == strings.ToLower(existingMember) {
				exists = true
				break
			}
		}
		if !exists {
			_, err := s.google.Members.Insert(groupUpdate.Before.Email, &admin.Member{Email: member}).Do()
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
			err := s.google.Members.Delete(groupUpdate.Before.Email, existingMember).Do()
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
			_, err := s.google.Groups.Aliases.Insert(groupUpdate.Before.Email, &admin.Alias{Alias: alias}).Do()
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
			err := s.google.Groups.Aliases.Delete(groupUpdate.Before.Email, existingAlias).Do()
			if err != nil {
				return err
			}
		}
	}

	_, err := s.google.Groups.Update(groupUpdate.Before.Email, &newGroup).Do()
	return err
}

func (s googleService) AddGroup(group goldapps.Group) error {
	newGroup := admin.Group{
		Email: group.Email,
	}

	_, err := s.google.Groups.Insert(&newGroup).Do()
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 10)

	// Add members
	for _, member := range group.Members {
		_, err := s.google.Members.Insert(group.Email, &admin.Member{Email: member}).Do()
		if err != nil {
			return err
		}
	}

	// Add Aliases
	for _, alias := range group.Aliases {
		_, err := s.google.Groups.Aliases.Insert(group.Email, &admin.Alias{Alias: alias}).Do()
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

func (s googleService) AddUser(user goldapps.User) error {

	usr := buildGoldappsUser(user, s.domain)
	usr.Password = user.PasswordHash
	usr.HashFunction = user.HashFunction

	_, err := s.google.Users.Insert(usr).Do()
	if err != nil {
		return err
	}

	// Google needs time for the addition to propagate
	time.Sleep(time.Second)

	// Add alias for nick@example.ex
	return s.addUserAlias(fmt.Sprintf("%s@%s", user.Cid, s.domain), fmt.Sprintf("%s@%s", user.Nick, s.domain))
}

func (s googleService) UpdateUser(update goldapps.UserUpdate) error {
	_, err := s.google.Users.Update(
		fmt.Sprintf("%s@%s", update.Before.Cid, s.domain),
		buildGoldappsUser(update.After, s.domain),
	).Do()
	if err != nil {
		return err
	}

	// Add alias for nick@example.ex
	return s.addUserAlias(fmt.Sprintf("%s@%s", update.After.Cid, s.domain), fmt.Sprintf("%s@%s", update.After.Nick, s.domain))
}

func (s googleService) DeleteUser(user goldapps.User) error {
	admin := fmt.Sprintf("%s@%s", s.admin, s.domain)
	userId := fmt.Sprintf("%s@%s", user.Cid, s.domain)
	if admin == userId {
		fmt.Printf("Skipping andmin user: %s\n", admin)
	}

	err := s.google.Users.Delete(userId).Do()
	return err
}

func (s googleService) GetUsers() ([]goldapps.User, error) {
	adminUsers, err := s.getGoogleUsers(googleCustomer)
	if err != nil {
		return nil, err
	}
	users := make([]goldapps.User, len(adminUsers)-1)

	admin := fmt.Sprintf("%s@%s", s.admin, s.domain)

	i := 0
	for _, adminUser := range adminUsers {
		if admin != adminUser.PrimaryEmail { // Don't list admin account
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
			gdpr := !(adminUser.Suspended && adminUser.SuspensionReason == gdprSuspensionText)

			users[i] = goldapps.User{
				Cid:           cid,
				FirstName:     firstName,
				SecondName:    adminUser.Name.FamilyName,
				Nick:          nick,
				GdprEducation: gdpr,
			}
			i++
		}
	}

	return users, err
}

func (s googleService) getGoogleGroups(customer string) ([]admin.Group, error) {
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

func (s googleService) getGoogleGroupMembers(email string) ([]string, error) {
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

func (s googleService) getGoogleUsers(customer string) ([]admin.User, error) {
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

func (s googleService) addUserAlias(userKey string, alias string) error {
	_, err := s.google.Users.Aliases.Insert(userKey, &admin.Alias{
		Alias: alias,
	}).Do()
	if err != nil {
		if err.Error() == googleDuplicateEntryError {
			fmt.Printf("Warning: Could not add alias for %s. It already exists. \n", alias)
		} else {
			return err
		}
	}
	return nil
}

func buildGoldappsUser(user goldapps.User, domain string) *admin.User {
	return &admin.User{
		Name: &admin.UserName{
			FamilyName: user.SecondName,
			GivenName:  fmt.Sprintf("%s / %s", user.Nick, user.FirstName),
		},
		IncludeInGlobalAddressList: true,
		PrimaryEmail:               fmt.Sprintf("%s@%s", user.Cid, domain),
		Suspended:                  !user.GdprEducation,
		SuspensionReason:           gdprSuspensionText,
	}
}
