package email

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type SesEmailClient struct {
	AwsCfg aws.Config
}
type SendEmailParams struct {
	ReceivingAddress string
	SendingAddress   string
	Subject          string
	FormattedHtml    *bytes.Buffer
}

func (c SesEmailClient) SendEmail(params SendEmailParams) error {
	if params.FormattedHtml == nil ||
		params.ReceivingAddress == "" ||
		params.SendingAddress == "" {
		return errors.New("missing required email params")
	}

	client := sesv2.NewFromConfig(c.AwsCfg)

	var email bytes.Buffer
	writer := multipart.NewWriter(&email)

	headers := textproto.MIMEHeader{}
	headers.Add("Content-Type", "text/html")
	headers.Set("MIME-Version", "1.0")
	// headers.Add("From", params.SendingAddress)
	headers.Add("To", params.ReceivingAddress)
	headers.Add("Subject", params.Subject)

	for k, v := range headers {
		fmt.Fprintf(&email, "%s: %s\r\n", k, v)
	}

	emailQP := quotedprintable.NewWriter(&email)
	_, err := emailQP.Write(params.FormattedHtml.Bytes())
	if err != nil {
		return errors.New("could't write html to email buffer: " + err.Error())
	}
	emailQP.Close()
	writer.Close()
	fmt.Println("Whole Email: ", email.String())

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
	output, err := client.SendEmail(context.Background(), &sesEmailParams)
	if err != nil {
		return err
	}
	fmt.Println("email output: ", output)

	return nil
}
