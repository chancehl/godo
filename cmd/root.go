package cmd

import (
	"fmt"
	"os"

	"github.com/chancehl/godo/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godo",
	Short: "A simple command-line TODO editor",
}

func init() {
	config.InitializeDotDir()
}

func Execute() {
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
