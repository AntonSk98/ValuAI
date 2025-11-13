package mail

// SendMailCommand represents the necessary information to send an email.
type SendMailCommand struct {
	To    string
	Title string
	Body  string
}
