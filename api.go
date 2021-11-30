package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3" // for local testing
	"log"
	"net/http"
	"os"
	"strconv"
)

type Employee struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

type Employees struct {
	Employees []Employee `json:"employees"`
}

func main() {
	db := getDBHandle()
	r := mux.NewRouter()
	r.HandleFunc("/employees", createEmployee(db)).Methods("POST").Headers("Content-Type", "application/x-www-form-urlencoded")
	r.HandleFunc("/employees", getEmployees(db)).Methods("GET")
	poisoned := false

	r.HandleFunc("/poisonPill", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		poisonedParam := req.PostFormValue("Poisoned")
		if poisonedParam == "" {
			poisonedParam = "1"
		}
		if poisonedParam == "1" {
			poisoned = true
		} else if poisonedParam == "0" {
			poisoned = false
		} else {
			http.Error(w, "Invalid 'Poisoned' param", 400)
			return
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, "Set poisoned to %v\n", poisoned)
	})

	// Healthcheck endpoint for kubernetes
	r.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		status := http.StatusOK
		responseStr := "OK"
		if poisoned {
			status = http.StatusInternalServerError
			responseStr = "UNHEALTHY"
		}
		w.WriteHeader(status)
		if _, err := fmt.Fprintf(w, "%s\n", responseStr); err != nil {
			log.Printf("Failed to write healthcheck response: %s\n", err)
		}
	}).Methods("GET")
    log.Fatal(http.ListenAndServe(":9090", r))
}

func getDBHandle() *sql.DB {
	dbConnStr := os.Getenv("DB_CONNECTION_STRING")
	if dbConnStr == "" {
		log.Fatal("DB connection string must be specified!")
	}
	driverName := os.Getenv("DB_DRIVER_NAME")
	if driverName == "" {
		log.Fatal("DB driver name must be specified!")
	}
	db, err := sql.Open(driverName, dbConnStr)
	if err != nil {
		log.Fatalf("Error opening DB connection: %s", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func createEmployee(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r * http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		name := r.PostFormValue("Name")
		if name == "" {
			http.Error(w, "'Name' form param must be specified", http.StatusBadRequest)
			return
		}
		age := r.PostFormValue("Age")
		if age == "" {
			http.Error(w, "'Age' form param must be specified", http.StatusBadRequest)
			return
		}
		ageI, err := strconv.Atoi(age)
		if err != nil {
			http.Error(w, "'Age' must be a valid number", http.StatusBadRequest)
			return
		}
		if ageI < 18 || ageI > 150 {
			http.Error(w, "'Age' must be between 18 and 150", http.StatusBadRequest)
			return
		}
		if _, err = db.ExecContext(r.Context(), `INSERT INTO employees (name, age) VALUES ($1, $2)`, name, ageI); err != nil {
			http.Error(w, fmt.Sprintf("DB error: %s", err), http.StatusInternalServerError)
			return
		}
		emp := Employee{
			Name: name,
			Age:  ageI,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(emp); err != nil {
			log.Fatalf("Failed to encode json, programming error: %s", err)
		}
	}
}

func getEmployees(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.QueryContext(r.Context(), `SELECT name, age FROM employees LIMIT 10`)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying database: %s", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var result []Employee
		for rows.Next() {
			emp := Employee{}
			if err := rows.Scan(&emp.Name, &emp.Age); err != nil {
				http.Error(w, fmt.Sprintf("Error scanning rows: %s", err), http.StatusInternalServerError)
				return
			}
			result = append(result, emp)
		}
		employees := Employees{result}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(employees); err != nil {
			log.Fatalf("Failed to json encode employees: %s", err)
		}
	}
}
