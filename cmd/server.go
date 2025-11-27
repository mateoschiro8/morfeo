package cmd

import (
	"github.com/mateoschiro8/morfeo/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Printf("ngrokurl: %v\n", GetNgrokUrl())
		server.StartServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

/*
POST /qrs
{
	id= del aula 5
}


server contesta > "/qrs/{elksjrfgvhn}"
*/
