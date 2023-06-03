package gamma

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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

	err := gammaReq(s, "/api/admin/groups", &groups)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return groups.Groups, nil
}

//Fetches all active groups
func getActiveGroups(s *GammaService) ([]FKITGroup, error) {
	groups := struct {
		GetFKITGroupResponse []FKITGroup `json:"groups"`
	}{}
	err := gammaReq(s, "/api/admin/groups", &groups)
	activeGroups := []FKITGroup{}
	for i := range groups.GetFKITGroupResponse {
		group := groups.GetFKITGroupResponse[i]
		if group.SuperGroup.Type != "ALUMNI" {
			activeGroups = append(activeGroups, group)
		}
	}

	return activeGroups, err
}
