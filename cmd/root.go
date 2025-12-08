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

var rootCmd = &cobra.Command{
	Use: "morfeo",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		godotenv.Load()
		serverURL = os.Getenv("SERVERURL")
	},
}

var (
	serverURL string
	msg       string
	extra  string
	in        string
	out       string
	chat 	  string
)

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	_ = rootCmd.Execute()
}

func CreateToken(msg string, extra string, chat string) string {

	data := types.UserInput{
		Msg:      msg,
		Extra: extra,
		Chat: chat,
	}
	
	body, err := json.Marshal(data)
	if err != nil{
		panic(err)
	}
	resp, err := http.Post(serverURL+"/tokens",
                        	"application/json",
                        	bytes.NewBuffer(body))
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil{
		panic(err)
	}
	tokenID := string(respBytes)
	return tokenID
	// body, err := json.Marshal(data)
	// if err != nil {
	// 	panic(err)
	// }

	// resp, err := http.Post(serverURL+"/tokens", "application/json", bytes.NewBuffer(body))
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// respBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }

	// respString := string(respBytes)
	// fmt.Println(respString)
	// return respString
}
