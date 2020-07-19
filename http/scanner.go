package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/freedomofpress/sdstatus/securedrop"
	"github.com/freedomofpress/sdstatus/status"
)

// Scanner allows to perform status checks across multiple HTTP services concurrently.
type Scanner struct {
	client *http.Client
}

// Scan performs status checks across multiple HTTP services concurrently.
func (s *Scanner) Scan(format string, onionURLs []string) error {

	i := 0

	results := make([]*securedrop.Instance, 0)

	ch := make(chan status.Information)

	// For each address we are creating a goroutine
	for _, v := range onionURLs {
		url := strings.TrimSpace(v)

		if url != "" {
			go func() {
				instance := securedrop.NewInstance(url)

				sc := securedrop.NewStatusChecker(s.client)
				err := sc.Check(instance)
				if err == nil {
					instance.Available = true
				}

				ch <- instance
			}()
			i = i + 1
		}

	}

	// Now wait for all the results
	for {
		result := <-ch

		if result != nil {

			if format == "csv" {
				fmt.Println(result.CSV())
			}

			results = append(results, result.(*securedrop.Instance))
			i = i - 1
		}
		if i == 0 {
			break
		}
	}

	if format == "json" {
		bits, err := json.MarshalIndent(results, "", "\t")
		if err == nil {
			fmt.Println(string(bits))
		}
	} else if format == "csv" {
	} else {
		log.Fatal(fmt.Sprintf("Invalid format: %s", format))
	}

	return nil
}

// NewScanner returns a new scanner configured to use a given HTTP client.
func NewScanner(c *http.Client) *Scanner {
	return &Scanner{
		client: c,
	}
}
