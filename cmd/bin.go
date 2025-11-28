package cmd

import (
	"github.com/spf13/cobra"
)

var binCmd = &cobra.Command{
	Use:   "bin",
	Short: "Genera un honeytoken a partir de un binario",
	Run:   generateBinaryWrapper,
}

func init() {
	binCmd.Flags().StringVar(&in, "in", "", "Path al binario a wrappear")
	binCmd.MarkFlagRequired("in")
	rootCmd.AddCommand(binCmd)
}

func generateBinaryWrapper(cmd *cobra.Command, args []string) {

}
