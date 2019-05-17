package ldap

import (
	"crypto/tls"

	"fmt"
	"github.com/cthit/goldapps/internal/pkg/model"
	"gopkg.in/ldap.v2"
	"strings"
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

func (s ServiceLDAP) GetUsers() ([]model.User, error) {
	return s.getUsers()
}

// Collect all users who are members of a committee
func (s ServiceLDAP) getUsers() ([]model.User, error) {
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

	// Create an empty model.Group slice
	privilegedUsers := make(model.Users, 0)

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
					if !privilegedUsers.Contains(user.GetAttributeValue("uid")) {
						if user.GetAttributeValue("gdprEducated") == "TRUE" { // only add user if he's gdpr educated
							privilegedUsers = append(privilegedUsers, model.User{
								Cid:        user.GetAttributeValue("uid"),
								Nick:       user.GetAttributeValue("nickname"),
								FirstName:  user.GetAttributeValue("givenName"),
								SecondName: user.GetAttributeValue("sn"),
								Mail:       user.GetAttributeValue("mail"),
							})
						}
					}
				}
			}
		}
	}

	return privilegedUsers, nil
}

// Recursively parse member tree and return users
func parsePrivilegedGroupMember(memberDN string, users []*ldap.Entry, groups []*ldap.Entry) []*ldap.Entry {
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
// model.Group slice.
func (s ServiceLDAP) GetGroups() ([]model.Group, error) {
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

	// Creates an empty model.Group slice
	groups := make([]model.Group, 0)

	// Creates a model.Group with appropriate mails and members
	for _, entry := range committees.Entries {

		// Creates a model.Group with it's mail
		committee := model.Group{
			Email:   entry.GetAttributeValue("mail"),
			Type:    entry.GetAttributeValue("type"),
			Members: nil,
		}

		if committee.Email == "" {
			continue
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

	positionGroups, err := s.getPositionGroups()
	if err != nil {
		return nil, err
	}
	groups = append(groups, positionGroups...)

	chairmenGroupMembers, err := s.getRoleInGroups("ordf", false)
	if err != nil {
		return nil, err
	}
	groups = append(groups, model.Group{
		Email:   "ordforanden@chalmers.it",
		Members: chairmenGroupMembers,
	})

	chairmenInCommitteesGroupMembers, err := s.getRoleInGroups("ordf", true)
	if err != nil {
		return nil, err
	}
	groups = append(groups, model.Group{
		Email:   "ordforanden.kommitteer@chalmers.it",
		Members: chairmenInCommitteesGroupMembers,
	})

	treasurersGroupMembers, err := s.getRoleInGroups("kassor", false)
	if err != nil {
		return nil, err
	}
	groups = append(groups, model.Group{
		Email:   "kassorer@chalmers.it",
		Members: treasurersGroupMembers,
	})

	treasurersInCommitteesGroupMembers, err := s.getRoleInGroups("kassor", true)
	if err != nil {
		return nil, err
	}
	groups = append(groups, model.Group{
		Email:   "kassorer.kommitteer@chalmers.it",
		Members: treasurersInCommitteesGroupMembers,
	})

	accounts, err := s.getUsers()
	if err != nil {
		return nil, err
	}

	for i, group := range groups {
		if group.Type == "Committee" {
			for _, comitteeGroupMemberEmail := range group.Members {

				for j := range groups {
					if groups[j].Email == comitteeGroupMemberEmail {
						groups[j] = replaceWithAccountEmail(groups[j], accounts)
					}
				}
			}
		} else if groups[i].Type == "CommitteeDirect" {
			groups[i] = replaceWithAccountEmail(groups[i], accounts)
		}
	}

	return groups, nil
}

func replaceWithAccountEmail(group model.Group, users model.Users) model.Group {
	for i := 0; i < len(group.Members); i++ {
		replacementFound := false
		for _, user := range users {
			if user.Mail == group.Members[i] {
				replacementFound = true
				group.Members[i] = user.Cid + "@chalmers.it"
			}
		}
		if !replacementFound {
			fmt.Printf("WARNING: no replacement could be found for %s in %s \n", group.Members[i], group.Email)

			//Remove member
			group.Members = append(group.Members[:i], group.Members[i+1:]...)
			i--
		}
	}
	return group
}

func (s ServiceLDAP) GetCustomGroups() ([]model.Group, error) {
	users, err := s.users()
	if err != nil {
		return nil, err
	}

	customGroups := make([]model.Group, 0)

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
					// FIXME: The %childRDN% is only necessary since year groups (e.g. snit14) are the same type as their Committee/Society.
					strings.Replace(entry.ParentFilter, "%childRDN%", getRDN(member.DN), -1), // The filter to apply
					entry.Attributes, // A list attributes to retrieve
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
				localMembers := member.GetAttributeValues("member")
				// Check if the found entry has a mail associated with it
				if mail == "" { // if not it should have members which do
					for _, localMember := range localMembers {
						mail = findEntry(users, localMember).GetAttributeValue("mail")
						//fmt.Println(mail)
						members = append(members, mail)
					}
				} else {
					mail := member.GetAttributeValue("mail")
					if mail != "" {
						members = append(members, mail)
					}
				}

			}
		}

		group := model.Group{
			Email:   entry.Mail,
			Members: members,
		}

		customGroups = append(customGroups, group)
	}

	return customGroups, nil
}

func (s ServiceLDAP) getPositionGroups() ([]model.Group, error) {
	users, err := s.users()
	if err != nil {
		return nil, err
	}

	var positionGroups []model.Group

	searchRequest := ldap.NewSearchRequest(
		"ou=fkit,ou=groups,dc=chalmers,dc=it", // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=itPosition))", // The filter to apply
		[]string{"cn", "member"},      // A list attributes to retrieve
		nil,
	)

	result, err := s.Connection.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	for _, entry := range result.Entries {
		groupType, err := dnPositionType(s, entry.DN)
		DnSplit := strings.SplitN(entry.DN, ",", 3)
		pos := DnSplit[0][3:]
		grp := DnSplit[1][3:]
		if err != nil {
			return nil, err
		}

		posGroup := model.Group{
			Email: fmt.Sprintf("%s.%s@chalmers.it", pos, grp),
			Type:  groupType + "Direct",
		}

		for _, member := range entry.GetAttributeValues("member") {
			ent := findEntry(users, member)
			// TODO detta krashar om användaren ej finns
			mail := ent.GetAttributeValue("mail")
			posGroup.Members = append(posGroup.Members, mail)
		}

		positionGroups = append(positionGroups, posGroup)

	}
	return positionGroups, nil

}

func (s ServiceLDAP) getRoleInGroups(role string, onlyCommittees bool) ([]string, error) {
	searchRequest := ldap.NewSearchRequest(
		"ou=fkit,ou=groups,dc=chalmers,dc=it", // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=itPosition)(cn=%s))", role), // The filter to apply
		[]string{"cn"}, // A list attributes to retrieve
		nil,
	)

	result, err := s.Connection.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var treasurersInCommitteeGroup []string

	for _, entry := range result.Entries {
		gtype, err := dnPositionType(s, entry.DN)
		if err != nil {
			return nil, err
		}
		if !onlyCommittees || gtype == "Committee" {
			dnSplit := strings.SplitN(entry.DN, ",", 3)
			treasurersInCommitteeGroup = append(treasurersInCommitteeGroup, fmt.Sprintf("%s.%s@chalmers.it", role, dnSplit[1][3:]))
		}

	}
	return treasurersInCommitteeGroup, nil
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

func dnPositionType(s ServiceLDAP, DN string) (string, error) {
	newDN := strings.SplitN(DN, ",", 2)[1] // Creates the dn for the group
	sr := ldap.NewSearchRequest(
		newDN, // The base dn to search
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=*))", // The filter to apply
		[]string{"type"},     // A list attributes to retrieve
		nil,
	)

	result, err := s.Connection.Search(sr)
	if err != nil {
		return "", err
	}

	return result.Entries[0].GetAttributeValue("type"), nil
}
