package richframe

import (
	"encoding/csv"
	"io"
	"os"
)

func LoadCSV(csvpath string, keys []string) (RichFrame, error) {

	csvFile, err := os.Open(csvpath)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()
	csvReader := csv.NewReader(csvFile)

	hasTitles := true
	if keys != nil {
		hasTitles = false
	}

	index := 0

	rf := RichFrame{}

	for {
		row, err := csvReader.Read()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}

		if index == 0 && hasTitles {
			keys = row
		} else {
			columnCount := len(keys)
			if len(keys) > len(row) {
				columnCount = len(row)
			}

			rm := RichMap{}
			for i := 0; i < columnCount; i++ {
				rm[keys[i]] = row[i]
			}
			rf = append(rf, rm)
		}

		index++

	}

	return rf, nil

}

func SaveCSV(csvpath string, rf RichFrame, keys []string, isAppend bool) error {

	hasData := true

	if keys == nil {
		keys = []string{}

		if len(rf) < 1 {
			hasData = false
		} else {
			row0 := rf[0]
			for k, _ := range row0 {
				keys = append(keys, k)
			}

		}
	}

	var csvFile *os.File
	var err error

	if isAppend {
		csvFile, err = os.OpenFile(csvpath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
	} else {
		csvFile, err = os.Create(csvpath)
		if err != nil {
			return err
		}
	}

	w := csv.NewWriter(csvFile)

	if hasData {
		w.Write(keys)

		for _, rm := range rf {
			row := make([]string, 0)
			for _, key := range keys {
				row = append(row, rm[key].(string))
			}
			w.Write(row)

		}
		w.Flush()
	}

	return nil

}
