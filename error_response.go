package apns2

// ErrorResponse contains reason, timestamp
type ErrorResponse struct {
	Reason    string `json:"reason,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}
