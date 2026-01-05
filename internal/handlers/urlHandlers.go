package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "github.com/Vhyron/url-shortener/internal/repository"
    "github.com/Vhyron/url-shortener/internal/utils"
    "github.com/gorilla/mux"
)

type URLHandler struct {
    repo    *repository.URLRepository
    baseURL string
}

func NewURLHandler(repo *repository.URLRepository, baseURL string) *URLHandler {
    return &URLHandler{repo: repo, baseURL: baseURL}
}

type CreateURLRequest struct {
    URL string `json:"url"`
}

type CreateURLResponse struct {
    ShortCode string `json:"short_code"`
    ShortURL  string `json:"short_url"`
    LongURL   string `json:"long_url"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, err := json.Marshal(payload)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal Server Error"))
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, ErrorResponse{Error: message})
}

func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
    var req CreateURLRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    if req.URL == "" {
        respondWithError(w, http.StatusBadRequest, "URL is required")
        return
    }

    if _, err := url.ParseRequestURI(req.URL); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid URL format")
        return
    }

    var shortCode string
    var err error
    maxAttempts := 5
    for i := 0; i < maxAttempts; i++ {
        shortCode, err = utils.GenerateShortCode()
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Failed to generate short code")
            return
        }
        existing, err := h.repo.GetByShortCode(shortCode)
        if err != nil {
            respondWithError(w, http.StatusInternalServerError, "Database error")
            return
        }
        if existing == nil {
            break
        }
    }

    _, err = h.repo.Create(shortCode, req.URL)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to create short URL")
        return
    }

    shortURL := fmt.Sprintf("%s/%s", h.baseURL, shortCode)
    respondWithJSON(w, http.StatusCreated, CreateURLResponse{
        ShortCode: shortCode,
        ShortURL:  shortURL,
        LongURL:   req.URL,
    })
}

func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    shortCode := vars["shortCode"]

    urlData, err := h.repo.GetByShortCode(shortCode)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Database error")
        return
    }

    if urlData == nil {
        respondWithError(w, http.StatusNotFound, "Short URL not found")
        return
    }

    if err := h.repo.IncrementClicks(shortCode); err != nil {
        fmt.Printf("Failed to increment clicks: %v\n", err)
    }

    http.Redirect(w, r, urlData.OriginalURL, http.StatusMovedPermanently)
}

func (h *URLHandler) GetAllURLs(w http.ResponseWriter, r *http.Request) {
    urls, err := h.repo.GetAll()
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to fetch URLs")
        return
    }
    respondWithJSON(w, http.StatusOK, urls)
}

func (h *URLHandler) GetURLStats(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    shortCode := vars["shortCode"]

    urlData, err := h.repo.GetByShortCode(shortCode)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Database error")
        return
    }

    if urlData == nil {
        respondWithError(w, http.StatusNotFound, "Short URL not found")
        return
    }

    respondWithJSON(w, http.StatusOK, urlData)
}