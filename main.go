package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"bitbucket.org/kibernick/parrot/parrot"
)

// errToMap is a cheap way to turn an error into a JSONable object.
func errToMap(e error) map[string]interface{} {
	return map[string]interface{}{"error": e.Error()}
}

// jsonResponse returns pretty JSON errors.
func jsonResponse(w http.ResponseWriter, payload map[string]interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "YA DUN GOOFED", http.StatusInternalServerError)
	}
}

// SendSMSHandler accepts SMS messages submitted via a POST request containing a JSON object.
func SendSMSHandler(w http.ResponseWriter, r *http.Request, p *parrot.Parrot) {
	if r.Method != "POST" {
		err := fmt.Errorf("Method %s not allowed.", r.Method)
		jsonResponse(w, errToMap(err), http.StatusMethodNotAllowed)
		return
	}

	var req parrot.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, errToMap(err), http.StatusBadRequest)
		return
	}
	if valErrors := req.Validate(); len(valErrors) != 0 {
		jsonResponse(w, valErrors, http.StatusBadRequest)
		return
	}

	n, err := p.PrepareSMS(req.Recipient, req.Originator, req.Message)
	if err != nil {
		jsonResponse(w, errToMap(err), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{"message": "SMS(s) sent!", "sent": n}
	jsonResponse(w, response, http.StatusCreated)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	parrot := parrot.NewParrot("config.json")
	go parrot.StartSending()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		SendSMSHandler(w, r, parrot)
	})
	fmt.Println("Parrot landed on port 8000...")
	http.ListenAndServe(":8000", nil)
}
