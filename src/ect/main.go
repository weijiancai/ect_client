package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"syscall"
	"unsafe"

	_ "github.com/denisenkom/go-mssqldb"
)

type DllCall struct {
	DllName   string        `json:"DllName"`
	FuncName  string        `json:"FuncName"`
	FuncParam []interface{} `json:"FuncParam"`
}

// var times int64

// func init() {
// 	flag.Int64Var(&times, "t", 100, "-t times")
// 	flag.Parse()
// }

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
	//host := flag.String("-h", "", "db host")
	//port := flag.Int("-P", 1433, "db port")
	//userName := flag.String("-u", "", "userName")
	//password := flag.String("-p", "", "password")
	//database := flag.String("-d", "", "database")
	//flag.Parse()
	// 读取当前目录配置文件config.properties，如果不存在则创建一个

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
	error := http.ListenAndServe(":8001", nil)
	if error != nil {
		log.Println(error)
	}
}
