package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"database/sql"
	"fmt"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	ID     	  int    `json:"id,omitempty"`
	FirstName string `json:"firstname,omstempty"`
  	LastName  string `json:"lastname,omitempty"`
  	Email 	  string `json:"email,omitempty"`
}

var (
	Frontport  = os.Getenv("FRONT_PORT")
	host       = os.Getenv("HOST")
	port       = os.Getenv("MYSQL_PORT")
	user       = os.Getenv("MYSQL_USER")
	password   = os.Getenv("MYSQL_PASSWORD")
	dbname     = os.Getenv("MYSQL_DATABASE")
)

var db *sql.DB

func databaseConnection() {
	db, err := sql.Open("mysql", user +":" + password + "@tcp(" + host + ":" + port + ")/" + dbname )

	    // if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
    	}
	err = db.Ping()
	if err != nil {
		panic(err)
		return
	}
	//query :="CREATE TABLE IF NOT EXISTS databases"+
	//"(" +
	//"id INT AUTO_INCRMENT, " +
	//"firstname VARCHAR(255), " +
	//"lastname VARCHAR(255), " +
	//"email VARCHAR(255) " +
	//");"
	query := "CREATE TABLE IF NOT EXISTS users (" +
	"id INT AUTO_INCREMENT PRIMARY KEY," +
	"FirstName VARCHAR(255) NOT NULL," +
	"LastName VARCHAR(255) NOT NULL," +
	"Email VARCHAR(255) NOT NULL" +
	")  ENGINE=INNODB;" 
	var errorOnCreate error
	//query := `create table if not exists databases(id int primary key auto_increment, firstname text not null, lastname text not null, email text unique not null);`
	_, errorOnCreate = db.Exec(query)

	fmt.Println("Successfully created!")
	if errorOnCreate != nil {
		_, errorOnGetRows := db.Query("select id, FirstName , LastName , Email from users")

		if errorOnGetRows != nil {
			panic(errorOnCreate)
			return
		}
	}

	fmt.Println("Successfully connected!")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://" + host + ":" + Frontport)
	fmt.Printf("Rota getAcessada")
	registers, errorOnGetRows := db.Query("select id, FirstName , LastName , Email from users")

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
	_, execError := db.Exec("insert into users (FirstName, LastName , Email) values (?, ?, ?);", newUser.FirstName , newUser.LastName , newUser.Email)

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
	router.HandleFunc("/api/users/", postUser).Methods("POST")
	router.HandleFunc("/api/users/", getUser).Methods("GET")
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
