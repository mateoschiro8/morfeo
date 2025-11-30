package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	url_original = ""
)

var cssCmd = &cobra.Command{
	Use:   "css",
	Short: "Genera el honeytoken de css para paginas clonadas",
	Run:   runCss,
}

func init() {
	cssCmd.Flags().StringVar(&in, "in", "", "Archivo de CSS")
	cssCmd.Flags().StringVar(&out, "out", "", "Archivo de CSS modificado, de no proveerse nada se crea en la carpeta out con el mismo nombre")
	cssCmd.Flags().StringVar(&url_original, "dominio", "", "Dominio del sitio original")
	cssCmd.Flags().StringVar(&msg, "mensaje", "", "Mensaje que debe mostrar el servidor Canary")
	cssCmd.MarkFlagRequired("in")
	cssCmd.MarkFlagRequired("dominio")
	cssCmd.MarkFlagRequired("mensaje")

	rootCmd.AddCommand(cssCmd)
}

func runCss(cmd *cobra.Command, args []string) {
	if out == "" {
		out = "out/" + in
	}

	var id = CreateToken(msg, url_original)
	createCss(id)

}

func createCss(id string) {

	var ngrokurl = GetNgrokUrl()

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

	var cssContent = fmt.Sprintf(" \nbody {\n    background: url(%s) !important; \n}\n", ngrokurl+"/fondo/"+id)

	_, err = outFile.WriteString(cssContent)
	if err != nil {
		panic(fmt.Errorf("error agregando el css: %w", err))
	}

	fmt.Print("===================== Agregamos la linea siguiente linea =====================\n")
	fmt.Println(cssContent)
	fmt.Print("==============================================================================\n")

	fmt.Printf("Ahora debes cambiar el css original(%v) por el nuevo: %v\n", in, out)
}
