package licence

import (
	"errors"
	"time"

	database "github.com/2bitburrito/mst-infra/db/sqlc"
	"github.com/google/uuid"
)

// licenceisValid is a helper method to check whether
// the licence is: paid or will expire later
// than now
func licenceIsValid(licence database.Licence) bool {
	now := time.Now()

	if licence.LicenceType.LicenceTypeEnum == "paid" {
		return true
	} else if licence.Expiry.Time.After(now) {
		return true
	}
	return false
}

// CheckForValid goes through a slice of licences and returns the licence that
// is either unused (with no machineid) or is least recently used
// or matches provided machineID
func CheckForValid(machineID string, licences []database.Licence) (bool, database.Licence, error) {
	if len(licences) == 0 {
		return false, database.Licence{}, errors.New("no licences found in the database matching")
	}
	userHasPaidLicence := false
	for _, licence := range licences {
		if licence.LicenceType.LicenceTypeEnum == "paid" {
			userHasPaidLicence = true
		}
	}

	var oldestLicence database.Licence
	var expiredTrial database.Licence

	for _, licence := range licences {
		// First check whether licence is a plan or within expiry
		if !licenceIsValid(licence) {
			if licence.LicenceType.LicenceTypeEnum == "trial" {
				expiredTrial = licence
			}
			continue
		}
		if licence.LicenceType.LicenceTypeEnum == "trial" &&
			userHasPaidLicence {
			continue
		}

		// Track the oldest licence:
		if !oldestLicence.LastUsedAt.Valid {
			oldestLicence = licence
		}
		if oldestLicence.LastUsedAt.Time.After(licence.LastUsedAt.Time) {
			oldestLicence = licence
		}
		// If licence doesn't have a machine ID attached then this is new licence
		if !licence.MachineID.Valid {
			return true, licence, nil
		} else if licence.MachineID.String == machineID {
			return true, licence, nil
		}
	}

	if oldestLicence.UserID == uuid.Nil {
		// if only have expired trial
		if expiredTrial.UserID != uuid.Nil {
			return false, expiredTrial, nil
		}
		return false, database.Licence{}, errors.New("couldn't find a valid licence")
	}
	// If nothing found then defaulting back to oldest available licence
	return true, oldestLicence, nil
}
