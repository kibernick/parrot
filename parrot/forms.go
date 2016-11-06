package parrot

// UserRequest contains the parsed user input.
type UserRequest struct {
	Recipient  string `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
}

const (
	//errMsgInt = "must be a non-zero integer"
	errMsgStr = "must be a non-empty string"
)

// Validate will validate user input and return any errors that may have occurred.
func (r UserRequest) Validate() map[string]interface{} {
	valErrs := make(map[string]interface{})

	if r.Recipient == "" {
		valErrs["Recipient"] = errMsgStr
	}
	if r.Originator == "" {
		valErrs["Originator"] = errMsgStr
	}
	if r.Message == "" {
		valErrs["Message"] = errMsgStr
	}
	return valErrs
}
