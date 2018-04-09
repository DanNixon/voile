package tui

import (
	"io/ioutil"
	"os"
	"os/exec"
)

func FindEditor() string {
	for _, envVar := range []string{"VISUAL", "EDITOR"} {
		// Find a set environment variable
		bin := os.Getenv(envVar)
		if len(bin) == 0 {
			continue
		}

		// Find a valid executable
		path, err := exec.LookPath(bin)
		if err != nil {
			continue
		}

		return path
	}

	return "vi"
}

func EditFile(filename string) error {
	cmd := exec.Command(FindEditor(), filename)

	// Connect console IO
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start process
	err := cmd.Start()
	if err != nil {
		return err
	}

	// Wait for termination
	return cmd.Wait()
}

func EditText(text string) (string, error) {
	// Create temporary file
	txtFile, err := ioutil.TempFile("", "voile")
	if err != nil {
		return "", err
	}

	// Ensure file is deleted
	defer os.Remove(txtFile.Name())

	// Write test to temp file
	if _, err = txtFile.WriteString(text); err != nil {
		return "", err
	}
	if err = txtFile.Close(); err != nil {
		return "", err
	}

	err = EditFile(txtFile.Name())
	if err != nil {
		return "", err
	}

	// Read text back from temp file
	var b []byte
	b, err = ioutil.ReadFile(txtFile.Name())
	result := string(b)

	return result, err
}
