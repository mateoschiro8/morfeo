package cmd

import (
	"encoding/base64"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var binCmd = &cobra.Command{
	Use:   "bin",
	Short: "Genera un honeytoken a partir de un binario",
	Run:   generateBinaryWrapper,
}

func init() {
	binCmd.Flags().StringVar(&msg, "msg", "", "Identificador del qr")
	binCmd.Flags().StringVar(&in, "in", "", "Path al binario a wrappear")
	binCmd.Flags().StringVar(&out, "out", "out", "Path al binario de salida")
	binCmd.MarkFlagRequired("msg")
	binCmd.MarkFlagRequired("in")
	rootCmd.AddCommand(binCmd)
}

func generateBinaryWrapper(cmd *cobra.Command, args []string) {

	tokenID := CreateToken(msg, "")

	data, err := os.ReadFile(in)
	if err != nil {
		panic(err)
	}

	b64 := base64.StdEncoding.EncodeToString(data)

	code := strings.ReplaceAll(wrapperTemplate, "{{B64}}", b64)
	code = strings.ReplaceAll(code, "{{Endpoint}}", serverURL+"/bins/"+tokenID)

	os.WriteFile("wrapper_tmp.go", []byte(code), 0644)

	outCmd := exec.Command("go", "build", "-o", out, "wrapper_tmp.go")
	outCmd.Stdout = os.Stdout
	outCmd.Stderr = os.Stderr
	outCmd.Run()

	os.Remove("wrapper_tmp.go")
}

const wrapperTemplate = `
package main

import (
    "encoding/base64"
    "net/http"
    "os"
    "os/exec"
    "io/ioutil"
    "path/filepath"
)

const encoded = "{{B64}}"
const endpoint = "{{Endpoint}}"

func sendAlert() {
    if endpoint == "" {
        return
    }
    _, _ = http.Get(endpoint)
}

func main() {
    
    sendAlert()

    data, _ := base64.StdEncoding.DecodeString(encoded)
    tmpDir, _ := os.MkdirTemp("", "honey-*")
    real := filepath.Join(tmpDir, "realbin")
    ioutil.WriteFile(real, data, 0755)

    cmd := exec.Command(real, os.Args[1:]...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin  = os.Stdin

    err := cmd.Run()
    if err != nil {
        if e, ok := err.(*exec.ExitError); ok {
            os.Exit(e.ExitCode())
        }
        panic(err)
    }
}
`
