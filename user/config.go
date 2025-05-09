package user

import "time"

type (
	Config struct {
		JWTsecret string        `env:"JWT_SECRET"`
		JWTexp    time.Duration `env:"JWT_EXP"`
	}
)
