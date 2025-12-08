package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mateoschiro8/morfeo/server/types"
	"github.com/spf13/cobra"
)

var (
	serverURL string

	msg  string
	chat string

	in    string
	out   string
	extra string
)

var rootCmd = &cobra.Command{
	Use: "morfeo",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		godotenv.Load()
		serverURL = os.Getenv("SERVERURL")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&msg, "msg", "", "Identificador del token")
	rootCmd.PersistentFlags().StringVar(&chat, "chat", "", "Chat ID al cual enviar la alerta al ser activado")

	_ = rootCmd.MarkPersistentFlagRequired("msg")
	_ = rootCmd.MarkPersistentFlagRequired("chat")
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	_ = rootCmd.Execute()
}

func CreateToken(msg string, extra string, chat string) string {

	data := types.UserInput{
		Msg:   msg,
		Extra: extra,
		Chat:  chat,
	}

	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(serverURL+"/tokens", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	tokenID := string(respBytes)

	return tokenID
}
