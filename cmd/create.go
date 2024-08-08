package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use: "create",
	Run: execute,
}

func execute(cmd *cobra.Command, args []string) {
	fmt.Println("I AM CREATING?")
}
