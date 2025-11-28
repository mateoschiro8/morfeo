package cmd

import (
	"fmt"
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/spf13/cobra"
)

var qrCmd = &cobra.Command{
	Use:   "qr",
	Short: "Genera el honeytoken de qr",
	Run:   generateQRCode,
}

func init() {
	qrCmd.Flags().StringVar(&msg, "msg", "", "Identificador del qr")
	qrCmd.Flags().StringVar(&redirect, "redirect", "http://www.google.com", "Sitio al cual redirigir")
	qrCmd.MarkFlagRequired("msg")
	rootCmd.AddCommand(qrCmd)
}

func generateQRCode(cmd *cobra.Command, args []string) {

	tokenID := CreateToken(msg, redirect)

	data := serverURL + "/qrs/" + tokenID

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
