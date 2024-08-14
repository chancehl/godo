package cmd

import (
	"fmt"

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

	items := []model.GodoItem{
		{Name: "TODO1"},
		{Name: "TODO2"},
		{Name: item},
	}

	fmt.Println(items)
}
