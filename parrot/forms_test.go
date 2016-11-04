package parrot

import "testing"

type forms_testpair struct {
	req UserRequest
	err interface{}
}

var forms_tests = []forms_testpair{
	{
		UserRequest{Recipient: 1234, Originator: "Ground Control", Message: "Commencing countdown"},
		nil,
	},
	{
		UserRequest{Originator: "Ground Control", Message: "Commencing countdown"},
		ValidationErrors{Recipient: errMsgInt},
	},
	{
		UserRequest{Recipient: 999},
		ValidationErrors{Originator: errMsgStr, Message: errMsgStr},
	},
}

func TestUserRequest_Validate(t *testing.T) {
	for _, pair := range forms_tests {
		err := pair.req.Validate()
		if err != pair.err {
			t.Errorf("Got: %s, Expected: %s", err, pair.err)
		}
	}
}
