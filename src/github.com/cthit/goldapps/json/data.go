package json

import "../../goldapps"

type data struct {
	Groups []goldapps.Group `json:"groups"`
	Users  []goldapps.User  `json:"users"`
}
