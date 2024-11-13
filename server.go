package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

func send_cmd(w http.ResponseWriter, r *http.Request) {
	command := "whoami"
	obfusCMD := base64.StdEncoding.EncodeToString([]byte(command))

	fmt.Fprint(w, obfusCMD)
	fmt.Print("Server send command:", command)
}

func main() {
	http.HandleFunc("/getCMD", send_cmd)
	fmt.Println("Server Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error to server", err)
	}
}
