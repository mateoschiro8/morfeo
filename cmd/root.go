package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use: "morfeo",
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	_ = rootCmd.Execute()
}
