package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/chancehl/godo/internal/clients/github"
	"github.com/chancehl/godo/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use: "list",
	Run: executeList,
}

func executeList(cmd *cobra.Command, args []string) {
	gistID, err := config.ReadGistIdFile()
	if err != nil {
		cmd.PrintErrln("Failed to read gist id from local config: ", err)
	}

	items, err := github.ReadGist(gistID)
	if err != nil {
		cmd.PrintErrln("Failed to fetch godo items: ", err)
	}

	// Create a new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.DiscardEmptyColumns)

	// Print table header
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tCREATED ON\tCOMPLETED ON\t")

	// Print each client in a new row
	for _, item := range items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", item.ID, item.Name, item.Status, item.CreatedOn)
	}

	w.Flush()
}
