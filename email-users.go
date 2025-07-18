package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/2bitburrito/mst-infra/email"
	"github.com/2bitburrito/mst-infra/email/html"
	"github.com/2bitburrito/mst-infra/utils"
)

type emailBetaUsersRequest struct {
	Emails []string `json:"emails"`
}

func (api *API) emailBetaUsers(w http.ResponseWriter, r *http.Request) {
	var req emailBetaUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		returnJsonError(w, "error de-serialising json: "+err.Error(), 400)
		return
	}
	betaMessageBytes, err := os.ReadFile("./email/messages/beta-message.html")
	if err != nil {
		returnJsonError(w, "error reading beta-message.html"+err.Error(), 500)
		return
	}
	betaMessage := string(betaMessageBytes)
	for _, emailAddr := range req.Emails {
		go func() {
			user, err := api.queries.GetUserFromEmail(r.Context(), emailAddr)
			if err != nil {
				returnJsonError(w, fmt.Sprintf("error getting user from email: %s\n%s", emailAddr, err.Error()), 400)
				return
			}
			emailData := html.GenericEmailData{
				HighlightWord: utils.StrPtr("Beta"),
				FirstName:     utils.StrPtr(user.FullName),
				MainMessage:   &betaMessage,
				CtaText:       utils.StrPtr("Download Now"),
				ClosingMessage: utils.StrPtr(`Note that the authors of this application will take no 
					liability for any damages caused by usage of this application.`),
				PreferencesUrl: utils.StrPtr("https.beta.metasoundtools.com/profile"),
			}
			genericTemplate, err := os.ReadFile("./email/html/templates/generic-template.html")
			if err != nil {
				returnJsonError(w, "error reading from html template"+err.Error(), 500)
				return
			}
			genericTemplateReader := bytes.NewReader(genericTemplate)
			htmlEmail, err := html.TemplateEmail(genericTemplateReader, emailData)
			if err != nil {
				returnJsonError(w, "error generating html template"+err.Error(), 500)
				return
			}
			fmt.Println(htmlEmail)
			params := email.SendEmailParams{
				ReceivingAddress: emailAddr,
				SendingAddress:   "hello@metasoundtools.com",
				Subject:          "Meta Sound Tools Beta Program",
				FormattedHtml:    htmlEmail,
			}
			api.config.EmailClient.SendEmail(params)
		}()
	}
}
