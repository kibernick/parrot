package parrot

import (
	"encoding/json"
	"os"
)

type Config struct {
	ApiKey string
}

// ReadConfig will read config JSON from the given filepath.
func ReadConfig(file string) (config Config, err error) {
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

var readConfig = ReadConfig
