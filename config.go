package main

import (
	"io/ioutil"
	"encoding/json"
)

type LDAPConfig struct {
	URL          string `json:url`
	User         string `json:user`
	Password     string `json:password`
	BaseDN       string `json:basedn`
	Filter       string `json:filter`
	Attributes []string `json:attributes`
}

func readLDAPConfig(filename string) (*LDAPConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config LDAPConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
