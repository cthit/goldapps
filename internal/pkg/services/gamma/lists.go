package gamma

import "github.com/cthit/goldapps/internal/pkg/model"

//Both SuperGroupList and NormalGroupList follows a linked list structure

type SuperGroupList struct {
	Next         *SuperGroupList
	MemberGroups *NormalGroupList
	Kit          bool
	model.Group
}

type NormalGroupList struct {
	Next   *NormalGroupList
	Active bool
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
	if group.Email == li.Email {
		li.Members = append(li.Members, getMembers(group)...)
		return li
	}

	li.Next = li.Next.insert(group)
	return li
}

//Creates a group item which contains the group email, members and a pointer to the next item
func (next *NormalGroupList) newListItem(group *FKITGroup) *NormalGroupList {
	return &NormalGroupList{
		Next:   next,
		Active: group.Active || group.SuperGroup.Type == "ALUMNI",
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

	if li.Active {
		return append(active, li.Group), inactive
	}

	return active, append(inactive, li.Group)
}

func (li *SuperGroupList) insert(group *FKITGroup) *SuperGroupList {
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
		return emptyGroup("fkit"), emptyGroup("kit"), []model.Group{}
	}

	activeGroups, inactiveGroups := li.MemberGroups.toGroups()
	superGroup := li.Group

	for _, group := range activeGroups {
		superGroup.Members = append(superGroup.Members, group.Email)
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
