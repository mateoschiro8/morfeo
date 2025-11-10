package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	pdfFile string
)

var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "Genera el honeytoken de pdf",
	Run:   runsdfk,
}

func init() {
	pdfCmd.Flags().StringVar(&pdfFile, "file", "", "Ruta al archivo")
	rootCmd.AddCommand(pdfCmd)
}

func runsdfk(cmd *cobra.Command, args []string) {
	fmt.Println(pdfFile)
}
