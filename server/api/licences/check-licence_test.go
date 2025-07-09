package licence

import (
	"database/sql"
	"testing"
	"time"

	database "github.com/2bitburrito/mst-infra/db/sqlc"
	"github.com/google/uuid"
)

type licenceCheckTests struct {
	licence      database.Licence
	licenceTable []database.Licence
	expect       string
}

func TestCheckLicence(t *testing.T) {
	t.Parallel()
	var (
		userID   = uuid.MustParse("00ea65bc-d60a-4b10-acff-87288446c031")
		now      = time.Now()
		past     = now.Add(-48 * time.Hour)
		future   = now.Add(48 * time.Hour)
		machine1 = "machine-1"
		machine2 = "machine-2"
	)

	testLicences := []database.Licence{
		{
			LicenceKey: "unused-valid",
			UserID:     userID,
			MachineID:  sql.NullString{Valid: false},
			CreatedAt:  sql.NullTime{Valid: true, Time: past},
			LastUsedAt: sql.NullTime{Valid: false},
			LicenceType: database.NullLicenceTypeEnum{
				LicenceTypeEnum: "paid",
				Valid:           true,
			},
			Expiry: sql.NullTime{Valid: true, Time: future},
			Jti:    uuid.NullUUID{Valid: false},
		},
		{
			LicenceKey: "used-valid-machine1",
			UserID:     userID,
			MachineID:  sql.NullString{Valid: true, String: machine1},
			CreatedAt:  sql.NullTime{Valid: true, Time: past},
			LastUsedAt: sql.NullTime{Valid: true, Time: now.Add(-24 * time.Hour)},
			LicenceType: database.NullLicenceTypeEnum{
				LicenceTypeEnum: "trial",
				Valid:           true,
			},
			Expiry: sql.NullTime{Valid: true, Time: future},
			Jti:    uuid.NullUUID{Valid: true, UUID: uuid.New()},
		},
		{
			LicenceKey: "used-beta-machine2",
			UserID:     userID,
			MachineID:  sql.NullString{Valid: true, String: machine2},
			CreatedAt:  sql.NullTime{Valid: true, Time: past},
			LastUsedAt: sql.NullTime{Valid: true, Time: now.Add(-36 * time.Hour)},
			LicenceType: database.NullLicenceTypeEnum{
				LicenceTypeEnum: "beta",
				Valid:           true,
			},
			Expiry: sql.NullTime{Valid: true, Time: past},
			Jti:    uuid.NullUUID{Valid: true, UUID: uuid.New()},
		},
		{
			LicenceKey: "paid-licence",
			UserID:     userID,
			MachineID:  sql.NullString{Valid: true, String: "machine-3"}, CreatedAt: sql.NullTime{Valid: true, Time: past}, LastUsedAt: sql.NullTime{Valid: true, Time: now.Add(-12 * time.Hour)},
			LicenceType: database.NullLicenceTypeEnum{
				LicenceTypeEnum: "paid",
				Valid:           true,
			},
			Expiry: sql.NullTime{Valid: true, Time: past},
			Jti:    uuid.NullUUID{Valid: true, UUID: uuid.New()},
		},
	}

	table := []licenceCheckTests{
		{
			licence: database.Licence{
				LicenceKey: "1",
				UserID:     userID,
				MachineID:  sql.NullString{Valid: true, String: machine1},
			},
			expect:       "pass",
			licenceTable: testLicences,
		},
		{
			licence: database.Licence{
				LicenceKey: "2",
				UserID:     userID,
				MachineID:  sql.NullString{Valid: true, String: machine1},
			},
			expect:       "pass",
			licenceTable: testLicences,
		},
		{
			licence: database.Licence{
				LicenceKey: "3",
				UserID:     userID,
				MachineID:  sql.NullString{Valid: true, String: machine2},
			},
			expect:       "fail",
			licenceTable: []database.Licence{},
		},
	}
	for _, test := range table {
		rtnLicence, err := Check(test.licence.MachineID.String, test.licenceTable)
		if test.expect == "pass" {
			if err != nil {
				t.Error(err)
			}
			if rtnLicence.UserID != test.licence.UserID {
				t.Errorf("Returned user id's don't match:\n got: %s \n wanted: %s",
					rtnLicence.UserID,
					test.licence.UserID)
			}
		}
	}
}
