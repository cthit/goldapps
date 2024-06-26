package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cthit/goldapps/internal/pkg/model"
)

type AuthService struct {
	apiKey string
	url    string
}

// Creates a auth service which has the url to auth and the pre-shared key
func CreateAuthService(apiKey string, url string) (AuthService, error) {
	return AuthService{
		apiKey: apiKey,
		url:    url,
	}, nil
}

// Executes a generic get request with api key
func request(s *AuthService, endpoint string, response interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", s.url, endpoint), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("pre-shared %s", s.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Printf("Request sent to: %s [key: %s] status %d\n", endpoint, s.apiKey, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s AuthService) GetGroups() ([]model.Group, error) {
	var groups AuthSuperGroups

	err := request(&s, "/api/account-scaffold/v1/supergroups", &groups)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return groups.ToGroups(), nil
}

func (s AuthService) GetUsers() ([]model.User, error) {
	var users AuthUsers

	err := request(&s, "/api/account-scaffold/v1/users", &users)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return users.ToUsers(), nil
}
