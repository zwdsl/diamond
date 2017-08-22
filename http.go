package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// http返回通用接口
type RespResult struct {
	Code    uint32      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PriceData struct {
	Price   float64 `json:"price"`
	carat   float64 `json:"carat"`
	color   uint8   `json:"color"`
	clarity string  `json:"clarity"`
}

func makeRespData(code uint32, msg string, data interface{}) []byte {
	resp := &RespResult{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	respdata, _ := json.Marshal(resp)
	return respdata
}

// 初始化DB
func InitDB(w http.ResponseWriter, req *http.Request) {
	err := CalcData2DB()
	if err == nil {
		w.Write(makeRespData(0, "", nil))
		return
	}
	w.Write(makeRespData(100, err.Error(), nil))
}

// 取数据
func GetData(w http.ResponseWriter, r *http.Request) {
	data := make([]byte, 0)

	defer func() {
		if len(data) == 0 {
			data = makeRespData(101, "处理异常", nil)
		}
		w.Write(data)
	}()
	if r.Method != "GET" {
		return
	}
	var (
		err     error   = nil
		carat   float64 = 0.0
		color   uint8   = 0
		clarity string  = ""
	)

	query := r.URL.Query()

	//拿carat 必须
	qstr, _ := query["carat"]
	if len(qstr) == 0 {
		data = makeRespData(102, "carat参数错误", nil)
		fmt.Println("parma err:", qstr)
		return
	}
	qqstr := qstr[0]
	carat, err = strconv.ParseFloat(qqstr, 64)
	if err != nil {
		data = makeRespData(102, "参数错误", nil)
		fmt.Println("parma err:", qstr)
		return
	}

	//拿color
	qcolor, _ := query["color"]
	if len(qcolor) != 0 {
		qqcolor := qcolor[0]
		bcolor := []byte(qqcolor)
		if len(bcolor) != 0 {
			color = uint8(bcolor[0])
		}
	}

	//拿clarity
	qclarity, _ := query["clarity"]
	if len(qclarity) != 0 {
		clarity = qclarity[0]
	}

	arg, argerr := DBManagerIns().GetPrice(carat, color, clarity)
	if argerr != nil {
		fmt.Println("GetPrice err:", argerr.Error())
		return
	}

	data = makeRespData(0, "", arg)

	fmt.Println("GetData carat:", carat, " color:", color, " clarity:", clarity)
}
