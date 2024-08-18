package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chancehl/godo/internal/clients/github"
	"github.com/chancehl/godo/internal/config"
	"github.com/chancehl/godo/internal/model"
	"github.com/chancehl/godo/internal/utils/cli"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:  "list",
	RunE: executeList,
}

func executeList(cmd *cobra.Command, args []string) error {
	gistID, err := config.ReadGistIdFile()
	if err != nil {
		return cli.CmdError(cmd, "Failed to read gist id from local config: ", err)
	}

	items, err := github.GetGodos(gistID)
	if err != nil {
		return cli.CmdError(cmd, "Failed to fetch godo items: ", err)
	}

	printList(items)

	return nil
}

func printList(items []model.GodoItem) {
	// Create a new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.DiscardEmptyColumns)

	// Print table header
	fmt.Fprintln(w, "TASK. NO\tID\tNAME\tSTATUS\tCREATED ON\tCOMPLETED ON\t")

	for index, item := range items {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t\n", index+1, item.ID, item.Name, item.Status, item.CreatedOn, item.CompletedOn)
	}

	w.Flush()
}
