package types

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/DisgoOrg/log"
	"github.com/DisgoOrg/snowflake"
)

func LoadConfig(log log.Logger) (*Config, error) {
	file, err := os.Open("config.json")
	if os.IsNotExist(err) {
		if file, err = os.Create("config.json"); err != nil {
			return nil, err
		}
		var data []byte
		if data, err = json.Marshal(Config{}); err != nil {
			return nil, err
		}
		if _, err = file.Write(data); err != nil {
			return nil, err
		}
		return nil, errors.New("config.json not found, created new one")
	} else if err != nil {
		return nil, err
	}

	var cfg Config
	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	log.Info("Loaded config")
	return &cfg, nil
}

type Config struct {
	DevMode    bool                `json:"dev_mode"`
	DevGuildID snowflake.Snowflake `json:"dev_guild_id"`
	LogLevel   log.Level           `json:"log_level"`
	Token      string              `json:"token"`
	PublicKey  string              `json:"public_key"`
	Database   Database            `json:"database"`
}

type Database struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}
