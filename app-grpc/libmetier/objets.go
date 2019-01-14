package libmetier

import "time"

// MessageSocial a common social msg
type MessageSocial struct {
	Data   string    `json:"data"`
	User   string    `json:"user"`
	Source string    `json:"source"`
	Date   time.Time `json:"timestamp"`
}
