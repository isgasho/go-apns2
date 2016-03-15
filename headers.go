package apns2

import (
	"net/http"
	"strconv"
	"time"
)

// Header request headers for apple push notification
type Headers struct {
	ID          string
	Expiration  time.Time
	LowPriority bool
	Topic       string
}

func (h *Headers) Set(header http.Header) {
	header.Set(ApnsID, h.ID)
	header.Set(ApnsExpiration, strconv.FormatInt(h.Expiration.Unix(), 10))
	header.Set(ApnsPriority, "5")
	header.Set(ApnsTopic, h.Topic)
}
