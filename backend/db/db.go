package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Connection *sql.DB
}

type DatabaseConfig struct {
	Database DatabaseDetails `json:"database"`
}

type DatabaseDetails struct {
	Name   string           `json:"name"`
	Tables map[string]Table `json:"tables"`
}

type Table struct {
	Columns map[string]Column `json:"columns"`
}

type Column struct {
	Type   string   `json:"type"`
	Extras []string `json:"extras"`
}

func Db_schema() DatabaseConfig {
	schema, err := os.ReadFile("./db/db_schema.json")
	if err != nil {
		log.Printf("Error reading db_schema.json: %v", err)
		return DatabaseConfig{}
	}
	var config DatabaseConfig
	err = json.Unmarshal(schema, &config)
	if err != nil {
		log.Printf("Error unmarshalling db_schema.json: %v", err)
		return DatabaseConfig{}
	}
	return config
}

func NewDatabase(config DatabaseConfig) *Database {
	db, err := sql.Open("sqlite3", config.Database.Name+".db")
	if err != nil {
		log.Fatal(err)
	}
	return &Database{
		Connection: db,
	}
}

func (db *Database) CreateTable(tableName string, table Table) {
	var columnDefs []string
	for columnName, column := range table.Columns {

		columnDef := fmt.Sprintf("%s %s", columnName, column.Type)
		if len(column.Extras) > 0 {
			columnDef += " " + strings.Join(column.Extras, " ")
		}
		columnDefs = append(columnDefs, columnDef)
	}
	createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);",
		tableName, strings.Join(columnDefs, ", "))

	log.Println("Create Table SQL:", createTableSQL)
	_, err := db.Connection.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

type Record struct {
	ID        int
	Name      string
	Frequency string
}

func (db *Database) FetchFromTable(tableName string) (*[]Record, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []Record

	for rows.Next() {
		var r Record
		err := rows.Scan(&r.ID, &r.Name, &r.Frequency)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	log.Println("Records Data:", records)

	return &records, nil
}
