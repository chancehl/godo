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
		return fmt.Errorf("could not read gist ID from config file (%s)", err)
	}

	godos, err := github.GetGodos(gistID)
	if err != nil {
		return fmt.Errorf("could not fetch godos from GitHub (%s)", err)
	}

	return clearItems(clearAll, gistID, godos)
}

func clearItems(all bool, gistID string, godos []model.GodoItem) error {
	if all {
		return clearAllItems(gistID, godos)
	}

	return clearCompletedItems(gistID, godos)
}

func clearAllItems(gistID string, godos []model.GodoItem) error {
	confirmationPrompt := "This action will delete ALL godo items. This is permanent and cannot be undone."

	if !cli.ConfirmAction(confirmationPrompt) {
		fmt.Println("clear all aborted by user")
		return nil
	}

	if err := github.UpdateGodos(gistID, []model.GodoItem{}); err != nil {
		return fmt.Errorf("failed to update godos (%s)", err)
	}

	fmt.Printf("🧹 removed %d godo items from list\n", len(godos))
	return nil
}

func clearCompletedItems(gistID string, godos []model.GodoItem) error {
	var updatedGodos []model.GodoItem

	for _, godo := range godos {
		if godo.Status != "COMPLETE" {
			updatedGodos = append(updatedGodos, godo)
		}
	}

	updated := len(godos) - len(updatedGodos)
	if updated == 0 {
		fmt.Println("no completed items to clear")
		return nil
	}

	if !cli.ConfirmAction("This action will delete all completed godo items.") {
		return nil
	}

	if err := github.UpdateGodos(gistID, updatedGodos); err != nil {
		return fmt.Errorf("failed to update godos (%s)", err)
	}

	fmt.Printf("🧹 cleared %d completed godo item(s) from list\n", updated)
	return nil
}
