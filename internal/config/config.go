package config

import "os"

// Config holds application configuration
type Config struct {
	Port       string
	RPCURL     string
	Network    string
	APIKey     string
	WebhookURL string
}

// Load loads configuration from environment variables
func Load() Config {
	return Config{
		Port:       getEnv("PORT", "3000"),
		RPCURL:     getEnv("SOLANA_RPC_URL", "https://api.devnet.solana.com"),
		Network:    getEnv("SOLANA_NETWORK", "devnet"),
		APIKey:     os.Getenv("API_KEY"),
		WebhookURL: os.Getenv("WEBHOOK_URL"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
