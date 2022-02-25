package middleware

import (
	"encoding/json"
	"fmt"
	"gsoc-api/models"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func PutOrg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// variables
	var org models.Organization
	var id uuid.UUID
	var err error
	var flag int = 0

	// database setup
	db := SetupDB()
	defer db.Close()

	// get the id from the url
	vars := mux.Vars(r)
	id, err = uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.JsonResponse{Type: "error", Message: "Invalid ID."})
		return
	}

	// decode the json
	err = json.NewDecoder(r.Body).Decode(&org)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.JsonResponse{Type: "error", Message: "Invalid JSON."})
		return
	}

	// update the organization
	_, err = db.Exec(fmt.Sprintf("UPDATE organizations SET name = '%s', year = %d, tech_stack = '%s', topics = '%s', short_desc = '%s', link = '%s', img_url = '%s', website = '%s', updated_at = NOW() WHERE id = '%s'", org.Name, org.Year, *org.TechStack, *org.Topics, *org.ShortDesc, *org.Link, *org.ImgUrl, *org.Website, id))
	flag = internalError(w, err, "Error updating organization.")

	if flag == -1 {
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.JsonResponse{Type: "success", Message: "Organization updated."})
}
