package config

import (
	"encoding/json"
	"os"
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

// Config is a config :)
type Config struct {
	BunVerbose			bool 	`envconfig:"BUN_VERBOSE"`
	LogLevel            int8 	`envconfig:"LOG_LEVEL"`
	DBURI               string 	`envconfig:"DB_URI"`
	APIAddr            	string 	`envconfig:"API_ADDR"`
	SessionDuration     string 	`envconfig:"SESSION_DURATION"`
	TokenDuration     	string 	`envconfig:"TOKEN_DURATION"`
	SignedToken			string 	`envconfig:"SIGNED_TOKEN"`
}

var (
	config Config
	once   sync.Once
)

// Get reads config from environment. Once.
func Get() *Config {
	once.Do(func() {
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatal(err)
		}
		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		infoLog := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)
		infoLog.Println("Configuration:\n", string(configBytes))
	})
	return &config
}