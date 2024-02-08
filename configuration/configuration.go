package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Mail_psw string `json:"mail_psw"`
}

var configuration = Configuration{}

func InitConfiguration() Configuration {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	fmt.Println(configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}
