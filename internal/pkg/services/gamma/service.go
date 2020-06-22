package gamma

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cthit/goldapps/internal/pkg/model"
)

type GammaService struct {
	apiKey   string
	gammaUrl string
}

type GammaData struct {
	Groups []model.Group `json:"groups"`
	Users  []model.User  `json:"users"`
}

func CreateGammaService(apiKey string, url string) (GammaService, error) {
	return GammaService{
		apiKey:   apiKey,
		gammaUrl: url,
	}, nil
}

func getGammaData(s *GammaService) (GammaData, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/admin/goldapps", s.gammaUrl), nil)
	if err != nil {
		log.Println(err)
		return GammaData{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("pre-shared %s", s.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return GammaData{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return GammaData{}, err
	}

	var groupsAndUsers GammaData

	err = json.Unmarshal(body, &groupsAndUsers)
	if err != nil {
		log.Println(err)
		return GammaData{}, err
	}

	return groupsAndUsers, nil
}

func (s GammaService) GetGroups() ([]model.Group, error) {
	data, err := getGammaData(&s)
	if err != nil {
		return nil, err
	}
	return data.Groups, nil
}

func (s GammaService) GetUsers() ([]model.User, error) {
	data, err := getGammaData(&s)
	if err != nil {
		return nil, err
	}
	return data.Users, nil
}
