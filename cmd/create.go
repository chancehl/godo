package cmd

import (
	"time"

	"github.com/chancehl/godo/internal/clients"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:  "create [item]",
	Run:  executeCreate,
	Args: cobra.ExactArgs(1),
}

func executeCreate(cmd *cobra.Command, args []string) {
	item := args[0]

	id, _ := config.ReadGistIdFile()
	godos, _ := clients.ReadGist(id)

	newGodo := model.GodoItem{
		Name:      item,
		Status:    "TODO",
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
	}
	godos = append(godos, newGodo)

	clients.UpdateGist(id, godos)
}
