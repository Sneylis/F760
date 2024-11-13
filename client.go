package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	response, err := http.Get("http://localhost:8080/getCMD")
	if err != nil {
		fmt.Println("Error to connect the server", err)

		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Eerror on Read Response", err)
		return
	}

	command, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		fmt.Println("Error to decode body CMD", err)
		return
	}

	fmt.Println("Received command from server:", string(command))
}
