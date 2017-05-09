package ldap

import (
	"crypto/tls"
	"fmt"

	glsync "github.com/hulthe/google-ldap-sync"
	"github.com/spf13/viper"
	"gopkg.in/ldap.v2"
)

func NewLDAPService(url string, serverName string, userName string, password string) (*ServiceLDAP, error) {

	l, err := ldap.DialTLS("tcp", url, &tls.Config{ServerName: serverName})
	if err != nil {
		return nil, err
	}
	defer l.Close()

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

func (s ServiceLDAP) Groups() ([]glsync.Group, error) {
	searchRequest := ldap.NewSearchRequest(
		viper.GetString("ldap.basedn"), // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		viper.GetString("ldap.filter"),          // The filter to apply
		viper.GetStringSlice("ldap.attributes"), // A list attributes to retrieve
		nil,
	)

	result, err := s.Connection.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	for _, entry := range result.Entries {
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn"))
	}

	return nil, nil
}
