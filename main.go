package main

import (
	"database/sql"
	"encoding/json"
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
	TechStack *string    `json:"tech_stack"`
	Topics    *string    `json:"topics"`
	ShortDesc *string    `json:"short_desc"`
	Link      *string    `json:"link"`
	ImgUrl    *string    `json:"img_url"`
	Website   *string    `json:"website"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
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

func getOrgs(w http.ResponseWriter, r *http.Request) {
	// variables
	var orgs []Organization
	var response []byte

	// database setup
	db := setupDB()
	defer db.Close()

	// getting all organizations
	data, err := db.Query("SELECT * FROM organizations")
	checkErr(err)

	// looping through all organizations and storing them in an array
	for data.Next() {
		// variables for each column
		var name string
		var techStack, topics, shortDesc, link, imgUrl, website *string
		var year int
		var id uuid.UUID
		var createdAt, updatedAt time.Time
		var deletedAt *time.Time

		// copying data from each row to the corresponding variables
		err = data.Scan(&id, &name, &year, &techStack, &topics, &shortDesc, &link, &imgUrl, &website, &createdAt, &updatedAt, &deletedAt)
		checkErr(err)

		// appending the data to the array
		orgs = append(orgs, Organization{Id: id, Name: name, Year: year, TechStack: techStack, Topics: topics, ShortDesc: shortDesc, Link: link, ImgUrl: imgUrl, Website: website})
	}

	// converting the array to json
	response, err = json.Marshal(JsonResponse{Type: "success", Message: "Organizations retrieved successfully.", Count: len(orgs), Data: orgs})
	checkErr(err)

	// writing the json to the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func router() (*mux.Router, *http.Server) {
	r := mux.NewRouter()
	r.HandleFunc("/orgs", getOrgs).Methods("GET")

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
