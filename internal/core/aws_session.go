package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AwsSession struct {
	Session *session.Session
}

type AwsSessionOptions struct {
	Region          string
	AccessKeyId     string
	SecretAccessKey string
}

func NewAwsSession(options AwsSessionOptions) (*AwsSession, error) {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(options.Region),
		Credentials: credentials.NewStaticCredentials(options.AccessKeyId, options.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}

	return &AwsSession{
		Session: session,
	}, nil
}
