package stripeAPI

import (
	"fmt"

	"github.com/2bitburrito/mst-infra/config"
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
