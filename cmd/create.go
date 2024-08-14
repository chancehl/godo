package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v50/github"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

type GodoItem struct {
	Name string
}

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use: "create",
	Run: execute,
}

func execute(cmd *cobra.Command, args []string) {
	fmt.Println("I AM CREATING?")

	ctx := context.Background()

	envFileData, _ := godotenv.Read(".env")

	accessToken := envFileData["GITHUB_ACCESS_TOKEN"]

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tokenClient := oauth2.NewClient(ctx, tokenSource)

	githubClient := github.NewClient(tokenClient)

	gistContent, err := json.Marshal(GodoItem{Name: "TODO?"})

	if err != nil {
		fmt.Println("Could not create JSON from godo item", err)
		return
	}

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
		return
	}

	fmt.Println("Gist created at URL:", *createdGist.HTMLURL)
}
