package main

import (
	"log"
	"fmt"
	"net/http"

	handler "github.com/getshelf/backend/internal/handler/rest"
)

func main() {
	fmt.Println("OK")
	router := handler.RootRouter()

	fmt.Println("Server up and running: http://localhost:3001")

	if err := http.ListenAndServe(":3001", router); err != nil {
		log.Fatal(err)
	}
}
