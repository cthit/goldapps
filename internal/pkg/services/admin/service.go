package admin

import (
	"github.com/cthit/goldapps/internal/pkg/services"

	"google.golang.org/api/admin/directory/v1" // Imports as admin
	"google.golang.org/api/gmail/v1"           // Imports as gmail

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"io/ioutil"
	"strings"
)

const googleDuplicateEntryError = "googleapi: Error 409: Entity already exists., duplicate"

// my_customer seems to work...
const googleCustomer = "my_customer"

type googleService struct {
	adminService *admin.Service
	mailService  *gmail.Service
	admin        string
	domain       string
}

func NewGoogleService(keyPath string, adminMail string) (services.UpdateService, error) {

	jsonKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	// Parse jsonKey
	config, err := google.JWTConfigFromJSON(jsonKey, Scopes()...)
	if err != nil {
		return nil, err
	}

	// Why do I need this??
	config.Subject = adminMail

	// Create a http client
	client := config.Client(context.Background())

	service, err := admin.New(client)
	if err != nil {
		return nil, err
	}

	mailService, err := gmail.New(client)
	if err != nil {
		return nil, err
	}

	// Extract account and mail
	s := strings.Split(adminMail, "@")
	admin := s[0]
	domain := s[1]

	gs := googleService{
		adminService: service,
		mailService:  mailService,
		admin:        admin,
		domain:       domain,
	}

	return gs, nil
}
