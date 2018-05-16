package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Invoke Git",
	Long:  `Invoke Git in the directory containing your bookmarks file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !IsBookmarksFileInGitRepository() {
			fmt.Println("Bookmarks file is not stored in a Git directory")
			return
		}

		repoDirArgs := []string{"-C", GetBookmarksFileParentDirectory()}
		gitArgs := append(repoDirArgs, args...)

		gitCmd := exec.Command("git", gitArgs...)

		// Connect console IO
		gitCmd.Stdin = os.Stdin
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		// Run process
		err := gitCmd.Run()
		CheckError(err)
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)
}
