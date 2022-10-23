package config

import (
	"context"
)

var configContextKey = struct{}{}

// Возвращает конфиг из конекста
func FromContext(ctx context.Context) *Config {
	cfgRaw := ctx.Value(configContextKey)
	cfg, ok := cfgRaw.(*Config)
	if ok {
		return cfg
	}
	return nil
}

// Обогащает контекст конфигом
func WrapContext(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, configContextKey, cfg)
}
