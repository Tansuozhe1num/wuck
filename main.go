package main

import (
	"log"
	"os"

	"luangao/biu"
	"luangao/server/httpserver"
)

func main() {
	biu.MustInit()

	r := httpserver.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server started on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
