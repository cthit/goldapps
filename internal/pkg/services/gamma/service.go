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

//Creates a gamma service which has the url to gamma and the pre-shared key
func CreateGammaService(apiKey string, url string) (GammaService, error) {
	return GammaService{
		apiKey:   apiKey,
		gammaUrl: url,
	}, nil
}

//Executes a generic get request to Gamma
func gammaReq(s *GammaService, endpoint string, response interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", s.gammaUrl, endpoint), nil)
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

	body, err := ioutil.ReadAll(resp.Body)
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

//Fetches all the fkit-groups from Gamma
func getGammaGroups(s *GammaService) ([]FKITGroup, error) {
	var groups struct {
		Groups []FKITGroup `json:"groups"`
	}

	err := gammaReq(s, "/api/groups", &groups)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return groups.Groups, nil
}

//Fetches all super groups from Gamma
func getSuperGroups(s *GammaService) ([]FKITSuperGroup, error) {
	var superGroups []FKITSuperGroup

	err := gammaReq(s, "/api/superGroups", &superGroups)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return superGroups, nil
}

//Returns all the Email addresses from the member of a group
func getMembers(group FKITGroup) []string {
	members := make([]string, len(group.GroupMembers))
	for k, v := range group.GroupMembers {
		members[k] = v.Email
	}
	return members
}

//Fetches the emails of the groups with the specified super group
func getGroupEmails(superGroupId string, normalGroups []FKITGroup) []string {
	emails := make([]string, 0)
	for _, v := range normalGroups {
		if superGroupId == v.SuperGroup.ID {
			emails = append(emails, v.Email)
		}
	}

	return emails
}

//Formats a list of super groups to a list of model.Group
func formatSuperGroups(superGroups []FKITSuperGroup, normalGroups []FKITGroup) (groups []model.Group) {
	groups = make([]model.Group, len(superGroups))
	for k, v := range superGroups {
		groups[k].Email = v.Email
		groups[k].Type = v.Type
		groups[k].Members = getGroupEmails(v.ID, normalGroups)
		groups[k].Expendable = false
	}

	return groups
}

//Formats a list of FKIT-groups to a list of model.Group
func formatGroups(gammaGroups []FKITGroup) (groups []model.Group) {
	groups = make([]model.Group, len(gammaGroups))
	for k, v := range gammaGroups {
		groups[k].Email = v.Email
		groups[k].Type = v.SuperGroup.Type
		groups[k].Members = getMembers(v)
		groups[k].Expendable = false
	}
	return groups
}

func (s GammaService) GetGroups() ([]model.Group, error) {
	groups, _ := getGammaGroups(&s)
	superGroups, _ := getSuperGroups(&s)

	formattedGroups := append(formatGroups(groups), formatSuperGroups(superGroups, groups)...)

	return formattedGroups, nil
}

func (s GammaService) GetUsers() ([]model.User, error) {
	return []model.User{}, nil
}
