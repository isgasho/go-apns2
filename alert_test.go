package apns2_test

import (
	"testing"

	"github.com/sger/go-apns2"
)

func TestAlert(t *testing.T) {

	payload := apns2.Payload{
		Alert: apns2.Alert{Body: "Hello World"},
	}

	alert := apns2.Alert{
		Body: "Hello World",
	}

	if alert.Body != payload.Alert.Body {
		t.Errorf("Expected %s, got %s", alert.Body, payload.Alert.Body)
	}
}
