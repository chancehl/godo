package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chancehl/godo/internal/clients"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use: "init",
	Run: executeInit,
}

func executeInit(cmd *cobra.Command, args []string) {
	configExists, _ := config.CheckIfGistIdFileExists()

	if configExists {
		if !handleExistingGist() {
			color.Yellow("Initialization aborted by the user.")
			return
		}
	}

	id, gistURL, err := clients.CreateGist([]model.GodoItem{})
	if err != nil {
		log.Fatalf("Error creating gist: %v", err)
	}

	gistIDFilePath, err := config.WriteGistIdFile(id)
	if err != nil {
		log.Fatalf("Error writing gist ID file: %v", err)
	}

	displayInitializationDetails(id, gistURL, gistIDFilePath)
	color.Green("Successfully initialized godo")
}

func handleExistingGist() bool {
	existingID, err := config.ReadGistIdFile()
	if err != nil {
		log.Fatalf("Error reading existing gist ID: %v", err)
	}

	prompt := fmt.Sprintf(`
You have already initialized godo with the following gist ID:

%s

Are you sure you want to reinitialize? (y/n): `, existingID)

	return confirmAction(prompt)
}

func confirmAction(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
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
