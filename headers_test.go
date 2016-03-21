package apns2_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/sger/go-apns2"
)

func TestHeaders(t *testing.T) {
	headers := apns2.Headers{
		ID:          "12345",
		Expiration:  time.Unix(1458565017, 0),
		LowPriority: true,
		Topic:       "com.test.app-id",
	}

	requestHeader := http.Header{}
	headers.Set(requestHeader)

	testHeader(t, requestHeader, "apns-id", headers.ID)
	testHeader(t, requestHeader, "apns-expiration", "1458565017")
	testHeader(t, requestHeader, "apns-priority", "5")
	testHeader(t, requestHeader, "apns-topic", headers.Topic)
}

func testHeader(t *testing.T, requestHeader http.Header, key string, expected string) {
	headerKey := requestHeader.Get(key)
	if headerKey != expected {
		t.Errorf("Expected %s %q, got %q.", key, expected, headerKey)
	}
}
