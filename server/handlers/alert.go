package handlers

import (
	"fmt"
	"strings"
)

func Alert(msg string) {
	fmt.Println("ALERTA! Fue activado el token: " + strings.ToLower(msg))
}
