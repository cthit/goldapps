package gamma

import (
	"testing"

	"github.com/cthit/goldapps/internal/pkg/model"
	"github.com/magiconair/properties/assert"
)

func TestShouldHaveMail(t *testing.T) {
	userA, userB := FKITUser{
		Gdpr: true,
	}, FKITUser{
		Gdpr: false,
	}

	groupA, groupB, groupC, groupD := FKITGroup{
		SuperGroup: FKITSuperGroup{
			Type: "COMMITTEE",
		},
		Active:       false,
		GroupMembers: []FKITUser{userA, userB},
	}, FKITGroup{
		SuperGroup: FKITSuperGroup{
			Type: "COMMITTEE",
		},
		Active:       true,
		GroupMembers: []FKITUser{userA, userB},
	}, FKITGroup{
		SuperGroup: FKITSuperGroup{
			Type: "BOARD",
		},
		Active:       false,
		GroupMembers: []FKITUser{userA, userB},
	}, FKITGroup{
		SuperGroup: FKITSuperGroup{
			Type: "BOARD",
		},
		Active:       true,
		GroupMembers: []FKITUser{userA, userB},
	}

	assert.Equal(t, shouldHaveMail(&groupA, &userA), false)
	assert.Equal(t, shouldHaveMail(&groupA, &userB), false)
	assert.Equal(t, shouldHaveMail(&groupB, &userA), true)
	assert.Equal(t, shouldHaveMail(&groupB, &userB), false)
	assert.Equal(t, shouldHaveMail(&groupC, &userA), false)
	assert.Equal(t, shouldHaveMail(&groupC, &userB), false)
	assert.Equal(t, shouldHaveMail(&groupD, &userA), true)
	assert.Equal(t, shouldHaveMail(&groupD, &userB), false)
}

func TestGetMembers(t *testing.T) {
	userA, userB := FKITUser{
		Cid:   "usera",
		Email: "usera@gmail.com",
		Gdpr:  true,
	}, FKITUser{
		Cid:   "userb",
		Email: "userb@gmail.com",
		Gdpr:  true,
	}

	groupA, groupB := FKITGroup{
		SuperGroup: FKITSuperGroup{
			Type: "COMMITTEE",
		},
		Active:       true,
		GroupMembers: []FKITUser{userA, userB},
	}, FKITGroup{
		SuperGroup: FKITSuperGroup{
			Type: "SOCIETY",
		},
		Active:       true,
		GroupMembers: []FKITUser{userA, userB, userB},
	}

	assert.Equal(t, getMembers(&groupA), []string{"usera@chalmers.it", "userb@chalmers.it"})
	assert.Equal(t, getMembers(&groupB), []string{"usera@gmail.com", "userb@gmail.com"})
}

func TestGetGroups(t *testing.T) {
	groups := []FKITGroup{
		{
			Active: true,
			SuperGroup: FKITSuperGroup{
				Type:  "COMMITTEE",
				Email: "digit@chalmers.it",
			},
		},
	}

	want := []model.Group{{"digit@chalmers.it", "COMMITTEE", []string{"ita.styrit@chalmers.it"}, nil, false},
		{"grupper@chalmers.it", "", []string{"digit@chalmers.it"}, nil, false},
		{"kommitteer@chalmers.it", "", []string{"digit@chalmers.it"}, nil, false}}

	assert.Equal(t, getGroups(groups), want)
}

func TestGetPostMails(t *testing.T) {
	userA := FKITUser{
		Cid: "usera",
		Post: Post{
			EmailPrefix: "ordf",
		},
		Gdpr:  true,
		Email: "usera@gamil.com",
	}

	groupA := FKITGroup{
		Active: true,
		SuperGroup: FKITSuperGroup{
			Type: "COMMITTEE",
			Name: "supergroup",
		},
		GroupMembers: []FKITUser{userA},
	}

	assert.Equal(t, getPostMails([]FKITGroup{groupA}), []model.Group{{"ordf.supergroup@chalmers.it", "", []string{"usera@chalmers.it"}, nil, false}})
}

func TestList(t *testing.T) {
	userA, userB := FKITUser{
		Cid:   "usera",
		Email: "usera@gmail.com",
		Gdpr:  true,
	}, FKITUser{
		Cid:   "userb",
		Email: "userb@gmail.com",
		Gdpr:  true,
	}

	groupA, groupB, groupC := FKITGroup{
		Email: "digit20@chalmers.it",
		SuperGroup: FKITSuperGroup{
			Type:  "COMMITTEE",
			Email: "digIT@chalmers.it",
		},
		Active:       true,
		GroupMembers: []FKITUser{userA, userB},
	}, FKITGroup{
		Email: "drawit20@chalmers.it",
		SuperGroup: FKITSuperGroup{
			Type:  "SOCIETY",
			Email: "drawit@chalmers.it",
		},
		Active:       true,
		GroupMembers: []FKITUser{userA, userB, userB},
	}, FKITGroup{
		Email: "drawit19@chalmers.it",
		SuperGroup: FKITSuperGroup{
			Type:  "SOCIETY",
			Email: "drawit@chalmers.it",
		},
		Active:       false,
		GroupMembers: []FKITUser{userA, userB, userB},
	}

	want := []model.Group{{"drawit19@chalmers.it", "SOCIETY", []string{"usera@gmail.com", "userb@gmail.com"}, []string{}, false},
		{"drawit20@chalmers.it", "SOCIETY", []string{"usera@gmail.com", "userb@gmail.com"}, []string{}, false},
		{"drawit@chalmers.it", "SOCIETY", []string{"drawit20@chalmers.it"}, []string{}, false},
		{"digit20@chalmers.it", "COMMITTEE", []string{"usera@chalmers.it", "userb@chalmers.it"}, []string{}, false},
		{"digIT@chalmers.it", "COMMITTEE", []string{"digit20@chalmers.it"}, []string{}, false},
		{"fkit@chalmers.it", "", []string{"drawit@chalmers.it", "digIT@chalmers.it"}, []string{}, false},
		{"kit@chalmers.it", "", []string{"digIT@chalmers.it"}, []string{}, false}}

	assert.Equal(t, getGroups([]FKITGroup{groupA, groupB, groupC}), want)
}
