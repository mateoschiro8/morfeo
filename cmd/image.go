package cmd

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Genera un honeytoken de imagen desde una imagen existente o crea una nueva",
	Run:   generateImageToken,
}

func init() {
	imageCmd.Flags().StringVar(&msg, "msg", "", "Identificador del honeytoken")
	imageCmd.Flags().StringVar(&chat, "chat", "", "Chat ID al cual enviar mensaje al activarse")
	imageCmd.Flags().StringVar(&in, "in", "", "Path a la imagen de entrada (opcional, si no hay se crea una imagen vacia)")
	imageCmd.Flags().StringVar(&out, "out", "honeytoken_image.html", "Path al archivo HTML de salida")
	imageCmd.MarkFlagRequired("msg")
	imageCmd.MarkFlagRequired("chat")
	rootCmd.AddCommand(imageCmd)
}

func generateImageToken(cmd *cobra.Command, args []string) {

	tokenID := CreateToken(msg, "")

	imageURL := serverURL + "/track?id=" + tokenID

	if out == "" {
		out = "honeytoken_image.html"
	}

	var htmlContent string
	var svgContent string

	if in != "" {
		fmt.Printf("Generando honeytoken desde imagen: %s\n", in)

		if _, err := os.Stat(in); err != nil {
			panic(fmt.Errorf("error: la imagen de entrada no existe: %w", err))
		}

		inputFile, err := os.Open(in)
		if err != nil {
			panic(fmt.Errorf("error abriendo imagen para leer dimensiones: %w", err))
		}
		imgConfig, _, err := image.DecodeConfig(inputFile)
		inputFile.Close()

		width := 800
		height := 600
		if err == nil {
			width = imgConfig.Width
			height = imgConfig.Height
		}
		// creo html
		htmlContent = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
* { margin: 0; padding: 0; }
body { overflow: hidden; background: #2b2b2b; }
img { display: block; width: 100%%; height: 100vh; object-fit: contain; }
</style>
</head>
<body>
<img src="%s" alt="">
<img src="%s" alt="" style="position:absolute;width:1px;height:1px;opacity:0;">
</body>
</html>`, filepath.Base(in), imageURL)

		// creo svg
		svgContent = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 %d %d" preserveAspectRatio="xMidYMid meet" style="width:100%%;height:100vh;background:#2b2b2b">
  <image href="%s" width="%d" height="%d"/>
  <image href="%s" width="1" height="1" opacity="0"/>
</svg>`, width, height, filepath.Base(in), width, height, imageURL)

		fmt.Printf("Imagen referenciada: %s\n", in)

	} else {
		// si no hay imagen creo html vacio
		fmt.Println("Generando tracking pixel invisible")

		htmlContent = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
* { margin: 0; padding: 0; }
body { background: #fff; }
</style>
</head>
<body>
<img src="%s" alt="" style="position:absolute;width:1px;height:1px;opacity:0;">
</body>
</html>`, imageURL)

		// creo svg vacio
		svgContent = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="1" height="1" viewBox="0 0 1 1">
  <image href="%s" width="1" height="1" opacity="0"/>
</svg>`, imageURL)
	}

	// Guardar el archivo HTML
	err := os.WriteFile(out, []byte(htmlContent), 0644)
	if err != nil {
		panic(fmt.Errorf("error creando archivo HTML: %w", err))
	}

	// Guardar el archivo SVG
	svgPath := out[:len(out)-len(filepath.Ext(out))] + ".svg"
	err = os.WriteFile(svgPath, []byte(svgContent), 0644)
	if err != nil {
		panic(fmt.Errorf("error creando archivo SVG: %w", err))
	}

	fmt.Printf("Token ID: %s\n", tokenID)
	fmt.Printf("URL del tokenn: %s\n", imageURL)
	fmt.Printf("Archivo HTML generado: %s\n", out)
	fmt.Printf("Archivo SVG generado: %s\n", svgPath)
	if in != "" {
		fmt.Printf("Imagen original referenciada: %s\n", filepath.Base(in))
	}
}

// generateBlankImage crea una imagen de 1x1 pixel transparente
func generateBlankImage() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{0, 0, 0, 0})

	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}
