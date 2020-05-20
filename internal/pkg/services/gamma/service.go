package gamma

import (
	"github.com/cthit/goldapps/internal/pkg/model"
)

type GammaService struct {
	apiKey   string
	gammaUrl string
}

func CreateGammaService() GammaService {
	return GammaService{
		apiKey:   "key",
		gammaUrl: "http://localhost:8081",
	}
}

func (s GammaService) GetGroups() ([]model.Group, error) {
	return []model.Group{}, nil
}

func (s GammaService) GetUsers() ([]model.User, error) {
	return []model.User{}, nil
}
