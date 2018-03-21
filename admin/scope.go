package admin

import "google.golang.org/api/admin/directory/v1"

func Scopes() []string {
	return []string{
		admin.AdminDirectoryGroupScope,
		admin.AdminDirectoryUserScope,
	}
}
// https://www.googleapis.com/auth/admin.directory.group, https://www.googleapis.com/auth/admin.directory.user