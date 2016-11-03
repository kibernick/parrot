package parrot

type SMSMessage struct {
	Recipient  int
	Originator string
	Message    [][]byte
}

func NewSMSMessage(recipient int, originator, message string) SMSMessage {
	messageLine := []byte("This is a line of text.")
	myMessage := [][]byte{}
	myMessage = append(myMessage, messageLine)
	myMessage = append(myMessage, []byte(message))
	return SMSMessage{
		Recipient:  recipient,
		Originator: originator,
		Message:    myMessage,
	}
}

func (m SMSMessage) Send() (string, error) {
	return "OK", nil
}
