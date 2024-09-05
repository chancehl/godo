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
		return fmt.Errorf("invalid item id: %d", itemID)
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

	return github.UpdateGodos(gistID, updatedGodos)
}

func checkIfIdExists(id int, items []model.GodoItem) bool {
	for _, item := range items {
		if item.ID == id {
			return true
		}
	}
	return false
}
