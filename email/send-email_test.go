package email

import (
	"context"
	"html/template"
	"os"
	"strings"
	"testing"

	"github.com/2bitburrito/mst-infra/email/html"
	"github.com/2bitburrito/mst-infra/utils"
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
		MainMessage:    template.HTML("<p>Hello and welcome to this email test...</p>"),
		FirstName:      utils.StrPtr("Test"),
		CtaText:        utils.StrPtr("Click Here"),
		CtaLink:        utils.StrPtr("https://metasoundtools.com"),
		ClosingMessage: template.HTML("<p>Thanks for agreeing to be a part of this beta program...</p>"),
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
		ReceivingAddress: "palmerhap@gmail.com",
		SendingAddress:   "hello@metasoundtools.com",
		Subject:          "Test Email",
		FormattedHtml:    html,
	}

	if err = sesClient.SendEmail(emailParams); err != nil {
		t.Fatalf("couldn't send email: %s", err)
	}
}
