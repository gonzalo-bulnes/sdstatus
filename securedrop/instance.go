package securedrop

import (
	"fmt"
	"net/http"
)

const instanceMetadataURLPattern = "http://%s/metadata"

// Metadata stores JSON metadata about a SecureDrop instance.
type Metadata struct {
	Fingerprint string `json:"gpg_fpr"`
	Version     string `json:"sd_version"`
}

// Instance represents a SecureDrop instance.
type Instance struct {
	Available bool
	Info      Metadata
	URL       string `json:"Url"`
}

// CSV implements the sdstatus.Information interface.
func (i *Instance) CSV() string {
	return fmt.Sprintf("%s,%s,%s", i.URL, i.Info.Version, i.Info.Fingerprint)
}

// NewMetadataRequest returns an http.Request to get the instance metadata.
func (i *Instance) NewMetadataRequest() (r *http.Request, err error) {
	metadataURL := fmt.Sprintf(instanceMetadataURLPattern, i.URL)
	r, err = http.NewRequest("GET", metadataURL, nil)
	if err != nil {
		err = fmt.Errorf("status request creation failed: %w", err)
	}
	return
}

// NewStatusRequest implements the status.Service interface.
func (i *Instance) NewStatusRequest() (r *http.Request, err error) {
	return i.NewMetadataRequest()
}

// NewInstance returns a new SecureDrop instance.
func NewInstance(url string) *Instance {
	return &Instance{
		URL: url,
	}
}
