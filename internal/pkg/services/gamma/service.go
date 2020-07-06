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

	if !isKit(group) || !group.Active {
		for _, member := range group.GroupMembers {
			members[getMemberEmail(group, &member)] = true
		}
	} else {
		for _, member := range group.GroupMembers {
			if shouldHaveMail(group, &member) {
				members[getMemberEmail(group, &member)] = true
			}
		}
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

//Makes a group which contains all the emails of groups with a specified prefix
func groupOfGroups(mailPrefix string, requiredPrefix string, groups []model.Group) model.Group {
	newGroup := emptyGroup(mailPrefix)
	for _, group := range groups {
		if strings.HasPrefix(group.Email, requiredPrefix) {
			newGroup.Members = append(newGroup.Members, group.Email)
		}
	}

	return newGroup
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

//Creates all mail groups connected to posts
func getPostMails(fkitGroups []FKITGroup) []model.Group {
	groupList := &PostGroupList{}
	for _, group := range fkitGroups {
		for _, member := range group.GroupMembers {
			groupList = groupList.insert(&group, &member)
		}
	}

	kit, other := groupList.toGroups()

	ordforanden := groupOfGroups("ordforanden", "ordf", append(kit, other...))
	kitOrdforanden := groupOfGroups("ordforanden.kommitteer", "ordf", kit)
	kassorer := groupOfGroups("kassorer", "kassor", append(kit, other...))
	kitKassorer := groupOfGroups("kassorer.kommitteer", "kassor", kit)

	return append(append(kit, ordforanden, kitOrdforanden, kassorer, kitKassorer), other...)
}

func (s GammaService) GetGroups() ([]model.Group, error) {
	groups, err := getGammaGroups(&s)
	if err != nil {
		log.Println("Failed to fetch all groups from Gamma")
		panic(err)
	}
	activeGroups, err := getActiveGroups(&s)
	if err != nil {
		log.Println("Failed to fetch active groups")
		panic(err)
	}

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
