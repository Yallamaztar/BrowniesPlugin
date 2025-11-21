package discord

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Yallamaztar/BrowniesPlugin/helpers"
)

type embed struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Color       int    `json:"color,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
}

type webhookPayload struct {
	Content string  `json:"content,omitempty"`
	Embeds  []embed `json:"embeds,omitempty"`
}

func SendWebhook(data webhookPayload) {
	url := os.Getenv("DISCORD_WEBHOOK")
	if url == "" {
		return
	}

	go func(u string, p webhookPayload) {
		b, err := json.Marshal(p)
		if err != nil {
			return
		}
		req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(b))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		io.Copy(io.Discard, resp.Body)
	}(url, data)
}

func WinWebhook(player string, amount int64) {
	SendWebhook(webhookPayload{
		Embeds: []embed{{
			Title:       "Gamble Win 🎉",
			Description: "**" + player + "** won **$" + helpers.FormatMoney(amount) + "**",
			Color:       0x00ff00,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		}},
	})
}

func LossWebhook(player string, amount int64) {
	SendWebhook(webhookPayload{
		Embeds: []embed{{
			Title:       "Gamble Loss 😿",
			Description: "**" + player + "** lost **$" + helpers.FormatMoney(amount) + "**",
			Color:       0xff0000,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		}},
	})
}

func GamblingEnabledWebhook(enabled bool) {
	status := "disabled 😿"
	color := 0xff0000

	if enabled {
		status = "enabled 🎉"
		color = 0x00ff00
	}

	SendWebhook(webhookPayload{
		Embeds: []embed{{
			Title:       "Gambling Status Changed",
			Description: "Gambling has been **" + status + "**",
			Color:       color,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		}},
	})
}

func GamblingMaxBetWebhook(amount int64) {
	desc := "Max bet has been **disabled 🎉**"
	color := 0xff0000

	if amount > 0 {
		desc = "Max bet set to **$" + helpers.FormatMoney(amount) + " 😿**"
		color = 0x00ff00
	}

	SendWebhook(webhookPayload{
		Embeds: []embed{{
			Title:       "Gambling Max Bet Changed",
			Description: desc,
			Color:       color,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		}},
	})
}

func OnReadyWebhook() {
	SendWebhook(webhookPayload{
		Embeds: []embed{{
			Title:       "Brownies Plugin is Online ✅",
			Description: "The Brownies Plugin has started and is now online",
			Color:       0x00ff00,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		}},
	})
}

func OnShutdownWebhook() {
	SendWebhook(webhookPayload{
		Embeds: []embed{{
			Title:       "Brownies Plugin is Offline ❌",
			Description: "The Brownies Plugin has been shut down and is now offline",
			Color:       0xff0000,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		}},
	})
}
