package config

import (
	"context"
	"os"
)

var contentServiceConfig = ContentServiceConfig{
	Host: "content-service",
}

type ContentServiceConfig struct {
	Host string
}

func InitContentServiceConfig(_ context.Context) {
	if value, exists := os.LookupEnv("CONTENT_SERVICE_HOST"); exists {
		contentServiceConfig.Host = value
	}
}

func GetContentServiceConfig() ContentServiceConfig {
	return contentServiceConfig
}
