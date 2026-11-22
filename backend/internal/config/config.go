package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Port           string
	DBDSN          string
	JWTSecret      string
	WebOrigin      string
	CookieSameSite string // "lax" untuk lokal, "none" saat frontend & backend beda domain
}

func Load() Config {
	return Config{
		Port:           env("PORT", "8080"),
		DBDSN:          dsn(),
		JWTSecret:      env("JWT_SECRET", "dev-secret-ganti-di-prod"),
		WebOrigin:      env("WEB_ORIGIN", "http://localhost:3000"),
		CookieSameSite: env("COOKIE_SAMESITE", "lax"),
	}
}

func dsn() string {
	if v := os.Getenv("DB_DSN"); v != "" {
		return v
	}
	// Railway/Render menyediakan MYSQL_URL (mysql://user:pass@host:port/db).
	if v := os.Getenv("MYSQL_URL"); v != "" {
		return dsnDariURL(v)
	}
	return "root:klinik@tcp(127.0.0.1:3307)/klinikgo?parseTime=true"
}

func dsnDariURL(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	port := parsed.Port()
	if port == "" {
		port = "3306"
	}
	pass, _ := parsed.User.Password()
	tambahan := "parseTime=true&multiStatements=true"
	// Aiven/TiDB mewajibkan TLS dan menandainya lewat ssl-mode/ssl_mode di URL.
	q := strings.ToLower(parsed.RawQuery)
	if strings.Contains(q, "ssl-mode") || strings.Contains(q, "ssl_mode") {
		tambahan += "&tls=true"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		parsed.User.Username(), pass, parsed.Hostname(), port, strings.TrimPrefix(parsed.Path, "/"), tambahan)
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
