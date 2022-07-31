// @Title Movie API
// @Version 1.0
// @Description Creating a watchlist of movies
// @ContactName Sakshi
// @ContactEmail sakshirads@email.com
// @Server localhost:4000
// @BasePath  /api/v1

package main

import (
	"fmt"
	"log"
	router "mongoapi/Router"
	"net/http"
)

func main() {
	fmt.Println("MongoDB API")
	router := router.Router()
	fmt.Println("Server Starting...")
	log.Fatal(http.ListenAndServe(":4000", router))
	fmt.Println("Listening at port 4000 ...")
}
