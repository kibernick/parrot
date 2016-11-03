package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kibernick/parrot/parrot"
)

type userRequest struct {
	Recipient  int    `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
}

type userRequestErrors struct {
	Recipient  string
	Originator string
	Message    string
}

func (e userRequestErrors) Error() string {
	return fmt.Sprintf("%#v", e)
}

func (r userRequest) Validate() error {
	var err = userRequestErrors{}

	if r.Recipient == 0 {
		err.Recipient = "must be a non-zero integer"
	}
	if r.Originator == "" {
		err.Originator = "must be a non-empty string"
	}
	if r.Message == "" {
		err.Message = "must be a non-empty string"
	}

	if (userRequestErrors{}) != err {
		return err
	}
	return nil
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	var reqData userRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := reqData.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := parrot.NewSMSMessage(reqData.Recipient, reqData.Originator, reqData.Message)
	//fmt.Fprintf(w, "%s\n", m)
	//fmt.Println(m)
	res, err := m.Send()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s\n", res)
}

func main() {
	http.HandleFunc("/", postHandler)

	fmt.Println("Now serving on port 8000: parrot!")
	http.ListenAndServe(":8000", nil)
}
