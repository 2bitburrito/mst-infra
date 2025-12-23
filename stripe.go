package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	database "github.com/2bitburrito/mst-infra/db/sqlc"
	"github.com/2bitburrito/mst-infra/internal/stripeAPI"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type checkoutPayload struct {
	StripeCustomerID string              `json:"stripeCustomerID"`
	UserID           string              `json:"userID"`
	Products         []stripeAPI.Product `json:"products"`
}

type createCustomerPayload struct {
	UserID uuid.UUID `json:"userID"`
}
type getStripeOrderPayload struct {
	SessionID string `json:"checkoutSessionID"`
}

func (api *API) createStripeCustomer(w http.ResponseWriter, r *http.Request) {
	var payload createCustomerPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		returnJsonError(w, "Couldn't decode response body"+err.Error(), 500)
		return
	}
	defer r.Body.Close()

	// if user already has stripe id get from db:
	user, err := api.queries.GetUser(r.Context(), payload.UserID)
	if err != nil {
		returnJsonError(w, "Error getting stripe ID from db"+err.Error(), 500)
		return
	}
	stripeCustomerID := user.StripeID.String
	// else create new user:
	if len(stripeCustomerID) == 0 {
		params := &stripe.CustomerCreateParams{
			Name:  stripe.String(user.FullName),
			Email: stripe.String(user.Email),
		}
		result, err := api.config.StripeClient.V1Customers.Create(r.Context(), params)
		if err != nil {
			returnJsonError(w, "Couldn't create a new stripe user"+err.Error(), 500)
			return
		}
		stripeCustomerID = result.ID
		err = api.queries.AddStripeIDtoUser(r.Context(), database.AddStripeIDtoUserParams{
			ID: user.ID,
			StripeID: sql.NullString{
				Valid:  true,
				String: stripeCustomerID,
			},
		})
		if err != nil {
			returnJsonError(w, "Couldn't add stripe user to db: "+err.Error(), 500)
			return
		}
	}
	returnJson := map[string]string{
		"stripeID": stripeCustomerID,
		"userID":   user.ID.String(),
	}
	respondWithJSON(w, http.StatusOK, returnJson)
}

func (api *API) createStripeCheckout(w http.ResponseWriter, r *http.Request) {
	var payload checkoutPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		returnJsonError(w, "Couldn't decode response body "+err.Error(), 500)
		return
	}
	defer r.Body.Close()

	var baseURL string
	if api.config.Env == "dev" {
		baseURL = "http://localhost:3000"
	} else {
		baseURL = "https://metasoundtools.com"
	}

	successURL := fmt.Sprintf("%s/success?session_id={CHECKOUT_SESSION_ID}", baseURL)
	params := &stripe.CheckoutSessionCreateParams{
		SuccessURL: stripe.String(successURL),
		LineItems:  []*stripe.CheckoutSessionCreateLineItemParams{},
		Customer:   stripe.String(payload.StripeCustomerID),
		InvoiceCreation: &stripe.CheckoutSessionCreateInvoiceCreationParams{
			Enabled: stripe.Bool(true),
		},
		Mode:     stripe.String("payment"),
		Metadata: map[string]string{"userID": payload.UserID},
	}

	for _, prod := range payload.Products {
		priceCode := stripeAPI.GetProductPrice(prod, api.config.Env)
		params.LineItems = append(params.LineItems, &stripe.CheckoutSessionCreateLineItemParams{
			AdjustableQuantity: &stripe.CheckoutSessionCreateLineItemAdjustableQuantityParams{
				Enabled: stripe.Bool(true),
			},
			Price:    stripe.String(priceCode),
			Quantity: stripe.Int64(1),
		})
	}

	result, err := api.config.StripeClient.V1CheckoutSessions.Create(r.Context(), params)
	if err != nil {
		returnJsonError(w, "Couldn't start new Stripe sessionObj: "+err.Error(), 500)
		return
	}
	returnJson := map[string]string{
		"checkoutURL": result.URL,
	}
	respondWithJSON(w, http.StatusOK, returnJson)
}

func (api *API) getStripeOrder(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("session_id")

	params := &stripe.CheckoutSessionRetrieveParams{}
	params.AddExpand("line_items")
	sessionDetails, err := api.config.StripeClient.V1CheckoutSessions.Retrieve(r.Context(), sessionID, params)
	if err != nil {
		returnJsonError(w, "Error while getting session from stripe: "+err.Error(), 500)
		return
	}
	respondWithJSON(w, http.StatusOK, sessionDetails)
}

func (api *API) handleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are accepted", http.StatusMethodNotAllowed)
		return
	}
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		returnJsonError(w, "Error reading request body: "+err.Error(), 503)
		return
	}
	sigHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sigHeader, api.config.StripeEndpointSecret)
	if err != nil {
		returnJsonError(w, "Couldn't construct new Webhook event or verify signature: "+err.Error(), 400)
		return
	}
	handlerType, exists := stripeAPI.EventHandler[event.Type]
	if !exists {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	exists, err = api.queries.StripeEventExists(r.Context(), event.ID)
	if err != nil {
		returnJsonError(w, "Couldn't check stripe event in DB: "+err.Error(), 500)
		return
	}
	if exists {
		returnJsonError(w, "Event already handled", 409)
		return
	}
	err = api.queries.AddStripeEvent(r.Context(), database.AddStripeEventParams{
		ID:   event.ID,
		Type: string(event.Type),
	})
	if err != nil {
		log.Println("ERROR: couldn't set stripe event in DB: ", err)
	}

	switch handlerType {
	case "PAYMENT_SUCCESS":
		var successObj *paymentSuccess

		switch event.Type {
		case "checkout.session.completed":
			successObj, err = api.handleCheckoutSuccess(r.Context(), event)
			if err != nil {
				returnJsonError(w, "Couldn't get checkout details "+err.Error(), 500)
				return
			}

		case "invoice.paid":
			successObj, err = api.handleInvoicePaid(r.Context(), event)
			if err != nil {
				returnJsonError(w, "Couldn't get invoice details "+err.Error(), 500)
				return
			}
		}
		if successObj == nil {
			returnJsonError(w, "Couldn't get required information from stripe", 500)
			return
		}

		// HACK: This whole transaction is a problem if we have multiple products...
		// Fine for now as we only sell a standard-licence but update if plugins or something
		// else comes along
		tx, err := api.db.Begin()
		if err != nil {
			returnJsonError(w, "Couldn't create db tx"+err.Error(), 500)
			return
		}
		defer tx.Rollback()
		qtx := api.queries.WithTx(tx)
		if err := qtx.IncrementLicences(r.Context(), database.IncrementLicencesParams{
			ID:               successObj.userID.UUID,
			NumberOfLicences: int32(successObj.lineItems[0].Quantity),
		}); err != nil {
			returnJsonError(w, "Couldn't increment licences: "+err.Error(), 500)
			return
		}
		newLicences := make([]string, 0, len(successObj.lineItems))
		for range successObj.lineItems {
			licence, err := qtx.AddPaidLicence(r.Context(), successObj.userID.UUID)
			if err != nil {
				returnJsonError(w, "Couldn't create new licence: "+err.Error(), 500)
				return
			}
			newLicences = append(newLicences, licence)
		}
		if err := tx.Commit(); err != nil {
			returnJsonError(w, "Error in DB Transaction: "+err.Error(), 500)
			return
		}
	case "PAYMENT_FAILED":
		fmt.Printf("Payment Failed for event: %s\n", event.Data.Object["id"])
	case "PAYMENT_CANCELED":
		// I assume we just email them...?
		fmt.Println("Payment Canceled")
	case "ACTION_REQUIRED":
		// Here we email them a prompt to take action to continue...
		fmt.Println("Action required")
	}
}

type paymentSuccess struct {
	lineItems []*stripe.LineItem
	userID    uuid.NullUUID
}

func (api *API) handleCheckoutSuccess(ctx context.Context, event stripe.Event) (*paymentSuccess, error) {
	eventID := event.Data.Object["id"].(string)
	checkoutParams := &stripe.CheckoutSessionRetrieveParams{}
	checkoutParams.AddExpand("line_items")
	sessionObj, err := api.config.StripeClient.V1CheckoutSessions.Retrieve(ctx, eventID, checkoutParams)
	if err != nil {
		return &paymentSuccess{}, err
	}

	lineItemList := sessionObj.LineItems.Data
	userID := sessionObj.Metadata["userID"]
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Println("Couldn't parse uuid from checkout sesssion metadata", err)
	}

	nullUUID := uuid.NullUUID{
		Valid: true,
		UUID:  userUUID,
	}

	// Remove user's entry from trial table
	go func() {
		err := api.queries.RemoveFromTrialMachines(ctx, nullUUID)
		if err != nil {
			log.Println("Couldn't remove user from trial table", err)
			sentry.CaptureException(fmt.Errorf("couldn't remove user from trial table %v", err))
		}
	}()

	returnObj := paymentSuccess{
		lineItems: lineItemList,
		userID:    nullUUID,
	}
	return &returnObj, nil
}

func (api *API) handleInvoicePaid(ctx context.Context, event stripe.Event) (*paymentSuccess, error) {
	return &paymentSuccess{}, nil
}
