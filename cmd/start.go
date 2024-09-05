package cmd

import (
	"fmt"
	"strconv"

	"github.com/chancehl/godo/internal/clients/github"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startcmd)
}

var startcmd = &cobra.Command{
	Use:  "start [item]",
	Args: cobra.ExactArgs(1),
	RunE: executeStart,
}

func executeStart(cmd *cobra.Command, args []string) error {
	itemID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Could not convert item ID to integer", err)
	}

	gistID, err := config.ReadGistIdFile()
	if err != nil {
		return fmt.Errorf("Could not read gist ID from config file", err)
	}

	godos, err := github.GetGodos(gistID)
	if err != nil {
		return fmt.Errorf("Could not fetch godos from GitHub", err)
	}

	var updatedGodos []model.GodoItem

	for index, godo := range godos {
		if index+1 == itemID {
			godo.Status = "IN_PROGRESS"

			if deleteGodoItem {
				continue
			}
		}
		updatedGodos = append(updatedGodos, godo)
	}

	if err := github.UpdateGodos(gistID, updatedGodos); err != nil {
		return fmt.Errorf("Failed to update godos", err)
	}

	return nil
}
