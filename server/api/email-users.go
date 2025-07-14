package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/2bitburrito/mst-infra/email"
	"github.com/2bitburrito/mst-infra/email/html"
)

type emailBetaUsersRequest struct {
	emails []string
}

func (api *API) emailBetaUsers(w http.ResponseWriter, r *http.Request) {
	var req emailBetaUsersRequest
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		returnJsonError(w, "error de-serialising json: "+err.Error(), 400)
		return
	}
	for _, emailAddr := range req.emails {
		emailData := html.GenericEmailData{}
		genericTemplate, err := os.ReadFile("../../email/html/templates/generic-template.html")
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
		params := email.SendEmailParams{
			ReceivingAddress: emailAddr,
			SendingAddress:   "hello@metasoundtools.com",
			Subject:          "Meta Sound Tools Beta Program",
			FormattedHtml:    htmlEmail,
		}
		api.config.EmailClient.SendEmail(params)
	}
}
