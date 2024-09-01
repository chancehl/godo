package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/chancehl/godo/internal/clients/github"
	"github.com/chancehl/godo/internal/config"
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

	// Create a new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.DiscardEmptyColumns)

	headers := []string{"TASK. NO", "ID", "NAME", "NOTES", "STATUS", "CREATED ON", "COMPLETED ON"}
	header := strings.Join(headers, "\t")

	// Print table header
	fmt.Fprintln(w, header)

	for index, item := range items {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t\n", index+1, item.ID, item.Name, item.Notes, item.Status, item.CreatedOn, item.CompletedOn)
	}

	w.Flush()

	return nil
}
