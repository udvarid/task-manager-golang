package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Mail_psw      string `json:"mail_psw"`
	Mail_from     string `json:"mail_from"`
	FlyIo         string `json:"flyIo"`
	Environment   string `json:"environment"`
	RemoteAddress string `json:"remote_address"`
}

var configuration = Configuration{}

func InitConfiguration(configFile string) Configuration {
	file, _ := os.Open(configFile)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}
