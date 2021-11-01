package messages

import "encoding/json"

type HelloPayload struct {
	Greeting string `json:"greeting"`
	Name     string `json:"name"`
}

func (p *HelloPayload) MarshalPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p *HelloPayload) UnmarshalPayload(payload []byte) error {
	return json.Unmarshal(payload, p)
}

func NewHelloPayload(greeting string, name string) *HelloPayload {
	return &HelloPayload{
		Greeting: greeting,
		Name:     name,
	}
}
