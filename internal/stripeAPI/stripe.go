package stripeAPI

import (
	"fmt"

	"github.com/2bitburrito/mst-infra/config"
	"github.com/stripe/stripe-go/v82"
)

type Product string

const (
	StandardLicence Product = "standard-licence"
)

var productMap = map[Product]string{
	StandardLicence: "price_1Rwvld281c2JDPUFIhaFAPSs",
}

var sandboxProductMap = map[Product]string{
	StandardLicence: "price_1RvqyXS3CWbfShtyfXZitBQJ",
}

func GetProductPrice(product Product, env config.Env) string {
	var priceCode string
	fmt.Println(env)
	if env == "dev" {
		priceCode = sandboxProductMap[product]
	} else {
		priceCode = productMap[product]
	}
	return priceCode
}

var AllowedEvents = []stripe.EventType{
	"checkout.session.completed",
	"customer.subscription.created",
	"customer.subscription.updated",
	"customer.subscription.deleted",
	"customer.subscription.paused",
	"customer.subscription.resumed",
	"customer.subscription.pending_update_applied",
	"customer.subscription.pending_update_expired",
	"customer.subscription.trial_will_end",
	"invoice.paid",
	"invoice.payment_failed",
	"invoice.payment_action_required",
	"invoice.upcoming",
	"invoice.marked_uncollectible",
	"invoice.payment_succeeded",
	"payment_intent.succeeded",
	"payment_intent.payment_failed",
	"payment_intent.canceled",
}

const (
	PAYMENT_SUCCESS = "PAYMENT_SUCCESS"
	PAYMENT_FAILED  = "PAYMENT_FAILED"
	PAYMENT_PENDING = "PAYMENT_PENDING"
	ACTION_REQUIRED = "ACTION_REQUIRED"
)

var EventHandler = map[stripe.EventType]string{
	"checkout.session.completed":      "PAYMENT_SUCCESS",
	"invoice.paid":                    "PAYMENT_SUCCESS",
	"invoice.payment_failed":          "PAYMENT_FAILED",
	"invoice.payment_action_required": "ACTION_REQUIRED",
	"invoice.upcoming":                "PAYMENT_PENDING",
	"invoice.marked_uncollectible":    "PAYMENT_FAILED",
	"invoice.payment_succeeded":       "PAYMENT_SUCCESS",
	"payment_intent.succeeded":        "PAYMENT_SUCCESS",
	"payment_intent.payment_failed":   "PAYMENT_FAILED",
	"payment_intent.canceled":         "PAYMENT_CANCELED",
}
