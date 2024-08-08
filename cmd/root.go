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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
