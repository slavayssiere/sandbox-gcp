package libmetier

// MessageSocial a common social msg
type MessageSocial struct {
	Data   string `json:"data"`
	User   string `json:"user"`
	Source string `json:"source"`
}
