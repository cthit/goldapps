package admin

import (
	"fmt"
	"github.com/cthit/goldapps"
	"google.golang.org/api/admin/directory/v1" // Imports as admin
	"strings"
	"time"
)

func (s googleService) AddUser(user goldapps.User) error {

	usr := buildGoldappsUser(user, s.domain)

	password := newPassword()

	usr.Password = password
	usr.ChangePasswordAtNextLogin = true

	_, err := s.adminService.Users.Insert(usr).Do()
	if err != nil {
		return err
	}

	err = s.sendPassword(user.Mail, password)
	if err != nil {
		return err
	}

	// Google needs time for the addition to propagate
	time.Sleep(time.Second)

	// Add alias for nick@example.ex
	return s.addUserAlias(fmt.Sprintf("%s@%s", goldapps.SanitizeEmail(user.Cid), s.domain), fmt.Sprintf("%s@%s", goldapps.SanitizeEmail(user.Nick), s.domain))
}

func (s googleService) UpdateUser(update goldapps.UserUpdate) error {
	_, err := s.adminService.Users.Update(
		fmt.Sprintf("%s@%s", goldapps.SanitizeEmail(update.Before.Cid), s.domain),
		buildGoldappsUser(update.After, s.domain),
	).Do()
	if err != nil {
		return err
	}

	// Add alias for nick@example.ex
	return s.addUserAlias(fmt.Sprintf("%s@%s", goldapps.SanitizeEmail(update.After.Cid), s.domain), fmt.Sprintf("%s@%s", goldapps.SanitizeEmail(update.After.Nick), s.domain))
}

func (s googleService) DeleteUser(user goldapps.User) error {
	admin := fmt.Sprintf("%s@%s", s.admin, s.domain)
	userId := fmt.Sprintf("%s@%s", goldapps.SanitizeEmail(user.Cid), s.domain)
	if admin == userId {
		fmt.Printf("Skipping andmin user: %s\n", admin)
	}

	err := s.adminService.Users.Delete(userId).Do()
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

			users[i] = goldapps.User{
				Cid:        cid,
				FirstName:  firstName,
				SecondName: adminUser.Name.FamilyName,
				Nick:       nick,
			}
			i++
		}
	}

	return users, err
}

func (s googleService) getGoogleUsers(customer string) ([]admin.User, error) {
	users, err := s.adminService.Users.List().Customer(customer).Do()
	if err != nil {
		return nil, err
	}

	for users.NextPageToken != "" {
		newUsers, err := s.adminService.Users.List().Customer(customer).PageToken(users.NextPageToken).Do()
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
	_, err := s.adminService.Users.Aliases.Insert(userKey, &admin.Alias{
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
