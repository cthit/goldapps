package auth

import (
	"github.com/cthit/goldapps/internal/pkg/model"
	"slices"
	"strings"
)

type AuthSuperGroup struct {
	Name              string      `json:"name"`
	PrettyName        string      `json:"prettyName"`
	Type              string      `json:"type"`
	Groups            []AuthGroup `json:"groups"`
	UseManagedAccount bool        `json:"useManagedAccount"`
}

type AuthGroup struct {
	Name       string `json:"name"`
	PrettyName string `json:"prettyName"`
	Members    []struct {
		Post struct {
			PostID      string `json:"postId"`
			SvText      string `json:"svText"`
			EnText      string `json:"enText"`
			EmailPrefix string `json:"emailPrefix"`
		} `json:"post"`
		User AuthUser `json:"user"`
	}
}

type AuthUser struct {
	Email     string `json:"email"`
	Cid       string `json:"cid"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Nick      string `json:"nick"`
}

type AuthSuperGroups []AuthSuperGroup

type AuthUsers []AuthUser

func (user AuthUser) ToUser() model.Users {
	return model.Users{
		model.User{
			Cid:        user.Cid,
			FirstName:  user.FirstName,
			SecondName: user.LastName,
			Nick:       user.Nick,
			Mail:       strings.ToLower(user.Cid) + "@chalmers.it",
		},
	}
}

func (users AuthUsers) ToUsers() model.Users {
	usersList := model.Users{}
	for _, user := range users {
		usersList = append(usersList, user.ToUser()...)
	}
	return usersList
}

func (superGroup AuthSuperGroup) ToGroups() model.Groups {
	groups := model.Groups{}
	outerGroup := model.Group{
		Email:   strings.ToLower(superGroup.Name) + "@chalmers.it",
		Type:    superGroup.Type,
		Aliases: []string{},
	}

	superGroupPostGroups := make(map[string]model.Group)

	for _, group := range superGroup.Groups {
		memberGroup := model.Group{
			Email:   strings.ToLower(group.Name) + "@chalmers.it",
			Type:    superGroup.Type,
			Aliases: []string{},
		}
		memberGroupPostGroups := make(map[string]model.Group)

		if group.Members == nil || len(group.Members) == 0 {
			continue
		}

		for _, member := range group.Members {

			var memberEmail string
			if superGroup.UseManagedAccount {
				memberEmail = member.User.Cid + "@chalmers.it"
			} else {
				memberEmail = member.User.Email
			}

			if member.Post.EmailPrefix != "" {
				memberGroupPostGroup, exists := memberGroupPostGroups[member.Post.EmailPrefix]
				if !exists {
					memberGroupPostGroup = model.Group{
						Email:   strings.ToLower(member.Post.EmailPrefix + "." + group.Name + "@chalmers.it"),
						Type:    superGroup.Type,
						Aliases: []string{},
					}
				}
				memberGroupPostGroup.Members = append(memberGroupPostGroup.Members, memberEmail)
				memberGroupPostGroups[member.Post.EmailPrefix] = memberGroupPostGroup

				postSuperGroup, exists := superGroupPostGroups[member.Post.EmailPrefix]
				if !exists {
					postSuperGroup = model.Group{
						Email:   strings.ToLower(member.Post.EmailPrefix + "." + superGroup.Name + "@chalmers.it"),
						Type:    superGroup.Type,
						Aliases: []string{},
					}
				}
				if !slices.Contains(postSuperGroup.Members, memberGroupPostGroup.Email) {
					postSuperGroup.Members = append(postSuperGroup.Members, memberGroupPostGroup.Email)
					superGroupPostGroups[member.Post.EmailPrefix] = postSuperGroup
				}
			}

			memberGroup.Members = append(memberGroup.Members, memberEmail)
		}
		outerGroup.Members = append(outerGroup.Members, memberGroup.Email)
		for _, memberGroupPostGroup := range memberGroupPostGroups {
			groups = append(groups, memberGroupPostGroup)
		}
		groups = append(groups, memberGroup)
	}
	for _, postGroup := range superGroupPostGroups {
		groups = append(groups, postGroup)
	}
	groups = append(groups, outerGroup)
	return groups
}

func (superGroups AuthSuperGroups) ToGroups() model.Groups {
	groupsList := model.Groups{}
	for _, superGroup := range superGroups {
		groupsList = append(groupsList, superGroup.ToGroups()...)
	}
	return groupsList
}
