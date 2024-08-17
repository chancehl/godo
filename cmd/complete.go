package cmd

import (
	"strconv"
	"time"

	"github.com/chancehl/godo/internal/clients/github"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	completeCmd.Flags().BoolVarP(&deleteGodoItem, "delete", "d", false, "Delete an item as you complete it")

	rootCmd.AddCommand(completeCmd)
}

var deleteGodoItem bool

var completeCmd = &cobra.Command{
	Use:  "complete [item]",
	Run:  executeComplete,
	Args: cobra.ExactArgs(1),
}

func executeComplete(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.PrintErrln("Error: Missing item argument")
		return
	}

	item := args[0]

	itemID, err := strconv.Atoi(item)
	if err != nil {
		cmd.PrintErrln("Could not convert item ID to integer")
		return
	}

	gistID, err := config.ReadGistIdFile()
	if err != nil {
		cmd.PrintErrln("Could not read gist id from config file: ", err)
		return
	}

	godos, err := github.GetGodos(gistID)
	if err != nil {
		cmd.PrintErrln("Could not fetch godods from github: ", err)
		return
	}

	newGodos := []model.GodoItem{}

	for index, godo := range godos {
		if index+1 == itemID {
			completedOn := time.Now().UTC().Format(time.RFC3339)

			godo.CompletedOn = completedOn
			godo.Status = "COMPLETE"

			if !deleteGodoItem {
				newGodos = append(newGodos, godo)
			}
		} else {
			newGodos = append(newGodos, godo)
		}
	}

	if err := github.UpdateGist(gistID, newGodos); err != nil {
		cmd.PrintErrln("Failed to update godos: ", err)
	}
}
