package admin

import (
	"fmt"
	"github.com/cthit/goldapps/internal/pkg/model"
	"google.golang.org/api/admin/directory/v1" // Imports as admin
)

func buildGoldappsUser(user model.User, domain string) *admin.User {
	return &admin.User{
		Name: &admin.UserName{
			FamilyName: user.SecondName,
			GivenName:  fmt.Sprintf("%s / %s", user.Nick, user.FirstName),
		},
		IncludeInGlobalAddressList: true,
		PrimaryEmail:               fmt.Sprintf("%s@%s", model.SanitizeEmail(user.Cid), domain),
	}
}
