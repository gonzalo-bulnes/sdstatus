package securedrop

import "fmt"

// Metadata stores JSON metadata about a SecureDrop instance
type Metadata struct {
	Fingerprint string `json:"gpg_fpr"`
	Version     string `json:"sd_version"`
}

// Instance represents a SecureDrop instance
type Instance struct {
	Available bool
	Info      Metadata
	URL       string
}

// CSV implements the sdstatus.Information interface
func (i Instance) CSV() string {
	return fmt.Sprintf("%s,%s,%s", i.URL, i.Info.Version, i.Info.Fingerprint)
}

// NewInstance returns a new SecureDrop instance
func NewInstance(url string) Instance {
	return Instance{
		URL: url,
	}
}
