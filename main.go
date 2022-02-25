// TODO: Add Caching
// TODO: Add GET route for individual organizations
// TODO: Add PATCH route for individual organizations
// TODO: Add DELETE route for individual organizations

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Organization struct {
	Id        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Year      int        `json:"year"`
	Link      *string    `json:"link"`
	Website   *string    `json:"website"`
	CreatedAt *time.Time `json:"created_at"`
}

type JsonResponse struct {
	Type    string         `json:"type"`
	Message string         `json:"message"`
	Count   int            `json:"count"`
	Data    []Organization `json:"data"`
}

func getEnv(key string) string {
	// set config file path
	viper.SetConfigFile(".env")

	// find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// get the key
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Error getting key: %s", key)
	}
	return value
}

func setupDB() *sql.DB {
	// setup database
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable connect_timeout=15", getEnv("DB_USER"), getEnv("DB_PASSWORD"), getEnv("DB_NAME"), getEnv("HOST"), getEnv("DB_PORT"))
	db, err := sql.Open(getEnv("DATABASE"), connectionString)
	checkErr(err)
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func router() (*mux.Router, *http.Server) {
	r := mux.NewRouter()
	r.HandleFunc("/orgs", getOrgs).Methods("GET")
	r.HandleFunc("/orgs/{id}", getOrg).Methods("GET")
	r.HandleFunc("/orgs/{id}", patchOrg).Methods("PATCH")
	r.HandleFunc("/orgs/{id}", deleteOrg).Methods("DELETE")

	server := &http.Server{
		Handler:      r,
		Addr:         getEnv("WEB_PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return r, server
}

func main() {
	// start the server
	_, server := router()
	fmt.Printf("Server started at port %s\n", getEnv("WEB_PORT"))
	log.Fatal(server.ListenAndServe())
}
