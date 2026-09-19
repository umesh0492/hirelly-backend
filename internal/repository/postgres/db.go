package postgres

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"hirelly-backend/internal/config"
)

// Database manages connection pooling and Supabase REST fallback.
type Database struct {
	Pool        *pgxpool.Pool
	SupabaseURL string
	SupabaseKey string
	HTTPClient  *http.Client
}

// NewDatabase initializes PostgreSQL connection pool and HTTP client.
func NewDatabase(ctx context.Context, cfg *config.Config) (*Database, error) {
	db := &Database{
		SupabaseURL: cfg.SupabaseURL,
		SupabaseKey: cfg.SupabaseKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	if cfg.PostgresURL != "" {
		poolConfig, err := pgxpool.ParseConfig(cfg.PostgresURL)
		if err != nil {
			log.Printf("⚠️ Warning: could not parse POSTGRES_URL: %v. Will rely on Supabase REST.", err)
		} else {
			poolConfig.MaxConns = 15
			poolConfig.MinConns = 2
			poolConfig.MaxConnLifetime = 30 * time.Minute
			poolConfig.MaxConnIdleTime = 5 * time.Minute

			poolCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			pool, err := pgxpool.NewWithConfig(poolCtx, poolConfig)
			if err != nil {
				log.Printf("⚠️ Warning: could not initialize Postgres pool: %v. Falling back to Supabase REST.", err)
			} else if err := pool.Ping(poolCtx); err != nil {
				log.Printf("⚠️ Warning: Postgres ping failed: %v. Falling back to Supabase REST.", err)
				pool.Close()
			} else {
				log.Println("✅ PostgreSQL connection pool initialized and verified.")
				db.Pool = pool
			}
		}
	} else {
		log.Println("ℹ️ No POSTGRES_URL configured. Operating in Supabase REST Mode.")
	}

	return db, nil
}

// Close gracefully terminates open database pools.
func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}

// QuerySupabase executes a PostgREST API request against the Supabase instance.
func (d *Database) QuerySupabase(ctx context.Context, method, endpoint string, body any) ([]byte, error) {
	if d.SupabaseURL == "" || d.SupabaseKey == "" {
		return nil, fmt.Errorf("supabase URL or Key not configured")
	}

	url := fmt.Sprintf("%s/rest/v1/%s", strings.TrimRight(d.SupabaseURL, "/"), endpoint)
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal payload error: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("apikey", d.SupabaseKey)
	req.Header.Set("Authorization", "Bearer "+d.SupabaseKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http execute error: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("supabase status %d: %s", resp.StatusCode, string(respBytes))
	}

	return respBytes, nil
}
