package apiserver

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	DatabaseUrl string `toml:"database_url"`
	Origin []string `toml:"origin"`
	ContentDir string `toml:"content_dir"`
}

func NewConfig() *Config{
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}