package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Healthy")
}

func host(w http.ResponseWriter, _ *http.Request) {
	node := os.Getenv("MY_NODE_NAME")
	podIP := os.Getenv("MY_POD_IP")

	fmt.Fprintf(w,"NODE: %v, POD IP:%v",node, podIP)
}

func dataHandler(w http.ResponseWriter, _ *http.Request) {
	db := CreateCon()

	err := db.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "Database Connected")
	}
}

func main() {
	http.HandleFunc("/", hello)

	http.HandleFunc("/healthz", healthz)

	http.HandleFunc("/data", dataHandler)

	http.HandleFunc("/host", host)

	http.ListenAndServe("0.0.0.0:8080", nil)
}

/*Create sql database connection*/
func CreateCon() *sql.DB {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v?sslmode=disable", user, pass, host, port)

	fmt.Printf("Database Connection String: %v \n", connStr)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	return db
}
