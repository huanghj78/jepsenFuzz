package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "gitee.com/opengauss/openGauss-connector-go-pq"
)

func main() {
	connStr := "host=10.186.133.141 port=26000 user=testuser password=Chen0031 dbname=test sslmode=disable"
	db, err := sql.Open("opengauss", connStr)
	if err != nil {
		log.Fatal(err)
	}
	var date string
	err = db.QueryRow("select current_date ").Scan(&date)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(date)
}
