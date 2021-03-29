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

//Fetches all posts which has a mail prefix
func getMailPosts(s *GammaService) ([]Post, error) {
	posts := []Post{}
	err := gammaReq(s, "/api/groups/posts", &posts)
	if err != nil {
		return nil, err
	}

	mailPosts := []Post{}
	for _, post := range posts {
		if post.EmailPrefix != "" {
			mailPosts = append(mailPosts, post)
		}
	}

	return mailPosts, nil
}

//Fetches all active groups
func getActiveGroups(s *GammaService) ([]FKITGroup, error) {
	groups := struct {
		GetFKITGroupResponse []FKITGroup `json:"getFKITGroupResponse"`
	}{}
	err := gammaReq(s, "/api/groups/active", &groups)
	return groups.GetFKITGroupResponse, err
}
