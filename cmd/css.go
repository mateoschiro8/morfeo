package cmd

import (
	"fmt"
	"os"
	"io"
	"github.com/spf13/cobra"
)

var (
	host = "http://localhost:8000/"
)

var cssCmd = &cobra.Command{
	Use:   "css",
	Short: "Genera el honeytoken de css para paginas clonadas",
	Run:   crearCss,
}

func init() {
	cssCmd.Flags().StringVar(&inputPath, "in", "", "Archivo de CSS")
	cssCmd.Flags().StringVar(&outputPath, "out", "", "Archivo de CSS modificado")
	cssCmd.MarkFlagRequired("in")
	cssCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(cssCmd)
}

func crearCss(cmd *cobra.Command, args []string) {
	
	inFile, err := os.Open(inputPath)
	if err != nil {
		panic(fmt.Errorf("error abriendo input: %w", err))
	}
	defer inFile.Close()
	
	outFile, err := os.Create(outputPath)
	if err != nil {
		panic(fmt.Errorf("error creando output: %w", err))
	}
	defer outFile.Close()
	
	_, err = io.Copy(outFile, inFile)
	if err != nil {
		panic(fmt.Errorf("error clonando input: %w", err))
	}
	
	var cssContent = fmt.Sprintf(" \nbody {\n    background: url(%s) !important; \n}\n", host+"fondo")
	fmt.Println(cssContent)
	
	_, err = outFile.WriteString(cssContent)
	if err!= nil{
		panic(fmt.Errorf("error agregando el css: %w", err))
	}

}