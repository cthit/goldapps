package gamma

import (
	"fmt"
	"log"
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

//Determins if the specified user in the specified group should have a gsuit account
func shouldHaveMail(group *FKITGroup, member *FKITUser) bool {
	return isKit(group) && member.Gdpr && group.Active
}

//Determins if specified group is a member of KIT
func isKit(group *FKITGroup) bool {
	return (group.SuperGroup.Type == "COMMITTEE" || group.SuperGroup.Type == "BOARD" || group.SuperGroup.Type == "FUNCTIONARIES")
}

//Returns the email which should be used for a specific user
func getMemberEmail(group *FKITGroup, member *FKITUser) string {
	if shouldHaveMail(group, member) {
		return fmt.Sprintf("%s@chalmers.it", member.Cid)
	}

	return member.Email
}

//Returns all the Email addresses from the member of a group
func getMembers(group *FKITGroup) []string {
	members := make(map[string]bool, len(group.GroupMembers))
	for _, v := range group.GroupMembers {
		members[getMemberEmail(group, &v)] = true
	}

	membersMail := []string{}
	for mail := range members {
		membersMail = append(membersMail, mail)
	}
	return membersMail
}

//Creates an empty group with a specific email
func emptyGroup(emailPrefix string) model.Group {
	return model.Group{
		Email:      fmt.Sprintf("%s@chalmers.it", emailPrefix),
		Type:       "",
		Members:    []string{},
		Aliases:    nil,
		Expendable: false,
	}
}

//Creates a map with of <post, <committee>, mailGroup> from all posts
func createMailPostMap(posts []Post) map[string]map[string]model.Group {
	mailPostMap := make(map[string]map[string]model.Group)
	for _, post := range posts {
		mailPostMap[post.EmailPrefix] = make(map[string]model.Group)
		mailPostMap[post.EmailPrefix]["kommitteer"] = emptyGroup(post.EmailPrefix + ".kommitteer")
	}
	return mailPostMap
}

//Inserts a new member to the mail group
func appendMember(mailPostMap *map[string]map[string]model.Group, post string, groupName string, member string) {
	tmpGroup := (*mailPostMap)[post][groupName]
	tmpGroup.Members = append(tmpGroup.Members, member)
	(*mailPostMap)[post][groupName] = tmpGroup
}

//Populates the map of post-mail-groups with the members for each post
func insertPostUsers(groups []FKITGroup, mailPostMap *map[string]map[string]model.Group) {
	var prefix string
	var groupName string
	var mailPrefix string

	for _, group := range groups {
		for _, member := range group.GroupMembers {
			prefix = member.Post.EmailPrefix
			groupName = group.SuperGroup.Name

			if prefix == "" {
				continue
			}

			if _, ok := (*mailPostMap)[prefix][groupName]; !ok {
				mailPrefix = fmt.Sprintf("%s.%s", prefix, groupName)
				(*mailPostMap)[prefix][groupName] = emptyGroup(mailPrefix)

				if isKit(&group) {
					appendMember(mailPostMap, prefix, "kommitteer", mailPrefix+"@chalmers.it")
				}
			}

			appendMember(mailPostMap, prefix, groupName, getMemberEmail(&group, &member))
		}
	}
}

//Converts the map of post-mail-groups to an array of post-mail-groups
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

//Creates all mail groups from the fkit groups except post specific groups
func getGroups(fkitGroups []FKITGroup) []model.Group {
	groupList := &SuperGroupList{}
	for _, group := range fkitGroups {
		groupList = groupList.insert(&group)
	}

	fkit, kit, groups := groupList.toGroups()
	return append(groups, fkit, kit)
}

func getPostMails(fkitGroups []FKITGroup) []model.Group {
	groupList := &PostGroupList{}
	for _, group := range fkitGroups {
		for _, member := range group.GroupMembers {
			groupList = groupList.insert(&group, &member)
		}
	}

	kit, other := groupList.toGroups()
	return append(kit, other...)
}

func (s GammaService) GetGroups() ([]model.Group, error) {
	groups, err := getGammaGroups(&s)
	if err != nil {
		log.Println("Failed to fetch all groups from Gamma")
		panic(err)
	}
	/*posts, err := getMailPosts(&s)
	if err != nil {
		log.Println("Failed to fetch all posts from Gamma")
		panic(err)
	}*/
	activeGroups, err := getActiveGroups(&s)
	if err != nil {
		log.Println("Failed to fetch active groups")
		panic(err)
	}

	/*mailPostMap := createMailPostMap(posts)
	insertPostUsers(activeGroups, &mailPostMap)*/

	formattedGroups := getGroups(groups)
	formattedGroups = append(formattedGroups, getPostMails(activeGroups)...)

	return formattedGroups, nil
}

//Fetches all the users in the specified groups
func extractUsers(groups []FKITGroup) []model.User {
	userFound := make(map[string]bool)
	users := []model.User{}

	for _, group := range groups {
		for _, member := range group.GroupMembers {
			if shouldHaveMail(&group, &member) && !userFound[member.Cid] {
				users = append(users, member.toUser(&group))
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
