package configs

import (
    "os"
    
    _ "github.com/denisenkom/go-mssqldb" // SQL Server driver
)

type Config struct {
    DatabaseURL    string
    DatabaseDriver string
    DatabaseHost   string
    DatabasePort   string
    DatabaseName   string
    DatabaseUser   string
    DatabasePass   string
    ServerPort     string
    JWTSecret      string
    SMTPHost       string
    SMTPPort       string
    SMTPUser       string
    SMTPPass       string
    SMTPFrom       string
    Environment    string
}

func Load() *Config {
    return &Config{
        DatabaseDriver: getEnv("DB_DRIVER", "mssql"),           // เปลี่ยนจาก postgres เป็น mssql
        DatabaseHost:   getEnv("DB_HOST", "localhost"),
        DatabasePort:   getEnv("DB_PORT", "1433"),             // เปลี่ยนจาก 5432 เป็น 1433
        DatabaseName:   getEnv("DB_NAME", "tpacaredb"),
        DatabaseUser:   getEnv("DB_USER", "dcsnewcoretpa"),
        DatabasePass:   getEnv("DB_PASSWORD", "TPA@mindcs!2"), // ใช้ DB_PASSWORD แทน DB_PASS
        ServerPort:     getEnv("PORT", "8777"),                // ใช้ PORT แทน SERVER_PORT
        JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
        SMTPHost:       getEnv("SMTP_HOST", "smtp.gmail.com"),
        SMTPPort:       getEnv("SMTP_PORT", "587"),
        SMTPUser:       getEnv("SMTP_USERNAME", ""),           // ใช้ SMTP_USERNAME
        SMTPPass:       getEnv("SMTP_PASSWORD", ""),           // ใช้ SMTP_PASSWORD
        SMTPFrom:       getEnv("SMTP_FROM", "noreply@company.com"),
        Environment:    getEnv("ENVIRONMENT", "development"),
    }
}

func (c *Config) GetDatabaseURL() string {
    if c.DatabaseURL != "" {
        return c.DatabaseURL
    }
    
    // SQL Server connection string format
    return "server=" + c.DatabaseHost + 
           ";port=" + c.DatabasePort +
           ";user id=" + c.DatabaseUser +
           ";password=" + c.DatabasePass +
           ";database=" + c.DatabaseName +
           ";encrypt=disable;connection timeout=30"
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}