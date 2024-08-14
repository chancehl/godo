package clients

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chancehl/godo/internal/model"
	"github.com/google/go-github/v50/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func GetGithubClient(ctx context.Context) *github.Client {
	envFileData, _ := godotenv.Read(".env")

	accessToken := envFileData["GITHUB_ACCESS_TOKEN"]

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tokenClient := oauth2.NewClient(ctx, tokenSource)

	githubClient := github.NewClient(tokenClient)

	return githubClient
}

func CreateGist(items []model.GodoItem) *string {
	ctx := context.Background()

	githubClient := GetGithubClient(ctx)

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
		fmt.Println("Error creating gist:", err)
		return nil
	}

	return createdGist.HTMLURL
}
