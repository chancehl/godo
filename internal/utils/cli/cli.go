package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func CmdErrorS(cmd *cobra.Command, msg string) error {
	cmd.PrintErrln(msg)
	return errors.New(msg)
}

func CmdError(cmd *cobra.Command, msg string, err error) error {
	cmd.PrintErrln(msg+":", err)
	return err
}

func ConfirmAction(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s Would you like to continue? (y/n): ", prompt)

	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
