package main

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type cfg struct {
	SecretKey string `envconfig:"AWS_SECRET_KEY"`
	AccessKey string `envconfig:"AWS_ACCESS_KEY"`
	Bucket    string `envconfig:"AWS_BUCKET"`
	Region    string `envconfig:"AWS_REGION"`
}

var c = new(cfg)

func init() {
	_ = godotenv.Overload(".env", ".env.local")
	_ = envconfig.Process("", c)
}
