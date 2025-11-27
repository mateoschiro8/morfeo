package cmd

import (
	"fmt"
	"os"
	"io"
	"github.com/spf13/cobra"
)

var (
	url_original = ""
)

var cssCmd = &cobra.Command{
	Use:   "css",
	Short: "Genera el honeytoken de css para paginas clonadas",
	Run:   run_Css,
}

func init() {
	cssCmd.Flags().StringVar(&inputPath, "in", "", "Archivo de CSS")
	cssCmd.Flags().StringVar(&outputPath, "out", "", "Archivo de CSS modificado, de no proveerse nada se crea en la carpeta out con el mismo nombre")
	cssCmd.Flags().StringVar(&url_original, "dominio", "", "Dominio del sitio original")
	cssCmd.MarkFlagRequired("in")
	cssCmd.MarkFlagRequired("url")

	rootCmd.AddCommand(cssCmd)
}

func run_Css(cmd *cobra.Command, args []string){
	if(outputPath == ""){
		outputPath = "out/" + inputPath
	}

	crearCss(inputPath, outputPath)

	cargarCssToken(outputPath)
}

func crearCss(input string, output string) {

	inFile, err := os.Open(input)
	if err != nil {
		panic(fmt.Errorf("error abriendo input: %w", err))
	}
	defer inFile.Close()
	
	outFile, err := os.Create(output)
	if err != nil {
		panic(fmt.Errorf("error creando output: %w", err))
	}
	defer outFile.Close()
	
	_, err = io.Copy(outFile, inFile)
	if err != nil {
		panic(fmt.Errorf("error clonando input: %w", err))
	}
	
	var cssContent = fmt.Sprintf(" \nbody {\n    background: url(%s) !important; \n}\n", ngrokurl+"fondo")
	
	_, err = outFile.WriteString(cssContent)
	if err!= nil{
		panic(fmt.Errorf("error agregando el css: %w", err))
	}

	fmt.Print("===================== Agregamos la linea siguiente linea =====================\n")
	fmt.Println(cssContent)
	fmt.Print("==============================================================================\n")

	fmt.Printf("Ahora debes cambiar el css original(%v) por el nuevo: %v\n",input, output)
}

func cargarCssToken(cssFile string){
	fmt.Printf("Tenemos que decidir como hacer el post \n")
	//Hay que mandar la url original al server para que cree el canary de dicha pagina
	//Si falla deberiamos borrar el contenido en cssFile
}