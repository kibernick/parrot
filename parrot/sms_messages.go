package parrot

// SMSMessage is a container for SMS message data that will be sent to MessageBird.
type SMSMessage struct {
	Recipient  int
	Originator string
	Message    [][]rune
}

// MaxMessageSize refers to the character limit for a single SMS text message transmission.
const MaxMessageSize = 160

// SplitMessageIntoParts splits a simple message string into multiple parts, ready to be
// transmitted as concatenated SMS.
func SplitMessageIntoParts(message string) [][]rune {
	msg := []rune(message)
	var msgParts [][]rune

	for len(msg) > MaxMessageSize {
		msgParts = append(msgParts, msg[:MaxMessageSize])
		msg = msg[MaxMessageSize:]
	}
	msgParts = append(msgParts, msg)
	return msgParts
}

// Send will start the process of sending an SMS message via MessageBird.
func (m SMSMessage) Send() (string, error) {
	return "OK", nil
}
