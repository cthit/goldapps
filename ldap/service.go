package ldap

import (
	"crypto/tls"

	"strings"

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

	result, err := s.Connection.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	groups := make([]glsync.Group, len(result.Entries))
	groupIndex := 0

	for _, entry := range result.Entries {

		cn := entry.GetAttributeValue("cn")
		dn := entry.DN

		if strings.TrimPrefix(dn, "cn="+cn+",") == "ou="+cn+","+baseDN {

			committé := glsync.Group{
				Name:    entry.GetAttributeValue("displayName"),
				Email:   entry.GetAttributeValue("mail"),
				Members: nil,
				Alias:   nil,
			}

			members := make([]glsync.Member, len(users))
			memberIndex := 0

			for _, subCommitté := range entry.GetAttributeValues("member") {

				for _, entry2 := range result.Entries {
					if subCommitté == entry2.DN {
						committéMembers := entry2.GetAttributeValues("member")

						for _, member := range committéMembers {
							for _, user := range users {
								if user.DN == member {
									members[memberIndex] = glsync.Member{
										Email: user.GetAttributeValue("mail"),
									}
									memberIndex++
								}
							}
						}
					}
				}
			}

			membersSlice := members[0:memberIndex]
			committé.Members = &(membersSlice)

			groups[groupIndex] = committé
			groupIndex++
		}
	}

	return groups[0:groupIndex], nil
}
