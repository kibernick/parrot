package parrot

import "fmt"

type UserRequest struct {
	Recipient  int    `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
}

type UserRequestErrors struct {
	Recipient  string
	Originator string
	Message    string
}

func (e UserRequestErrors) Error() string {
	return fmt.Sprintf("%#v", e)
}

func (r UserRequest) Validate() error {
	var err = UserRequestErrors{}

	if r.Recipient == 0 {
		err.Recipient = "must be a non-zero integer"
	}
	if r.Originator == "" {
		err.Originator = "must be a non-empty string"
	}
	if r.Message == "" {
		err.Message = "must be a non-empty string"
	}

	if (UserRequestErrors{}) != err {
		return err
	}
	return nil
}
