package config

import (
    "os"
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
        DatabaseDriver: getEnv("DB_DRIVER", "postgres"),
        DatabaseHost:   getEnv("DB_HOST", "localhost"),
        DatabasePort:   getEnv("DB_PORT", "5432"),
        DatabaseName:   getEnv("DB_NAME", "provider_detail_db"),
        DatabaseUser:   getEnv("DB_USER", "postgres"),
        DatabasePass:   getEnv("DB_PASS", "password"),
        ServerPort:     getEnv("SERVER_PORT", "8080"),
        JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
        SMTPHost:       getEnv("SMTP_HOST", "smtp.gmail.com"),
        SMTPPort:       getEnv("SMTP_PORT", "587"),
        SMTPUser:       getEnv("SMTP_USER", ""),
        SMTPPass:       getEnv("SMTP_PASS", ""),
        SMTPFrom:       getEnv("SMTP_FROM", "noreply@company.com"),
        Environment:    getEnv("ENVIRONMENT", "development"),
    }
}

func (c *Config) GetDatabaseURL() string {
    if c.DatabaseURL != "" {
        return c.DatabaseURL
    }
    return "postgres://" + c.DatabaseUser + ":" + c.DatabasePass + 
           "@" + c.DatabaseHost + ":" + c.DatabasePort + 
           "/" + c.DatabaseName + "?sslmode=disable"
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}