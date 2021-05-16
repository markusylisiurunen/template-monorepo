package payloads

type Payload interface {
	UnmarshalJSON([]byte) error
}
