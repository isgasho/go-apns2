package apns2

// ApnsResponse contains apns-id, reason, status code, status code description.
type ApnsResponse struct {
	StatusCode            int
	StatusCodeDescription string
	ApnsID                string `json:"apns-id,omitempty"`
	Reason                string `json:"reason,omitempty"`
}
