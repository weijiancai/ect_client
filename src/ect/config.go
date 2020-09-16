package ect

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
	ClientKey string `json:"client_key"`
	ErpType string `json:"erp_type"`
	ErpStation string `json:"erp_station"`
	LocalPort int `json:"local_port"`
}

var ConfigInfo Config

func init() {
	curPath := GetCurrentPath()
	log.Println(curPath)
	configPath := GetCurrentPath() + "\\config.json"
	var f *os.File
	var err error

	ConfigInfo = Config{}
	if checkFileIsExist(configPath) {
		f, err = os.OpenFile(configPath, os.O_RDONLY, 0666)
		Error(err)
		bytes,err := ioutil.ReadAll(f)
		Error(err)
		log.Println(string(bytes))

		err = json.Unmarshal(bytes,&ConfigInfo)
		Error(err)

		fmt.Printf("CONFIG:%v\n", ConfigInfo)

	} else {
		// 创建文件
		f, err = os.Create(configPath)
		Error(err)
		// 写入内容
		str, err :=json.MarshalIndent(ConfigInfo,"","\t");
		if err != nil {
			Error(err)
		}
		io.WriteString(f, string(str))
	}
}

func GetCurrentPath() string {
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

func CheckConfig() {
	if ConfigInfo.DbType == "" {
		configError("数据库类型")
	}
}

func configError(param string) {
	ErrorStr("参数错误：" + param)
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
