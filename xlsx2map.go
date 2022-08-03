package richframe

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Options struct {
}

func Marshal(outFilePath string, input interface{}, def *XlsxFileDef) error {

	f := excelize.NewFile()

	deleteDefaultSheet := true

	if data, ok := input.(map[string]*RichFrame); ok {
		// fmt.Println("map[string]map[string]interface{}")
		for _, sheetDef := range def.SheetDefs {
			if "sheet1" == strings.ToLower(sheetDef.GetTitle()) {
				deleteDefaultSheet = false
			}
			f.NewSheet(sheetDef.GetTitle())

			for colIndex, fieldDef := range sheetDef.FieldDefs {
				columnName, columnErr := excelize.ColumnNumberToName(colIndex + 1)
				if columnErr != nil {
					return columnErr
				}
				f.SetCellValue(sheetDef.GetTitle(), columnName+"1", fieldDef.GetTitle())

			}

			sheetData := data[sheetDef.Key]
			if sheetData == nil {
				return fmt.Errorf("no data for key: %v", sheetDef.Key)
			}

			for i := 0; i < len(sheetData.RichMaps); i++ {
				rowData := sheetData.RichMaps[i]
				for colIndex, fieldDef := range sheetDef.FieldDefs {
					columnName, columnErr := excelize.ColumnNumberToName(colIndex + 1)
					if columnErr != nil {
						return columnErr
					}
					f.SetCellValue(sheetDef.GetTitle(), columnName+strconv.Itoa(i+2), rowData[fieldDef.Key])
				}

			}

			// fmt.Println(data, index)

		}

		if deleteDefaultSheet {
			f.DeleteSheet("sheet1")
		}

		f.SetActiveSheet(0)
		if err := f.SaveAs(outFilePath); err != nil {
			return err
		}

	} else {
		return errors.New("not supported data type")
	}

	return nil

}

func Unmarshal(xslxFile string, result interface{}, def *XlsxFileDef, opts *Options) error {

	f, err := excelize.OpenFile(xslxFile)
	if err != nil {
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			return
		}
	}()

	for _, sheetName := range f.GetSheetList() {
		sheetDef := def.GetSheetDef(sheetName)
		if sheetDef != nil {
			sheetMap, err := parseSheet(f, sheetName, sheetDef)
			if err != nil {
				return err
			}

			switch v := result.(type) {
			case map[string]*RichFrame:
				v[sheetDef.Key] = sheetMap
			default:
				return fmt.Errorf("not a sheepmap: %v", sheetMap)
			}

			// xlsxMaps[sheetDef.Key] = sheetMap
			// fmt.Println(sheetMap)

		}

	}

	return nil
}

func parseSheet(f *excelize.File, sheet string, sheetDef *SheetDef) (*RichFrame, error) {
	rows, err := f.GetRows(sheet, excelize.Options{RawCellValue: true})
	if err != nil {
		return nil, err
	}

	var columns *Columns = nil
	results := &RichFrame{
		RichMaps: []RichMap{},
	}
	for i, row := range rows {
		if i == 0 {
			columns = PrepareColumns(row, sheetDef)
		} else {
			data := PrepareRow(row, columns)
			results.RichMaps = append(results.RichMaps, data)
		}

	}
	return results, nil
}

func PrepareColumns(titles []string, sheetDef *SheetDef) *Columns {
	columns := &Columns{FieldDefs: make(map[int]*FieldDef)}
	for index, title := range titles {
		if fieldDef := sheetDef.GetFieldDef(title); fieldDef != nil {
			columns.AddColumns(index, fieldDef)
		}
	}
	return columns
}

func PrepareRow(values []string, columns *Columns) map[string]interface{} {
	data := make(map[string]interface{})
	for index, value := range values {
		fieldDef := columns.GetFieldDef(index)
		if fieldDef != nil && fieldDef.Key != "" {
			v, err := fieldDef.ParseValue(value)
			if err != nil {
				data[fieldDef.Key] = err
			} else {
				data[fieldDef.Key] = v
			}

		}

	}
	return data
}

func LoadFromFile(excelFile, excelDefFile string, opts *Options) (map[string][]map[string]interface{}, error) {
	def := &XlsxFileDef{}
	file, err := os.Open(excelDefFile)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	loadErr := LoadXlsxFileDef(file, def)
	if loadErr != nil {
		return nil, loadErr
	}

	xlsxMaps := make(map[string][]map[string]interface{})

	err = Unmarshal(excelFile, xlsxMaps, def, nil)
	if err != nil {
		return nil, err
	}

	return xlsxMaps, nil

}

func ExportToFile(data map[string][]map[string]interface{}, outExcelFile, excelDefFile string, opts *Options) error {
	def := &XlsxFileDef{}
	file, err := os.Open(excelDefFile)

	if err != nil {
		return err
	}

	defer file.Close()

	loadErr := LoadXlsxFileDef(file, def)
	if loadErr != nil {
		return loadErr
	}

	return Marshal(outExcelFile, data, def)

}

func LoadRichFrames(excelFile, excelDefFile string, opts *Options) (map[string]*RichFrame, error) {
	def := &XlsxFileDef{}
	file, err := os.Open(excelDefFile)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	loadErr := LoadXlsxFileDef(file, def)
	if loadErr != nil {
		return nil, loadErr
	}

	frames := map[string]*RichFrame{}

	err = Unmarshal(excelFile, frames, def, nil)
	if err != nil {
		return nil, err
	}

	return frames, nil

}

func ExportRichFrames(data map[string]*RichFrame, outExcelFile, excelDefFile string, opts *Options) error {
	def := &XlsxFileDef{}
	file, err := os.Open(excelDefFile)

	if err != nil {
		return err
	}

	defer file.Close()

	loadErr := LoadXlsxFileDef(file, def)
	if loadErr != nil {
		return loadErr
	}

	return Marshal(outExcelFile, data, def)

}
