package cli

import "github.com/spf13/cobra"

func CmdError(cmd *cobra.Command, msg string, err error) error {
	cmd.PrintErrln(msg+":", err)
	return err
}
