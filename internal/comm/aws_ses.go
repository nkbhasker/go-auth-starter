package comm

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type awsSES struct {
	svc    *ses.SES
	sender string
}

func NewSES(session *session.Session, sender string) EmailClient {
	svc := ses.New(session)
	return &awsSES{
		svc:    svc,
		sender: sender,
	}
}

func (s *awsSES) Send(
	recipients []string,
	subject string,
	body string,
) error {
	toAddressess := make([]*string, len(recipients))
	for i, r := range recipients {
		toAddressess[i] = aws.String(r)
	}
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: toAddressess,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(body),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(s.sender),
	}

	_, err := s.svc.SendEmail(input)
	if err != nil {
		return err
	}

	return nil
}
