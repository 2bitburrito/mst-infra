package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	database "github.com/2bitburrito/mst-infra/db/sqlc"
	"github.com/2bitburrito/mst-infra/utils"
)

type BinaryInfo struct {
	Version string `json:"version"`
	Files   []struct {
		URL    string `json:"url"`
		Sha512 string `json:"sha512"`
		Size   int64  `json:"size"`
	} `json:"files"`
	Path        string    `json:"path"`
	Sha512      string    `json:"sha512"`
	ReleaseDate time.Time `json:"releaseDate"`
	Arch        string    `json:"arch"`
	Platform    string    `json:"platform"`
}

type GetBinaries struct {
	Arch     string `json:"arch"`
	Platform string `json:"platform"`
}
type ReturnURL struct {
	URL           string `json:"url"`
	LatestVersion string `json:"latest_version"`
}

func (api *API) insertLatestBinaries(w http.ResponseWriter, req *http.Request) {
	var request BinaryInfo
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		returnJsonError(w, "Error while inserting Binaries: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()
	utils.PrintPretty(request)

	if len(request.Files) == 0 {
		returnJsonError(w, "No files provided", http.StatusBadRequest)
		return
	}
	splitPath := strings.Split(request.Path, ".")

	err = api.queries.UnsetIsLatest(req.Context(), database.UnsetIsLatestParams{
		Platform: request.Platform,
		Architecture: sql.NullString{
			Valid:  true,
			String: request.Arch,
		},
	})
	if err != nil {
		log.Println("Error while unsetting is latest in app_releases: " + err.Error())
	}

	err = api.queries.AddNewReleaseData(req.Context(), database.AddNewReleaseDataParams{
		Platform: request.Platform,
		Architecture: sql.NullString{
			Valid:  true,
			String: request.Arch,
		},
		ReleaseVersion: request.Version,
		UrlFilename:    splitPath[0],
		FileSize: sql.NullInt64{
			Valid: true,
			Int64: int64(request.Files[0].Size),
		},
		ReleaseDate: sql.NullTime{
			Valid: true,
			Time:  time.Time(request.ReleaseDate),
		},
	})
	if err != nil {
		returnJsonError(w, "Error Writing New Release Data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Release data inserted successfully",
	})
}

func (api *API) getLatestBinaries(w http.ResponseWriter, req *http.Request) {
	log.Println("Getting Latest Binaries")
	platform := req.PathValue("platform")
	arch := req.PathValue("arch")

	log.Printf("received binaries request for platform: %v, arch: %v", platform, arch)
	latestRelease, err := api.queries.GetLatestBinary(req.Context(), database.GetLatestBinaryParams{
		Architecture: sql.NullString{
			Valid:  true,
			String: arch,
		},
		Platform: platform,
	})
	if err != nil {
		returnJsonError(w, "Error while fetching latest binaries from db"+err.Error(), http.StatusInternalServerError)
		return
	}

	returnBody := ReturnURL{
		URL:           latestRelease.UrlFilename,
		LatestVersion: latestRelease.ReleaseVersion,
	}
	dat, err := json.Marshal(returnBody)
	if err != nil {
		returnJsonError(w, "Error marshalling respoonse", http.StatusInternalServerError)
		return
	}
	log.Println("Successfully retreived Binaries data")
	log.Println("Latest Version details:")
	utils.PrintPretty(returnBody)

	w.Write(dat)
}
