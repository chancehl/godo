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

	return clearItems(clearAll, cmd, gistID, godos)
}

func clearItems(all bool, cmd *cobra.Command, gistID string, godos []model.GodoItem) error {
	if all {
		return clearAllItems(cmd, gistID, godos)
	}

	return clearCompletedItems(cmd, gistID, godos)
}

func clearAllItems(cmd *cobra.Command, gistID string, godos []model.GodoItem) error {
	confirmationPrompt := "This action will delete ALL godo items. This is permanent and cannot be undone."

	if !cli.ConfirmAction(confirmationPrompt) {
		fmt.Println("Clear all aborted by user")
		return nil
	}

	if err := github.UpdateGodos(gistID, []model.GodoItem{}); err != nil {
		return cli.CmdError(cmd, "Failed to update godos", err)
	}

	fmt.Printf("Removed %d godo items from list\n", len(godos))
	return nil
}

func clearCompletedItems(cmd *cobra.Command, gistID string, godos []model.GodoItem) error {
	var updatedGodos []model.GodoItem

	for _, godo := range godos {
		if godo.Status != "COMPLETE" {
			updatedGodos = append(updatedGodos, godo)
		}
	}

	if (len(godos) - len(updatedGodos)) == 0 {
		fmt.Println("No completed items to clear.")
		return nil
	}

	if !cli.ConfirmAction("This action will delete all completed godo items.") {
		return nil
	}

	if err := github.UpdateGodos(gistID, updatedGodos); err != nil {
		return cli.CmdError(cmd, "Failed to update godos", err)
	}

	fmt.Printf("Cleared all completed godo item(s) from list\n")
	return nil
}
