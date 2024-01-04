package update

import (
	"backend/db"
	goUtils "backend/goutilspkg"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type NewHabit struct {
	Id        int
	Name      string   `json:"name"`
	Frequency []string `json:"frequency"`
}

type GetHabit struct {
	Id int `json:"id,omitempty"`
}

// UpdateHandler handles the update habit endpoint.
func UpdateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		var habit NewHabit
		err := json.NewDecoder(r.Body).Decode(&habit)
		if err != nil {
			http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		createdHabit, err := CreateHabit(context.Background(), &habit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(createdHabit)
	}

	if r.Method == http.MethodGet {
		var habit GetHabit
		err := json.NewDecoder(r.Body).Decode(&habit)
		if errors.Is(err, io.EOF) {
			http.Error(w, "Request body is empty, at least specify {}", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if habit.Id == 0 {
			getAllHabits, err := GetAllHabits(context.Background())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(getAllHabits)
		} else {
			getHabit, err := GetSpecifcHabit(context.Background(), habit.Id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(getHabit)
		}
	}

}

func InitializeDatabase() (db.DatabaseConfig, *db.Database, string) {
	config := db.Db_schema()
	database := db.NewDatabase(config)
	if database.Connection == nil {
		log.Fatal("Failed to establish a database connection")
	}
	tableName := "habits"

	if _, exists := config.Database.Tables[tableName]; !exists {
		log.Fatal("Habits table not found in configuration")
	}

	return config, database, tableName
}

// CreateHabit creates a habit in the database & handles logic.
func CreateHabit(ctx context.Context, create *NewHabit) (string, error) {
	validFrequency := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	create.Frequency = goUtils.RemoveListDups(create.Frequency)
	if len(create.Frequency) > 7 || len(create.Frequency) <= 0 {
		return "", errors.New("frequency length is invalid: " + strconv.Itoa(len(create.Frequency)))
	}
	if len(create.Name) > 100 || len(create.Name) <= 0 {
		return "", errors.New("name length is invalid: " + strconv.Itoa(len(create.Name)))
	}
	for _, day := range create.Frequency {
		if !goUtils.IsInList(day, validFrequency) {
			return "", errors.New(day + " is not a valid day of the week & unsupported")
		}
	}
	response, err := InsertData(*create)
	if err != nil {
		return "", errors.New("Error inserting data: " + err.Error())
	}

	return response, nil
}

// InsertData inserts data into the database.
func InsertData(createHabit NewHabit) (status string, err error) {
	_, db, tableName := InitializeDatabase()
	defer db.Connection.Close()

	frequencyJSON, err := json.Marshal(createHabit.Frequency)
	if err != nil {
		return "Error marshalling frequency", err
	}
	if tableName != "" {
		sqlStatement := fmt.Sprintf(`INSERT INTO %s (id, name, frequency) VALUES (?, ?, ?)`, tableName)
		fmt.Println("SQL Statement:", sqlStatement)
		fmt.Println("Habit Values:", nil, createHabit.Name, string(frequencyJSON))

		_, err = db.Connection.Exec(sqlStatement, nil, createHabit.Name, string(frequencyJSON))
		if err != nil {
			return "Error inserting data", err
		}
	} else {
		fmt.Println("Habit table or columns not found in configuration")
	}

	return "Successfully loaded data into database", nil
}

// GetAllHabits gets all habits from the database.
func GetAllHabits(ctx context.Context) ([]NewHabit, error) {
	_, db, tableName := InitializeDatabase()
	defer db.Connection.Close()

	tableData, err := db.FetchFromTable(tableName)
	if err != nil {
		return nil, err
	}

	var habits []NewHabit
	for _, row := range *tableData {
		var frequency []string
		errFreq := json.Unmarshal([]byte(row.Frequency), &frequency)
		if errFreq != nil {
			return nil, fmt.Errorf("failed to parse frequency JSON: %v", errFreq)
		}

		var habit NewHabit

		habit.Id = row.ID
		habit.Name = row.Name
		habit.Frequency = frequency

		habits = append(habits, habit)
	}
	if habits == nil {
		return nil, errors.New("no habits found")
	}

	return habits, nil
}

// GetSpecifcHabit gets a specific habit from the database.
func GetSpecifcHabit(ctx context.Context, id int) (*NewHabit, error) {
	_, db, tableName := InitializeDatabase()
	defer db.Connection.Close()

	tableData, err := db.FetchFromTable(tableName)
	if err != nil {
		return nil, err
	}
	var habit NewHabit
	found := false
	for _, row := range *tableData {

		var frequency []string
		errFreq := json.Unmarshal([]byte(row.Frequency), &frequency)
		if errFreq != nil {
			return nil, fmt.Errorf("failed to parse frequency JSON: %v", errFreq)
		}
		if row.ID == id {
			log.Println("Found habit with ID:", row.ID)
			habit.Id = row.ID
			habit.Name = row.Name
			habit.Frequency = frequency
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("no habits found")
	}

	return &habit, nil
}
