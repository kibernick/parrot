package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/kibernick/parrot/parrot"
)

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	var reqData parrot.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := reqData.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := parrot.NewSMSMessage(reqData.Recipient, reqData.Originator, reqData.Message)
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
