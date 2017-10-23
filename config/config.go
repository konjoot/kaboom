package config

// New - provides user with fresh config
func New() *Config {
	return &Config{Listen: "localhost:20000"}
}

// Config is a structured config
type Config struct {
	Listen string
}
