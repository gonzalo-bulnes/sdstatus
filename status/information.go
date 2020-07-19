package status

// Information represents data that can be serialized as CSV
type Information interface {
	CSV() string
}
