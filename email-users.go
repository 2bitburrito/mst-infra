package main

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	database "github.com/2bitburrito/mst-infra/db/sqlc"
	"github.com/2bitburrito/mst-infra/email"
	"github.com/2bitburrito/mst-infra/email/html"
	"github.com/2bitburrito/mst-infra/utils"
)

type emailBetaUsersRequest struct {
	Emails []string `json:"emails"`
}
type sendInviteParams struct {
	betaRows []database.BetaLicence
	isTest   bool
}

//go:embed email/messages/beta-message.html
var betaMessage string

//go:embed email/html/templates/generic-template.html
var genericTemplate []byte

func (api *API) emailSelectBetaUsers(w http.ResponseWriter, r *http.Request) {
	var req emailBetaUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		returnJsonError(w, "error de-serialising json: "+err.Error(), 400)
		return
	}
	var sendInviteParams sendInviteParams
	for _, email := range req.Emails {
		sendInviteParams.betaRows = append(sendInviteParams.betaRows, database.BetaLicence{
			Email: sql.NullString{
				Valid: true, String: email,
			},
		})
	}

	err := api.sendBetaInvites(r.Context(), sendInviteParams)
	if err != nil {
		returnJsonError(w, "error sending emails: "+err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"emails sent successfully"}`))
}

func (api *API) testEmails(w http.ResponseWriter, r *http.Request) {
	var req emailBetaUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		returnJsonError(w, "error de-serialising json: "+err.Error(), 400)
		return
	}

	var sendInviteParams sendInviteParams
	for _, email := range req.Emails {
		sendInviteParams.betaRows = append(sendInviteParams.betaRows, database.BetaLicence{
			Email: sql.NullString{
				Valid: true, String: email,
			},
			Name: sql.NullString{
				Valid:  true,
				String: "Test",
			},
		})
	}
	sendInviteParams.isTest = true
	err := api.sendBetaInvites(r.Context(), sendInviteParams)
	if err != nil {
		returnJsonError(w, "error sending emails: "+err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"emails sent successfully"}`))
}

func (api *API) emailAllBetaUsers(w http.ResponseWriter, r *http.Request) {
	betaRows, err := api.queries.GetAllBetaEmails(r.Context())
	if err != nil {
		returnJsonError(w, "error getting licences from db"+err.Error(), 500)
	}
	log.Printf("Sending emails to: %+v", betaRows)

	err = api.sendBetaInvites(r.Context(), sendInviteParams{betaRows: betaRows})
	if err != nil {
		returnJsonError(w, "error sending emails: "+err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"emails sent successfully"}`))
}

func (api *API) sendBetaInvites(ctx context.Context, params sendInviteParams) error {
	var wg sync.WaitGroup
	wg.Add(len(params.betaRows))
	errCh := make(chan error, len(params.betaRows))

	for _, row := range params.betaRows {
		go func(row database.BetaLicence) {
			defer wg.Done()
			name := row.Name.String
			if name == "" {
				user, err := api.queries.GetNameFromBetaList(ctx, sql.NullString{Valid: true, String: row.Email.String})
				if err != nil {
					errCh <- fmt.Errorf("error getting user's name from email address: %s: %s", row.Email.String, err)
					return
				}
				name = user.String
			}
			var closingMessage template.HTML
			if params.isTest {
				closingMessage = template.HTML("<p>mlrch-5f37951e982acdd9b393429</p>")
			} else {
				closingMessage = template.HTML(`<p> Please treat this tool like any professional-grade software: keep backups, test workflows on non-critical files first, and make sure you’re confident before using batch operations on large projects.</p> 
					<p>By using this application, you acknowledge that you are solely responsible for your data. The developers and authors of this software accept no liability for data loss, corruption, or any other damages resulting from its use.</p> `)
			}
			emailData := html.GenericEmailData{
				HighlightWord:  utils.StrPtr("Beta"),
				FirstName:      utils.StrPtr(name),
				MainMessage:    template.HTML(betaMessage),
				CtaLink:        utils.StrPtr("https://beta.metasoundtools.com"),
				CtaText:        utils.StrPtr("Download Now"),
				ClosingMessage: closingMessage,
				PreferencesUrl: utils.StrPtr("https.beta.metasoundtools.com/profile"),
				ExtraTags:      false,
			}
			genericTemplateReader := bytes.NewReader(genericTemplate)
			htmlEmail, err := html.TemplateEmail(genericTemplateReader, emailData)
			if err != nil {
				errCh <- fmt.Errorf("error generating html template %v", err)
				return
			}
			params := email.SendEmailParams{
				ReceivingAddress: row.Email.String,
				SendingAddress:   "Hugh <hugh@metasoundtools.com>",
				Subject:          "Meta Sound Tools Beta Program",
				FormattedHtml:    htmlEmail,
			}
			api.config.EmailClient.SendEmail(params)
			log.Printf("Successfully sent email to: %s", row.Email.String)
		}(row)
	}

	wg.Wait()
	close(errCh)

	var errs []string
	for err := range errCh {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}
	return nil
}
