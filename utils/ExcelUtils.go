package utils

import (
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
)

func ExportExcel(heads []string, content [][]string) (*xlsx.File, error) {
	excelFile := xlsx.NewFile()
	sheet, e := excelFile.AddSheet("Sheet1")
	if e != nil {
		return nil, e
	}
	if len(heads) == 0 {
		return nil, errors.New("empty head")
	}

	if len(content) == 0 {
		return nil, errors.New("empty content")
	}

	headRow := sheet.AddRow()
	for _, head := range heads {
		cell := headRow.AddCell()
		cell.Value = head
	}

	for _, con := range content {
		row := sheet.AddRow()
		for _, val := range con {
			cell := row.AddCell()
			cell.Value = val
		}
	}
	return excelFile, nil
}
