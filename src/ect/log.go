package ect

import (
	"log"
	"io/ioutil"
	"os"
	"strings"
	"path/filepath"
	"time"
	"runtime/debug"
	"fmt"
)

var LogMaxDate string
var logger *log.Logger

func init() {
	path := GetCurrentPath()
	logFiles, err :=ListDir(path, "log")
	check(err)
	var hasFile = false
	var logFile *os.File
	for _, file := range logFiles {
		var dateStr = strings.Replace(file, "log_", "", -1)
		dateStr = strings.Replace(dateStr, ".log", "", -1)
		if LogMaxDate == "" {
			LogMaxDate = time.Now().Format("2006-01-02")
		} else {
			if dateAfter(dateStr, LogMaxDate) {
				LogMaxDate = dateStr
			}
		}

		if dateStr == LogMaxDate {
			hasFile = true
			fmt.Println("file = ", file)
			fmt.Println("path = ", path)
			logFilePath := path + "/" + file
			logfile, err := os.OpenFile(logFilePath, os.O_APPEND, 0666)
			check(err)
			logFile = logfile
		}
	}

	fmt.Println("maxDate = ", LogMaxDate)

	if !hasFile {
		logFilePath := path + "\\log_" + LogMaxDate + ".log"
		logfile, err := os.Create(logFilePath)
		check(err)
		logFile = logfile
	}

	logger = log.New(logFile, "", log.LstdFlags)
}

func dateAfter(first string, second string) bool {
	var fdate, _ = time.Parse("2006-01-02", first)
	var sdate, _ = time.Parse("2006-01-02", second)
	fmt.Println("first = ", fdate)
	fmt.Println("second = ", sdate)
	return fdate.After(sdate)
}

func Error(err error) {
	if err == nil || logger == nil{
		return
	}
	logger.SetPrefix("ERROR ")
	logger.Println(err, "\r\n", fmt.Sprintf("%s", debug.Stack()))
	dl := DllCall{DllName: "user32.dll", FuncName: "MessageBoxA", FuncParam: []interface{}{0, err.Error(), "Text", 0}}
	//dl := new(DllCall)
	//dl.DllName = "user32.dll"
	//dl.FuncName = "MessageBoxA"
	//dl.FuncParam = []interface{}{0, err.Error(), "Text", 0}

	dl.Call()
}
func ErrorStr(err string) {
	if err == "" || logger == nil{
		return
	}
	logger.SetPrefix("ERROR ")
	logger.Print(err, "\r\n")
	logger.Println(err, "\r\n", fmt.Sprintf("%s", debug.Stack()))
	dl := DllCall{DllName: "user32.dll", FuncName: "MessageBoxA", FuncParam: []interface{}{0, err, "Text", 0}}
	//dl := new(DllCall)
	//dl.DllName = "user32.dll"
	//dl.FuncName = "MessageBoxA"
	//dl.FuncParam = []interface{}{0, err.Error(), "Text", 0}

	dl.Call()
}

func Info(info ...string) {
	if logger == nil{
		return
	}
	logger.SetFlags(log.LstdFlags)
	logger.SetPrefix("INFO  ")
	logger.Println(info)
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	//PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			//files = append(files, dirPth+PthSep+fi.Name())
			files = append(files, fi.Name())
		}
	}
	return files, nil
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)                                                     //忽略后缀匹配的大小写
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		//if err != nil { //忽略错误
		// return err
		//}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}