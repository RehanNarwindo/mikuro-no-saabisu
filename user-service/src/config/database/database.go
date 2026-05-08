package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func DatabaseConfig() error {
	return DatabaseConfigWithContext(context.Background())
}

func DatabaseConfigWithContext(ctx context.Context) error {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system env")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || port == "" || user == "" || dbname == "" {
		return errors.New("missing required database configuration")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		errChan <- db.PingContext(ctxWithTimeout)
	}()

	select {
	case <-ctxWithTimeout.Done():
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Warning: error closing database connection: %v", closeErr)
		}
		return fmt.Errorf("database connection timeout after 5 seconds: %w", ctxWithTimeout.Err())
	case err := <-errChan:
		if err != nil {
			if closeErr := db.Close(); closeErr != nil {
				log.Printf("Warning: error closing database connection: %v", closeErr)
			}
			return fmt.Errorf("database not accessible: %w", err)
		}
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	fmt.Println("===  Database connected successfully ===")
	DB = db

	return nil
}

func DatabaseConfigWithRetry(maxRetries int) error {
	return DatabaseConfigWithRetryContext(context.Background(), maxRetries)
}

func DatabaseConfigWithRetryContext(ctx context.Context, maxRetries int) error {
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Attempting database connection (%d/%d)...", attempt, maxRetries)

		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

		err := DatabaseConfigWithContext(attemptCtx)
		cancel()

		if err == nil {
			if attempt > 1 {
				log.Printf("✅ Database connected successfully on attempt %d", attempt)
			}
			return nil
		}

		lastErr = err
		log.Printf("❌ Attempt %d failed: %v", attempt, err)

		if attempt < maxRetries {
			backoff := time.Duration(attempt*attempt) * time.Second
			log.Printf("Retrying in %v...", backoff)

			select {
			case <-time.After(backoff):
				continue
			case <-ctx.Done():
				return fmt.Errorf("context canceled: %w", ctx.Err())
			}
		}
	}

	return fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, lastErr)
}

func PingWithTimeout(timeout time.Duration) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return DB.PingContext(ctx)
}

func HealthCheck(ctx context.Context) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return DB.PingContext(ctx)
}

func QueryWithTimeout(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if DB == nil {
		return nil, errors.New("database not initialized")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}

	return DB.QueryContext(ctx, query, args...)
}

func ExecWithTimeout(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if DB == nil {
		return nil, errors.New("database not initialized")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	return DB.ExecContext(ctx, query, args...)
}

func CloseDatabase() error {
	if DB != nil {
		log.Println("Closing database connection...")
		return DB.Close()
	}
	return nil
}
