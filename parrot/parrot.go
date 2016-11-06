package parrot

import (
	"errors"
	"fmt"
	"time"

	"github.com/messagebird/go-rest-api"
)

const (
	SendRate          = time.Second
	MaxMessageSize    = 160
	ConcatSMSPartSize = 153 // CSMS is to split the message into 153 7-bit character parts (134 octets)
)

type SMSProvider interface {
	Client() *messagebird.Client
	CreateMessage(SMSMessage) (*messagebird.Message, error)
}

type MessageBirdSMSProvider struct{ ApiKey string }

// Client returns the MessageBird SDK client.
func (mb MessageBirdSMSProvider) Client() *messagebird.Client {
	return messagebird.New(mb.ApiKey)
}

// CreateMessage will create a new message on the MessageBird server.
func (mb MessageBirdSMSProvider) CreateMessage(sms SMSMessage) (*messagebird.Message, error) {
	params := &messagebird.MessageParams{}
	if len(sms.UDHHeader) > 0 {
		params.Type = "binary"
		params.TypeDetails = messagebird.TypeDetails{"udh": sms.UDHHeader}
	}
	return mb.Client().NewMessage(sms.Originator, []string{sms.Recipient}, sms.Message, params)
}

// Parrot contains the logic for sending SMS messages via an SMSProvider. It uses a worker
// goroutine that receives SMSMessage via a channel.
type Parrot struct {
	mbp  SMSProvider
	cfg  *Config
	Work chan SMSMessage
}

// NewParrot is a convenience function for instantiating a Parrot instance along with its configuration,
// SMSProvider interface and a channel for sending SMSMessages.
func NewParrot(filename string) (*Parrot, error) {
	if filename == "" {
		return nil, errors.New("missing filename")
	}
	cfg, err := readConfig(filename) // e.g."config.json"
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s\n", err)
	}
	return &Parrot{
		mbp:  MessageBirdSMSProvider{cfg.ApiKey},
		cfg:  &cfg,
		Work: make(chan SMSMessage),
	}, nil
}

// prepareConcatenatedSMSs will prepare a number of binary SMSMessages with a UDH header, ready to
// be sent as CSMS.
func (p Parrot) prepareConcatenatedSMSs(recipient, originator, message string) ([]SMSMessage, error) {
	var smsMsgs []SMSMessage
	if recipient == "" || originator == "" || message == "" {
		return nil, errors.New("SMS arguments cannot be empty strings")
	}
	msgParts, err := splitMessageIntoParts(message, ConcatSMSPartSize)
	if err != nil {
		return nil, err
	}
	nParts := len(msgParts)
	udhRef := generateUDHRef()

	for i, msgPart := range msgParts {
		udhHeader, err := generateUDHHeader(i, nParts, udhRef)
		if err != nil {
			return nil, err
		}

		msg := SMSMessage{
			Recipient:  recipient,
			Originator: originator,
			Message:    fmt.Sprintf("%x", msgPart),
			UDHHeader:  udhHeader,
		}
		smsMsgs = append(smsMsgs, msg)
	}
	return smsMsgs, nil
}

// PrepareSMS will prepare oneor more SMSMessages to the sender worker channel. It will break up
// a long message into concatenated SMSs if necessary. The actual sending of messages sends to be
// done by pushing SMSMessages to parrot's worker channel.
func (p Parrot) PrepareSMS(recipient, originator, message string) ([]SMSMessage, error) {
	if recipient == "" || originator == "" || message == "" {
		return nil, errors.New("SMS arguments cannot be empty strings")
	}
	if len(message) > MaxMessageSize {
		smsMsgs, err := p.prepareConcatenatedSMSs(recipient, originator, message)
		if err != nil {
			return nil, err
		}
		return smsMsgs, nil
	} else {
		msg := SMSMessage{Recipient: recipient, Originator: originator, Message: message}
		return []SMSMessage{msg}, nil
	}
}

// sendSMS creates a new message using the SMSProvider.
func (p Parrot) sendSMS(sms SMSMessage) error {
	m, err := p.mbp.CreateMessage(sms)
	if err != nil {
		// messagebird.ErrResponse means custom JSON errors.
		if err == messagebird.ErrResponse {
			for _, mbError := range m.Errors {
				fmt.Printf("New SMS Error:\n%#v\n", mbError)
			}
		}
	} else {
		fmt.Printf("New SMS Success:\n%+v\n", m)
	}
	return err
}

// StartSending listens to a channel for new SMSMessages to be dispatched and sends them. It then waits for a period
// to provide for the (theoretical) throughput to MessageBird.
func (p Parrot) StartSending() {
	for {
		select {
		case sms := <-p.Work:
			p.sendSMS(sms)
			time.Sleep(SendRate)
		}
	}
}
