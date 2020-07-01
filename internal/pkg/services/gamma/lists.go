package gamma

import (
	"fmt"

	"github.com/cthit/goldapps/internal/pkg/model"
)

//Both SuperGroupList and NormalGroupList follows a linked list structure

type SuperGroupList struct {
	SuperGroupId string
	Next         *SuperGroupList
	MemberGroups *NormalGroupList
	Kit          bool
	model.Group
}

type NormalGroupList struct {
	GroupId string
	Next    *NormalGroupList
	Active  bool
	model.Group
}

type PostGroupList struct {
	Next        *PostGroupList
	GroupName   string
	EmailPrefix string
	Kit         bool
	model.Group
}

func (li *NormalGroupList) insert(group *FKITGroup) *NormalGroupList {
	//If you have reached the last item of the chain
	//create a new item connected to a group
	if li.Next == nil {
		return li.newListItem(group)
	}

	//If you have reached the correct group
	//append the member emails
	if group.ID == li.GroupId {
		li.Members = append(li.Members, getMembers(group)...)
		return li
	}

	li.Next = li.Next.insert(group)
	return li
}

//Creates a group item which contains the group email, members and a pointer to the next item
func (next *NormalGroupList) newListItem(group *FKITGroup) *NormalGroupList {
	return &NormalGroupList{
		GroupId: group.ID,
		Next:    next,
		Active:  group.Active || group.SuperGroup.Type == "ALUMNI",
		Group: model.Group{
			Email:      group.Email,
			Type:       group.SuperGroup.Type,
			Members:    getMembers(group),
			Aliases:    nil,
			Expendable: false,
		},
	}
}

//Returns all active and inactive groups in the list
func (li *NormalGroupList) toGroups() ([]model.Group, []model.Group) {
	if li.Next == nil {
		return []model.Group{}, []model.Group{}
	}

	active, inactive := li.Next.toGroups()

	if len(li.Members) <= 0 {
		return active, inactive
	}

	if li.Active {
		return append(active, li.Group), inactive
	}

	return active, append(inactive, li.Group)
}

func (li *SuperGroupList) insert(group *FKITGroup) *SuperGroupList {
	//Ignoring admin group
	if group.SuperGroup.Type == "ADMIN" {
		return li
	}

	//If you have reached the last item of the chain
	//create a new item connected to a super group
	if li.Next == nil {
		return li.newListItem(group)
	}

	//If you have reached the correct super group
	//save the normal group as a group member
	if group.SuperGroup.Email == li.Email {
		li.MemberGroups = li.MemberGroups.insert(group)
		return li
	}

	li.Next = li.Next.insert(group)
	return li
}

//Creates a super group item which contains the group email, members and a pointer to the next item
func (next *SuperGroupList) newListItem(group *FKITGroup) *SuperGroupList {
	memberGroups := &NormalGroupList{}
	return &SuperGroupList{
		SuperGroupId: group.SuperGroup.ID,
		Next:         next,
		MemberGroups: memberGroups.insert(group),
		Kit:          isKit(group),
		Group: model.Group{
			Email:      group.SuperGroup.Email,
			Type:       group.SuperGroup.Type,
			Members:    []string{},
			Aliases:    nil,
			Expendable: false,
		},
	}
}

//Returns the fkit group and all the rest of the groups
func (li *SuperGroupList) toGroups() (model.Group, model.Group, []model.Group) {
	if li.Next == nil {
		return emptyGroup("grupper"), emptyGroup("kommitteer"), []model.Group{}
	}

	activeGroups, inactiveGroups := li.MemberGroups.toGroups()
	superGroup := li.Group

	for _, group := range activeGroups {
		superGroup.Members = append(superGroup.Members, group.Email)
	}

	if len(superGroup.Members) <= 0 {
		superGroup.Members = []string{"ita.styrit@chalmers.it"}
	}

	fkit, kit, groups := li.Next.toGroups()
	groups = append(groups, inactiveGroups...)
	groups = append(groups, activeGroups...)

	if li.Kit {
		kit.Members = append(kit.Members, li.Email)
	}

	if li.Type != "ALUMNI" {
		fkit.Members = append(fkit.Members, li.Email)
	}

	return fkit, kit, append(groups, superGroup)
}

//Creates a new post mail group
func (pl *PostGroupList) newListItem(group *FKITGroup, member *FKITUser) *PostGroupList {
	return &PostGroupList{
		Next:        pl,
		EmailPrefix: member.Post.EmailPrefix,
		GroupName:   group.SuperGroup.Name,
		Kit:         isKit(group),
		Group: model.Group{
			Email:      fmt.Sprintf("%s.%s", member.Post.EmailPrefix, group.SuperGroup.Email),
			Members:    []string{getMemberEmail(group, member)},
			Aliases:    nil,
			Expendable: false,
		},
	}
}

//Creates a post email for active member if their posts should have email
func (pl *PostGroupList) insert(group *FKITGroup, member *FKITUser) *PostGroupList {
	if !group.Active || member.Post.EmailPrefix == "" || (!member.Gdpr && isKit(group)) {
		return pl
	}

	if pl.Next == nil {
		return pl.newListItem(group, member)
	}

	if pl.GroupName == group.SuperGroup.Name &&
		pl.EmailPrefix == member.Post.EmailPrefix {
		pl.Members = append(pl.Members, getMemberEmail(group, member))
		return pl
	}

	pl.Next = pl.Next.insert(group, member)
	return pl
}

//Returns all post groups in kit and post groups which is not in kit
func (pl *PostGroupList) toGroups() ([]model.Group, []model.Group) {
	if pl.Next == nil {
		return []model.Group{}, []model.Group{}
	}

	kit, other := pl.Next.toGroups()

	if pl.Kit {
		kit = append(kit, pl.Group)
	} else {
		other = append(other, pl.Group)
	}

	return kit, other
}
