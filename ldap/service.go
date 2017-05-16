package ldap

import (
	"crypto/tls"

	glsync "github.com/hulthe/google-ldap-sync"
	"github.com/spf13/viper"
	"gopkg.in/ldap.v2"
)

func NewLDAPService(url string, serverName string, userName string, password string) (*ServiceLDAP, error) {

	l, err := ldap.DialTLS("tcp", url, &tls.Config{ServerName: serverName})
	if err != nil {
		return nil, err
	}
	// FIXME: Close connection on garbage collection
	//defer l.Close()

	err = l.Bind(userName, password)
	if err != nil {
		return nil, err
	}

	ld := &ServiceLDAP{
		Connection: l,
	}

	return ld, nil

}

type ServiceLDAP struct {
	Connection *ldap.Conn
}

func (s ServiceLDAP) users() ([]*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(
		viper.GetString("ldap.users.basedn"), // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		viper.GetString("ldap.users.filter"),          // The filter to apply
		viper.GetStringSlice("ldap.users.attributes"), // A list attributes to retrieve
		nil,
	)

	result, err := s.Connection.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	return result.Entries, nil
}

func findEntry(ldapEntries []*ldap.Entry, DN string) *ldap.Entry {
	for _, entry := range ldapEntries {
		if entry.DN == DN {
			return entry
		}
	}
	return nil
}

func dnIsUser(DN string) bool {
	return len(DN) >= 4 && DN[0:4] == "uid="
}

func (s ServiceLDAP) Groups() ([]glsync.Group, error) {
	users, err := s.users()
	if err != nil {
		return nil, err
	}

	baseDN := viper.GetString("ldap.groups.basedn")
	searchRequest := ldap.NewSearchRequest(
		baseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		viper.GetString("ldap.groups.filter"),          // The filter to apply
		viper.GetStringSlice("ldap.groups.attributes"), // A list attributes to retrieve
		nil,
	)

	committés, err := s.Connection.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	groups := make([]glsync.Group, len(committés.Entries))
	groupIndex := 0

	for _, entry := range committés.Entries {

		committé := glsync.Group{
			Name:    entry.GetAttributeValue("displayName"),
			Email:   entry.GetAttributeValue("mail"),
			Members: nil,
			Alias:   nil,
		}

		members := make([]glsync.Member, len(users))
		memberIndex := 0

		for _, member := range entry.GetAttributeValues("member") {

			var m *ldap.Entry

			if dnIsUser(member) {
				m = findEntry(users, member)
			} else {
				m = findEntry(committés.Entries, member)
			}

			if m != nil {
				mail := m.GetAttributeValue("mail")
				if mail != "" {
					members[memberIndex] = glsync.Member{Email: mail}
					memberIndex++
				}
			}
		}

		membersSlice := members[0:memberIndex]
		committé.Members = &(membersSlice)

		groups[groupIndex] = committé
		groupIndex++
	}

	return groups[0:groupIndex], nil
}
