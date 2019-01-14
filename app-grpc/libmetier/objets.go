package libmetier

import (
	"encoding/json"
	"log"
	"time"
)

// MessageSocial a common social msg
type MessageSocial struct {
	Data   string    `json:"data"`
	User   string    `json:"user"`
	Source string    `json:"source"`
	Date   time.Time `json:"timestamp"`
}

func (ms MessageSocial) toMessageSocial(mstpl []byte) {
	err := json.Unmarshal(mstpl, &ms)
	if err != nil {
		log.Println(err)
	}
}

func (ms MessageSocial) toByteArray() []byte {
	b, err := json.Marshal(ms)
	if err != nil {
		log.Println(err)
	}
	return []byte(b)
}
