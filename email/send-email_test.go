package email

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/2bitburrito/mst-infra/emails/html"
	"github.com/2bitburrito/mst-infra/server/api/utils"
	"github.com/aws/aws-sdk-go-v2/config"
)

func TestEmail(t *testing.T) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-1"))
	if err != nil {
		t.Fatal("Couldn't configure aws: ", err)
	}
	sesClient := SesEmailClient{
		AwsCfg: awsCfg,
	}

	emailData := html.GenericEmailData{
		HighlightWord:  utils.StrPtr("Beta"),
		MainMessage:    utils.StrPtr("Hello and welcome to this email test..."),
		FirstName:      utils.StrPtr("Test"),
		CtaText:        utils.StrPtr("Click Here"),
		CtaLink:        utils.StrPtr("https://metasoundtools.com"),
		ClosingMessage: utils.StrPtr("Thanks for agreeing to be a part of this beta program..."),
		PreferencesUrl: utils.StrPtr("https://beta.metasoundtools.com/profile"),
	}

	file, err := os.ReadFile("./html/templates/generic-template.html")
	if err != nil {
		t.Fatal(err)
	}
	htmlTmpl := strings.NewReader(string(file))
	html, err := html.TemplateEmail(htmlTmpl, emailData)
	if err != nil {
		t.Fatal(err)
	}
	emailParams := SendEmailParams{
		ReceivingAddress: "hughandelsa@gmail.com",
		SendingAddress:   "hello@metasoundtools.com",
		Subject:          "Test Email",
		FormattedHtml:    html,
	}

	if err = sesClient.SendEmail(emailParams); err != nil {
		t.Fatalf("couldn't send email: %s", err)
	}
}
