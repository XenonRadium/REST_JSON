package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

var db *sql.DB

//	Struct for Versions
type Versions struct {
	Version     string
	Date        string
	Description string
}

func insertInDatabase(data Versions) error {
	_, err := db.Exec("INSERT versions(version, date, description) VALUES (?,?,?)",
		data.Version, data.Date, data.Description)

	return err
}

func getFromDatabase(version string, w http.ResponseWriter) error {
	dataRetrieved := Versions{}

	//	db.Query returns >1 rows, db.QueryRow returns 1 row
	err := db.QueryRow("SELECT * FROM versions WHERE version=?", version).Scan(&dataRetrieved.Version, &dataRetrieved.Date, &dataRetrieved.Description)

	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(&dataRetrieved)

	return err
}

func deleteFromDatabase(version string, w http.ResponseWriter) error {
	dataDeleted := Versions{}

	//db.QueryRow
	err := db.QueryRow("SELECT * FROM versions WHERE version=?", version).Scan(&dataDeleted.Version, &dataDeleted.Date, &dataDeleted.Description)
	if err != nil {
		return err
	} else {
		json.NewEncoder(w).Encode(&dataDeleted)
		db.Exec("DELETE FROM versions WHERE version=?", version)
	}
	return err
}

func updateInDatabase(version string, data Versions, w http.ResponseWriter) error {
	dataUpdated := Versions{}

	//db.Query
	err := db.QueryRow("SELECT * FROM versions WHERE version=?", version).Scan(&dataUpdated.Version, &dataUpdated.Date, &dataUpdated.Description)
	if err != nil {
		return err
	} else {
		json.NewEncoder(w).Encode(&dataUpdated)
		db.Exec("UPDATE versions SET date=?,description=? WHERE version=?", data.Date, data.Description, version)
	}
	return err

}

func userAddHandler(w http.ResponseWriter, r *http.Request) {
	reqbody, _ := ioutil.ReadAll(r.Body)

	var version Versions
	json.Unmarshal(reqbody, &version)

	//insert into database
	err := insertInDatabase(version)

	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte(`{"message":"success"}`))
}

func userGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	version := vars["version"]

	err := getFromDatabase(version, w)
	if err != nil {
		fmt.Println(err)
	}
}

func userDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	version := vars["version"]

	err := deleteFromDatabase(version, w)
	if err != nil {
		log.Print(err)
	}
}

func userUpdateHandler(w http.ResponseWriter, r *http.Request) {
	reqbody, _ := ioutil.ReadAll(r.Body)

	var body Versions
	json.Unmarshal(reqbody, &body)

	vars := mux.Vars(r)
	version := vars["version"]
	err := updateInDatabase(version, body, w)
	if err != nil {
		log.Print(err)
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/getMessage/{version}", userGetHandler).Methods("GET")
	myRouter.HandleFunc("/addMessage", userAddHandler).Methods("POST")
	myRouter.HandleFunc("/deleteMessage/{version}", userDeleteHandler).Methods("DELETE")
	myRouter.HandleFunc("/updateMessage/{version}", userUpdateHandler).Methods("PATCH")
	http.ListenAndServe(":8090", myRouter)
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:bernard@/messages")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	handleRequests()
}
