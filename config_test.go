package main

import "testing"

func TestConfig(t *testing.T) {

	loadConfig()

}

func TestServiceKeyLoading(t *testing.T) {

	_, err := getGoogleServiceKey()
	if err != nil{
		t.Fail()
	}

}