package db

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DSN         string
	MaxRetries  int
	RetryDelay  time.Duration
	ConnTimeout time.Duration
}

// InitPostgres opens a gorm DB with retries and ping verification.
// Returns a ready-to-use *gorm.DB or an error.
func InitPostgres(cfg Config) (*gorm.DB, error) {
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 5
	}
	if cfg.RetryDelay <= 0 {
		cfg.RetryDelay = 3 * time.Second
	}
	if cfg.ConnTimeout <= 0 {
		cfg.ConnTimeout = 5 * time.Second
	}

	var db *gorm.DB
	var err error

	for attempt := 1; attempt <= cfg.MaxRetries; attempt++ {
		// Open driver (this is cheap). We still verify with ping below.
		db, err = gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
		if err == nil {
			// Verify connection with timeout
			sqlDB, sqlErr := db.DB()
			if sqlErr != nil {
				err = sqlErr
			} else {
				ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnTimeout)
				defer cancel()
				// Ping using context
				pingCh := make(chan error, 1)
				go func() {
					pingCh <- sqlDB.Ping()
				}()
				select {
				case perr := <-pingCh:
					if perr == nil {
						// success
						return db, nil
					}
					err = perr
				case <-ctx.Done():
					err = ctx.Err()
				}
			}
		}

		log.Printf("DB connect attempt %d/%d failed: %v — retrying in %s", attempt, cfg.MaxRetries, err, cfg.RetryDelay)
		time.Sleep(cfg.RetryDelay)
		// exponential backoff
		cfg.RetryDelay *= 2
	}

	// final check
	if err == nil {
		err = errors.New("failed to init db: unknown error")
	}
	return nil, err
}
