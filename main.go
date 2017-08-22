package main

import (
	"fmt"
	"net/http"
)

//carat扩大的倍数
const CaratMuti = 1000

func main() {

	// 初始化dbmgr
	dbmgr := DBManagerIns()
	dbmgr.Init()

	http.HandleFunc("/initdb", InitDB)
	http.HandleFunc("/get", GetData)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
