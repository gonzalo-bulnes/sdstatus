package securedrop

import (
	"testing"
)

func TestInstanceCSV(t *testing.T) {
	t.Run("known values", func(t *testing.T) {
		i := Instance{
			Info: Metadata{
				Fingerprint: "F6E0E2901B787C3721E1C0BF4BD6284B525A3DF4",
				Version:     "1.4.1",
			},
			URL: "zdf4nikyuswdzbt6.onion",
		}
		expected := "zdf4nikyuswdzbt6.onion,1.4.1,F6E0E2901B787C3721E1C0BF4BD6284B525A3DF4"

		if result := i.CSV(); result != expected {
			t.Errorf("Expected '%s', got %s", expected, result)
		}
	})
}

func TestInstanceNewMetadataRequest(t *testing.T) {
	t.Run("returns an http.Request to get instance metadata", func(t *testing.T) {
		i := Instance{
			URL: "zdf4nikyuswdzbt6.onion",
		}
		expected := struct {
			method string
			url    string
		}{
			method: "GET",
			url:    "http://zdf4nikyuswdzbt6.onion/metadata",
		}

		r, err := i.NewMetadataRequest()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if method := r.Method; method != expected.method {
			t.Errorf("Expected %s method, got %s", expected.method, method)
		}
		if url := r.URL.String(); url != expected.url {
			t.Errorf("Expected %s URL, got %s", expected.url, url)
		}
	})
}
