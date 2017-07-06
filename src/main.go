package main

import "net/http"

func DllCall(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello"))
}

func  main()  {
	http.HandleFunc("/dllcall", DllCall)
	http.ListenAndServe(":8001", nil)
}