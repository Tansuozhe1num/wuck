package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type User struct {
	ID           string        `json:"id"`
	Username     string        `json:"username"`
	Email        string        `json:"email"`
	PasswordHash string        `json:"-"`
	ClickHistory []ClickRecord `json:"clickHistory"`
	CreatedAt    time.Time     `json:"createdAt"`
}

type ClickRecord struct {
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Source    string    `json:"source"`
	Category  string    `json:"category"`
	ClickedAt time.Time `json:"clickedAt"`
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(dbPath string) (*UserStore, error) {
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite serialized mode — one writer at a time

	s := &UserStore{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return s, nil
}

func (s *UserStore) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS clicks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			url TEXT NOT NULL,
			title TEXT NOT NULL,
			source TEXT NOT NULL,
			category TEXT NOT NULL,
			clicked_at TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_clicks_user ON clicks(user_id);
		CREATE INDEX IF NOT EXISTS idx_clicks_url ON clicks(user_id, url);
	`)
	return err
}

func (s *UserStore) Close() error {
	return s.db.Close()
}

func (s *UserStore) Create(username, email, passwordHash string) (*User, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	_, err = s.db.Exec(
		"INSERT INTO users (id, username, email, password_hash, created_at) VALUES (?, ?, ?, ?, ?)",
		id, username, email, passwordHash, now.Format(time.RFC3339),
	)
	if err != nil {
		if isUniqueConstraint(err) {
			return nil, errors.New("该邮箱已被注册")
		}
		return nil, err
	}

	return &User{
		ID:           id,
		Username:     username,
		Email:        email,
		ClickHistory: []ClickRecord{},
		CreatedAt:    now,
	}, nil
}

func (s *UserStore) FindByEmail(email string) *User {
	row := s.db.QueryRow("SELECT id, username, email, password_hash, created_at FROM users WHERE email = ?", email)
	return s.scanUser(row)
}

func (s *UserStore) FindByID(id string) *User {
	row := s.db.QueryRow("SELECT id, username, email, password_hash, created_at FROM users WHERE id = ?", id)
	user := s.scanUser(row)
	if user == nil {
		return nil
	}
	user.ClickHistory = s.getClicks(id)
	return user
}

func (s *UserStore) GetClicksByUserID(userID string) []ClickRecord {
	return s.getClicks(userID)
}

func (s *UserStore) AddClick(userID string, record ClickRecord) error {
	_, err := s.db.Exec(
		"INSERT INTO clicks (user_id, url, title, source, category, clicked_at) VALUES (?, ?, ?, ?, ?, ?)",
		userID, record.URL, record.Title, record.Source, record.Category, time.Now().Format(time.RFC3339),
	)
	return err
}

func (s *UserStore) GetClickedURLs(userID string) map[string]bool {
	rows, err := s.db.Query("SELECT DISTINCT url FROM clicks WHERE user_id = ?", userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	urls := make(map[string]bool)
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err == nil {
			urls[url] = true
		}
	}
	return urls
}

func (s *UserStore) GetCategoryPreference(userID string) map[string]float64 {
	rows, err := s.db.Query("SELECT category, COUNT(*) as cnt FROM clicks WHERE user_id = ? GROUP BY category", userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	prefs := make(map[string]float64)
	total := 0
	type cnt struct {
		cat string
		n   int
	}
	var counts []cnt
	for rows.Next() {
		var cat string
		var n int
		if err := rows.Scan(&cat, &n); err == nil {
			counts = append(counts, cnt{cat, n})
			total += n
		}
	}

	if total == 0 {
		return nil
	}

	for _, c := range counts {
		prefs[c.cat] = float64(c.n) / float64(total)
	}
	return prefs
}

// ── Internal helpers ─────────────────────────────

func (s *UserStore) scanUser(row *sql.Row) *User {
	var u User
	var createdAtStr string
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &createdAtStr)
	if err != nil {
		return nil
	}
	u.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	u.ClickHistory = []ClickRecord{}
	return &u
}

func (s *UserStore) getClicks(userID string) []ClickRecord {
	rows, err := s.db.Query(
		"SELECT url, title, source, category, clicked_at FROM clicks WHERE user_id = ? ORDER BY clicked_at DESC LIMIT 200",
		userID,
	)
	if err != nil {
		return []ClickRecord{}
	}
	defer rows.Close()

	var clicks []ClickRecord
	for rows.Next() {
		var c ClickRecord
		var t string
		if err := rows.Scan(&c.URL, &c.Title, &c.Source, &c.Category, &t); err == nil {
			c.ClickedAt, _ = time.Parse(time.RFC3339, t)
			clicks = append(clicks, c)
		}
	}
	if clicks == nil {
		clicks = []ClickRecord{}
	}
	return clicks
}

func isUniqueConstraint(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || contains(err.Error(), "UNIQUE constraint")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func generateID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
