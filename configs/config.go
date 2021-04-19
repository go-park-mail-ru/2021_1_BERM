package configs

type Config struct {
	BindAddr    string   `toml:"bind_addr"`
	LogLevel    string   `toml:"log_level"`
	DatabaseURL string   `toml:"database_url"`
	Origin      []string `toml:"origin"`
	ContentDir  string   `toml:"content_dir"`
	HTTPS       bool     `toml:"https"`
	DSN         string   `toml:"dsn"`
	LogFile     string   `toml:"log_file"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}
