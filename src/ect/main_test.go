package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestDllCall(t *testing.T) {
	fmt.Println("Test1")
	dll := DllCall{DllName: "user32.dll", FuncName: "MessageBoxA", FuncParam: []interface{}{0, "HelloWorld", "Text", 0}}
	jsonStr, _ := json.Marshal(dll)
	println(string(jsonStr))
	res, err := http.Post("http://localhost:8001/dllcall", "application/x-www-form-urlencoded", strings.NewReader(string(jsonStr)))
	if err != nil {
		t.Error(err.Error())
	}

	t.Log(res)
}
