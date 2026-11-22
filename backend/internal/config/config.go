package config

import "os"

type Config struct {
	Port      string
	DBDSN     string
	JWTSecret string
	WebOrigin string
}

func Load() Config {
	return Config{
		Port:      env("PORT", "8080"),
		DBDSN:     env("DB_DSN", "root:klinik@tcp(127.0.0.1:3307)/klinikgo?parseTime=true"),
		JWTSecret: env("JWT_SECRET", "dev-secret-ganti-di-prod"),
		WebOrigin: env("WEB_ORIGIN", "http://localhost:3000"),
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
