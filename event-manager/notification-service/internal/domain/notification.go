package domain

type NotificationMessage struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}
