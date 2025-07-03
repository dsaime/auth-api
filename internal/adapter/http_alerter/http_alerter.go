package httpAlerter

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"net"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type HttpAlerter struct {
	Webhook string
}

func (a *HttpAlerter) UnknownIPRefreshSessionAlert(sessionID uuid.UUID, oldIP, newIP net.IP) {
	if a.Webhook == "" {
		return
	}

	t, err := template.New("").Parse(a.Webhook)
	if err != nil {
		slog.Error("http_alerter.UnknownIPRefreshSessionAlert: " + err.Error())
		return
	}

	var readWriter bufio.ReadWriter
	if err = t.Execute(&readWriter, map[string]any{
		"SessionID": sessionID,
		"OldIP":     oldIP,
		"NewIP":     newIP,
	}); err != nil {
		slog.Error("http_alerter.UnknownIPRefreshSessionAlert: " + err.Error())
		return
	}

	url, err := io.ReadAll(readWriter)
	if err != nil {
		slog.Error("http_alerter.UnknownIPRefreshSessionAlert: " + err.Error())
		return
	}

	_, _, errs := fiber.Post(string(url)).Bytes()
	if len(errs) > 0 {
		slog.Error("http_alerter.UnknownIPRefreshSessionAlert: " + errors.Join(errs...).Error())
	}
}
