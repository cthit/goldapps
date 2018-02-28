package ldap

import (
	"crypto/tls"

	"github.com/cthit/goldapps"
	"gopkg.in/ldap.v2"
)

type ServiceLDAP struct {
	Connection         *ldap.Conn
	DBConfig           ServerConfig
	GroupsConfig       EntryConfig
	UsersConfig        EntryConfig
	CustomEntryConfigs []CustomEntryConfig
}

type ServerConfig struct {
	Url        string
	ServerName string
}

type EntryConfig struct {
	BaseDN     string
	Filter     string
	Attributes []string
}

type CustomEntryConfig struct {
	BaseDN     string
	Filter     string
	Attributes []string
	Mail       string
}

type LoginConfig struct {
	UserName string
	Password string
}

func NewLDAPService(dbConfig ServerConfig, login LoginConfig, usersConfig EntryConfig, groupsConfig EntryConfig, customEntryConfigs []CustomEntryConfig) (*ServiceLDAP, error) {

	l, err := ldap.DialTLS("tcp", dbConfig.Url, &tls.Config{ServerName: dbConfig.ServerName})
	if err != nil {
		return nil, err
	}
	// FIXME: Close connection on garbage collection
	//defer l.Close()

	err = l.Bind(login.UserName, login.Password)
	if err != nil {
		return nil, err
	}

	ld := &ServiceLDAP{
		Connection:   l,
		DBConfig:     dbConfig,
		UsersConfig:  usersConfig,
		GroupsConfig: groupsConfig,
		CustomEntryConfigs: customEntryConfigs,
	}

	return ld, nil

}

// Collects all users from LDAP as a slice of *ldap.Entry's
func (s ServiceLDAP) users() ([]*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(
		s.UsersConfig.BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		s.UsersConfig.Filter,     // The filter to apply
		s.UsersConfig.Attributes, // A list attributes to retrieve
		nil,
	)

	result, err := s.Connection.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	return result.Entries, nil
}

// Collects all committees from LDAP and then creates a
// goldapps.Group slice.
func (s ServiceLDAP) GetGroups() ([]goldapps.Group, error) {
	users, err := s.users()
	if err != nil {
		return nil, err
	}

	// Creates a search request to collect all committees from LDAP
	searchRequest := ldap.NewSearchRequest(
		s.GroupsConfig.BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		s.GroupsConfig.Filter,     // The filter to apply
		s.GroupsConfig.Attributes, // A list attributes to retrieve
		nil,
	)

	// Collects the committee entries
	committees, err := s.Connection.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	// Creates an empty goldapps.Group slice
	groups := make([]goldapps.Group, 0)

	// Creates a goldapps.Group with appropriate mails and members
	for _, entry := range committees.Entries {

		// Creates a goldapps.Group with it's mail
		committee := goldapps.Group{
			Email:   entry.GetAttributeValue("mail"),
			Members: nil,
		}

		// Creates an empty members slice
		members := make([]string, 0) // len(users) might break if we have all users and some groups in the members field

		// Fills the members slice with data
		for _, member := range entry.GetAttributeValues("member") {
			var m *ldap.Entry

			if dnIsUser(member) {
				m = findEntry(users, member)
			} else {
				m = findEntry(committees.Entries, member)
			}

			if m != nil {
				mail := m.GetAttributeValue("mail")
				if mail != "" {
					members = append(members, mail)
				}
			}
		}

		committee.Members = members

		groups = append(groups, committee)
	}

	customGroups, err := s.GetCustomGroups()
	if err != nil {
		return nil, err
	}

	groups = append(groups, customGroups...)

	return groups, nil
}

func (s ServiceLDAP) GetCustomGroups() ([]goldapps.Group, error) {
	customGroups := make([]goldapps.Group, 0)

	for _, entry := range s.CustomEntryConfigs {
		// Creates a search request to collect all committees from LDAP
		searchRequest := ldap.NewSearchRequest(
			entry.BaseDN, // The base dn to search
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			entry.Filter,     // The filter to apply
			entry.Attributes, // A list attributes to retrieve
			nil,
		)

		result, err := s.Connection.Search(searchRequest)
		if err != nil {
			return nil, err
		}

		members := make([]string, 0) // len(users) might break if we have all users and some groups in the members field

		for _, member := range result.Entries {
			mail := member.GetAttributeValue("mail")
			if mail != "" {
				members = append(members, mail)
			}
		}

		group := goldapps.Group{
			Email:   entry.Mail,
			Members: members,
		}

		customGroups = append(customGroups, group)
	}

	return customGroups, nil
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
