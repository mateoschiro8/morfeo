package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "morfeo",
}

var (
	ngrokurl = GetNgrokUrl()
)

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	_ = rootCmd.Execute()
}

type NgrokTunnel struct {
	PublicURL string `json:"public_url"`
}

type NgrokResponse struct {
	Tunnels []NgrokTunnel `json:"tunnels"`
}

func GetNgrokUrl() string{

	//Se loopea por las dudas de que falle la respuesta de ngrok
	for i := 0; i < 5; i++ {
		resp, err := http.Get("http://127.0.0.1:4040/api/tunnels") //Se monta siempre en esta IP con ese puerto
		if err != nil {
			fmt.Println("Esperando a ngrok... (asegúrate de correr 'ngrok http 8000')")
			time.Sleep(2 * time.Second)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic("Error al leer respuesta")
		}

		var data NgrokResponse
		json.Unmarshal(body, &data)

		if data.Tunnels[0].PublicURL != "" {
			return data.Tunnels[0].PublicURL
		}
		
	}

	panic("NO ENCONTRE PUBLIC URL")
}