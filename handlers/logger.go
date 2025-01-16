package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"strings"
)

// LogDetails representa os detalhes de um log com chaves e valores.
type LogDetails map[string]interface{}

// formatDetails formata a string dos detalhes para o log.
func formatDetails(details LogDetails) string {
	if len(details) == 0 {
		return ""
	}

	var fDetails strings.Builder
	for key, value := range details {
		fDetails.WriteString(fmt.Sprintf("%s: %v\n", key, value))
	}

	return fDetails.String()
}

// LogMessage realiza o registro de tela do log da requisição.
func LogMessage(c echo.Context, message string, details LogDetails) {
	log.Printf(
		"[%s | %s] %s\n%s",
		c.Path(),
		c.Request().Method,
		message,
		formatDetails(details),
	)
}
