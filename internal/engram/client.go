package engram

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Client provides access to Engram memory store
type Client struct {
	db     *sql.DB
	dbPath string
}

// Observation represents a memory observation
type Observation struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Scope       string    `json:"scope"`
	TopicKey    string    `json:"topic_key,omitempty"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SearchResult represents a search result
type SearchResult struct {
	Observation
	Rank float64 `json:"rank"`
}

// NewClient creates a new Engram client
func NewClient() (*Client, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home dir: %w", err)
	}

	dbPath := filepath.Join(home, ".ayrton", "engram.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	c := &Client{db: db, dbPath: dbPath}
	if err := c.initSchema(); err != nil {
		return nil, fmt.Errorf("init schema: %w", err)
	}

	return c, nil
}

// initSchema initializes the database schema with FTS5
func (c *Client) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS observations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		type TEXT NOT NULL,
		scope TEXT NOT NULL DEFAULT 'project',
		topic_key TEXT,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_observations_topic ON observations(topic_key);
	CREATE INDEX IF NOT EXISTS idx_observations_type ON observations(type);
	CREATE INDEX IF NOT EXISTS idx_observations_scope ON observations(scope);
	CREATE INDEX IF NOT EXISTS idx_observations_created ON observations(created_at);

	-- FTS5 virtual table for full-text search
	CREATE VIRTUAL TABLE IF NOT EXISTS observations_fts USING fts5(
		title, content, type, scope, topic_key,
		content='observations', content_rowid='id'
	);

	-- Triggers to keep FTS in sync
	CREATE TRIGGER IF NOT EXISTS observations_ai AFTER INSERT ON observations BEGIN
		INSERT INTO observations_fts(rowid, title, content, type, scope, topic_key)
		VALUES (new.id, new.title, new.content, new.type, new.scope, new.topic_key);
	END;

	CREATE TRIGGER IF NOT EXISTS observations_ad AFTER DELETE ON observations BEGIN
		INSERT INTO observations_fts(observations_fts, rowid, title, content, type, scope, topic_key)
		VALUES ('delete', old.id, old.title, old.content, old.type, old.scope, old.topic_key);
	END;

	CREATE TRIGGER IF NOT EXISTS observations_au AFTER UPDATE ON observations BEGIN
		INSERT INTO observations_fts(observations_fts, rowid, title, content, type, scope, topic_key)
		VALUES ('delete', old.id, old.title, old.content, old.type, old.scope, old.topic_key);
		INSERT INTO observations_fts(rowid, title, content, type, scope, topic_key)
		VALUES (new.id, new.title, new.content, new.type, new.scope, new.topic_key);
	END;
	`

	_, err := c.db.Exec(schema)
	return err
}

// Save saves an observation
func (c *Client) Save(ctx context.Context, obs *Observation) (int64, error) {
	now := time.Now()
	obs.CreatedAt = now
	obs.UpdatedAt = now

	result, err := c.db.ExecContext(ctx, `
		INSERT INTO observations (title, type, scope, topic_key, content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, obs.Title, obs.Type, obs.Scope, obs.TopicKey, obs.Content, obs.CreatedAt, obs.UpdatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	obs.ID = id
	return id, nil
}

// Update updates an existing observation
func (c *Client) Update(ctx context.Context, obs *Observation) error {
	obs.UpdatedAt = time.Now()
	_, err := c.db.ExecContext(ctx, `
		UPDATE observations SET title=?, type=?, scope=?, topic_key=?, content=?, updated_at=?
		WHERE id=?
	`, obs.Title, obs.Type, obs.Scope, obs.TopicKey, obs.Content, obs.UpdatedAt, obs.ID)
	return err
}

// SaveOrUpdate saves or updates an observation by topic_key (upsert)
func (c *Client) SaveOrUpdate(ctx context.Context, obs *Observation) (int64, error) {
	if obs.TopicKey == "" {
		return c.Save(ctx, obs)
	}

	// Check if exists
	var existing Observation
	err := c.db.QueryRowContext(ctx, `
		SELECT id, title, type, scope, topic_key, content, created_at, updated_at
		FROM observations WHERE topic_key=? AND scope=?
	`, obs.TopicKey, obs.Scope).Scan(
		&existing.ID, &existing.Title, &existing.Type, &existing.Scope,
		&existing.TopicKey, &existing.Content, &existing.CreatedAt, &existing.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.Save(ctx, obs)
	}
	if err != nil {
		return 0, err
	}

	// Update existing
	existing.Title = obs.Title
	existing.Type = obs.Type
	existing.Content = obs.Content
	existing.UpdatedAt = time.Now()
	return existing.ID, c.Update(ctx, &existing)
}

// Get retrieves an observation by ID
func (c *Client) Get(ctx context.Context, id int64) (*Observation, error) {
	obs := &Observation{}
	err := c.db.QueryRowContext(ctx, `
		SELECT id, title, type, scope, topic_key, content, created_at, updated_at
		FROM observations WHERE id=?
	`, id).Scan(&obs.ID, &obs.Title, &obs.Type, &obs.Scope, &obs.TopicKey, &obs.Content, &obs.CreatedAt, &obs.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return obs, err
}

// Search performs full-text search using FTS5
func (c *Client) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	// Sanitize query for FTS5 - escape special characters
	query = sanitizeFTS5Query(query)

	rows, err := c.db.QueryContext(ctx, `
		SELECT o.id, o.title, o.type, o.scope, o.topic_key, o.content, o.created_at, o.updated_at,
		       bm25(observations_fts) as rank
		FROM observations_fts
		JOIN observations o ON o.id = observations_fts.rowid
		WHERE observations_fts MATCH ?
		ORDER BY rank
		LIMIT ?
	`, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var r SearchResult
		err := rows.Scan(&r.ID, &r.Title, &r.Type, &r.Scope, &r.TopicKey, &r.Content, &r.CreatedAt, &r.UpdatedAt, &r.Rank)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

// ListByTopic retrieves observations by topic_key
func (c *Client) ListByTopic(ctx context.Context, topicKey, scope string) ([]Observation, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT id, title, type, scope, topic_key, content, created_at, updated_at
		FROM observations WHERE topic_key=? AND scope=?
		ORDER BY created_at DESC
	`, topicKey, scope)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Observation
	for rows.Next() {
		var o Observation
		err := rows.Scan(&o.ID, &o.Title, &o.Type, &o.Scope, &o.TopicKey, &o.Content, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, o)
	}
	return results, rows.Err()
}

// ListRecent retrieves recent observations
func (c *Client) ListRecent(ctx context.Context, scope string, limit int) ([]Observation, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT id, title, type, scope, topic_key, content, created_at, updated_at
		FROM observations WHERE scope=?
		ORDER BY created_at DESC
		LIMIT ?
	`, scope, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Observation
	for rows.Next() {
		var o Observation
		err := rows.Scan(&o.ID, &o.Title, &o.Type, &o.Scope, &o.TopicKey, &o.Content, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, o)
	}
	return results, rows.Err()
}

// Close closes the database connection
func (c *Client) Close() error {
	return c.db.Close()
}

// sanitizeFTS5Query sanitizes user input for FTS5 MATCH queries
func sanitizeFTS5Query(query string) string {
	// FTS5 special characters that need escaping
	special := []string{`"`, `'`, `-`, `+`, `*`, `(`, `)`, `:`}
	result := query
	for _, char := range special {
		result = strings.ReplaceAll(result, char, " ")
	}
	// Collapse multiple spaces
	result = strings.Join(strings.Fields(result), " ")
	// Wrap in quotes for phrase search if multiple terms
	if strings.Contains(result, " ") {
		return `"` + result + `"`
	}
	return result
}

// ToJSON serializes observation to JSON
func (o *Observation) ToJSON() (string, error) {
	b, err := json.MarshalIndent(o, "", "  ")
	return string(b), err
}