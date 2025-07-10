package emails

import (
	"bytes"
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
		params.SendingAddress == "" ||
		params.SendingAddress == "" {
		return errors.New("missing required email params")
	}

	client := sesv2.NewFromConfig(c.AwsCfg)

	var email bytes.Buffer
	writer := multipart.NewWriter(&email)
	boundary := writer.Boundary()
	fmt.Println("The Boundary: ", boundary)

	headers := textproto.MIMEHeader{}
	headers.Add("Content-Type", "text/html")
	headers.Set("MIME-Version", "1.0")
	headers.Add("From", params.SendingAddress)
	headers.Add("To", params.ReceivingAddress)
	headers.Add("Subject", params.Subject)

	for k, v := range headers {
		fmt.Fprintf(&email, "%s: %s\r\n", k, v)
	}
	fmt.Println("Email Headers: ", email)

	emailQP := quotedprintable.NewWriter(&email)
	_, err := emailQP.Write(params.FormattedHtml.Bytes())
	if err != nil {
		return errors.New("could't write html to email buffer: " + err.Error())
	}
	emailQP.Close()
	writer.Close()
	fmt.Println("Whole Email: ", email)

	sesEmailParams := sesv2.SendEmailInput{
		Content: &types.EmailContent{
			&types.RawMessage{
				Data: email.Bytes(),
			},
		},
	}

	return nil
}
