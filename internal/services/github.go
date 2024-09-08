package services

import (
	"context"
	"encoding/json"
	"fmt"

	appContext "github.com/chancehl/godo/internal/context"
	"github.com/chancehl/godo/internal/model"
	"github.com/google/go-github/v50/github"
)

const GistFileDescription = "godo"
const GistFilePublic = false
const GistFileName = "godo.json"

type GithubService struct {
	GithubClient       *github.Client
	Context            *context.Context
	ApplicationContext *appContext.ApplicationContext
}

func (service *GithubService) CreateGist(godos []model.GodoItem) (string, string, error) {
	gistContent, _ := json.Marshal(godos)

	gist := &github.Gist{
		Description: github.String(GistFileDescription),
		Public:      github.Bool(GistFilePublic),
		Files: map[github.GistFilename]github.GistFile{
			GistFileName: {
				Content: github.String(string(gistContent)),
			},
		},
	}

	createdGist, _, err := service.GithubClient.Gists.Create(*service.Context, gist)
	if err != nil {
		return "", *createdGist.HTMLURL, fmt.Errorf("failed to create gist file (err=%s)", err)
	}

	return *createdGist.ID, *createdGist.HTMLURL, nil
}

func (service *GithubService) GetGists() ([]*github.Gist, *github.Response, error) {
	return service.GithubClient.Gists.List(*service.Context, "", &github.GistListOptions{})
}

func (service *GithubService) GetGist(id string) (*github.Gist, *github.Response, error) {
	return service.GithubClient.Gists.Get(*service.Context, id)
}

func (service *GithubService) EditGist(id string, gist *github.Gist) (*github.Gist, *github.Response, error) {
	return service.GithubClient.Gists.Edit(*service.Context, id, gist)
}
