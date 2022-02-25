// TODO: Add Caching

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
	Id        uuid.UUID  `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Year      int        `json:"year,omitempty"`
	TechStack *string    `json:"tech_stack,omitempty"`
	Topics    *string    `json:"topics,omitempty"`
	ShortDesc *string    `json:"short_desc,omitempty"`
	Link      *string    `json:"link,omitempty"`
	ImgUrl    *string    `json:"img_url,omitempty"`
	Website   *string    `json:"website,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type JsonResponse struct {
	Type    string         `json:"type,omitempty"`
	Message string         `json:"message,omitempty"`
	Count   int            `json:"count,omitempty"`
	Data    []Organization `json:"data,omitempty"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
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

	// ping the database
	err = db.Ping()
	checkErr(err)

	return db
}

func router() (*mux.Router, *http.Server) {
	r := mux.NewRouter()
	r.HandleFunc("/orgs", getOrgs).Methods("GET")
	r.HandleFunc("/orgs/{id}", getOrg).Methods("GET")
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
