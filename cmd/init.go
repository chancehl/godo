package cmd

import (
	"fmt"

	"github.com/chancehl/godo/internal/clients/github"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/chancehl/godo/internal/utils/cli"
	"github.com/spf13/cobra"
)

type GodoGistData struct {
	id  string
	url string
}

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize godo with a new or existing gist",
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	if err := handleExistingConfig(); err != nil {
		return err
	}

	existingGistData, err := findExistingGodoGist()
	if err != nil {
		return err
	}

	return createOrUpdateGistFile(existingGistData)
}

func createOrUpdateGistFile(existing *GodoGistData) error {
	var gistID string
	var gistURL string

	if existing != nil && confirmUserWantsToUseExistingGistFile(existing.id) {
		gistID = existing.id
		gistURL = existing.url
	} else {
		newGistID, newGistURL, err := github.CreateGist([]model.GodoItem{})
		if err != nil {
			return fmt.Errorf("failed to create gist: %w", err)
		}

		gistID = newGistID
		gistURL = newGistURL
	}

	gistIDFilePath, err := config.WriteGistIdFile(gistID)
	if err != nil {
		return fmt.Errorf("failed to write gist ID file: %w", err)
	}

	displayInitDetails(gistID, gistURL, gistIDFilePath)

	fmt.Println("⚡ successfully initialized godo")
	return nil
}

func handleExistingConfig() error {
	configExists, err := config.CheckIfGistIdFileExists()
	if err != nil {
		return fmt.Errorf("error checking gist ID file: %w", err)
	}

	if !configExists {
		return nil
	}

	existingGistFileID, err := config.ReadGistIdFile()
	if err != nil {
		return fmt.Errorf("error reading existing gist ID: %w", err)
	}

	if !confirmOverwrite(existingGistFileID) {
		fmt.Println("Initialization aborted by the user.")
		return fmt.Errorf("user canceled")
	}

	return nil
}

func findExistingGodoGist() (*GodoGistData, error) {
	gists, _, err := github.GetGists()
	if err != nil {
		return nil, fmt.Errorf("error fetching gists: %w", err)
	}

	for _, gist := range gists {
		if gist.Description != nil && *gist.Description == "godo" {
			return &GodoGistData{id: *gist.ID, url: gist.GetHTMLURL()}, nil
		}
	}

	return nil, nil
}

func confirmUserWantsToUseExistingGistFile(id string) bool {
	prompt := fmt.Sprintf(`When initializing godo found a godo gist in your Github account with the following ID:

%s

Would you like to use this existing gist file?`, id)
	return cli.ConfirmAction(prompt)

}

func confirmOverwrite(id string) bool {
	prompt := fmt.Sprintf(`You have already initialized godo with the following gist ID:

%s

Are you sure you want to reinitialize?`, id)
	return cli.ConfirmAction(prompt)
}

func displayInitDetails(id, gistURL, gistIDFilePath string) {
	fmt.Printf(`
GitHub Details:

* Gist ID: %s
* Gist URL: %s

Local Details:

* Gist ID File Path: %s
`, id, gistURL, gistIDFilePath)
}
