package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"database/sql"
	"fmt"

	"github.com/gorilla/mux"

	"github.com/rs/cors"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Database struct {
	ID     int      `json:"id,omitempty"`
	Dbname string   `json:"dbname,omitempty"`
	Images []string `json:"images,omitempty"`
}

const (
	host     = "localhost"
	port     = 3001
	user     = "docker"
	password = "docker"
	dbname   = "docker"
)

var db *sql.DB

func databaseConnection() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error

	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
		return
	}
	var errorOnCreate error

	_, errorOnCreate = db.Exec(
		"CREATE TABLE DATABASES (" +
			"ID serial PRIMARY KEY," +
			"Dbname VARCHAR ( 50 ) UNIQUE NOT NULL," +
			"images TEXT []" +
			");")

	if errorOnCreate != nil {
		_, errorOnGetRows := db.Query("SELECT ID, Dbname , images  FROM DATABASES")

		if errorOnGetRows != nil {
			panic(errorOnCreate)
			return
		}
	}

	fmt.Println("Successfully connected!")
}

func getDataBase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	fmt.Printf("Rota getAcessada")
	registers, errorOnGetRows := db.Query("SELECT ID, Dbname , images  FROM DATABASES")

	if errorOnGetRows != nil {
		panic(errorOnGetRows)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var databases []Database = make([]Database, 0)

	for registers.Next() {
		var database Database
		scanErorr := registers.Scan(&database.ID, &database.Dbname, pq.Array(&database.Images))
		if database.Images == nil {
			databaseImages := []string{}
			database.Images = databaseImages
		}
		if scanErorr != nil {
			panic(scanErorr)
			continue
		}

		databases = append(databases, database)
	}

	closeRergistersError := registers.Close()

	if closeRergistersError != nil {
		panic(closeRergistersError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Content-Type","text/html")
	json.NewEncoder(w).Encode(databases)
}

func postDataBase(w http.ResponseWriter, r *http.Request) {
	body, erro := ioutil.ReadAll(r.Body)

	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var newDataBase Database

	json.Unmarshal(body, &newDataBase)
	fmt.Println(fmt.Printf("Array de images = %v", newDataBase.Images))
	_, execError := db.Exec("INSERT INTO DATABASES (Dbname, images) VALUES ($1, $2);", newDataBase.Dbname, pq.Array(newDataBase.Images))

	if execError != nil {
		panic(execError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newDataBase)
}

func deleteDataBase(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	id, _ := strconv.Atoi(vars["databaseID"])

	registers := db.QueryRow("SELECT ID FROM DATABASES WHERE ID = $1", id)

	var database Database

	scanErorr := registers.Scan(&database.ID, &database.Dbname, pq.Array(&database.Images))

	w.Header().Add("Content-Type","text/html")   
	w.Header().Set("Content-Type", "application/json")
	if scanErorr != nil {
		panic(scanErorr)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, execError := db.Exec("DELETE FROM DATABASES WHERE ID = $1", database.Images, id)

	if execError != nil {
		panic(execError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func putDataBase(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	id, _ := strconv.Atoi(vars["databaseID"])

	registers := db.QueryRow("SELECT ID, Dbname , images  FROM DATABASES WHERE ID = $1", id)

	var database Database

	scanErorr := registers.Scan(&database.ID, &database.Dbname, pq.Array(&database.Images))

	w.Header().Add("Content-Type","text/html")
	w.Header().Set("Content-Type", "application/json")
	if scanErorr != nil {
		panic(scanErorr)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var mdifiedDatabase Database

	json.Unmarshal(body, &mdifiedDatabase)

	_, execError := db.Exec("UPDATE DATABASES SET Dbname = $1, images = $2 WHERE ID = $3", mdifiedDatabase.Dbname, pq.Array(mdifiedDatabase.Images), id)
	if execError != nil {
		panic(execError)
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(mdifiedDatabase)

}

func searchDataBase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["databaseID"])

	registers := db.QueryRow("SELECT ID, Dbname , images  FROM DATABASES WHERE ID = $1", id)

	var database Database

	scanErorr := registers.Scan(&database.ID, &database.Dbname, pq.Array(&database.Images))

	                                         
	 w.Header().Add("Content-Type","text/html")
	w.Header().Set("Content-Type", "application/json")
	
	if scanErorr != nil {
		panic(scanErorr)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(database)

}

func configureServer() {

	router := mux.NewRouter()
	router.HandleFunc("/api/databases/{databaseID}/", searchDataBase).Methods("GET")
	router.HandleFunc("/api/databases/", postDataBase).Methods("POST")
	router.HandleFunc("/api/databases/", getDataBase).Methods("GET")
	router.HandleFunc("/api/databases/{databaseID}/", putDataBase).Methods("PUT")
	router.HandleFunc("/api/databases/{databaseID}/", deleteDataBase).Methods("DELETE")
	fmt.Printf("Rota configurada")
c := cors.New(cors.Options{
	AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		ExposedHeaders:   []string{"Link", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":3003", handler))
}

func main() {
	databaseConnection()
	configureServer()
}
