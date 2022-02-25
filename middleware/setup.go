package middleware

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetEnv(key string) string {
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

func SetupDB() *sql.DB {
	// setup database
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable connect_timeout=15", GetEnv("DB_USER"), GetEnv("DB_PASSWORD"), GetEnv("DB_NAME"), GetEnv("HOST"), GetEnv("DB_PORT"))
	db, err := sql.Open(GetEnv("DATABASE"), connectionString)
	CheckErr(err)

	// ping the database
	err = db.Ping()
	CheckErr(err)

	return db
}
