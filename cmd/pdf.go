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
	Run:   createPDFTokenWithJs,
}

func init() {
	pdfCmd.Flags().StringVar(&in, "in", "", "Ruta al archivo de entrada")
	pdfCmd.Flags().StringVar(&out, "out", "", "Ruta al archivo de salida")
	pdfCmd.MarkFlagRequired("in")
	pdfCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(pdfCmd)
}

func createPDFTokenWithJs(cmd *cobra.Command, args []string) {

	tokenID := CreateToken(msg, "")

	url := serverURL + "/bins/" + tokenID

	js := "app.launchURL('" + url + "', true);"

	err := license.SetMeteredKey(offlineLicenseKey)
	if err != nil {
		log.Fatalf("Licencia inválida: %v", err)
	}

	f, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader, err := model.NewPdfReader(f)
	if err != nil {
		panic(err)
	}

	isEncrypted, err := reader.IsEncrypted()
	if err != nil {
		panic(err)
	}

	if isEncrypted {
		ok, err := reader.Decrypt([]byte(""))
		if err != nil || !ok {
			panic(fmt.Errorf("PDF encriptado y sin password"))
		}
	}

	nPages, err := reader.GetNumPages()
	if err != nil {
		panic(err)
	}

	writer := model.NewPdfWriter()

	// Copiar páginas
	for i := 1; i <= nPages; i++ {
		page, err := reader.GetPage(i)
		if err != nil {
			panic(err)
		}
		err = writer.AddPage(page)
		if err != nil {
			panic(err)
		}
	}

	// Crear diccionario JavaScript
	jsDict := core.MakeDict()
	jsDict.Set("S", core.MakeName("JavaScript"))
	jsDict.Set("JS", core.MakeString(js))

	// Setear OpenAction del documento
	err = writer.SetOpenAction(jsDict)
	if err != nil {
		panic(err)
	}

	out, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	writer.Write(out)
}

const offlineLicenseKey = `1b6e2b4d1bf6137cd0b6f7b7b00ebb54196700b748c51cc460bbfaca74a5de74`
