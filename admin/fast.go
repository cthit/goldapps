package admin

import (
	glsync "github.com/hulthe/google-ldap-sync"
	"google.golang.org/api/admin/directory/v1"
)

type memberResponse struct {
	Members *[]glsync.Member
	Error   error
}

type futureMembers struct {
	done     bool
	incoming chan memberResponse
	saved    memberResponse
}

func (f *futureMembers) Done() bool {
	return f.done
}

func (f *futureMembers) Members() (*[]glsync.Member, error) {
	if f.done == false {
		f.saved = <-f.incoming
		f.done = true
		f.incoming = nil
	}
	return f.saved.Members, f.saved.Error
}

func (f *futureMembers) Start(s *GoogleService, g *admin.Group) {

	f.done = false

	f.incoming = make(chan memberResponse, 1)

	go s.asyncMembers(g, f.incoming)

}
