package server

type Config struct {
	bind_addr string
}

func NewConfig() *Config {
	return &Config{bind_addr: ":4003"}
}
