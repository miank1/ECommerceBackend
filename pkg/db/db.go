package db

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

func InitDB(dsn string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // standard Go logger
		logger.Config{
			SlowThreshold:             time.Second,   // threshold for slow queries
			LogLevel:                  logger.Silent, // ✅ turn off all SQL logs
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
		if err == nil {
			log.Println("✅ Connected to DB")
			break
		}
		log.Printf("⏳ DB not ready, retrying... (%d/10)", i+1)
		time.Sleep(3 * time.Second)
	}
	if db == nil {
		log.Fatalf("❌ Could not connect to DB after 10 retries")
	}

	return db, nil
}
