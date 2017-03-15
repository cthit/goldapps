package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gopkg.in/ldap.v2"
)

func getLDAPGroups() {

	l, err := ldap.DialTLS("tcp", viper.GetString("ldap.url"), &tls.Config{ServerName: viper.GetString("ldap.servername")})
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	err = l.Bind(viper.GetString("ldap.user"), viper.GetString("ldap.assword"))
	if err != nil {
		log.Fatal(err)
	}

	searchRequest := ldap.NewSearchRequest(
		viper.GetString("ldap.basedn"), // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		viper.GetString("ldap.filter"),          // The filter to apply
		viper.GetStringSlice("ldap.Attributes"), // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range sr.Entries {
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn"))
	}
}
