package config

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/viper"
)

type Config struct {
	Env                    string
	Port                   string
	AwsCognitoUserPoolId   string
	AwsCognitoClientId     string
	AwsCognitoClientSecret string
	AwsConfig              aws.Config
	AwsTokenURL            string
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetDefault("ENV", "local")
	v.SetDefault("PORT", "8080")

	v.SetConfigFile(".env")
	// v.SetConfigFile("../../.env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		slog.Error("error on parsing .env file", "err", err)
		return nil, err
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration")
	}

	cfg := &Config{
		Env:                    v.GetString("ENV"),
		Port:                   v.GetString("PORT"),
		AwsCognitoUserPoolId:   v.GetString("AWS_COGNITO_USER_POOL_ID"),
		AwsCognitoClientId:     v.GetString("AWS_COGNITO_CLIENT_ID"),
		AwsCognitoClientSecret: v.GetString("AWS_COGNITO_CLIENT_SECRET"),
		AwsConfig:              awsCfg,
		AwsTokenURL:            v.GetString("AWS_COGNITO_TOKEN_URL"),
	}

	return cfg, nil
}
