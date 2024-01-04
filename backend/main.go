package main

import (
	"log"
	"net/http"

	"backend/db"
	"backend/help"
	"backend/update"

	"github.com/gorilla/mux"
)

func main() {
	config := db.Db_schema()
	r := mux.NewRouter()
	db := db.NewDatabase(config)
	defer db.Connection.Close()
	for tableName, tableDetails := range config.Database.Tables {
		db.CreateTable(tableName, tableDetails)
	}

	r.HandleFunc("/habits", update.UpdateHandler)
	r.HandleFunc("/help", help.HelpHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
