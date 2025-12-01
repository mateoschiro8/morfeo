package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type SendMessageReq struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func Alert(msg string, ip string) {

	alertText := "ALERTA! Fue activado el token " + strings.ToLower(msg) + " desde la IP: " + ip
	fmt.Println(alertText)

	token := os.Getenv("BOT_TOKEN")
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	body := SendMessageReq{
		ChatID: chatID,
		Text:   alertText,
	}

	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
