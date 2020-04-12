package tui

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

type MultiChoiceOption struct {
	Key  string
	Desc string
}

func Option(prompt string, options []MultiChoiceOption) (string, error) {
	// Prompt
	fmt.Print(prompt)
	for _, a := range options {
		fmt.Print(", (", a.Key, ") ", a.Desc)
	}
	fmt.Print(" ")

	// Read first line from console
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}

	// Check for response
	for _, a := range options {
		if bytes.ContainsAny(line, a.Key) {
			return a.Key, err
		}
	}

	return "", err
}
