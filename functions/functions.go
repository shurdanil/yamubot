package functions

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

var fields = []string{"Enter token:|token"}

func CreateConfig() (config struct {
	Token    string `json:"token"`
	Login    string `json:"login"`
	Password string `json:"password"`
}) {

	if !Exists("config.json") {
		for {
			var newAPIConfig = make(map[string]string)
			for _, field := range fields {
				fmt.Println(strings.Split(field, "|")[0])
				input := bufio.NewScanner(os.Stdin)
				input.Scan()
				newAPIConfig[strings.Split(field, "|")[1]] = input.Text()
			}

			data, _ := json.MarshalIndent(newAPIConfig, "", " ")
			err := os.WriteFile("config.json", data, 0644)
			if err != nil {
				continue
			}
			break
		}
	}

	byteValue, err := os.ReadFile("config.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(byteValue, &config)
	return
}
