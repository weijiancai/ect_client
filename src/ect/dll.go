package ect

import (
	"fmt"
	"reflect"
	"unsafe"
	"syscall"
)

type DllCall struct {
	DllName   string        `json:"DllName"`
	FuncName  string        `json:"FuncName"`
	FuncParam []interface{} `json:"FuncParam"`
}

func (dl *DllCall) Call() {
	var param = make([]uintptr, 0)
	for _, v := range dl.FuncParam {
		fmt.Println(v)
		fmt.Println(reflect.TypeOf(v))
		fmt.Println(reflect.TypeOf(v).Kind())
		if reflect.TypeOf(v).Kind() == reflect.Float64 {
			param = append(param, uintptr(v.(float64)))
		} else if reflect.TypeOf(v).Kind() == reflect.String {
			param = append(param, uintptr(unsafe.Pointer(syscall.StringBytePtr(v.(string)))))
		}
	}

	dll := syscall.NewLazyDLL(dl.DllName)
	proc := dll.NewProc(dl.FuncName)
	r1, r2, err := proc.Call(param...)
	fmt.Println(r1)
	fmt.Println(r2)
	fmt.Println(err)
}