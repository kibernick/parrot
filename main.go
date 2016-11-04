package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"bitbucket.org/kibernick/parrot/parrot"
)

type ParrotConfig struct {
	ApiKey string
}

// readConfig will read config JSON from the given filepath.
func readConfig(file string) (config ParrotConfig, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&config); err != nil {
		return
	}
	return
}

// postHandler accepts SMS messages submitted via a POST request containing a JSON object.
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

	m := parrot.SMSMessage{
		Recipient:  reqData.Recipient,
		Originator: reqData.Originator,
		Message:    parrot.SplitMessageIntoParts(reqData.Message),
	}
	res, err := m.Send()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s\n", res)
}

func main() {
	_, err := readConfig("config.json")
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", postHandler)

	fmt.Println("Now serving on port 8000: parrot!")
	http.ListenAndServe(":8000", nil)
}
