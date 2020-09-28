package main

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
)

type CSVData struct {
	Header []string
	Data   []map[string]string
}

func ReadCSV(csvFileName string, trimBlank bool) (*CSVData, error) {
	data := &CSVData{}
	csvFile, err := os.Open(csvFileName)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	bFirst := true
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return data, err
		}
		if bFirst {
			bFirst = false

			if trimBlank {
				for _, cell := range row {
					cell = strings.ReplaceAll(cell, " ", "")
					cell = strings.ReplaceAll(cell, "	", "")
					cell = strings.ReplaceAll(cell, "\t", "")
					cell = strings.ReplaceAll(cell, "\n", "")
					cell = strings.ReplaceAll(cell, "\r", "")
					data.Header = append(data.Header, cell)
				}
			} else {
				data.Header = row
			}

			continue
		}
		rowLen := len(row)
		newRow := make(map[string]string)
		for idx, columnName := range data.Header {
			if rowLen > idx {
				cell := row[idx]
				if trimBlank {
					cell = strings.ReplaceAll(cell, " ", "")
					cell = strings.ReplaceAll(cell, "	", "")
					cell = strings.ReplaceAll(cell, "\t", "")
					cell = strings.ReplaceAll(cell, "\n", "")
					cell = strings.ReplaceAll(cell, "\r", "")
				}
				newRow[columnName] = cell
			} else {
				newRow[columnName] = ""
			}
		}
		data.Data = append(data.Data, newRow)
	}

	return data, nil
}
