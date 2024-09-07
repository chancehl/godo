package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chancehl/godo/internal/model"
	"github.com/google/go-github/v50/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func GetGithubClient(ctx context.Context) (*github.Client, error) {
	envFileData, _ := godotenv.Read(".env")

	accessToken := envFileData["GITHUB_ACCESS_TOKEN"]

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tokenClient := oauth2.NewClient(ctx, tokenSource)

	githubClient := github.NewClient(tokenClient)

	return githubClient, nil
}

func CreateGist(items []model.GodoItem) (string, string, error) {
	ctx := context.Background()

	githubClient, err := GetGithubClient(ctx)

	if err != nil {
		return "", "", err
	}

	gistContent, _ := json.Marshal(items)

	gist := &github.Gist{
		Description: github.String("godo"),
		Public:      github.Bool(false),
		Files: map[github.GistFilename]github.GistFile{
			"godo.json": {
				Content: github.String(string(gistContent)),
			},
		},
	}

	createdGist, _, err := githubClient.Gists.Create(ctx, gist)
	if err != nil {
		return "", *createdGist.HTMLURL, err
	}

	return *createdGist.ID, *createdGist.HTMLURL, nil
}

func GetGodos(id string) ([]model.GodoItem, error) {
	ctx := context.Background()

	githubClient, err := GetGithubClient(ctx)

	if err != nil {
		return []model.GodoItem{}, err
	}

	gist, resp, err := githubClient.Gists.Get(ctx, id)

	if err != nil || resp.StatusCode != 200 {
		return []model.GodoItem{}, err
	}

	content := *gist.Files["godo.json"].Content

	var items []model.GodoItem

	if err := json.Unmarshal([]byte(content), &items); err != nil {
		fmt.Printf("Failed to parse godo items from gist content: %s, %s", content, err)
		return []model.GodoItem{}, nil
	}

	return items, nil
}

func UpdateGodos(id string, items []model.GodoItem) error {
	ctx := context.Background()

	content, err := json.Marshal(items)
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

	githubClient, err := GetGithubClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to update godos (%s)", err)
	}

	_, resp, err := githubClient.Gists.Edit(ctx, id, gist)

	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("failed to update gist: %s", err)
	}

	return nil
}

func GetGists() ([]*github.Gist, *github.Response, error) {
	ctx := context.Background()

	client, err := GetGithubClient(ctx)

	if err != nil {
		return nil, nil, err
	}

	return client.Gists.List(ctx, "", &github.GistListOptions{})
}
