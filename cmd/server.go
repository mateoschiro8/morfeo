package cmd

import (
	"fmt"

	"github.com/mateoschiro8/morfeo/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ngrokurl: %v\n", ngrokurl)
		server.StartServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
