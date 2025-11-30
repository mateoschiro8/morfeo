package handlers

import (
	"fmt"
	"strings"
)

func Alert(msg string) {
	fmt.Println("ALERTAA! Fue activado el token: " + strings.ToLower(msg))
}
