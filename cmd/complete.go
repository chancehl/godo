package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/chancehl/godo/internal/clients/github"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/spf13/cobra"
)

var deleteGodoItem bool

func init() {
	completeCmd.Flags().BoolVarP(&deleteGodoItem, "delete", "d", false, "Delete an item as you complete it")
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

	gistID, err := config.ReadGistIdFile()
	if err != nil {
		return fmt.Errorf("could not read gist id from config file (%s)", err)
	}

	godos, err := github.GetGodos(gistID)
	if err != nil {
		return fmt.Errorf("could not fetch godos from GitHub (%s)", err)
	}

	idExists := checkIfIdExists(itemID, godos)
	if !idExists {
		return fmt.Errorf("item id %d not found", itemID)
	}

	alreadyComplete := checkIfItemIsAlreadyComplete(itemID, godos)
	if alreadyComplete {
		return fmt.Errorf("item %d is already complete", itemID)
	}

	var updatedGodos []model.GodoItem

	for index, godo := range godos {
		if index+1 == itemID {
			godo.CompletedOn = time.Now().UTC().Format(time.RFC3339)
			godo.Status = "COMPLETE"

			if deleteGodoItem {
				continue
			}
		}
		updatedGodos = append(updatedGodos, godo)
	}

	if err := github.UpdateGodos(gistID, updatedGodos); err != nil {
		return fmt.Errorf("could not update godos (%s)", err)
	}

	fmt.Printf("✅ Completed godo item %d\n", itemID)
	return nil
}

func checkIfIdExists(id int, items []model.GodoItem) bool {
	for _, item := range items {
		if item.ID == id {
			return true
		}
	}
	return false
}

func checkIfItemIsAlreadyComplete(id int, items []model.GodoItem) bool {
	for _, item := range items {
		if item.ID == id && item.Status == "COMPLETE" {
			return true
		}
	}
	return false
}
