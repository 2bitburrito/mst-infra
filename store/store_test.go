package store

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

var tests = map[string]uuid.UUID{
	"1": uuid.New(),
	"2": uuid.New(),
	"3": uuid.New(),
	"4": uuid.New(),
}

func TestVerificationStore_NewDelete(t *testing.T) {
	store := CreateVerificationStore(1*time.Second, 3*time.Second)
	store.New(tests["1"])
	store.New(tests["2"])
	store.New(tests["3"])
	store.Delete(tests["2"])
	want := 2
	result := len(store.Map)
	fmt.Printf("Current Store has %v items\n", result)
	if result != want {
		t.Errorf("incorrect adding into store")
	}
}

func TestVerificationStore_Expire(t *testing.T) {
	store := CreateVerificationStore(1*time.Second, 5*time.Second)
	store.New(tests["1"])
	store.New(tests["2"])
	time.Sleep(1 * time.Second)
	store.New(tests["3"])
	tests := []struct {
		description string
		id          uuid.UUID
		want        bool
	}{
		{"1 Should be Expired", tests["1"], false},
		{"2 Should be Expired", tests["2"], false},
		{"3 Should be Valid", tests["3"], true},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			res, err := store.Get(test.id)
			if test.want {
				if err != nil {
					t.Errorf("Error getting variable: %v \n%v", test.id, err.Error())
				}
				if len(res) == 0 {
					t.Errorf("Should have had a return value from Expire %v", err.Error())
				}
			}
		})
	}
}

func TestVerificationStore_Reap(t *testing.T) {
	store := CreateVerificationStore(60*time.Second, 2*time.Second)
	store.New(tests["1"])
	time.Sleep(2 * time.Second)
	store.New(tests["2"])
	store.New(tests["3"])
	tests := []struct {
		description string
		id          uuid.UUID
		want        bool
	}{
		{"1 Should be Vaid", tests["1"], true},
		{"2 Should be Expired", tests["2"], false},
		{"3 Should be Expired", tests["3"], false},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			res, err := store.Get(test.id)
			if !test.want {
				if err != nil {
					t.Errorf("Error getting variable: %v \n%v", test.id, err.Error())
				}
				if len(res) == 0 {
					t.Errorf("Should have had a return value from Expire %v", err.Error())
				}
			}
		})
	}
}
