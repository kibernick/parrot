package parrot

import (
	"fmt"
	"os"
	"time"

	"github.com/messagebird/go-rest-api"
)

const (
	MaxMessageSize    = 160
	ConcatSMSPartSize = 153 // CSMS is to split the message into 153 7-bit character parts (134 octets)
)

type Parrot struct {
	client *messagebird.Client
	cfg    *Config
	chn    chan SMSMessage
}

func NewParrot(file string) *Parrot {
	cfg, err := ReadConfig(file) // "config.json"
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}
	return &Parrot{
		client: messagebird.New(cfg.ApiKey),
		cfg:    &cfg,
		chn:    make(chan SMSMessage),
	}
}

func (p Parrot) PrepareSMS(recipient, originator, message string) (int, error) {
	if len(message) > MaxMessageSize {
		msgParts, err := splitMessageIntoParts(message, ConcatSMSPartSize)
		if err != nil {
			return 0, err
		}
		nParts := len(msgParts)
		udhRef := generateUDHRef()

		for i, msgPart := range msgParts {
			udhHeader, err := generateUDHHeader(i, nParts, udhRef)
			if err != nil {
				return 0, err
			}

			msg := SMSMessage{
				Recipient:  recipient,
				Originator: originator,
				Message:    fmt.Sprintf("%x", msgPart),
				UDHHeader:  udhHeader,
			}
			p.chn <- msg
		}
		return len(msgParts), nil
	} else {
		msg := SMSMessage{Recipient: recipient, Originator: originator, Message: message}
		p.chn <- msg
		return 1, nil
	}
}

func (p Parrot) sendSMS(sms SMSMessage) error {
	params := &messagebird.MessageParams{}
	if len(sms.UDHHeader) > 0 {
		params.Type = "binary"
		params.TypeDetails = messagebird.TypeDetails{"udh": sms.UDHHeader}
		fmt.Printf("params:\n%+v\n", params)
	}
	m, err := p.client.NewMessage(sms.Originator, []string{sms.Recipient}, sms.Message, params)
	if err != nil {
		// messagebird.ErrResponse means custom JSON errors.
		if err == messagebird.ErrResponse {
			for _, mbError := range m.Errors {
				// todo remove

				fmt.Printf("SMS Error:\n%#v\n", mbError)
			}
		}
	} else {
		fmt.Printf("SMS Success:\n%+v\n", m)
	}
	return err
}

func (p Parrot) StartSending() {
	// TODO: improvement - split off a coroutine that would acutally send the SMS
	for {
		select {
		case sms := <-p.chn:
			fmt.Printf("%+v\n", sms)
			p.sendSMS(sms)
			time.Sleep(time.Second)
		}

		//<-p.chn
		////sms := <-p.chn
		//fmt.Println("I GOT A JOB!!!")
		////fmt.Println(sms)
		////p.sendSMS(sms)
		//time.Sleep(time.Second)
	}
}
