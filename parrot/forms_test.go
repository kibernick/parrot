package parrot

import (
	"reflect"
	"testing"
)

type forms_testargs struct {
	req UserRequest
	err interface{}
}

var forms_tests = []forms_testargs{
	{
		UserRequest{Recipient: "1234", Originator: "+3166666", Message: "Wubalubadubdub"},
		map[string]interface{}{},
	},
	{
		UserRequest{Originator: "Ground Control", Message: "Commencing countdown"},
		map[string]interface{}{"Recipient": errMsgStr},
	},
	{
		UserRequest{Recipient: "999"},
		map[string]interface{}{"Originator": errMsgStr, "Message": errMsgStr},
	},
}

func TestUserRequest_Validate(t *testing.T) {
	for _, args := range forms_tests {
		err := args.req.Validate()

		if !reflect.DeepEqual(err, args.err) {
			t.Errorf("got: %s, expected: %s", err, args.err)
		}
	}
}
