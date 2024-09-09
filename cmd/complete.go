package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/chancehl/godo/internal/model"
	"github.com/spf13/cobra"
)

const complete = "COMPLETE"

var deleteGodoItem bool

func init() {
	// register flags
	completeCmd.Flags().BoolVarP(&deleteGodoItem, "delete", "d", false, "Delete an item as you complete it")

	// register cmd
	rootCmd.AddCommand(completeCmd)
}

var completeCmd = &cobra.Command{
	Use:  "complete [item]",
	Args: cobra.ExactArgs(1),
	RunE: executeComplete, // Use RunE to handle errors better
}

func executeComplete(cmd *cobra.Command, args []string) error {
	itemID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("could not convert id to integer")
	}

	godos, err := godoService.GetGodos()
	if err != nil {
		return fmt.Errorf("could not fetch godos from GitHub (%s)", err)
	}

	if !doesIdExist(itemID, godos) {
		return fmt.Errorf("item id %d not found", itemID)
	}

	if isAlreadyComplete(itemID, godos) {
		return fmt.Errorf("item %d is already complete", itemID)
	}

	var updatedGodos []model.GodoItem

	for _, godo := range godos {
		if godo.ID == itemID {
			godo.Status = complete
			godo.CompletedOn = time.Now().UTC().Format(time.RFC3339)

			if deleteGodoItem {
				continue
			}
		}
		updatedGodos = append(updatedGodos, godo)
	}

	if err := godoService.UpdateGodos(updatedGodos); err != nil {
		return fmt.Errorf("could not update godos (%s)", err)
	}

	fmt.Printf("✅ Completed godo item %d\n", itemID)
	return nil
}

func doesIdExist(id int, items []model.GodoItem) bool {
	for _, item := range items {
		if item.ID == id {
			return true
		}
	}
	return false
}

func isAlreadyComplete(id int, items []model.GodoItem) bool {
	for _, item := range items {
		if item.ID == id && item.Status == complete {
			return true
		}
	}
	return false
}
