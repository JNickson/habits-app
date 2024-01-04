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
	// Delete after viewing db content from a table (habit)
	// query := "SELECT * FROM habit"
	// rows, _ := db.Connection.Query(query)
	// defer rows.Close()
	// cols, _ := rows.Columns()
	// fmt.Println("Columns:", cols)

	// for rows.Next() {
	// 	columns := make([]interface{}, len(cols))
	// 	columnPointers := make([]interface{}, len(cols))
	// 	for i := range columns {
	// 		columnPointers[i] = &columns[i]
	// 	}
	// 	if err := rows.Scan(columnPointers...); err != nil {
	// 		return
	// 	}
	// 	rowMap := make(map[string]interface{})
	// 	for i, colName := range cols {
	// 		val := columnPointers[i].(*interface{})
	// 		rowMap[colName] = *val
	// 	}

	// 	fmt.Println(rowMap)
	// }
	// to here

	r.HandleFunc("/habits", update.UpdateHandler)
	r.HandleFunc("/help", help.HelpHandler)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}
