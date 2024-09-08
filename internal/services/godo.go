package services

import (
	"encoding/json"
	"fmt"

	appContext "github.com/chancehl/godo/internal/context"
	"github.com/chancehl/godo/internal/model"
	"github.com/google/go-github/v50/github"
)

type GodoService struct {
	GithubService      *GithubService
	ApplicationContext *appContext.ApplicationContext
}

func (service *GodoService) GetGodos() ([]model.GodoItem, error) {
	gistID := service.ApplicationContext.GistID

	gist, _, err := service.GithubService.GetGist(gistID)

	if err != nil {
		return []model.GodoItem{}, fmt.Errorf("could not retrieve gist by id (err=%s)", err)
	}

	content := *gist.Files[GistFileName].Content

	var items []model.GodoItem

	if err := json.Unmarshal([]byte(content), &items); err != nil {
		return []model.GodoItem{}, fmt.Errorf("failed to parse godo items from gist content (content=%s, err=%s)", content, err)
	}

	return items, nil
}

func (service *GodoService) UpdateGodos(godos []model.GodoItem) error {
	gistID := service.ApplicationContext.GistID

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

	_, _, err = service.GithubService.EditGist(gistID, gist)
	if err != nil {
		return fmt.Errorf("failed to update gist: %s", err)
	}

	return nil
}
