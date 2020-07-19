package securedrop

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/freedomofpress/sdstatus/status"
)

// StatusChecker performs status checks for SecureDrop instances.
type StatusChecker struct {
	client    *http.Client
	doRequest func(*http.Request) (*http.Response, error)
}

// Check perform a request and makes the outcome available in the result channel.
func (c *StatusChecker) Check(instance status.Service) error {
	req, err := instance.NewStatusRequest()
	if err != nil {
		return fmt.Errorf("status check error: %w", err)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return fmt.Errorf("status check error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("instance responded with HTTP %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading status check response: %w", err)
	}

	if instance, ok := instance.(*Instance); ok {
		err = json.Unmarshal(body, &instance.Info)
		if err != nil {
			return fmt.Errorf("error deserializing status check response: %w", err)
		}
	} else {
		panic("Expected service to be a securedrop.Instance, but wasn't.")
	}

	return nil
}

// NewStatusChecker returns a status.Checker adequate for Securerop instances.
func NewStatusChecker(client *http.Client) status.Checker {
	c := &StatusChecker{
		client: client,
	}
	c.doRequest = client.Do
	return c
}
