// Package vault provides utility functions to interact with HashiCorp Vault to securely retrieve
// stored secrets, specifically database credentials. This package encapsulates the complexity of
// Vault operations including client initialization and secret retrieval.
package vault

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
)

// DatabaseCredentials represents the structure of the database credentials as stored in Vault.
// It maps the Vault secrets to a Go struct, facilitating the use of these credentials in the application.
type DatabaseCredentials struct {
	DatabaseName     string `mapstructure:"database_name" json:"database_name"`
	DatabasePassword string `mapstructure:"database_password" json:"database_password"`
	DatabasePort     int    `mapstructure:"database_port" json:"database_port"`
	DatabaseUrl      string `mapstructure:"database_url" json:"database_url"`
	DatabaseUsername string `mapstructure:"database_username" json:"database_username"`
}

type RedisCredentials struct {
	RedisUrl string `mapstructure:"redisUrl" json:"redisUrl"`
}

// createClient initializes and returns a new Vault client using the VAULT_ADDR environment variable.
// It returns an error if the VAULT_ADDR is not set or the client cannot be created.
func createClient() (*api.Client, error) {
	vaultAddress := os.Getenv("VAULT_ADDR")
	if vaultAddress == "" {
		return nil, fmt.Errorf("VAULT_ADDR environment variable is not set")
	}

	config := &api.Config{
		Address: vaultAddress,
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %v", err)
	}

	return client, nil
}

// getSecret retrieves a secret from Vault at a specified path and sets the client token using
// the VAULT_TOKEN environment variable. It returns the secret or an error if the secret cannot be read.
func getSecret(client *api.Client, path string) (*api.Secret, error) {
	client.SetToken(os.Getenv("VAULT_TOKEN"))
	secret, err := client.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret: %v", err)
	}

	return secret, nil
}

// GetDatabaseSecret initializes the Vault client, retrieves the secret containing database credentials,
// and decodes the secret into the DatabaseCredentials struct. This function loads environment variables from
// a .env file at runtime. It logs fatal errors and exits the application if it cannot complete any step,
// including environment loading, client creation, or secret retrieval.
//
// Returns:
//   - A pointer to a DatabaseCredentials struct populated with database credentials from Vault.
//   - Nil if there are any failures during the process, after logging the error.
func GetDatabaseSecret() *DatabaseCredentials {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	client, err := createClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	secret, err := getSecret(client, "kv/data/tpa-newcore")
	if err != nil {
		log.Fatalf("Error getting secret: %v", err)
	}

	if secret != nil && secret.Data != nil {
		data, ok := secret.Data["data"].(map[string]interface{})
		if ok {
			var creds DatabaseCredentials
			if err := mapstructure.Decode(data, &creds); err != nil {
				log.Fatalf("Failed to decode data into struct: %v", err)
			}
			return &creds
		} else {
			fmt.Println("Secret data type is not expected.")
		}
	} else {
		fmt.Println("No data returned.")
	}
	return nil
}

func GetRedisSecret() *RedisCredentials {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	client, err := createClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	secret, err := getSecret(client, "kv/data/redis")
	if err != nil {
		log.Fatalf("Error getting secret: %v", err)
	}

	if secret != nil && secret.Data != nil {
		data, ok := secret.Data["data"].(map[string]interface{})
		if ok {
			var redisData RedisCredentials
			if err := mapstructure.Decode(data, &redisData); err != nil {
				log.Fatalf("Failed to decode data into struct: %v", err)
			}
			return &redisData
		} else {
			fmt.Println("Secret data type is not expected.")
		}
	} else {
		fmt.Println("No data returned.")
	}
	return nil
}
