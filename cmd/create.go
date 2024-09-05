package cmd

import (
	"fmt"
	"time"

	"github.com/chancehl/godo/internal/clients/github"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/chancehl/godo/internal/utils/cli"
	"github.com/spf13/cobra"
)

func init() {
	createCmd.PersistentFlags().StringP("notes", "n", "", "The notes to associate with the godo item")
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:  "create [item]",
	RunE: createGodo,
	Args: cobra.ExactArgs(1),
}

func createGodo(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing item argument")
	}

	item := args[0]

	gistID, err := config.ReadGistIdFile()
	if err != nil {
		return fmt.Errorf("error reading gist id from file (%s)", err)
	}

	godos, err := github.GetGodos(gistID)
	if err != nil {
		return fmt.Errorf("error reading gist file from github (%s)", err)
	}

	if exists, existing := checkIfAlreadyExists(item, godos); exists && existing.Status != "COMPLETED" {
		prompt := fmt.Sprintf("Looks like you already have an incomplete godo item with the name \"%s\".", existing.Name)

		if !cli.ConfirmAction(prompt) {
			return nil
		}
	}

	notes, _ := cmd.Flags().GetString("notes")

	id := generateID(godos)

	newGodo := model.GodoItem{
		ID:        id,
		Name:      item,
		Status:    "TODO",
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		Notes:     notes,
	}

	godos = append(godos, newGodo)

	if err := github.UpdateGodos(gistID, godos); err != nil {
		return fmt.Errorf("error updating gist (%s)", err)
	}

	fmt.Printf("Added \"%s\" to your godo list.\n", item)
	return nil
}

func checkIfAlreadyExists(item string, items []model.GodoItem) (bool, *model.GodoItem) {
	for _, existing := range items {
		if existing.Name == item {
			return true, &existing
		}
	}
	return false, nil
}

func generateID(items []model.GodoItem) int {
	if len(items) == 0 {
		return 1
	}

	max := 0

	for _, item := range items {
		if item.ID > max {
			max = item.ID
		}
	}

	return max + 1
}
