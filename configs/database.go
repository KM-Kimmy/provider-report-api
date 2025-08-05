package configs

import (
    "database/sql"
    "fmt"
    "log"
    "time"
	"sync"

    "provider-report-api/pkg/vault"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

var db *sql.DB
var once sync.Once

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

// GetDB returns a singleton DB instance
func GetDB() *sql.DB {
	once.Do(func() {
		creds := vault.GetDatabaseSecret()
		if creds == nil {
			log.Fatalf("Secret data type is not expected.")
		}
		var err error
		connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
			creds.DatabaseUrl, creds.DatabaseUsername, creds.DatabasePassword, creds.DatabasePort, creds.DatabaseName)

		db, err = sql.Open("mssql", connString)
		if err != nil {
			log.Fatalf("Error creating database connection: %v", err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatalf("Error pinging database: %v", err)
		}
	})
	return db
}