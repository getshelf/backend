type Sender interface {
	sender(from string, to string, subject string, body string) bool
}