// api connected to mongodb
package main

import (
	"fmt"
	"log"
	"net/http"

	"vishsec.dev/goapi/router"
)


func main(){
	fmt.Println("MongoDb API")

	fmt.Println("server getting started...")

	r := router.Router()
	log.Fatal(http.ListenAndServe(":4004", r))
	fmt.Println("listening at port 4004")

}


