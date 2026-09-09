package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Ответ для корневого маршрута
type RootResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	RemoteIP  string `json:"remote_ip"`
}

// Ответ для /echo
type EchoResponse struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	resp := RootResponse{
		Message:   "Привет! Я Go-сервер из GitHub 🚀",
		Timestamp: time.Now().Format(time.RFC3339),
		RemoteIP:  r.RemoteAddr,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	bodyBytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	resp := EchoResponse{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: r.Header,
		Body:    string(bodyBytes),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/echo", echoHandler)

	port := ":3000"
	fmt.Printf("🚀 Сервер запущен на http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
