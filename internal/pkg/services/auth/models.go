package auth

import "github.com/cthit/goldapps/internal/pkg/model"

type AuthSuperGroup struct {
	Name       string `json:"name"`
	PrettyName string `json:"prettyName"`
	Type       string `json:"type"`
	Members    []struct {
		Post struct {
			PostID      string `json:"postId"`
			SvText      string `json:"svText"`
			EnText      string `json:"enText"`
			EmailPrefix string `json:"emailPrefix"`
		} `json:"post"`
		User AuthUser `json:"user"`
	} `json:"members"`
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
			Mail:       user.Email,
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

func (superGroup AuthSuperGroup) ToGroup() model.Group {
	group := model.Group{
		Email:   superGroup.Name + "@chalmers.it",
		Type:    superGroup.Type,
		Aliases: []string{},
	}
	for _, member := range superGroup.Members {
		group.Members = append(group.Members, member.User.Email)
	}
	return group
}

func (superGroups AuthSuperGroups) ToGroups() model.Groups {
	groupsList := model.Groups{}
	for _, superGroup := range superGroups {
		groupsList = append(groupsList, superGroup.ToGroup())
	}
	return groupsList
}
