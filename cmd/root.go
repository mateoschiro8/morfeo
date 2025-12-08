package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		godotenv.Load()
		serverURL = os.Getenv("SERVERURL")

		if cmd.Name() == "server" {
			return nil
		}

		msg, _ = cmd.Flags().GetString("msg")
		chat, _ = cmd.Flags().GetString("chat")

		var missingFlags []string

		if msg == "" {
			missingFlags = append(missingFlags, "msg")
		}

		if chat == "" {
			missingFlags = append(missingFlags, "chat")
		}

		if len(missingFlags) == 0 {
			return nil
		}

		return fmt.Errorf("required flag(s) missing: --%s", strings.Join(missingFlags, ", --"))
	},
}

func init() {
	rootCmd.PersistentFlags().String("msg", "", "Identificador del token")
	rootCmd.PersistentFlags().String("chat", "", "Chat ID al cual enviar la alerta al ser activado")
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
