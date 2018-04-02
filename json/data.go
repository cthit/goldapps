package json

import "github.com/cthit/goldapps"

type data struct {
	Groups []goldapps.Group `json:"groups"`
	Users  []goldapps.User  `json:"users"`
}
