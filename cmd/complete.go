package cmd

import (
	"strconv"
	"time"

	"github.com/chancehl/godo/internal/clients/github"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/chancehl/godo/internal/utils/cli"
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
		return cli.CmdError(cmd, "Could not convert item ID to integer", err)
	}

	gistID, err := config.ReadGistIdFile()
	if err != nil {
		return cli.CmdError(cmd, "Could not read gist ID from config file", err)
	}

	godos, err := github.GetGodos(gistID)
	if err != nil {
		return cli.CmdError(cmd, "Could not fetch godos from GitHub", err)
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
		return cli.CmdError(cmd, "Failed to update godos", err)
	}

	return nil
}
