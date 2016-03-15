package apns2_test

import (
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
	payload := []byte(`{ "aps" : { "alert" : "Hello World" } }`)
	apnsID := "922D9F1F-B82E-B337-EDC9-DB4FC8527676"

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
		if !reflect.DeepEqual(body, payload) {
			t.Errorf("Expected body %v, got %v", payload, body)
		}
		w.Header().Set("apns-id", apnsID)
	})

	client := apns2.Client{
		HTTPClient: http.DefaultClient,
		Host:       server.URL,
	}

	resp, err := client.Send(payload, deviceToken)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println(resp.Header.Get("apns-id"))
	}
}
