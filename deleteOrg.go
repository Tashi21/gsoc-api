// function to delete an organization
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/google/uuid"
)

func deleteOrg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// variables
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

	// delete the organization
	_, err = db.Exec(fmt.Sprintf("DELETE FROM organizations WHERE id = '%s'", id))
	checkErr(err)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JsonResponse{Type: "success", Message: "Organization deleted."})
}
