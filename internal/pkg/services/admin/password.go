package admin

import (
	"github.com/sethvargo/go-password/password"
	"math/rand"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"encoding/base64"
)

const passwordMailBody = "Action required! You are a member of a committee at the IT-section and have therefor been provided a google-account by the section. Login within the following week to setup two-factor-authentication or you might get locked out from your account. You can login on any google service such as gmail.google.com or drive.google.com with cid@chalmers.it and your provided password: %s"
const passwordMailSubject = "Login details for google services at chalmers.it"

func newPassword() string {
	numbers := rand.Intn(10) + 5
	//symbols := rand.Intn(10) + 5
	pass, err := password.Generate(64, numbers, 0, false, true)
	if err != nil {
		panic("Password generation failed")
	}
	return pass
}

func (s googleService) sendPassword(to string, password string) error {

	from := s.admin + "@" + s.domain
	body := fmt.Sprintf(passwordMailBody, password)

	msgRaw := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + passwordMailSubject + "\r\n\r\n" +
		body + "\r\n"

	msg := &gmail.Message{
		Raw: base64.StdEncoding.EncodeToString([]byte(msgRaw)),
	}
	_, err := s.mailService.Users.Messages.Send(from, msg).Do()

	return err
}
