package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"ect"
	_ "github.com/denisenkom/go-mssqldb"
	"strconv"
)

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
	js := new(ect.DllCall)

	json.Unmarshal(jsonstr, js)
	fmt.Println(js)

	js.Call()
}

func IndexHandle(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

func main() {
	ect.Info("系统启动..........")
	ect.Info("当前工作路径：" + ect.GetCurrentPath())

	if ect.ConfigInfo.LocalPort == 0 {
		ect.ConfigInfo.LocalPort = 12233
	}
	ect.Info("监听端口：" + strconv.Itoa(ect.ConfigInfo.LocalPort))

	http.HandleFunc("/dllcall", DllCallHandle)
	http.HandleFunc("/", IndexHandle)
	err := http.ListenAndServe(":" + strconv.Itoa(ect.ConfigInfo.LocalPort), nil)
	if err != nil {
		ect.Error(err)
	}
}
