package main

import (
	"fmt"
	"log"
	"gopkg.in/ldap.v2"
)

func getLDAPGroups() {
	config, err := readLDAPConfig("ldap.json")
	if err != nil {
		log.Fatal(err)
	}

	l, err := ldap.Dial("tcp", config.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	err = l.Bind(config.User, config.Password)
	if err != nil {
		log.Fatal(err)
	}

	searchRequest := ldap.NewSearchRequest(
		config.BaseDN,       // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		config.Filter,       // The filter to apply
		config.Attributes,   // A list attributes to retrieve
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
