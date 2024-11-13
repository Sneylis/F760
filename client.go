package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"

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
		}

		time.Sleep(10 * time.Second) // Пауза между запросами
	}
}
