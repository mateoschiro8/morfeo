package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var cssCmd = &cobra.Command{
	Use:   "css",
	Short: "Genera el honeytoken de css para paginas clonadas",
	Run:   runCss,
}

func init() {
	cssCmd.Flags().StringVar(&in, "in", "", "Archivo CSS")
	cssCmd.Flags().StringVar(&out, "out", "", "Nombre de archivo CSS modificado")
	cssCmd.Flags().StringVar(&extra, "dominio", "", "Dominio del sitio original")
	cssCmd.MarkFlagRequired("in")
	cssCmd.MarkFlagRequired("dominio")
	rootCmd.AddCommand(cssCmd)
}

func runCss(cmd *cobra.Command, args []string) {
	if out == "" {
		out = "new_" + in
	}

	var id = CreateToken(msg, extra, chat)
	createCss(id)

}

func createCss(id string) {

	inFile, err := os.Open(in)
	if err != nil {
		panic(fmt.Errorf("error abriendo input: %w", err))
	}
	defer inFile.Close()

	outFile, err := os.Create(out)
	if err != nil {
		panic(fmt.Errorf("error creando output: %w", err))
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, inFile)
	if err != nil {
		panic(fmt.Errorf("error clonando input: %w", err))
	}

	var cssContent = fmt.Sprintf(" \nbody {\n    background: url(%s) !important; \n}\n", serverURL+"/fondo/"+id)

	_, err = outFile.WriteString(cssContent)
	if err != nil {
		panic(fmt.Errorf("error agregando el css: %w", err))
	}

}
