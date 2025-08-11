package configs

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    "sync"

    "provider-report-api/pkg/vault"

    "github.com/jmoiron/sqlx"
    _ "github.com/denisenkom/go-mssqldb" // SQL Server driver
)

var db *sql.DB
var once sync.Once

func Initialize(databaseURL string) (*sqlx.DB, error) {
    log.Printf("Connecting to SQL Server database...")
    log.Printf("Connection string: %s", maskPassword(databaseURL))
    
    // เปลี่ยนจาก postgres เป็น mssql
    db, err := sqlx.Connect("mssql", databaseURL)
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

    log.Println("Successfully connected to SQL Server database")
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
        log.Println("Initializing database connection...")
        
        // ลองใช้ vault secret ก่อน
        creds := vault.GetDatabaseSecret()
        if creds != nil {
            log.Println("Using Vault credentials")
            connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=disable;connection timeout=30",
                creds.DatabaseUrl, creds.DatabaseUsername, creds.DatabasePassword, creds.DatabasePort, creds.DatabaseName)
            
            var err error
            db, err = sql.Open("mssql", connString)
            if err != nil {
                log.Printf("Error creating database connection with Vault: %v", err)
                db = nil
            } else {
                err = db.Ping()
                if err != nil {
                    log.Printf("Error pinging database with Vault: %v", err)
                    db.Close()
                    db = nil
                }
            }
        }
        
        // ถ้า vault ไม่ได้หรือเชื่อมต่อไม่ได้ ให้ใช้ environment variables
        if db == nil {
            log.Println("Vault failed, using environment variables")
            config := Load()
            connString := config.GetDatabaseURL()
            log.Printf("Connection string: %s", maskPassword(connString))
            
            var err error
            db, err = sql.Open("mssql", connString)
            if err != nil {
                log.Fatalf("Error creating database connection: %v", err)
            }

            err = db.Ping()
            if err != nil {
                log.Fatalf("Error pinging database: %v", err)
            }
            
            log.Println("Successfully connected to database using environment variables")
        }
    })
    return db
}

// Helper function to mask password in connection string for logging
func maskPassword(connStr string) string {
    // Simple password masking for logging
    return "server=***;user id=***;password=***;port=***;database=***;encrypt=disable"
}