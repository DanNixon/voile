package tui

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
)

func Confirm(prompt string) (bool, error) {
	// Prompt
	fmt.Print(prompt, " (y/N) ")

	result := false

	// Read first line from console
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		return result, err
	}

	// Check for response
	positiveResponse := bytes.ContainsAny(line, "yY")
	negativeResponse := bytes.ContainsAny(line, "nN")
	result = positiveResponse && !negativeResponse

	// Ensure valid and set error if not
	if positiveResponse && negativeResponse {
		err = errors.New("Ambiguous response")
	}

	return result, err
}
