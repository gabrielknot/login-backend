package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"database/sql"
	"fmt"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	ID     	  int    `json:"id,omitempty"`
	FirstName string `json:"id,omitempty"`
  	LastName  string `json:"id,omitempty"`
  	Email 	  string `json:"id,omitempty"`
}

const (
	Frontport = os.Getenv("FRONT_PORT")
	host      = os.Getenv("HOST")
	port      = os.Getenv("MYSQL_PORT")
	user      = os.Getenv("MYSQL_USER")
	password  = os.Getenv("MYSQL_PASSWORD")
	dbname  = os.Getenv("MYSQL_DATABASE")
)

var db *sql.DB

func databaseConnection() {
	mysqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error

	db, err = sql.Open("mysql", mysqlInfo)

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
			"FirstName  VARCHAR ( 50 ) UNIQUE NOT NULL," +
			"LastName VARCHAR ( 50 ) UNIQUE NOT NULL," +
			"Email VARCHAR ( 50 ) UNIQUE NOT NULL," +
			");")

	if errorOnCreate != nil {
		_, errorOnGetRows := db.Query("SELECT ID, FirstName , LastName , Email FROM DATABASES")

		if errorOnGetRows != nil {
			panic(errorOnCreate)
			return
		}
	}

	fmt.Println("Successfully connected!")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:" + Frontport)
	fmt.Printf("Rota getAcessada")
	registers, errorOnGetRows := db.Query("SELECT ID, FirstName , LastName , Email FROM DATABASES")

	if errorOnGetRows != nil {
		panic(errorOnGetRows)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var databases []Database = make([]Database, 0)

	for registers.Next() {
		var database Database
		scanErorr := registers.Scan(&database.ID, &database.FirstName,&database.LastName,&database.Email)
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

func postUser(w http.ResponseWriter, r *http.Request) {
	body, erro := ioutil.ReadAll(r.Body)

	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var newUser Database

	json.Unmarshal(body, &newUser)
	_, execError := db.Exec("INSERT INTO DATABASES (FirstName, LastName , Email) VALUES (?, ?, ?);", newUser.FirstName , newUser.LastName , newUser.Email)

	if execError != nil {
		panic(execError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newUser)
}

//func deleteUser(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:" + Frontport)
//	id, _ := strconv.Atoi(vars["databaseID"])
//
//	registers := db.QueryRow("SELECT ID FROM DATABASES WHERE ID = ?", id)
//
//	var database Database
//
//
//	scanErorr := registers.Scan(&database.ID, &database.FirstName,&database.LastName,&database.Email))
//	w.Header().Add("Content-Type","text/html")   
//	w.Header().Set("Content-Type", "application/json")
//	if scanErorr != nil {
//		panic(scanErorr)
//		w.WriteHeader(http.StatusNoContent)
//		return
//	}
//
//	w.WriteHeader(http.StatusNoContent)
//
//}
//
////func putUser(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:" + Frontport)
//	id, _ := strconv.Atoi(vars["databaseID"])
//
//	registers := db.QueryRow("SELECT ID, FirstName , LastName , Email  FROM DATABASES WHERE ID = ?", id)
//
//	var database Database
//
//
//	scanErorr := registers.Scan(&database.ID, &database.FirstName,&database.LastName,&database.Email))
//	w.Header().Add("Content-Type","text/html")
//	w.Header().Set("Content-Type", "application/json")
//	if scanErorr != nil {
//		panic(scanErorr)
//		w.WriteHeader(http.StatusNoContent)
//		return
//	}
//
//	body, _ := ioutil.ReadAll(r.Body)
//	var mdifiedDatabase Database
//
//	json.Unmarshal(body, &mdifiedDatabase)
//
//	_, execError := db.Exec("UPDATE DATABASES SET FirstName = ?, LastName , Email = ? WHERE ID = ?", mdifiedDatabase.FirstName , mdifiedDatabase.LastName , mdifiedDatabase.Email , id)
//	if execError != nil {
//		panic(execError)
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//
//	json.NewEncoder(w).Encode(mdifiedDatabase)
//
//}
//
////func searchUser(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:" + Frontport)
//	vars := mux.Vars(r)
//	id, _ := strconv.Atoi()
//
//	registers := db.QueryRow("SELECT ID, FirstName , LastName , Email  FROM DATABASES WHERE ID = ?", id)
//
//	var database Database
//
//	scanErorr := registers.Scan(&database.ID, &database.FirstName, &database.LastName)
//
//	                                         
//	 w.Header().Add("Content-Type","text/html")
//	w.Header().Set("Content-Type", "application/json")
//	
//	if scanErorr != nil {
//		panic(scanErorr)
//		w.WriteHeader(http.StatusNoContent)
//		return
//	}
//
//	w.WriteHeader(http.StatusFound)
//	json.NewEncoder(w).Encode(database)
//
//}
//
func configureServer() {

	router := mux.NewRouter()
	router.HandleFunc("/api/databases/", postUser).Methods("POST")
	router.HandleFunc("/api/databases/", getUser).Methods("GET")
	//router.HandleFunc("/api/databases/{databaseID}/", searchUser).Methods("GET")
	//router.HandleFunc("/api/databases/{databaseID}/", putUser).Methods("PUT")
	//router.HandleFunc("/api/databases/{databaseID}/", deleteUser).Methods("DELETE")
	fmt.Printf("Rota configurada")
c := cors.New(cors.Options{
	AllowedOrigins:   []string{"http://localhost:" + Frontport},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		ExposedHeaders:   []string{"Link", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":3080", handler))
}

func main() {
	databaseConnection()
	configureServer()
}
