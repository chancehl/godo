package cmd

import (
	"fmt"
	"time"

	"github.com/chancehl/godo/internal/clients/github"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/chancehl/godo/internal/utils/cli"
	"github.com/lithammer/shortuuid/v4"
	"github.com/spf13/cobra"
)

func init() {
	createCmd.PersistentFlags().StringP("notes", "n", "", "The notes to associate with the godo item")
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:  "create [item]",
	RunE: executeCreate,
	Args: cobra.ExactArgs(1),
}

func executeCreate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing item argument")
	}

	item := args[0]

	gistID, err := config.ReadGistIdFile()
	if err != nil {
		return cli.CmdError(cmd, "Error reading gist id: ", err)
	}

	godos, err := github.GetGodos(gistID)
	if err != nil {
		return cli.CmdError(cmd, "Error reading gist: ", err)
	}

	notes, _ := cmd.Flags().GetString("notes")
	newGodo := model.GodoItem{
		ID:        shortuuid.New()[:12],
		Name:      item,
		Status:    "TODO",
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		Notes:     notes,
	}

	godos = append(godos, newGodo)

	if err := github.UpdateGodos(gistID, godos); err != nil {
		return cli.CmdError(cmd, "Error updating gist: ", err)
	}

	fmt.Println("Created godo item:", item)
	return nil
}
