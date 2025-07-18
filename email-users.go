package main

import (
	"bytes"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/2bitburrito/mst-infra/email"
	"github.com/2bitburrito/mst-infra/email/html"
	"github.com/2bitburrito/mst-infra/utils"
)

type emailBetaUsersRequest struct {
	Emails []string `json:"emails"`
}

//go:embed email/messages/beta-message.html
var betaMessage string

//go:embed email/html/templates/generic-template.html
var genericTemplate []byte

func (api *API) emailBetaUsers(w http.ResponseWriter, r *http.Request) {
	var req emailBetaUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		returnJsonError(w, "error de-serialising json: "+err.Error(), 400)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(req.Emails))
	var eMu sync.Mutex
	var errs []string

	for _, emailAddr := range req.Emails {
		go func(emailAddr string) {
			defer wg.Done()
			user, err := api.queries.GetNameFromBetaList(r.Context(), sql.NullString{Valid: true, String: emailAddr})
			if err != nil {
				eMu.Lock()
				errs = append(errs, fmt.Sprintf("error getting user from email: %s: %s", emailAddr, err))
				eMu.Unlock()
				return
			}
			emailData := html.GenericEmailData{
				HighlightWord: utils.StrPtr("Beta"),
				FirstName:     utils.StrPtr(user.String),
				MainMessage:   template.HTML(betaMessage),
				CtaLink:       utils.StrPtr("https://beta.metasoundtools.com"),
				CtaText:       utils.StrPtr("Download Now"),
				ClosingMessage: template.HTML(`<p>Please treat this tool like any professional-grade software: keep backups, test workflows on non-critical files first, and make sure you’re confident before using batch operations on large projects.</p>
					<p>By using this application, you acknowledge that you are solely responsible for your data. The developers and authors of this software accept no liability for data loss, corruption, or any other damages resulting from its use.</p> `),
				PreferencesUrl: utils.StrPtr("https.beta.metasoundtools.com/profile"),
				ExtraTags:      false,
			}
			genericTemplateReader := bytes.NewReader(genericTemplate)
			htmlEmail, err := html.TemplateEmail(genericTemplateReader, emailData)
			if err != nil {
				eMu.Lock()
				errs = append(errs, fmt.Sprintf("error generating html template %v", err))
				eMu.Unlock()
				return
			}
			params := email.SendEmailParams{
				ReceivingAddress: emailAddr,
				SendingAddress:   "Hugh <hello@metasoundtools.com>",
				Subject:          "Meta Sound Tools Beta Program",
				FormattedHtml:    htmlEmail,
			}
			api.config.EmailClient.SendEmail(params)
		}(emailAddr)
	}
	wg.Wait()
	if len(errs) > 0 {
		returnJsonError(w, fmt.Sprintf("Errors occured in email Beta Users: %v", errs), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"emails sent successfully"}`))
}
