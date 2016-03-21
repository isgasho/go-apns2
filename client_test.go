package apns2_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sger/go-apns2"
)

func TestPush(t *testing.T) {
	deviceToken := "c7800a79efffe8ffc01b280717a936937cb69f8ca307545eb6983c60f12e167a"
	payload := apns2.Payload{
		Alert: apns2.Alert{
			Body: "Hello World"},
	}

	apnsID := "674EB1D5-7E7C-3DC9-B0F5-32A55E54960E"

	handler := http.NewServeMux()
	server := httptest.NewServer(handler)

	handler.HandleFunc("/3/device/", func(w http.ResponseWriter, r *http.Request) {
		expectURL := fmt.Sprintf("/3/device/%s", deviceToken)
		if r.URL.String() != expectURL {
			t.Errorf("Expected url %v, got %v", expectURL, r.URL)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(body, b) {
			t.Errorf("Expected body %v, got %v", payload, body)
		}
		w.Header().Set("apns-id", apnsID)
	})

	client := apns2.Client{
		HTTPClient: http.DefaultClient,
		Host:       server.URL,
	}

	resp, err := client.SendPush(payload, deviceToken, &apns2.Headers{})
	if err != nil {
		t.Error(err)
	}

	remoteApnsID := resp.ApnsID

	if remoteApnsID != apnsID {
		t.Errorf("Expected apns-id %q, but got %q ", apnsID, remoteApnsID)
	}
}
