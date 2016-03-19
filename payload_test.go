package apns2_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sger/go-apns2"
)

func TestPayload(t *testing.T) {
	var tests = []struct {
		input    apns2.Payload
		expected []byte
	}{
		{
			apns2.Payload{
				Alert: apns2.Alert{Body: "Hello World"},
			},
			[]byte(`{"aps":{"alert":"Hello World"}}`),
		},
		{
			apns2.Payload{
				Alert: apns2.Alert{
					Title: "My Title",
					Body:  "Hello APNS 2"},
			},
			[]byte(`{"aps":{"alert":{"title":"My Title","body":"Hello APNS 2"}}}`),
		},
	}

	for _, tt := range tests {
		testPayload(t, tt.input, tt.expected)
	}
}

func testPayload(t *testing.T, p interface{}, expected []byte) {
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal("Error", err)
	}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("Expected %s, got %s", expected, b)
	}
}
