package admin

import (
	"google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/gmail/v1"
)

func Scopes() []string {
	return []string{
		admin.AdminDirectoryGroupScope,
		admin.AdminDirectoryUserScope,
		gmail.GmailSendScope,
	}
}
// https://www.googleapis.com/auth/admin.directory.group, https://www.googleapis.com/auth/admin.directory.user