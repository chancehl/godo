package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chancehl/godo/internal/model"
	"github.com/google/go-github/v50/github"
)

type githubService interface {
	CreateGist(godos []model.GodoItem) (string, string, error)
	GetGists() ([]*github.Gist, *github.Response, error)
}

type GithubService struct {
	githubClient *github.Client
	context      context.Context
}

func NewGithubService(githubClient *github.Client, context context.Context) *GithubService {
	return &GithubService{githubClient, context}
}

func (service *GithubService) CreateGist(godos []model.GodoItem) (string, string, error) {
	ctx := context.Background()

	gistContent, _ := json.Marshal(godos)

	gist := &github.Gist{
		Description: github.String("godo"),
		Public:      github.Bool(false),
		Files: map[github.GistFilename]github.GistFile{
			"godo.json": {
				Content: github.String(string(gistContent)),
			},
		},
	}

	createdGist, _, err := service.githubClient.Gists.Create(ctx, gist)
	if err != nil {
		return "", *createdGist.HTMLURL, fmt.Errorf("failed to create gist file (err=%s)", err)
	}

	return *createdGist.ID, *createdGist.HTMLURL, nil
}

func (service *GithubService) GetGists() ([]*github.Gist, *github.Response, error) {
	return service.githubClient.Gists.List(service.context, "", &github.GistListOptions{})
}
