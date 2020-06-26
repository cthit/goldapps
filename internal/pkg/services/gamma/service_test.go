package gamma

import (
	"testing"

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

	assert.Equal(t, shouldHaveMail(groupA, userA), false)
	assert.Equal(t, shouldHaveMail(groupA, userB), false)
	assert.Equal(t, shouldHaveMail(groupB, userA), true)
	assert.Equal(t, shouldHaveMail(groupB, userB), false)
	assert.Equal(t, shouldHaveMail(groupC, userA), false)
	assert.Equal(t, shouldHaveMail(groupC, userB), false)
	assert.Equal(t, shouldHaveMail(groupD, userA), true)
	assert.Equal(t, shouldHaveMail(groupD, userB), false)
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

	assert.Equal(t, getMembers(groupA), []string{"usera@chalmers.it", "userb@chalmers.it"})
	assert.Equal(t, getMembers(groupB), []string{"usera@gmail.com", "userb@gmail.com"})
}
