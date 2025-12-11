package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	dominio = ""
)

var cssCmd = &cobra.Command{
	Use:   "css",
	Short: "Genera el honeytoken de css para paginas clonadas",
	Run:   runCss,
}

func init() {
	cssCmd.Flags().StringVar(&msg, "msg", "", "Mensaje que debe mostrar el servidor Canary")
	cssCmd.Flags().StringVar(&chat, "chat", "", "Chat ID al cual enviar mensaje al activarse")
	cssCmd.Flags().StringVar(&in, "in", "", "Archivo de CSS")
	cssCmd.Flags().StringVar(&out, "out", "", "Archivo de CSS modificado, de no proveerse nada se crea con el mismo nombre")
	cssCmd.Flags().StringVar(&dominio, "dominio", "", "Dominio del sitio original")
	cssCmd.MarkFlagRequired("in")
	cssCmd.MarkFlagRequired("dominio")
	rootCmd.AddCommand(cssCmd)
}

func runCss(cmd *cobra.Command, args []string) {
	
	formatIn()

	formatOut()

	var id = CreateToken(msg, dominio, chat)
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
