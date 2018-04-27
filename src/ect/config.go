package main

import (
	"os/exec"
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"log"
	"io"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	DbType string `json:"db_type"`
	DbHost string `json:"db_host"`
	DbPort int `json:"db_port"`
	DbUserName string `json:"db_username"`
	DbPassword string `json:"db_password"`
	DbDatabase string `json:"db_database"`
	ServerUrl string `json:"server_url"`
	ServerKey string `json:"server_key"`
	ErpType string `json:"erp_type"`
	ErpStation string `json:"erp_station"`
}

func init() {
	curPath := getCurrentPath()
	log.Println(curPath)
	configPath := getCurrentPath() + "\\config.json"
	var f *os.File
	var err error

	var config = Config{}
	if checkFileIsExist(configPath) {
		f, err = os.OpenFile(configPath, os.O_RDONLY, 0666)
		check(err)
		bytes,err := ioutil.ReadAll(f)
		check(err)
		log.Println(string(bytes))

		err = json.Unmarshal(bytes,&config)
		check(err)

		fmt.Printf("CONFIG:%v\n",config)

	} else {
		// 创建文件
		f, err = os.Create(configPath)
		check(err)
		// 写入内容
		str, err :=json.MarshalIndent(config,"","\t");
		if err != nil {
			check(err)
		}
		io.WriteString(f, string(str))
	}
}

func getCurrentPath() string {
	file, _ := exec.LookPath(os.Args[0])
	fmt.Println("file:", file)
	path, _ := filepath.Abs(file)
	fmt.Println("path:", path)
	splitstring := strings.Split(path, "\\")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	ret := strings.Replace(splitstring[0], "\\", "/", size-1)
	return ret
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func main() {
	log.Println("start")
}