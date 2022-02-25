// function to get details of one organization
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func getOrg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// variables
	var orgs []Organization
	var data *sql.Rows
	var id uuid.UUID
	var err error

	// database setup
	db := setupDB()
	defer db.Close()

	// get the id from the url
	vars := mux.Vars(r)
	id, err = uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{Type: "error", Message: "Invalid ID."})
		return
	}

	// select the given id
	data, err = db.Query(fmt.Sprintf("SELECT id, name, year, tech_stack, topics, short_desc, link, img_url, website, created_at, updated_at FROM organizations WHERE id = '%s'", id))
	checkErr(err)

	// looping through all organizations and storing them in an array
	for data.Next() {
		// variables for each column
		var id uuid.UUID
		var name string
		var year int
		var tech_stack, topics, short_desc, link, img_url, website *string
		var created_at, updated_at *time.Time

		// scan the columns
		err = data.Scan(&id, &name, &year, &tech_stack, &topics, &short_desc, &link, &img_url, &website, &created_at, &updated_at)
		checkErr(err)

		// storing the data in an array
		orgs = append(orgs, Organization{Id: id, Name: name, Year: year, TechStack: tech_stack, Topics: topics, ShortDesc: short_desc, Link: link, ImgUrl: img_url, Website: website, CreatedAt: created_at, UpdatedAt: updated_at})
	}

	// if no organization was found
	if len(orgs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JsonResponse{Type: "error", Message: "Organization not found.", Count: 0, Data: []Organization{}})
		return
	}

	// if one organization was found
	if len(orgs) == 1 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(JsonResponse{Type: "success", Message: "Organization found.", Count: 1, Data: orgs})
		return
	}

	// if more than one organization was found
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(JsonResponse{Type: "error", Message: "More than one organization found.", Count: len(orgs), Data: orgs})
}
