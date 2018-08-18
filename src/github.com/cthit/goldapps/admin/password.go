package admin

import (
	"github.com/sethvargo/go-password/password"
	"math/rand"
)

func newPassword() string {
	numbers := rand.Intn(10) + 5
	//symbols := rand.Intn(10) + 5
	pass, err := password.Generate(64, numbers, 0, false, true)
	if err != nil {
		panic("Password generation failed")
	}
	return pass
}
