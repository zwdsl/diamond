package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tealeg/xlsx"
)

type Price struct {
	CaratBegin int
	CaratEnd   int
	ColorBegin uint8
	ColorEnd   uint8
	Clarity    string
	Price      string
	TotalPrice float64 `xorm:"-"`
}

type CTLocaltion struct {
	CTRow  int
	ColEnd int
	RowEnd int
}

type ResultPrice struct {
	Carat   string
	Color   uint8
	Clarity string
	Price   uint32
}

func GetCTLocaltion(xlFile *xlsx.File) ([]*CTLocaltion, error) {
	ctlocal := make([]*CTLocaltion, 0)

	ctid := 0
	if len(xlFile.Sheets) == 0 {
		return nil, errors.New("Sheets == 0")
	}
	sheet := xlFile.Sheets[0]
	//fmt.Println(sheet.Name)
	for rowid, row := range sheet.Rows {
		if row.Cells[0].String() == "ct" {
			//记录ct位置
			ctl := &CTLocaltion{
				CTRow: rowid,
			}
			ctlocal = append(ctlocal, ctl)
			if ctid != 0 {
				ctlocal[ctid-1].RowEnd = rowid - 1
			}
		} else if row.Cells[0].String() == "data" {
			//计算列的项数 ColEnd
			colcnt := 0
			for _, cell := range row.Cells {
				if len(cell.String()) == 0 {
					break
				}
				colcnt++
			}
			ctlocal[ctid].ColEnd = colcnt
			//到这里再+1
			ctid++
		}
	}

	return ctlocal, nil
}

func GetDiamondData(xlFile *xlsx.File, ctlocal []*CTLocaltion) ([]*Price, error) {
	data := make([]*Price, 0)
	if len(xlFile.Sheets) == 0 {
		return nil, errors.New("Sheets == 0")
	}
	sheet := xlFile.Sheets[0]
	//	fmt.Println(sheet.Name)

	for _, v := range ctlocal {
		for i := v.CTRow; i <= v.RowEnd; i++ {
			if sheet.Cell(i, 0).String() == "ct" ||
				sheet.Cell(i, 0).String() == "data" {
				continue
			}
			for j := 0; j < v.ColEnd; j++ {
				if j == 0 {
					continue
				}
				bbyte := []byte(sheet.Cell(i, 0).String())
				var (
					caratbegin64 float64
					caratend64   float64
					colorbegin   uint8
					colorend     uint8
					err          error
				)
				caratbegin64, err = strconv.ParseFloat(sheet.Cell(v.CTRow, 1).String(), 64)
				if err != nil {
					return nil, err
				}
				caratend64, err = strconv.ParseFloat(sheet.Cell(v.CTRow, 2).String(), 64)
				if err != nil {
					return nil, err
				}
				if len(bbyte) != 3 {
					colorbegin = bbyte[0]
					colorend = bbyte[0]
				} else {
					colorbegin = bbyte[0]
					colorend = bbyte[2]
				}
				arg := &Price{
					CaratBegin: int(caratbegin64 * CaratMuti),
					CaratEnd:   int(caratend64 * CaratMuti),
					ColorBegin: colorbegin,
					ColorEnd:   colorend,
					Clarity:    sheet.Cell(v.CTRow+1, j).String(),
					Price:      sheet.Cell(i, j).String(),
				}
				data = append(data, arg)
			}

		}
	}
	return data, nil
}

func CalcData2DB() error {
	excelFileName := "xls.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	ctlocal, _ := GetCTLocaltion(xlFile)
	for _, l := range ctlocal {
		fmt.Println(l)
	}
	data, err := GetDiamondData(xlFile, ctlocal)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	for _, l := range data {
		fmt.Println(l)
	}
	return DBManagerIns().Inset2DB(data)
}
