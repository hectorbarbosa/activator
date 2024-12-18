package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	DbUrl          string `env:"DB_URL"`
	ServerAddr     string `env:"SERVER_ADDR"`
	LogLevel       int    `env:"LOG_LEVEL"`
	SmtpHost       string `env:"SMTP_HOST"`
	SmtpPort       int    `env:"SMTP_PORT"`
	ActivationPath string `env:"ACTIVATION_PATH"`
}

func NewConfig(path string) (Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
