package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {

	conndb, errdb := sql.Open("mssql", "server=127.0.0.1;user id=sa;password=P@ssw0rd")

	if errdb != nil {
		fmt.Println("Error !!!! ", errdb.Error())
	}

	var (
		sqlversion string
	)
	rows, err := conndb.Query("SELECT @@version")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err1 := rows.Scan(&sqlversion)
		if err1 != nil {
			log.Fatal(err1)
		}
		log.Println(sqlversion)
	}
	defer conndb.Close()
}
