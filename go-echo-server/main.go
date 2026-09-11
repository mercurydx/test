package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// Счётчик запросов — потокобезопасный
var requestCounter int64

// Список случайных фраз
var phrases = []string{
	"🚀 Полетели!",
	"🎉 Отличный день для кода!",
	"☕ Время для кофе и curl",
	"🔥 Сервер горяч, как мой процессор",
	"🌌 Бесконечность не предел",
	"🎸 Rock'n'roll is alive",
	"🐳 Docker — наше всё",
	"💡 Идея родилась в душе",
}

type Response struct {
	RequestID string `json:"request_id"`
	RequestNo int64  `json:"request_no"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Атомарно увеличиваем счётчик (потокобезопасно!)
	no := atomic.AddInt64(&requestCounter, 1)

	// Генерируем уникальный ID
	reqID := uuid.New().String()

	// Выбираем случайную фразу
	rand.Seed(time.Now().UnixNano())
	phrase := phrases[rand.Intn(len(phrases))]

	resp := Response{
		RequestID: reqID,
		RequestNo: no,
		Message:   phrase,
		Timestamp: time.Now().Format(time.RFC3339Nano),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/", handler)

	port := ":8080"
	fmt.Printf("🚀 Сервер запущен на http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
