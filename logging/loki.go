package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
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

type logEntry struct {
	line    string
	level   string
	service string
	labels  map[string]string
}

type LokiBatcher struct {
	client    *LokiClient
	queue     chan logEntry
	batchSize int
	interval  time.Duration
	stop      chan struct{}
	wg        sync.WaitGroup
}

func NewLokiBatcher(client *LokiClient, batchSize int, interval time.Duration) *LokiBatcher {
	b := &LokiBatcher{
		client:    client,
		queue:     make(chan logEntry, 1000),
		batchSize: batchSize,
		interval:  interval,
		stop:      make(chan struct{}),
		wg:        sync.WaitGroup{},
	}
	b.wg.Add(1)
	go b.run()
	return b
}

func (b *LokiBatcher) run() {
	defer b.wg.Done()
	batch := make([]logEntry, 0, b.batchSize)
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()
	for {
		select {
		case entry := <-b.queue:
			batch = append(batch, entry)
			if len(batch) >= b.batchSize {
				b.flush(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				b.flush(batch)
				batch = batch[:0]
			}
		case <-b.stop:
			if len(batch) > 0 {
				b.flush(batch)
			}
			return
		}
	}
}

func (b *LokiBatcher) flush(entries []logEntry) {
	if len(entries) == 0 {
		return
	}
	streams := []struct {
		Stream map[string]string `json:"stream"`
		Values [][2]string       `json:"values"`
	}{}
	for _, e := range entries {
		labels := e.labels
		if labels == nil {
			labels = map[string]string{
				"Language":     b.client.language,
				"source":       b.client.source,
				"service_name": e.service,
			}
		}
		streams = append(streams, struct {
			Stream map[string]string `json:"stream"`
			Values [][2]string       `json:"values"`
		}{
			Stream: labels,
			Values: [][2]string{{strconv.FormatInt(time.Now().UnixNano(), 10), e.line}},
		})
	}
	push := lokiPush{Streams: streams}
	b.client.sendBatch(push)
}

func (b *LokiBatcher) SendLog(line, level, service string, labels map[string]string) {
	select {
	case b.queue <- logEntry{line, level, service, labels}:
		// enqueued successfully
	default:
		fmt.Fprintf(os.Stdout, "[LokiBatcher] Dropping log: queue full (level=%s, service=%s)\n", level, service)
	}
}

func (b *LokiBatcher) Close() {
	close(b.stop)
	b.wg.Wait()
}

// Replace LokiClient.SendLog with async batcher
func (lc *LokiClient) sendBatch(push lokiPush) error {
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
