package api

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

const CREATE_TABLE = `CREATE TABLE employees (name varchar(256) NOT NULL, age int NOT NULL)`

func createSqliteDB() *sql.DB {
	f, err := os.Create("emp-api.sqlite")
	if err != nil {
		log.Fatalf("Failed to create sqlite file: %s", err)
	}
	if err := f.Close(); err != nil {
		log.Fatalf("Error closing sqlite file: %s", err)
	}
	db, err := sql.Open("sqlite3", "./emp-api.sqlite")
	if err != nil {
		log.Fatalf("Error opening sqlite source: %s", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to open sqlite source: %s", err)
	}
	if _, err := db.Exec(CREATE_TABLE); err != nil {
		log.Fatalf("Failed to create table: %s", err)
	}

	return db
}

func TestApi(t *testing.T) {
	db := createSqliteDB()
	defer db.Close()
	defer func() { os.Remove("emp-api.sqlite")}()


	form := url.Values{"Name": {"juggernaut"}, "Age": {"37"}}
	req, err := http.NewRequest("POST", "/employees", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createEmployee(db))
	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != 201 {
		t.Fatalf("Got status %d, expected 201", status)
	}

}
