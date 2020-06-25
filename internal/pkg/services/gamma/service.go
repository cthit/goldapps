package gamma

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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

func put(arr []string, value string) []string {
	if arr == nil {
		return []string{value}
	}

	for _, v := range arr {
		if v == value {
			return arr
		}
	}

	return append(arr, value)
}

func emptyGroup(emailPrefix string) model.Group {
	return model.Group{
		Email:      fmt.Sprintf("%s@chalmers.it", emailPrefix),
		Type:       "",
		Members:    []string{},
		Aliases:    nil,
		Expendable: false,
	}
}

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

func createMailPostMap(posts []Post) map[string]map[string]model.Group {
	mailPostMap := make(map[string]map[string]model.Group)
	for _, post := range posts {
		mailPostMap[post.EmailPrefix] = make(map[string]model.Group)
		mailPostMap[post.EmailPrefix]["kommitteer"] = emptyGroup(post.EmailPrefix + ".kommitteer")
	}
	return mailPostMap
}

func insertPostUsers(groups []FKITGroup, mailPostMap *map[string]map[string]model.Group) {
	var prefix string
	var groupName string
	var postMailPrefix string
	var tmpGroup model.Group

	for _, group := range groups {
		for _, member := range group.GroupMembers {
			prefix = member.Post.EmailPrefix
			groupName = group.SuperGroup.Name

			if prefix == "" || !member.Gdpr || group.SuperGroup.Type == "ALUMNI" {
				continue
			}

			if _, ok := (*mailPostMap)[prefix][groupName]; !ok {
				postMailPrefix = fmt.Sprintf("%s.%s", prefix, groupName)
				(*mailPostMap)[prefix][groupName] = emptyGroup(postMailPrefix)

				if group.SuperGroup.Type == "COMMITTEE" {
					tmpGroup = (*mailPostMap)[prefix]["kommitteer"]
					tmpGroup.Members = append(tmpGroup.Members, postMailPrefix+"@chalmers.it")
					(*mailPostMap)[prefix]["kommitteer"] = tmpGroup
				}
			}

			tmpGroup = (*mailPostMap)[prefix][groupName]
			tmpGroup.Members = append(tmpGroup.Members, member.Email)
			(*mailPostMap)[prefix][groupName] = tmpGroup
		}
	}
}

func convertPostMailGroups(mailPostMap *map[string]map[string]model.Group) []model.Group {
	mailGroups := []model.Group{}
	var postGroupMail model.Group

	for postName, postMap := range *mailPostMap {
		postGroupMail = emptyGroup(postName)
		for _, group := range postMap {
			mailGroups = append(mailGroups, group)
			if !strings.Contains(group.Email, "kommitteer") {
				postGroupMail.Members = append(postGroupMail.Members, group.Email)
			}
		}
		mailGroups = append(mailGroups, postGroupMail)
	}

	return mailGroups
}

func (s GammaService) GetGroups() ([]model.Group, error) {
	groups, err := getGammaGroups(&s)
	if err != nil {
		panic(err)
	}
	superGroups, err := getSuperGroups(&s)
	if err != nil {
		panic(err)
	}
	posts, err := getMailPosts(&s)
	if err != nil {
		panic(err)
	}

	mailPostMap := createMailPostMap(posts)
	insertPostUsers(groups, &mailPostMap)

	formattedGroups := append(formatGroups(groups), formatSuperGroups(superGroups, groups)...)
	formattedGroups = append(formattedGroups, convertPostMailGroups(&mailPostMap)...)

	return formattedGroups, nil
}

func getActiveGroups(s *GammaService) ([]FKITGroup, error) {
	groups := struct {
		GetFKITGroupResponse []FKITGroup `json:"getFKITGroupResponse"`
	}{}
	err := gammaReq(s, "/api/groups/active", &groups)
	return groups.GetFKITGroupResponse, err
}

func shouldHaveMail(group FKITGroup, member FKITUser) bool {
	return group.Active &&
		(group.SuperGroup.Type == "COMMITTEE" || group.SuperGroup.Type == "BOARD") &&
		member.Gdpr
}

func extractUsers(groups []FKITGroup) []model.User {
	userFound := make(map[string]bool)
	users := []model.User{}
	var newMember model.User

	for _, group := range groups {
		for _, member := range group.GroupMembers {
			if shouldHaveMail(group, member) && !userFound[member.Cid] {
				newMember = model.User{}
				newMember.Cid = member.Cid
				newMember.FirstName = member.FirstName
				newMember.SecondName = member.LastName
				newMember.Nick = member.Nick
				newMember.Mail = member.Email
				users = append(users, newMember)
				userFound[member.Cid] = true
			}
		}
	}

	return users
}

func (s GammaService) GetUsers() ([]model.User, error) {
	groups, err := getActiveGroups(&s)
	if err != nil {
		return nil, err
	}
	return extractUsers(groups), nil
}
