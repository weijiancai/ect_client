package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"
	"flag"
	"bytes"

	_ "github.com/denisenkom/go-mssqldb"
)

type DB struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Format   string
}

func (d *DB) GetDb() *sql.DB {
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s;encrypt=disable", d.Host, d.Port, d.User, d.Password, d.Database))
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


	if d.Format == "sql" {
		return d.toSql("table_name", list)
	}

	return d.toJson(list)
}

func (db *DB) toJson(list []map[string]interface{}) string {
	da, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err.Error)
		return err.Error()
	}

	return string(da)
}


func (db *DB) toSql(table string, list []map[string]interface{}) string {
	var buf bytes.Buffer
	var valBuf bytes.Buffer
	
	for _, data := range list {
		buf.WriteString("INSERT INTO ")
		buf.WriteString(table)
		buf.WriteString("(")
		var i int
		for key, value := range data {
			buf.WriteString(key)
			switch v := value.(type) {
			case nil:
				valBuf.WriteString("NULL")
			case []byte:
				valBuf.WriteString("NULL")
			case time.Time:
				valBuf.WriteString(v.Format("2016-01-02 15:05:05.999"))
			case int:
				valBuf.WriteString(fmt.Sprintf("%v", value))
			case int16:
				valBuf.WriteString(fmt.Sprintf("%v", value))
			case int32:
				valBuf.WriteString(fmt.Sprintf("%v", value))
			case int64:
				valBuf.WriteString(fmt.Sprintf("%v", value))
			case int8:
				valBuf.WriteString(fmt.Sprintf("%v", value))
			case float32:
				valBuf.WriteString(fmt.Sprintf("%v", value))
			case float64:
				valBuf.WriteString(fmt.Sprintf("%v", value))
			default:
				valBuf.WriteString(fmt.Sprintf("'%s'", value))
			}
			if i < len(data) - 1 {
				buf.WriteString(",")
				valBuf.WriteString(",")
			}
			
			i += 1
		}
		buf.WriteString(") VALUES (")
		buf.WriteString(valBuf.String())
		buf.WriteString(");")
	}

	return buf.String()
}

func dbHandler(w http.ResponseWriter, req *http.Request) {

}

func main() {
	host := flag.String("H", "127.0.0.1", "db host")
	port := flag.Int("P", 1433, "db port")
	userName := flag.String("u", "", "userName")
	password := flag.String("p", "", "password")
	database := flag.String("d", "", "database")
	sql := flag.String("s", "", "sql 语句")
	format := flag.String("f", "json", "格式化输出，默认json，支持json，sql")
	flag.Parse()

	db := DB{Host: *host, Port: *port, User: *userName, Password: *password, Database: *database, Format: *format}

	result := db.Query(*sql)
	// println(result)

	fmt.Fprintln(os.Stdout, result)
	// http.HandleFunc("/db", dbHandler)
	// http.ListenAndServe(":8001", nil)

}
