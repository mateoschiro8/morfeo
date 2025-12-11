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

func formatIn() {
	if strings.Split(in, "/")[0] != "input" {
		in = "input/" + in
	}
	fmt.Printf("El token se creara en base a: %v\n", in)
}

func formatOut() {
	if out == "" {
		var directories = strings.Split(in, "/")
		out = "output/" + directories[len(directories)-1]
	}else{
		var directories = strings.Split(out, "/")
		if(directories[0]=="tokens"){
			directories[0]="output"
			out = strings.Join(directories, "/")
		}else if(directories[0] != "output"){
			out = "output/" + out
		}
	}

	fmt.Printf("El token se creara en: %v\n", out)

}