package main

import (
	"database/sql"
	_ "ects"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"syscall"
	"time"
	"unsafe"

	_ "github.com/denisenkom/go-mssqldb"
)

type DllCall struct {
	DllName   string        `json:"DllName"`
	FuncName  string        `json:"FuncName"`
	FuncParam []interface{} `json:"FuncParam"`
}

var times int64

func init() {
	flag.Int64Var(&times, "t", 100, "-t times")
	flag.Parse()
}

func DllCallHandle(w http.ResponseWriter, req *http.Request) {
	fmt.Println("HelloWorld")
	// fmt.Println(string(req.Body))
	jsonstr, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(jsonstr))
	js := new(DllCall)

	json.Unmarshal(jsonstr, js)
	fmt.Println(js)
	var param []uintptr = make([]uintptr, 0)
	for _, v := range js.FuncParam {
		fmt.Println(v)
		fmt.Println(reflect.TypeOf(v))
		fmt.Println(reflect.TypeOf(v).Kind())
		if reflect.TypeOf(v).Kind() == reflect.Float64 {
			param = append(param, uintptr(v.(float64)))
		} else if reflect.TypeOf(v).Kind() == reflect.String {
			param = append(param, uintptr(unsafe.Pointer(syscall.StringBytePtr(v.(string)))))
		}

		// ptr :=1 unsafe.Pointer(&v)
		// s := (uintptr)(ptr)
		// param = append(param, uintptr(v))
	}

	// fmt.Println(param111111111111111111111111111111111111111111111111111111)
	dll := syscall.NewLazyDLL(js.DllName)
	proc := dll.NewProc(js.FuncName)
	r1, r2, err := proc.Call(param...)
	fmt.Println(r1)
	fmt.Println(r2)
	fmt.Println(err)

}

func IndexHandle(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

func main() {
	fmt.Println(times)
	fmt.Println("DB Test start......")

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

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute")
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(2)) //12 Red light

	fmt.Println("hello world")

	handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7)) //White dark
	CloseHandle := kernel32.NewProc("CloseHandle")
	CloseHandle.Call(handle)

	fmt.Printf("\x1b[0;%dmhello world 30: 黑 \n\x1b[0m", 0)
	fmt.Printf("\x1b[0;%dmhello world 31: 红 \n\x1b[0m", 1)
	fmt.Printf("\x1b[0;%dmhello world 32: 绿 \n\x1b[0m", 2)
	fmt.Printf("\x1b[0;%dmhello world 33: 黄 \n\x1b[0m", 3)
	fmt.Printf("\x1b[0;%dmhello world 34: 蓝 \n\x1b[0m", 4)
	fmt.Printf("\x1b[0;%dmhello world 35: 紫 \n\x1b[0m", 5)
	fmt.Printf("\x1b[0;%dmhello world 36: 深绿 \n\x1b[0m", 6)
	fmt.Printf("\x1b[0;%dmhello world 37: 白色 \n\x1b[0m", 7)

	fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 47: 白色 30: 黑 \n", 47, 30)
	fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 46: 深绿 31: 红 \n", 46, 31)
	fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 45: 紫   32: 绿 \n", 45, 32)
	fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 44: 蓝   33: 黄 \n", 44, 33)
	fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 43: 黄   34: 蓝 \n", 43, 34)
	fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 42: 绿   35: 紫 \n", 42, 35)
	fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 41: 红   36: 深绿 \n", 41, 36)
	fmt.Printf("\x1b[%d;%dmhello world \x1b[0m 40: 黑   37: 白色 \n", 40, 37)

	http.HandleFunc("/dllcall", DllCallHandle)
	http.HandleFunc("/", IndexHandle)
	http.ListenAndServe(":8001", nil)
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
