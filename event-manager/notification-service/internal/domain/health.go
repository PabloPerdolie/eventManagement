package domain

// Health represents the health status of the service
type Health struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	Service   string `json:"service"`
}
