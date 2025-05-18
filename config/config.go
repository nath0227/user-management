package config

import (
	"time"

	env "github.com/caarlos0/env/v10"
)

type AppConfig struct {
	HttpServer        HttpServer
	GrpcServer        GrpcServer
	Crypto            CryptoCredential
	MongoDB           MongoConfig
	UserCountInterval time.Duration `env:"USER_COUNT_INTERVAL" envDefault:"10s"`
}

type HttpServer struct {
	Port string `env:"HTTP_SERVER_PORT"`
}

type GrpcServer struct {
	Port string `env:"GRPC_SERVER_PORT"`
}

type CryptoCredential struct {
	JwtKey            string        `env:"CRYPTO_JWT_KEY"`
	JwtExpireDuration time.Duration `env:"CRYPTO_JWT_EXPIRE_DURATION"`
}

type MongoConfig struct {
	Uri            string `env:"MONGO_CONFIG_URI"`
	Username       string `env:"MONGO_CONFIG_USERNAME"`
	Password       string `env:"MONGO_CONFIG_PASSWORD"`
	Database       string `env:"MONGO_CONFIG_DATABASE"`
	UserCollection string `env:"MONGO_CONFIG_USER_COLLECTION"`
}

func NewAppConfig() (*AppConfig, error) {
	var cfg AppConfig

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
