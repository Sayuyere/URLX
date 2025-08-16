package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type LokiClient struct {
	url         string
	user        string
	apiKey      string
	language    string
	source      string
	serviceName string
	client      *http.Client
}

func NewLokiClient() *LokiClient {
	url := os.Getenv("GRAFANA_LOKI_URL")        // e.g. https://logs-prod-028.grafana.net/loki/api/v1/push
	user := os.Getenv("GRAFANA_LOKI_USER")      // e.g. 1306555
	apiKey := os.Getenv("GRAFANA_LOKI_API_KEY") // e.g. glc_...
	language := os.Getenv("GRAFANA_LOKI_LANGUAGE")
	if language == "" {
		language = "Go"
	}
	source := os.Getenv("GRAFANA_LOKI_SOURCE")
	if source == "" {
		source = "Code"
	}
	serviceName := os.Getenv("GRAFANA_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "urlx"
	}
	if url == "" || user == "" || apiKey == "" {
		panic("Loki URL, user, and API key must be set in environment variables")
	}
	return &LokiClient{
		url:         url,
		user:        user,
		apiKey:      apiKey,
		language:    language,
		source:      source,
		serviceName: serviceName,
		client:      &http.Client{Timeout: 5 * time.Second},
	}
}

type lokiPush struct {
	Streams []struct {
		Stream map[string]string `json:"stream"`
		Values [][2]string       `json:"values"`
	} `json:"streams"`
}

func (lc *LokiClient) SendLog(line string) error {
	entry := [2]string{
		strconv.FormatInt(time.Now().UnixNano(), 10),
		line,
	}
	streamLabels := map[string]string{
		"Language":     lc.language,
		"source":       lc.source,
		"service_name": lc.serviceName,
	}
	push := lokiPush{
		Streams: []struct {
			Stream map[string]string `json:"stream"`
			Values [][2]string       `json:"values"`
		}{
			{
				Stream: streamLabels,
				Values: [][2]string{entry},
			},
		},
	}
	b, _ := json.Marshal(push)
	req, err := http.NewRequest("POST", lc.url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(lc.user, lc.apiKey)
	resp, err := lc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("loki returned status: %s", resp.Status)
	}
	return nil
}
