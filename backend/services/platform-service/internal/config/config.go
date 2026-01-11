package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

const defaultConfigName = "default.json"

//go:embed *.json
var defaultConfig embed.FS

type Config struct {
	Server struct {
		GRPC struct {
			Addr string `json:"addr"`
		} `json:"grpc"`
		HTTP struct {
			Addr string `json:"addr"`
		} `json:"http"`
	} `json:"server"`
	ContentService struct {
		Addr string `json:"addr"`
	} `json:"content_service"`
	Logging struct {
		Level slog.Level `json:"level"`
	} `json:"logging"`
}

func readDefaultConfig() (*Config, error) {
	config := &Config{}

	data, err := defaultConfig.ReadFile(defaultConfigName)
	if err != nil {
		return nil, fmt.Errorf("failed read default config: %w", err)
	}

	if err = json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed parse default config: %w", err)
	}

	return config, nil
}

func readConfigWithOverride(config *Config, fileName string) (*Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed read config %q: %w", fileName, err)
	}

	if err = json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed parse config %q: %w", fileName, err)
	}

	return config, nil
}

func Load() (*Config, error) {
	config, err := readDefaultConfig()
	if err != nil {
		return nil, err
	}

	flags := readFlags(flagValues{
		config:  "",
		verbose: false,

		grpcAddr: config.Server.GRPC.Addr,
		httpAddr: config.Server.HTTP.Addr,
	})

	if flags.config != "" {
		config, err = readConfigWithOverride(config, flags.config)
		if err != nil {
			return nil, err
		}
	}

	config.Server.GRPC.Addr = flags.grpcAddr
	config.Server.HTTP.Addr = flags.httpAddr

	if flags.verbose {
		config.Logging.Level = slog.LevelDebug
	}

	return config, nil
}
