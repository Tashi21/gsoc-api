package middleware

import (
	"encoding/json"
	"fmt"
	"gsoc-api/models"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func DeleteOrg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// variables
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

	// delete the organization
	_, err = db.Exec(fmt.Sprintf("DELETE FROM organizations WHERE id = '%s'", id))
	flag = internalError(w, err, "Error deleting organization.")

	if flag == -1 {
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.JsonResponse{Type: "success", Message: "Organization deleted."})
}
