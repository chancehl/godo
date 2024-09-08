package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/chancehl/godo/internal/model"
	"github.com/spf13/cobra"
)

const DEFAULT_SORT = "ID"

func init() {
	listCommand.PersistentFlags().StringP("sort", "s", DEFAULT_SORT, "The property you want to sort by")
	rootCmd.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:  "list",
	RunE: executeList,
}

func executeList(cmd *cobra.Command, args []string) error {
	items, err := godoService.GetGodos()
	if err != nil {
		return fmt.Errorf("failed to fetch godo items (%s)", err)
	}

	sort, _ := cmd.Flags().GetString("sort")
	sortItems(items, sort)

	// Create a new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.DiscardEmptyColumns)

	headers := []string{"ID", "NAME", "NOTES", "STATUS", "CREATED ON", "COMPLETED ON"}
	header := strings.Join(headers, "\t")

	// Print table header
	fmt.Fprintln(w, header)

	for _, item := range items {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t\n", item.ID, item.Name, item.Notes, item.Status, item.CreatedOn, item.CompletedOn)
	}

	w.Flush()

	return nil
}

func sortItems(items []model.GodoItem, field string) {
	switch strings.ToUpper(field) {
	case "STATUS":
		sort.Slice(items, func(i, j int) bool {
			return items[i].Status < items[j].Status
		})
	case "ID":
	default:
		sort.Slice(items, func(i, j int) bool {
			return items[i].ID < items[j].ID
		})
	}

}
