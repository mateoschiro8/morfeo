package cmd

import (
	"fmt"
	"image/png"
	"os"

	"encoding/base64"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/spf13/cobra"
)

var qrCmd = &cobra.Command{
	Use:   "qr",
	Short: "Genera el honeytoken de pdf",
	Run:   generateQRCode,
}

var id string

func init() {
	qrCmd.Flags().StringVar(&id, "id", "", "Identificador del qr")
	rootCmd.AddCommand(qrCmd)
}

func generateQRCode(cmd *cobra.Command, args []string) {

	data := "http://localhost:8000" + "/qs?data=" + base64.RawURLEncoding.EncodeToString([]byte(id))

	qrCode, err := qr.Encode(data, qr.M, qr.Auto)
	if err != nil {
		panic(fmt.Errorf("failed to encode QR code: %w", err))
	}

	scaledQR, err := barcode.Scale(qrCode, 300, 300)
	if err != nil {
		panic(fmt.Errorf("failed to scale QR code: %w", err))
	}

	file, err := os.Create("qrcode.png")
	if err != nil {
		panic(fmt.Errorf("failed to create file: %w", err))
	}
	defer file.Close()

	err = png.Encode(file, scaledQR)
	if err != nil {
		panic(fmt.Errorf("failed to encode QR code image: %w", err))
	}
}
