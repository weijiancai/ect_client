package db

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	db, err := sql.Open("mssql", "server=weiyi1998.com;port=18888;user id=sa;password=123!@#qwe;database=TEST;encrypt=disable")
	if err != nil {
		log.Fatalf("Open database error: %s\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select top 2 * from db_product")
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

	// values := make([]sql.RawBytes, len(cols))
	// values := make([]interface{}, len(cols))
	// fmt.Println(values)

	for rows.Next() {
		// var h float64
		// if err := rows.Scan(&h); err != nil {
		// 	fmt.Println(err.Error())
		// }
		// fmt.Println(h)
		// map = make(map[string]interface{})
		rows.Scan(scanArgs...)
		PrintRow(scanArgs)
		// record := make(map[string]string)
		// for i, col := range values {
		// 	if col != nil {
		// 		record[cols[i]] = string(col.([]byte))
		// 	}
		// }
		// fmt.Println(record)
	}

	fmt.Println("finish")
	fmt.Println("DB Test end")
}

func PrintRow(colsdata []interface{}) {
	for _, val := range colsdata {
		fmt.Println(reflect.TypeOf(*(val.(*interface{}))))
		switch v := (*(val.(*interface{}))).(type) {
		case nil:
			fmt.Println("NULL")
		case bool:
			if v {
				fmt.Println("True")
			} else {
				fmt.Println("False")
			}
		case []byte:
			fmt.Println("string: " + string(v))
		case time.Time:
			fmt.Println(v.Format("2016-01-02 15:05:05.999"))
		default:
			fmt.Println(v)
		}
		fmt.Print("=================\n")
	}
	fmt.Println()
}
