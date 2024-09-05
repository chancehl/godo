package cmd

import (
	"fmt"
	"log"

	"github.com/chancehl/godo/internal/clients/github"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/chancehl/godo/internal/utils/cli"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:  "init",
	RunE: executeInit,
}

func executeInit(cmd *cobra.Command, args []string) error {
	configExists, _ := config.CheckIfGistIdFileExists()

	if configExists {
		if !handleExistingGist() {
			fmt.Println("Initialization aborted by the user.")
			return nil
		}
	}

	id, gistURL, err := github.CreateGist([]model.GodoItem{})
	if err != nil {
		return fmt.Errorf("error creating gist (%v)", err)
	}

	gistIDFilePath, err := config.WriteGistIdFile(id)
	if err != nil {
		log.Fatalf("Error writing gist ID file: %v", err)
	}

	displayInitializationDetails(id, gistURL, gistIDFilePath)

	fmt.Println("Successfully initialized godo")

	return nil
}

func handleExistingGist() bool {
	existingID, err := config.ReadGistIdFile()
	if err != nil {
		log.Fatalf("Error reading existing gist ID: %v", err)
	}

	prompt := fmt.Sprintf(`
You have already initialized godo with the following gist ID:

%s

Are you sure you want to reinitialize?`, existingID)

	return cli.ConfirmAction(prompt)
}

func displayInitializationDetails(id, gistURL, gistIDFilePath string) {
	message := fmt.Sprintf(`
GitHub Details:

* Gist ID: %s
* Gist URL: %s

Local Details:

* Gist ID File Path: %s
`, id, gistURL, gistIDFilePath)

	fmt.Println(message)
}
