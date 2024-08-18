package cmd

import (
	"fmt"

	"github.com/chancehl/godo/internal/clients/github"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/chancehl/godo/internal/utils/cli"
	"github.com/spf13/cobra"
)

var clearAll bool

func init() {
	clearCmd.Flags().BoolVarP(&clearAll, "all", "a", false, "Clear all items (instead of just completed ones)")
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

	if clearAll {
		return handleClearAllItems(cmd, gistID, godos)
	} else {
		return handleClearCompletedItems(cmd, gistID, godos)
	}
}

func handleClearAllItems(cmd *cobra.Command, gistID string, godos []model.GodoItem) error {
	confirmationPrompt := "This action will delete ALL godo items. Continue? (y/n): "

	if !cli.ConfirmAction(confirmationPrompt) {
		fmt.Println("Clear all aborted by user")
		return nil
	}

	if err := github.UpdateGodos(gistID, []model.GodoItem{}); err != nil {
		return cli.CmdError(cmd, "Failed to update godos", err)
	}

	fmt.Printf("Removed %d godo items from list", len(godos))
	return nil
}

func handleClearCompletedItems(cmd *cobra.Command, gistID string, godos []model.GodoItem) error {
	var updatedGodos []model.GodoItem

	for _, godo := range godos {
		if godo.Status != "COMPLETE" {
			updatedGodos = append(updatedGodos, godo)
		}
	}

	toBeDeleted := len(godos) - len(updatedGodos)

	if toBeDeleted == 0 {
		fmt.Println("No items were deleted")
		return nil
	}

	confirmationPrompt := fmt.Sprintf("This action will delete %d godo items. Continue? (y/n): ", toBeDeleted)

	if !cli.ConfirmAction(confirmationPrompt) {
		fmt.Println("Clear aborted by the user.")
		return nil
	}

	if err := github.UpdateGodos(gistID, updatedGodos); err != nil {
		return cli.CmdError(cmd, "Failed to update godos", err)
	}

	fmt.Printf("Removed %d godo item(s) from list\n", toBeDeleted)
	return nil
}
