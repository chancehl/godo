package cmd

import (
	"fmt"

	"github.com/chancehl/godo/internal/clients"
	"github.com/chancehl/godo/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:  "create [item]",
	Run:  execute,
	Args: cobra.ExactArgs(1),
}

func execute(cmd *cobra.Command, args []string) {
	item := args[0]

	items := []model.GodoItem{
		{Name: "TODO1"},
		{Name: "TODO2"},
		{Name: item},
	}

	gist := clients.CreateGist(items)

	fmt.Println("gist created at", gist)
}
