package parrot

import (
	"errors"
	"testing"

	"strings"

	"github.com/messagebird/go-rest-api"
)

func TestNewParrot(t *testing.T) {
	var tests = []struct {
		filename string
		wantErr  bool
	}{
		{"", true},
		{"validpath", false},
		{"/dev/null", true},
	}

	saved := readConfig
	defer func() { readConfig = saved }()

	readConfig = func(file string) (config Config, err error) {
		if file == "" || file != "validpath" {
			return Config{}, errors.New("!")
		}
		return Config{ApiKey: "12345"}, nil
	}

	for _, test := range tests {
		_, err := NewParrot(test.filename)
		if err != nil && !test.wantErr {
			t.Fatalf("error raised and not wanted for file %s", test.filename)
		}
		if err == nil && test.wantErr {
			t.Fatalf("error not raised for filename %q", test.filename)
		}
	}

}

type FakeSMSProvider struct{ Err error }

func (f FakeSMSProvider) Client() *messagebird.Client { return &messagebird.Client{} }

func (f FakeSMSProvider) CreateMessage(sms SMSMessage) (*messagebird.Message, error) {
	if f.Err != nil {
		if f.Err == messagebird.ErrResponse {
			err := messagebird.Error{}
			errMsg := messagebird.Message{}
			errMsg.Errors = append(errMsg.Errors, err)
			return &errMsg, f.Err
		} else {
			return nil, f.Err
		}
	}
	return &messagebird.Message{}, nil
}

func TestParrot_sendSMS(t *testing.T) {
	var tests = []struct {
		sms SMSMessage
		fp  FakeSMSProvider
		err error
	}{
		{
			SMSMessage{"123", "123", "Dobar dan!", ""},
			FakeSMSProvider{},
			nil,
		},
		{
			SMSMessage{},
			FakeSMSProvider{messagebird.ErrResponse},
			messagebird.ErrResponse,
		},
		{
			SMSMessage{},
			FakeSMSProvider{messagebird.ErrUnexpectedResponse},
			messagebird.ErrUnexpectedResponse,
		},
	}

	for _, test := range tests {
		fakeParrot := Parrot{mbp: test.fp}
		got := fakeParrot.sendSMS(test.sms)
		if got != test.err {
			t.Fatalf("expected %s but got %s", test.err, got)
		}
	}
}

func TestParrot_prepareConcatenatedSMSs(t *testing.T) {
	var tests = []struct {
		recipient  string
		originator string
		message    string
		expLen     int
		expErr     bool
	}{
		{"", "", "", 0, true},
		{"123", "123", "Yo", 1, false},
		{"123", "123", strings.Repeat("1234567890", 16) + "X", 2, false},
	}
	p := Parrot{}
	for i, test := range tests {
		msgs, err := p.prepareConcatenatedSMSs(test.recipient, test.originator, test.message)
		if err != nil && !test.expErr {
			t.Fatalf("error raised and not wanted for test #%v", i)
		}
		if err == nil && test.expErr {
			t.Fatalf("error not raised and wanted for test #%v", i)
		}
		if len(msgs) != test.expLen {
			t.Fatalf("expected %v messages but got %v", test.expLen, len(msgs))
		}
		for _, msg := range msgs {
			if msg.UDHHeader == "" {
				t.Fatal("UDH header was not set when preparing CSMS")
			}
		}
	}
}
