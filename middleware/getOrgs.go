package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gsoc-api/models"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func GetOrgs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// variables
	var orgs []models.Organization
	var data *sql.Rows
	var err error

	// database setup
	db := SetupDB()
	defer db.Close()

	params := r.URL.Query()
	n := params.Get("name")
	y := params.Get("year")

	checkParams := r.URL.Query()
	checkParams.Del("name")
	checkParams.Del("year")

	if len(checkParams) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.JsonResponse{Type: "error", Message: "Invalid parameters.", Count: 0, Data: []models.Organization{}})
		return
	}

	if n != "" && y != "" {
		data, err = db.Query(fmt.Sprintf("SELECT id, name, year, link, website, created_at FROM organizations WHERE name LIKE '%%%s%%' AND year = %s", n, y))
		CheckErr(err)
	} else if n != "" {
		// getting all organizations with name contaning given letters
		data, err = db.Query(fmt.Sprintf("SELECT id, name, year, link, website, created_at FROM organizations WHERE name LIKE '%%%s%%'", n))
		CheckErr(err)
	} else if y != "" {
		// getting all organizations for a specific year
		data, err = db.Query(fmt.Sprintf("SELECT id, name, year, link, website, created_at FROM organizations WHERE year = %s", y))
		CheckErr(err)
	} else {
		data, err = db.Query("SELECT id, name, year, link, website, created_at FROM organizations")
		CheckErr(err)
	}

	// looping through all organizations and storing them in an array
	for data.Next() {
		// variables for each column
		var id uuid.UUID
		var name string
		var year int
		var link, website *string
		var createdAt *time.Time

		// copying data from each row to the corresponding variables
		err = data.Scan(&id, &name, &year, &link, &website, &createdAt)
		CheckErr(err)

		// appending the data to the array
		orgs = append(orgs, models.Organization{Id: id, Name: name, Year: year, Link: link, Website: website, CreatedAt: createdAt})
	}

	// writing the json to the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.JsonResponse{Type: "success", Message: "Organizations retrieved successfully.", Count: len(orgs), Data: orgs})
}
