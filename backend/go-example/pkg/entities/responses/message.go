package responses

import (
	"encoding/json"
	"time"
)

type Message struct {
	ID        string
	CreatedAt time.Time
	Text      string
}

type jsonMessage struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
}

func (m Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonMessage{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		Text:      m.Text,
	})
}
