package apns2

import (
	"net/http"
	"strconv"
	"time"
)

// Headers Request headers for apple push notification
type Headers struct {
	ID          string
	Expiration  time.Time
	LowPriority bool
	Topic       string
}

// Set request headers for HTTPClient
func (h *Headers) Set(header http.Header) {

	header.Set("Content-Type", "application/json")

	if h.ID != "" {
		header.Set(ApnsID, h.ID)
	}

	if !h.Expiration.IsZero() {
		timestamp := strconv.FormatInt(h.Expiration.Unix(), 10)
		header.Set(ApnsExpiration, timestamp)
	}

	if h.LowPriority {
		header.Set(ApnsPriority, "5")
	}

	if h.Topic != "" {
		header.Set(ApnsTopic, h.Topic)
	}
}
