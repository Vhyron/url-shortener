package repository

import (
    "database/sql"
    "time"
    "github.com/Vhyron/url-shortener/internal/models"
    _ "github.com/mattn/go-sqlite3"
)

type URLRepository struct {
    db *sql.DB
}

func NewURLRepository(dbPath string) (*URLRepository, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    query := `
    CREATE TABLE IF NOT EXISTS urls (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        short_code TEXT UNIQUE NOT NULL,
        original_url TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        clicks INTEGER DEFAULT 0
    );
    CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
    `

    if _, err := db.Exec(query); err != nil {
        return nil, err
    }

    return &URLRepository{db: db}, nil
}

func (r *URLRepository) Create(shortCode, originalURL string) (*models.URL, error) {
    query := `INSERT INTO urls (short_code, original_url, created_at) VALUES (?, ?, ?)`
    now := time.Now()
    result, err := r.db.Exec(query, shortCode, originalURL, now)
    if err != nil {
        return nil, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }

    return &models.URL{
        ID:          int(id),
        ShortCode:   shortCode,
        OriginalURL: originalURL,
        CreatedAt:   now,
        Clicks:      0,
    }, nil
}

func (r *URLRepository) GetByShortCode(shortCode string) (*models.URL, error) {
    query := `SELECT id, short_code, original_url, created_at, clicks FROM urls WHERE short_code = ?`
    var url models.URL
    err := r.db.QueryRow(query, shortCode).Scan(
        &url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt, &url.Clicks,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &url, nil
}

func (r *URLRepository) IncrementClicks(shortCode string) error {
    query := `UPDATE urls SET clicks = clicks + 1 WHERE short_code = ?`
    _, err := r.db.Exec(query, shortCode)
    return err
}

func (r *URLRepository) GetAll() ([]models.URL, error) {
    query := `SELECT id, short_code, original_url, created_at, clicks FROM urls ORDER BY created_at DESC`
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var urls []models.URL
    for rows.Next() {
        var url models.URL
        if err := rows.Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt, &url.Clicks); err != nil {
            return nil, err
        }
        urls = append(urls, url)
    }
    return urls, nil
}

func (r *URLRepository) Close() error {
    return r.db.Close()
}