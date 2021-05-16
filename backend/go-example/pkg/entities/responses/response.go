package responses

type Response interface {
	MarshalJSON() ([]byte, error)
}
