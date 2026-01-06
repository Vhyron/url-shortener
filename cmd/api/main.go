package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/Vhyron/url-shortener/internal/handlers"
    "github.com/Vhyron/url-shortener/internal/repository"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    port := getEnv("PORT", "8080")
    dbPath := getEnv("DB_PATH", "./urls.db")
    baseURL := getEnv("BASE_URL", "http://localhost:8080")

    repo, err := repository.NewURLRepository(dbPath)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer repo.Close()

    urlHandler := handlers.NewURLHandler(repo, baseURL)
    router := mux.NewRouter()

    api := router.PathPrefix("/api").Subrouter()
    api.HandleFunc("/shorten", urlHandler.CreateShortURL).Methods("POST")
    api.HandleFunc("/urls", urlHandler.GetAllURLs).Methods("GET")
    api.HandleFunc("/stats/{shortCode}", urlHandler.GetURLStats).Methods("GET")

    router.HandleFunc("/{shortCode}", urlHandler.RedirectURL).Methods("GET")
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }).Methods("GET")

    addr := fmt.Sprintf(":%s", port)
    fmt.Printf("🚀 URL Shortener running on http://localhost:%s\n", port)
    fmt.Printf("📊 Database: %s\n", dbPath)
    fmt.Printf("🔗 Base URL: %s\n", baseURL)
    
    log.Fatal(http.ListenAndServe(addr, router))
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}