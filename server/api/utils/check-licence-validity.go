package utils

import (
	"time"

	database "github.com/2bitburrito/mst-infra/db/sqlc"
)

func LicenceIsValid(licence database.Licence) bool {
	now := time.Now()

	if licence.LicenceType.LicenceTypeEnum == "paid" {
		return true
	}
	if licence.Expiry.Time.After(now) {
		return true
	}
	return false
}
