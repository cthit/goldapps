package admin

import "google.golang.org/api/admin/directory/v1"

func Scopes() []string {
	return []string{
		admin.AdminDirectoryGroupScope,
	}
}
