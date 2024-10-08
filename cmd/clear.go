package cmd

import (
	"fmt"

	"github.com/chancehl/godo/internal/model"
	"github.com/chancehl/godo/internal/utils/cli"
	"github.com/spf13/cobra"
)

var clearAll bool

func init() {

	// register flag
	clearCmd.Flags().BoolVarP(&clearAll, "all", "a", false, "Clear all items (instead of just completed ones)")

	// register cmd
	rootCmd.AddCommand(clearCmd)
}

var clearCmd = &cobra.Command{
	Use:  "clear",
	RunE: executeClear,
}

func executeClear(cmd *cobra.Command, args []string) error {
	godos, err := godoService.GetGodos()
	if err != nil {
		return fmt.Errorf("could not fetch godos from GitHub (%s)", err)
	}
	return clearItems(clearAll, godos)
}

func clearItems(all bool, godos []model.GodoItem) error {
	if all {
		return clearAllItems(godos)
	}

	return clearCompletedItems(godos)
}

func clearAllItems(godos []model.GodoItem) error {
	if !cli.ConfirmAction(fmt.Sprintf("This action will delete ALL %d items in your godo list. This is permanent and cannot be undone. Continue?", len(godos))) {
		fmt.Println("clear all aborted by user")
		return nil
	}

	if err := godoService.UpdateGodos([]model.GodoItem{}); err != nil {
		return fmt.Errorf("failed to update godos (%s)", err)
	}

	fmt.Printf("🧹 removed %d godo items from list\n", len(godos))
	return nil
}

func clearCompletedItems(godos []model.GodoItem) error {
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

	if !cli.ConfirmAction("This action will delete all completed godo items. Continue?") {
		return nil
	}

	if err := godoService.UpdateGodos(updatedGodos); err != nil {
		return fmt.Errorf("failed to update godos (%s)", err)
	}

	fmt.Printf("🧹 cleared %d completed godo item(s) from list\n", updated)
	return nil
}
