package jwt

import "os"

type Config struct {
	JWTSecret string
}

func LoadConfigJWT() Config {
	return Config{
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
