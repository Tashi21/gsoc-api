// TODO: Add Caching
// TODO: Add GET route for individual organizations
// TODO: Add PATCH route for individual organizations
// TODO: Add DELETE route for individual organizations
// TODO: Add POST route for individual organizations
// TODO: Add params for all the above endpoints

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

func getOrgs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// variables
	var orgs []Organization
	var data *sql.Rows
	var err error

	// database setup
	db := setupDB()
	defer db.Close()

	params := r.URL.Query()
	n := params.Get("name")
	y := params.Get("year")

	checkParams := r.URL.Query()
	checkParams.Del("name")
	checkParams.Del("year")

	if len(checkParams) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{Type: "error", Message: "Invalid parameters.", Count: 0, Data: []Organization{}})
		return
	}

	if n != "" && y != "" {
		data, err = db.Query(fmt.Sprintf("SELECT id, name, year, link, website, created_at FROM organizations WHERE name LIKE '%%%s%%' AND year = %s", n, y))
		checkErr(err)
	} else if n != "" {
		// getting all organizations with name contaning given letters
		data, err = db.Query(fmt.Sprintf("SELECT id, name, year, link, website, created_at FROM organizations WHERE name LIKE '%%%s%%'", n))
		checkErr(err)
	} else if y != "" {
		// getting all organizations for a specific year
		data, err = db.Query(fmt.Sprintf("SELECT id, name, year, link, website, created_at FROM organizations WHERE year = %s", y))
		checkErr(err)
	} else {
		data, err = db.Query("SELECT id, name, year, link, website, created_at FROM organizations")
		checkErr(err)
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
		checkErr(err)

		// appending the data to the array
		// orgs = append(orgs, Organization{Id: id, Name: name, Year: year, TechStack: techStack, Topics: topics, ShortDesc: shortDesc, Link: link, ImgUrl: imgUrl, Website: website})
		orgs = append(orgs, Organization{Id: id, Name: name, Year: year, Link: link, Website: website, CreatedAt: createdAt})
	}

	// writing the json to the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JsonResponse{Type: "success", Message: "Organizations retrieved successfully.", Count: len(orgs), Data: orgs})
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
