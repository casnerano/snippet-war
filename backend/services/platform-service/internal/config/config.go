package config

import "context"

func Init(ctx context.Context) error {
	InitContentServiceConfig(ctx)
	
	return nil
}
