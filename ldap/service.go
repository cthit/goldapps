package ldap

import (
	"crypto/tls"

	"github.com/cthit/goldapps"
	"gopkg.in/ldap.v2"
	"strings"
	"fmt"
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
	BaseDN       string
	Filter       string
	ParentFilter string
	Attributes   []string
	Mail         string
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
		Connection:         l,
		DBConfig:           dbConfig,
		UsersConfig:        usersConfig,
		GroupsConfig:       groupsConfig,
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

// Collect all users who are members of a committee
func (s ServiceLDAP) GetUsers() ([]goldapps.User, error) {
	users, err := s.users()
	if err != nil {
		return nil, err
	}

	// Create a search request to collect all groups from LDAP
	searchRequest := ldap.NewSearchRequest(
		s.GroupsConfig.BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		s.GroupsConfig.Filter,     // The filter to apply
		s.GroupsConfig.Attributes, // A list attributes to retrieve
		nil,
	)

	// Collect the group entries
	groups, err := s.Connection.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	// Create an empty goldapps.Group slice
	privilegedUsers := make([]goldapps.User, 0)

	for _, group := range groups.Entries {
		// TODO: What qualified as a privileged group should be made configurable. See FIXME:s
		if group.GetAttributeValue("type") != "Committee" /* FIXME */ {
			continue // Only Committees are considered privileged groups
		}

		cn := group.GetAttributeValue("cn")
		// Check if RDN is the same as the groups parent. FIXME
		if strings.HasPrefix(group.DN, fmt.Sprintf("cn=%s,ou=%s", cn, cn)) {
			for _, member := range group.GetAttributeValues("member") {
				for _, user := range parsePrivilegedGroupMember(member, users, groups.Entries) {
					privilegedUsers = append(privilegedUsers, goldapps.User{
						// TODO: Make these attribute values configurable
						Cid:           user.GetAttributeValue("uid"),
						Nick:          user.GetAttributeValue("nickname"),
						FirstName:     user.GetAttributeValue("givenName"),
						SecondName:    user.GetAttributeValue("sn"),
						Mail:          user.GetAttributeValue("mail"),
						GdprEducation: user.GetAttributeValue("gdprEducated") == "TRUE",
					})
				}
			}
		}
	}

	return privilegedUsers, nil
}

// Recursively parse member tree and return users
func parsePrivilegedGroupMember(memberDN string, users []*ldap.Entry, groups []*ldap.Entry) ([]*ldap.Entry) {
	res := make([]*ldap.Entry, 0)
	if dnIsUser(memberDN) {
		for _, user := range users {
			if user.DN == memberDN {
				res = append(res, user)
				break
			}
		}
	} else {
		for _, group := range groups {
			if group.DN == memberDN {
				for _, subMember := range group.GetAttributeValues("member") {
					res = append(res, parsePrivilegedGroupMember(subMember, users, groups)...)
				}
				break
			}
		}
	}
	return res
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

			var parentResult *ldap.SearchResult = nil
			if entry.ParentFilter != "" {
				parentSearchRequest := ldap.NewSearchRequest(
					entry.BaseDN, // The base dn to search
					ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
					strings.Replace(entry.ParentFilter, "%childRDN%", getRDN(member.DN), -1), // The filter to apply
					entry.Attributes,                                                         // A list attributes to retrieve
					nil,
				)

				parentResult, err = s.Connection.Search(parentSearchRequest)
				if err != nil {
					return nil, err
				}
			}

			// If parent filter exists: check if member has a parent that matches
			var addMember = parentResult == nil
			if !addMember {
				for _, parent := range parentResult.Entries {
					if dnIsParentOf(parent.DN, member.DN) {
						addMember = true
						break
					}
				}
			}

			if addMember {
				mail := member.GetAttributeValue("mail")
				if mail != "" {
					members = append(members, mail)
				}
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

func getRDN(DN string) string {
	return strings.Split(strings.Split(DN, ",")[0], "=")[1]
}

func dnIsParentOf(parent string, node string) bool {
	return len(parent) != len(node) && strings.Contains(node, parent)
}

func dnIsUser(DN string) bool {
	return len(DN) >= 4 && DN[0:4] == "uid="
}
