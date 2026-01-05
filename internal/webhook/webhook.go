package webhook

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// Notifier sends events to external webhooks
type Notifier struct {
	url    string
	client *http.Client
}

// Event is the payload sent to webhooks
type Event struct {
	Event     string          `json:"event"`
	Timestamp string          `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// New creates a webhook notifier (nil if URL empty)
func New(url string) *Notifier {
	if url == "" {
		return nil
	}
	return &Notifier{
		url:    url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// Notify sends an event to the webhook (non-blocking)
func (n *Notifier) Notify(event string, data interface{}) {
	if n == nil {
		return
	}

	go func() {
		dataBytes, _ := json.Marshal(data)
		payload := Event{
			Event:     event,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Data:      dataBytes,
		}

		body, _ := json.Marshal(payload)
		resp, err := n.client.Post(n.url, "application/json", strings.NewReader(string(body)))
		if err != nil {
			log.Printf("Webhook error: %v", err)
			return
		}
		resp.Body.Close()
	}()
}
