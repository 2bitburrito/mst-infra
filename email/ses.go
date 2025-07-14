package email

import (
	"bytes"
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type SesEmailClient struct {
	AwsCfg       aws.Config
	SendingEmail string
}
type SendEmailParams struct {
	ReceivingAddress string
	SendingAddress   string
	Subject          string
	FormattedHtml    *bytes.Buffer
}

type EmailSender interface {
	SendEmail(params SendEmailParams) error
}

func (c SesEmailClient) SendEmail(params SendEmailParams) error {
	if params.FormattedHtml == nil ||
		params.ReceivingAddress == "" ||
		params.SendingAddress == "" {
		return errors.New("missing required email params")
	}

	client := sesv2.NewFromConfig(c.AwsCfg)

	sesEmailParams := sesv2.SendEmailInput{
		FromEmailAddress: &params.SendingAddress,
		Destination: &types.Destination{
			ToAddresses: []string{params.ReceivingAddress},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(params.Subject),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(params.FormattedHtml.String()),
					},
				},
			},
		},
	}
	// returns a MessageId and result metadata
	_, err := client.SendEmail(context.Background(), &sesEmailParams)
	if err != nil {
		return err
	}
	return nil
}
