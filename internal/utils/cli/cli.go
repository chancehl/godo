package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ConfirmAction(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("⚠️  %s (y/n): ", prompt)

	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error reading input:", err)
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
