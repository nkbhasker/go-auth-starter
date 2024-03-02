package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type SrvConfig struct {
	Host                       string
	Port                       string
	PostgresUrl                string
	RedisUrl                   string
	JwtPrivateKey              string
	AccessTokenExpiryInMinutes int
	OtpExpiryInMinutes         int
	OtpGenerateRateLimit       int
	OtpGenerateRateLimitWindow int
	OtpVerifyRateLimit         int
	OtpVerifyRateLimitWindow   int
	AwsRegion                  string
	AwsAccessKeyId             string
	AwsSecretAccessKey         string
	AwsSesSender               string
}

func InitSrvConfig() (*SrvConfig, error) {
	envErrors := []string{}
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	host := os.Getenv("HOST")
	if host == "" {
		host = "http://localhost:8080"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	postgresUrl := os.Getenv("POSTGRES_URL")
	if postgresUrl == "" {
		envErrors = append(envErrors, "postgres url is required")
	}
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		envErrors = append(envErrors, "redis url is required")
	}
	jwtPrivateKey := os.Getenv("JWT_BASE64_ENCODED_PRIVATE_KEY")
	if jwtPrivateKey == "" {
		envErrors = append(envErrors, "jwt private key is required")
	}
	accessTokenExpiryInMinutes, ok := parseInt(os.Getenv("ACCESS_TOKEN_EXPIRY_IN_MINUTES"))
	if !ok {
		accessTokenExpiryInMinutes = 86400
	}
	otpExpiryInMinutes, ok := parseInt(os.Getenv("OTP_EXPIRY_IN_MINUTES"))
	if !ok {
		otpExpiryInMinutes = 5
	}
	otpGenerateRateLimit, ok := parseInt(os.Getenv("OTP_VERIFY_RATE_LIMIT"))
	if !ok {
		otpGenerateRateLimit = 3
	}
	otpGenerateRateLimitWindow, ok := parseInt(os.Getenv("OTP_VERIFY_RATE_LIMIT_WINDOW"))
	if !ok {
		otpGenerateRateLimitWindow = 7200
	}
	otpVerifyRateLimit, ok := parseInt(os.Getenv("OTP_VERIFY_RATE_LIMIT"))
	if !ok {
		otpVerifyRateLimit = 5
	}
	otpVerifyRateLimitWindow, ok := parseInt(os.Getenv("OTP_VERIFY_RATE_LIMIT_WINDOW"))
	if !ok {
		otpVerifyRateLimitWindow = 86400
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "ap-south-1"
	}
	awsAccessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	if awsAccessKeyId == "" {
		envErrors = append(envErrors, "aws access key id is required")
	}
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if awsSecretAccessKey == "" {
		envErrors = append(envErrors, "aws secret access key is required")
	}
	awsSesSender := os.Getenv("AWS_SES_SENDER")
	if awsSesSender == "" {
		awsSesSender = "auth@elevatr.in"
	}
	if len(envErrors) != 0 {
		return nil, errors.New(strings.Join(envErrors, "\n"))
	}

	return &SrvConfig{
		Host:                       host,
		Port:                       port,
		PostgresUrl:                postgresUrl,
		RedisUrl:                   redisUrl,
		JwtPrivateKey:              jwtPrivateKey,
		AccessTokenExpiryInMinutes: accessTokenExpiryInMinutes,
		OtpExpiryInMinutes:         otpExpiryInMinutes,
		OtpGenerateRateLimit:       otpGenerateRateLimit,
		OtpGenerateRateLimitWindow: otpGenerateRateLimitWindow,
		OtpVerifyRateLimit:         otpVerifyRateLimit,
		OtpVerifyRateLimitWindow:   otpVerifyRateLimitWindow,
		AwsRegion:                  awsRegion,
		AwsAccessKeyId:             awsAccessKeyId,
		AwsSecretAccessKey:         awsSecretAccessKey,
		AwsSesSender:               awsSesSender,
	}, nil
}

func parseInt(str string) (int, bool) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return i, false
	}

	return i, true
}
