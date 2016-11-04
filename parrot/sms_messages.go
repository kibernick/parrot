package parrot

type SMSMessage struct {
	Recipient  int
	Originator string
	Message    [][]rune
}

const MaxMessageSize = 160

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

func NewSMSMessage(recipient int, originator, message string) SMSMessage {
	return SMSMessage{
		Recipient:  recipient,
		Originator: originator,
		Message:    SplitMessageIntoParts(message),
	}
}

func (m SMSMessage) Send() (string, error) {
	return "OK", nil
}
