package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	database "github.com/2bitburrito/mst-infra/db/sqlc"
	"github.com/2bitburrito/mst-infra/internal/stripeAPI"
	"github.com/2bitburrito/mst-infra/utils"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
)

type checkoutPayload struct {
	StripeCustomerID string              `json:"stripeCustomerID"`
	Products         []stripeAPI.Product `json:"products"`
}

type createCustomerPayload struct {
	UserID uuid.UUID `json:"userID"`
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
	}
	json.NewEncoder(w).Encode(returnJson)
}

func (api *API) createStripeCheckout(w http.ResponseWriter, r *http.Request) {
	var payload checkoutPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		returnJsonError(w, "Couldn't decode response body "+err.Error(), 500)
		return
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Print("Error: Couldn't close the body within create stripe customer")
		}
	}()
	baseURL := utils.GetBaseURL(r)
	successURL := fmt.Sprintf("%s/success/", baseURL)
	params := &stripe.CheckoutSessionCreateParams{
		SuccessURL: stripe.String(successURL),
		LineItems:  []*stripe.CheckoutSessionCreateLineItemParams{},
		Customer:   stripe.String(payload.StripeCustomerID),
		Mode:       stripe.String("payment"),
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
		returnJsonError(w, "Couldn't start new Stripe session: "+err.Error(), 500)
		return
	}
	returnJson := map[string]string{
		"checkoutURL": result.URL,
	}
	if err := json.NewEncoder(w).Encode(returnJson); err != nil {
		returnJsonError(w, "Couln't start new checkout session: "+err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
