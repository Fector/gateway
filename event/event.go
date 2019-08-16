package event

type Event struct {
	Message   string
	Error     error
	Timestamp int64
}
