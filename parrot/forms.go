package parrot

import "fmt"

// UserRequest contains the parsed user input.
type UserRequest struct {
	Recipient  int    `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
}

// ValidationErrors gives information on any errors that occurred when parsing user input.
type ValidationErrors struct {
	Recipient  string
	Originator string
	Message    string
}

const (
	errMsgInt = "must be a non-zero integer"
	errMsgStr = "must be a non-empty string"
)

func (e ValidationErrors) Error() string {
	return fmt.Sprintf("%#v", e)
}

// Validate will validate user input and return any errors that may have occurred.
func (r UserRequest) Validate() error {
	var err = ValidationErrors{}

	if r.Recipient == 0 {
		err.Recipient = errMsgInt
	}
	if r.Originator == "" {
		err.Originator = errMsgStr
	}
	if r.Message == "" {
		err.Message = errMsgStr
	}

	if (ValidationErrors{}) != err {
		return err
	}
	return nil
}
