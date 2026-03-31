package main

import (
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()

	api.Init()

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
   	 log.Printf("Server startup error: %v", err)
    	return 
	}
}