package router

import (
	"gsoc-api/middleware"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Router() (*mux.Router, *http.Server) {
	r := mux.NewRouter()
	r.HandleFunc("/orgs", middleware.GetOrgs).Methods("GET")
	r.HandleFunc("/orgs/{id}", middleware.GetOrg).Methods("GET")
	r.HandleFunc("/orgs/{id}", middleware.DeleteOrg).Methods("DELETE")

	server := &http.Server{
		Handler:      r,
		Addr:         middleware.GetEnv("WEB_PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return r, server
}
