package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AaravMalani/cvwo-forum-backend/internal/auth"
	"github.com/AaravMalani/cvwo-forum-backend/internal/database"
	"github.com/AaravMalani/cvwo-forum-backend/internal/env"
	"github.com/AaravMalani/cvwo-forum-backend/internal/router"
)

func main() {
	env.Setup()
	database.Setup()
	auth.Setup()
	r := router.Setup()
	fmt.Print("Listening on port 8000 at http://localhost:8000!")
	log.Fatalln(http.ListenAndServe(":8000", r))
}
