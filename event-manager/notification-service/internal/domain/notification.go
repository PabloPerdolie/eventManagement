package domain

// NotificationType defines the type of notification
type NotificationType string

const (
	EventCreated NotificationType = "event_created"
	TaskAssigned NotificationType = "task_assigned"
	ExpenseAdded NotificationType = "expense_added"
)

type NotificationResult struct {
	Success      bool   `json:"success"`
	MessageID    string `json:"message_id,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type EventNotification struct {
	EventID     int64  `json:"event_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	UserEmail   string `json:"user_email"`
}

type TaskNotification struct {
	TaskID    int64  `json:"task_id"`
	TaskName  string `json:"task_name"`
	EventName string `json:"event_name,omitempty"`
	UserEmail string `json:"user_email"`
}

type ExpenseNotification struct {
	ExpenseID   int64   `json:"expense_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description,omitempty"`
	EventName   string  `json:"event_name,omitempty"`
	UserEmail   string  `json:"user_email"`
}
