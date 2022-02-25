// TODO: Add Caching

package main

import (
	"fmt"
	"gsoc-api/middleware"
	"gsoc-api/router"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// start the server
	_, server := router.Router()
	fmt.Printf("Server started at port %s\n", middleware.GetEnv("WEB_PORT"))
	log.Fatal(server.ListenAndServe())
}
