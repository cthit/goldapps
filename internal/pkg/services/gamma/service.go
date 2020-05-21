package gamma

import (
	"github.com/cthit/goldapps/internal/pkg/model"
)

type GammaService struct {
	apiKey   string
	gammaUrl string
}

func CreateGammaService(apiKey string, url string) (GammaService, error) {
	return GammaService{
		apiKey:   apiKey,
		gammaUrl: url,
	}, nil
}

func (s GammaService) GetGroups() ([]model.Group, error) {
	return []model.Group{}, nil
}

func (s GammaService) GetUsers() ([]model.User, error) {
	return []model.User{}, nil
}
