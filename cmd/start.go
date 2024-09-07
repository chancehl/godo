package cmd

import (
	"fmt"
	"strconv"

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
		return fmt.Errorf("could not convert item ID to integer (%s)", err)
	}

	godos, err := appCtx.GodoService.GetGodos()
	if err != nil {
		return fmt.Errorf("could not fetch godos from GitHub (%s)", err)
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

	if err := appCtx.GodoService.UpdateGodos(updatedGodos); err != nil {
		return fmt.Errorf("failed to update godos (%s)", err)
	}

	return nil
}
