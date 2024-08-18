package cmd

import (
	"fmt"

	"github.com/chancehl/godo/internal/clients/github"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/chancehl/godo/internal/utils/cli"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clearCmd)
}

var clearCmd = &cobra.Command{
	Use:  "clear",
	RunE: executeClear,
}

func executeClear(cmd *cobra.Command, args []string) error {
	gistID, err := config.ReadGistIdFile()
	if err != nil {
		return cli.CmdError(cmd, "Could not read gist ID from config file", err)
	}

	godos, err := github.GetGodos(gistID)
	if err != nil {
		return cli.CmdError(cmd, "Could not fetch godos from GitHub", err)
	}

	var updatedGodos []model.GodoItem

	for _, godo := range godos {
		if godo.Status != "COMPLETE" {
			updatedGodos = append(updatedGodos, godo)
		}
	}

	toBeDeleted := len(godos) - len(updatedGodos)

	if toBeDeleted == 0 {
		color.Yellow("No items were deleted")
		return nil
	}

	confirmationPrompt := fmt.Sprintf("This action will delete %d godo items. Continue? (y/n): ", toBeDeleted)

	if !cli.ConfirmAction(confirmationPrompt) {
		color.Yellow("Clear aborted by the user.")
		return nil
	}

	if err := github.UpdateGodos(gistID, updatedGodos); err != nil {
		return cli.CmdError(cmd, "Failed to update godos", err)
	}

	color.Green(fmt.Sprintf("Removed %d godo items from list", toBeDeleted))
	return nil
}
