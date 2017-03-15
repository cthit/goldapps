package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gopkg.in/ldap.v2"
)

func getLDAPGroups() {

	l, err := ldap.DialTLS("tcp", viper.GetString("URL"), &tls.Config{ServerName: viper.GetString("ServerName")})
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	err = l.Bind(viper.GetString("User"), viper.GetString("Password"))
	if err != nil {
		log.Fatal(err)
	}

	searchRequest := ldap.NewSearchRequest(
		viper.GetString("BaseDN"), // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		viper.GetString("Filter"),          // The filter to apply
		viper.GetStringSlice("Attributes"), // A list attributes to retrieve
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
