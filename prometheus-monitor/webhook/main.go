package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

//ref: https://prometheus.io/docs/alerting/latest/configuration/#webhook_config
type AlertManagerWebhook struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []struct {
		Status       string            `json:"status"`
		Label        map[string]string `json:"labels"`
		Annotations  map[string]string `json:"annotations"`
		StartsAt     string            `json:"startsAt"`
		EndsAt       string            `json:"endsAt"`
		GeneratorURL string            `json:"generatorURL"`
		Fingerprint  string            `json:"fingerprint"`
	} `json:"alerts"`
}

//ref: https://discord.com/developers/docs/resources/webhook#execute-webhook
type (
	DiscordEmbedField struct {
		Name   string `json:"name"`
		Value  string `json:"value"`
		Inline bool   `json:"inline"`
	}
	DiscordEmbed struct {
		Title       string              `json:"title"`
		Description string              `json:"description"`
		Fields      []DiscordEmbedField `json:"fields"`
	}
	DiscordHook struct {
		Content string         `json:"content"`
		Embeds  []DiscordEmbed `json:"embeds"`
	}
)

func discordhook(ctx *gin.Context) {
	info := &AlertManagerWebhook{}
	json.NewDecoder(ctx.Request.Body).Decode(info)
	fmt.Printf("%+v\n", info)
	hook := DiscordHook{
		Content: "=== Alert ===",
		Embeds: []DiscordEmbed{
			{
				Title:       fmt.Sprintf("[%s] %s", strings.ToUpper(info.Status), info.CommonAnnotations["summary"]),
				Description: info.CommonAnnotations["description"],
			},
		},
	}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(hook)
	if _, err := http.Post(os.Getenv("WEBHOOK_URL"), "application/json", &buf); err != nil {
		log.Println(err)
	}
}

func main() {
	r := gin.Default()
	r.POST("/webhook", discordhook)
	r.Run(":8080")
}
