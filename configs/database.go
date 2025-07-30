package configs

import (
    "fmt"
    "log"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func Initialize(databaseURL string) (*sqlx.DB, error) {
    log.Printf("Connecting to database...")
    
    db, err := sqlx.Connect("postgres", databaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    // Set connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(time.Hour)

    log.Println("Successfully connected to database")
    return db, nil
}

// Health check function
func HealthCheck(db *sqlx.DB) error {
    var result int
    err := db.Get(&result, "SELECT 1")
    if err != nil {
        return fmt.Errorf("database health check failed: %w", err)
    }
    return nil
}