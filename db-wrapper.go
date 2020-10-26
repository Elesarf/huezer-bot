package main

import (
	"database/sql"
	"log"
	"strings"
)

type baseWord struct {
	key      string
	value    string
	username string
}

func insertIntoDB(key string, value string, db *sql.DB, user string) error {

	// default value. fuck go
	if user == "" {
		user = "system"
	}

	log.Println("Insert into db")
	_, err := db.Exec("insert into dictionary (key, value, username) values ($1, $2, $3)", key, value, user)
	return err
}

// get value from db
func getValue(key string, db *sql.DB) string {
	// some -> Some
	key = strings.Title(key)
	log.Println("Get value to key: " + key)
	// query first row from db
	row := db.QueryRow("select * from dictionary where key=$1", key)
	// cast row to baseWord
	word := baseWord{}
	row.Scan(&word.key, &word.value, &word.username)

	log.Println(word)
	return word.value
}

// base have key?
func checkKey(key string, db *sql.DB) bool {
	// some -> Some
	key = strings.Title(key)
	log.Println("Check key: _" + key + "_")
	// get first row
	row := db.QueryRow("select * from dictionary where key=$1", key)
	// cast to baseWord
	word := baseWord{}
	row.Scan(&word.key, &word.value, &word.username)
	log.Println(word)

	if word.key != "" {
		return true
	}

	return false
}

// get aliases to cities
func getAliases(city string, db *sql.DB, user string) []string {
	// SomE -> some
	value := strings.ToLower(city)
	// some -> Some
	value = strings.Title(value)

	if !IsLetter(city) {
		return []string{}
	}

	log.Println("Collect values to: _" + value + "_" + ": _" + user + "_")

	// if user not defined, get all aliases
	var rows *sql.Rows
	var err error
	if user != "" {
		rows, err = db.Query("select * from dictionary where value=$1 OR key=$1 OR username=$2", value, user)
		_check(err)
	} else {
		rows, err = db.Query("select * from dictionary where value=$1 OR key=$1", value)
		_check(err)
	}
	defer rows.Close()

	// make city stringList
	values := []string{}
	for rows.Next() {
		var str string
		var strBee string
		var uName string
		err := rows.Scan(&str, &strBee, &uName)
		_check(err)
		values = append(values, str+" : "+strBee)
	}

	return values
}

// something strange i wanted
func getTransfer(key string, db *sql.DB) string {
	return ""
}
