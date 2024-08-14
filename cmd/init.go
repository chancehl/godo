package cmd

import (
	"fmt"

	"github.com/chancehl/godo/internal/clients"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/chancehl/godo/internal/utils"
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
	exists := config.CheckIfGistIdFileExists()

	if exists {
		existingId, _ := config.ReadGistIdFile()

		prompt := fmt.Sprintf(`
You have already initialized godo in the past and have a stored gist id of:

%s

Are you sure you want to proceed?
`, existingId)

		confirmation := utils.ConfirmAction(prompt)

		if !confirmation {
			fmt.Println("Did not create gist_file_id file. Aborting.")
			return
		}
	}

	id, gist, err := clients.CreateGist(make([]model.GodoItem, 0))

	if err != nil {
		fmt.Printf("Encountered error when creating gist: %s", err)
	}

	gist_id_file, err := config.WriteGistIdFile(id)

	if err != nil {
		fmt.Printf("Failed to write gist id file: %s", err)
	}

	message := fmt.Sprintf(`
Github:

* %s
* %s

Local:

* %s
`, id, gist, gist_id_file)

	fmt.Println(message)
	fmt.Println("\033[32mSuccessfully initialized godo!\033[0m")
}
