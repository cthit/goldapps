package admin

import (
	"fmt"
	"github.com/cthit/goldapps"
	"google.golang.org/api/admin/directory/v1" // Imports as admin
)

func buildGoldappsUser(user goldapps.User, domain string) *admin.User {
	return &admin.User{
		Name: &admin.UserName{
			FamilyName: user.SecondName,
			GivenName:  fmt.Sprintf("%s / %s", user.Nick, user.FirstName),
		},
		IncludeInGlobalAddressList: true,
		PrimaryEmail:               goldapps.SanitizeEmail(fmt.Sprintf("%s@%s", user.Cid, domain)),
	}
}
