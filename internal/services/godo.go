package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chancehl/godo/internal/model"
	"github.com/google/go-github/v50/github"
)

type godoService interface {
	GetGodos(id string) ([]model.GodoItem, error)
	UpdateGodos(items []model.GodoItem) error
}

type GodoService struct {
	githubClient *github.Client
	context      context.Context
	gistID       string
}

func NewGodoService(githubClient *github.Client, context context.Context, gistID string) *GodoService {
	return &GodoService{githubClient, context, gistID}
}

func (service *GodoService) GetGodos() ([]model.GodoItem, error) {
	gist, resp, err := service.githubClient.Gists.Get(service.context, service.gistID)

	if err != nil || resp.StatusCode != 200 {
		return []model.GodoItem{}, err
	}

	content := *gist.Files[GistFileName].Content

	var items []model.GodoItem

	if err := json.Unmarshal([]byte(content), &items); err != nil {
		return []model.GodoItem{}, fmt.Errorf("failed to parse godo items from gist content (content=%s, err=%s)", content, err)
	}

	return items, nil
}

func (service *GodoService) UpdateGodos(godos []model.GodoItem) error {
	content, err := json.Marshal(godos)
	if err != nil {
		return fmt.Errorf("could not serialize items: %s", err)
	}

	gist := &github.Gist{
		Files: map[github.GistFilename]github.GistFile{
			"godo.json": {
				Content: github.String(string(content)),
			},
		},
	}

	_, _, err = service.githubClient.Gists.Edit(service.context, service.gistID, gist)
	if err != nil {
		return fmt.Errorf("failed to update gist: %s", err)
	}

	return nil
}
