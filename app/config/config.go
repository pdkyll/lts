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
	LogLevel            string 	`envconfig:"LOG_LEVEL"`
	PgURL               string 	`envconfig:"PG_URL"`
	InitDB	            bool 	`envconfig:"INIT_DB"`
	HTTPAddr            string 	`envconfig:"HTTP_ADDR"`
	FolderData          string 	`envconfig:"FOLDERDATA"`
	TokenExp	        string 	`envconfig:"TOKENEXP"`
	TokenKey	        string 	`envconfig:"TOKENKEY"`
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
		infoLog := log.New(os.Stdout, "INIT\t", log.Ldate|log.Ltime|log.LUTC)
		infoLog.Println("Configuration:", string(configBytes))
	})
	return &config
}