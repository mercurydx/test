package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

type Response struct {
	RequestID string `json:"request_id"`
	Quote     string `json:"quote"`
	Timestamp string `json:"timestamp"`
}

func initDB() {
	// Берем настройки из переменных окружения (с дефолтными значениями)
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "quotes_db")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	dbPool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	// Проверка соединения
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Не удалось пингануть БД: %v", err)
	}
	log.Println("✅ Успешное подключение к PostgreSQL")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func handler(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.New().String()

	var quote string
	// Берем случайную цитату из БД
	err := dbPool.QueryRow(context.Background(), "SELECT text FROM quotes ORDER BY RANDOM() LIMIT 1").Scan(&quote)
	if err != nil {
		http.Error(w, "Ошибка получения цитаты из БД", http.StatusInternalServerError)
		log.Printf("Ошибка БД: %v", err)
		return
	}

	resp := Response{
		RequestID: reqID,
		Quote:     quote,
		Timestamp: time.Now().Format(time.RFC3339Nano),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	initDB()
	defer dbPool.Close()

	http.HandleFunc("/", handler)

	port := getEnv("PORT", "3000")
	fmt.Printf("🚀 Сервер запущен на http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
