package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
)

var (
	imageURL string = "http://localhost:8000/track"
	uri      string
)

var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "Genera el honeytoken de pdf",
	Run:   createPDFTokenWith,
}

func init() {
	pdfCmd.Flags().StringVar(&in, "in", "", "Ruta al archivo de entrada")
	pdfCmd.Flags().StringVar(&out, "out", "", "Ruta al archivo de salida")
	pdfCmd.MarkFlagRequired("in")
	pdfCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(pdfCmd)
}

func checkError(err error){
	if err != nil {
        panic(err)
    }
}

func createPDFTokenWith(cmd *cobra.Command, args []string) {

	tokenID := CreateToken(msg, "")

	url := serverURL + "/bins/" + tokenID

	injectedCode := "app.launchURL('" + url + "', true);"

	err := license.SetMeteredKey(offlineLicenseKey)
	if err != nil {
		log.Fatalf("Licencia inválida: %v", err)
	}

	f, err := os.Open(in)
	checkError(err)
	defer f.Close()

	reader, err := model.NewPdfReader(f)
	checkError(err)

	isEncrypted, err := reader.IsEncrypted()
	checkError(err)

	if isEncrypted {
		ok, err := reader.Decrypt([]byte(""))
		if err != nil || !ok {
			panic(fmt.Errorf("PDF encriptado y sin password"))
		}
	}

	nPages, err := reader.GetNumPages()
	checkError(err)

	writer := model.NewPdfWriter()

	for i := 1; i <= nPages; i++ {
		page, err := reader.GetPage(i)
		checkError(err)
		err = writer.AddPage(page)
		checkError(err)
	}

	// Crear diccionario JavaScript
	dict := core.MakeDict()
	dict.Set("S", core.MakeName("JavaScript"))
	dict.Set("", core.MakeString(injectedCode))

	// Setear OpenAction del documento. Se dispara apenas se abre con acrobat reader
	err = writer.SetOpenAction(dict)
	checkError(err)

	out, err := os.Create(out)
	checkError(err)
	defer out.Close()

	writer.Write(out)
}

const offlineLicenseKey = `1b6e2b4d1bf6137cd0b6f7b7b00ebb54196700b748c51cc460bbfaca74a5de74`
