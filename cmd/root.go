package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/chancehl/godo/internal/config"
	appContext "github.com/chancehl/godo/internal/context"
	"github.com/chancehl/godo/internal/services"
	"github.com/google/go-github/v50/github"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var applicationContext *appContext.ApplicationContext
var godoService *services.GodoService
var githubService *services.GithubService

var rootCmd = &cobra.Command{
	Use:   "godo",
	Short: "A simple command-line TODO editor",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		envFileData, _ := godotenv.Read(".env")

		accessToken := envFileData["GITHUB_ACCESS_TOKEN"]

		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
		tokenClient := oauth2.NewClient(ctx, tokenSource)

		githubClient := github.NewClient(tokenClient)

		gistID, err := config.ReadGistIdFile()
		if err != nil {
			return fmt.Errorf("could not read local gist id file (err=%s)", err)
		}

		applicationContext = &appContext.ApplicationContext{
			GistID: gistID,
		}

		githubService = &services.GithubService{
			Context:            &ctx,
			GithubClient:       githubClient,
			ApplicationContext: applicationContext,
		}

		godoService = &services.GodoService{
			GithubService:      githubService,
			ApplicationContext: applicationContext,
		}

		return nil
	},
}

func init() {
	config.InitializeDotDir()
}

func Execute() {
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
