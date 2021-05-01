/*
-----UserDB.go------
all functions related to user db in mysql is hosted here
*/
package main

import (
	"database/sql"
	"fmt"
)

func queryUser(db *sql.DB, username string) (user, bool) {
	results, _ := db.Query("SELECT * FROM VenueDB.Users where Username=?", username)
	if results.Next() {
		var person user
		_ = results.Scan(&person.UserName, &person.Password)
		return person, true
	} else {
		return user{}, false
	}
}

func insertRecord(db *sql.DB, username, password string) int {
	results, err := db.Exec("INSERT INTO VenueDB.Users VALUES (?,?)", username, password)
	if err != nil {
		// fmt.Println(err)
		return 0
	} else {
		rows, _ := results.RowsAffected()
		return int(rows)
	}
}

func deleteUserDB(db *sql.DB, username string) bool {
	_, err := db.Exec("DELETE FROM VenueDB.Users WHERE Username= ?", username)
	if err != nil {
		return false
	} else {
		fmt.Println("Successfully Deleted!!")
		return true
	}
}
