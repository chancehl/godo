package cmd

import (
	"time"

	"github.com/chancehl/godo/internal/clients"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:  "create [item]",
	Run:  executeCreate,
	Args: cobra.ExactArgs(1),
}

func executeCreate(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.PrintErrln("Error: Missing item argument")
		return
	}

	item := args[0]

	id, err := config.ReadGistIdFile()
	if err != nil {
		cmd.PrintErrln("Error reading Gist ID:", err)
		return
	}

	godos, err := clients.ReadGist(id)
	if err != nil {
		cmd.PrintErrln("Error reading Gist:", err)
		return
	}

	newGodo := model.GodoItem{
		Name:      item,
		Status:    "TODO",
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
	}

	godos = append(godos, newGodo)

	if err := clients.UpdateGist(id, godos); err != nil {
		cmd.PrintErrln("Error updating Gist:", err)
	} else {
		color.Green("Successfully created godo item")
	}
}
