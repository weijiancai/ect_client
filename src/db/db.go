package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func (d *DB) GetDb() *sql.DB {
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s;encrypt=disable", d.Host, d.Port, d.User, d.Password, d.Database))
	if err != nil {
		log.Fatalf("Open database error: %s\n", err)
	}
	return db
}

func (d *DB) Query(sql string) string {
	db := d.GetDb()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(sql)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	scanArgs := make([]interface{}, len(cols))
	data := make([]interface{}, len(cols))
	for i, _ := range scanArgs {
		scanArgs[i] = &data[i]
	}

	list := []map[string]interface{}{}
	for rows.Next() {
		rows.Scan(scanArgs...)
		record := make(map[string]interface{})
		for i, val := range scanArgs {
			switch v := (*(val.(*interface{}))).(type) {
			case nil:
				continue
			case bool:
				if v {
					record[cols[i]] = true
				} else {
					record[cols[i]] = false
				}
			case []byte:
				record[cols[i]] = string(v)
			case time.Time:
				record[cols[i]] = v.Format("2016-01-02 15:05:05.999")
			default:
				record[cols[i]] = v
			}
		}

		list = append(list, record)
	}

	da, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err.Error)
	}

	return string(da)
}

func dbHandler(w http.ResponseWriter, req *http.Request) {

}

func main() {
	db := DB{Host: "weiyi1998.com", Port: "18888", User: "sa", Password: "123!@#qwe", Database: "TEST"}
	result := db.Query("select top 2 * from db_product")
	println(result)

	http.HandleFunc("/db", dbHandler)
	http.ListenAndServe(":8001", nil)
}
