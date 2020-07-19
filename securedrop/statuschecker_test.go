package securedrop

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestStatusCheck(t *testing.T) {

	t.Run("returns no errors when service response is OK", func(t *testing.T) {

		checker := &StatusChecker{}
		checker.doRequest = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
				StatusCode: http.StatusOK,
			}, nil
		}

		instance := &Instance{}

		err := checker.Check(instance)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})

	t.Run("returns an error unless the SecureDrop instance responds with HTTP 200 OK", func(t *testing.T) {

		testcases := []struct {
			responseStatusCode int
			expectedError      bool
		}{
			{http.StatusOK, false},
			{http.StatusNotFound, true},
			{http.StatusServiceUnavailable, true},
			{http.StatusRequestTimeout, true},
			{http.StatusBadRequest, true},
			{http.StatusGone, true},
		}

		for _, tc := range testcases {

			checker := &StatusChecker{}
			checker.doRequest = func(*http.Request) (*http.Response, error) {
				return &http.Response{
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
					StatusCode: tc.responseStatusCode,
				}, nil
			}

			instance := &Instance{}

			err := checker.Check(instance)

			if tc.expectedError && err == nil {
				t.Errorf("Expected error for response HTTP %d, got none", tc.responseStatusCode)
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		}
	})
}
