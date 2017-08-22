package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var onceNewRM sync.Once
var dbManagerInstance *DBManager

func DBManagerIns() *DBManager {
	onceNewRM.Do(func() {
		if dbManagerInstance == nil {
			dbManagerInstance = new(DBManager)
		}
	})
	return dbManagerInstance
}

type DBManager struct {
	engine *xorm.Engine
}

func (rm *DBManager) Init() {
	var err error
	rm.engine, err = xorm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3306)/diamond?charset=utf8&loc=Local")
	//rm.engine, err = xorm.NewEngine("mysql", "root:123456@tcp(192.168.163.162:3306)/diamond?charset=utf8&loc=Local")
	if err != nil {
		panic(err)
	}
	rm.engine.ShowSQL()
}

func (rm *DBManager) Inset2DB(data []*Price) error {
	if len(data) == 0 {
		return errors.New("data nil err")
	}
	var err error = nil

	truncateSql := `truncate table price`
	_, err = rm.engine.Exec(truncateSql)
	if err != nil {
		return err
	}

	_, err = rm.engine.Insert(data)
	if err != nil {
		return err
	}
	return nil
}

func (rm *DBManager) GetPrice(carat float64, color uint8, clarity string) (interface{}, error) {
	prices := make([]*Price, 0)
	caratSql := ""
	if carat != 0.0 {
		caratSql = fmt.Sprintf(`and carat_end>=%v and carat_begin <=%v`, int(carat*CaratMuti), int(carat*CaratMuti))
	}
	colorSql := ""
	if color != 0 {
		colorSql = fmt.Sprintf(`and color_end>=%v and color_begin <=%v`, color, color)
	}
	claritySql := ""
	if len(clarity) != 0 {
		claritySql = fmt.Sprintf(`and clarity="%v"`, clarity)
	}

	querySql := `select carat_begin,carat_end,color_begin,color_end,clarity,price from price where id > 0 %v %v %v`

	price := new(Price)
	priceRows, err := rm.engine.SQL(fmt.Sprintf(querySql, caratSql, colorSql, claritySql)).Rows(price)
	if err != nil {
		return prices, err
	}
	defer priceRows.Close()
	for priceRows.Next() {
		price := new(Price)
		priceRows.Scan(price)
		priceparma, _ := strconv.ParseFloat(price.Price, 64)
		price.TotalPrice = carat * 100 * priceparma
		fmt.Println(price)
		prices = append(prices, price)
	}
	return prices, nil
}
