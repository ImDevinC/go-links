package config

import "os"

type Config struct {
	StaticPath string
}

func FromEnv() Config {
	cfg := Config{
		StaticPath: os.Getenv("STATIC_PATH"),
	}
	return cfg
}
