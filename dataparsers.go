package richframe

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"
)

type ParseDataFuncFactory struct {
	ParseDataFuncs map[string]ParseDataFunc
}

func GetParseDataFuncFactory() *ParseDataFuncFactory {
	factory := &ParseDataFuncFactory{make(map[string]ParseDataFunc)}
	factory.AddFunc("string", ParseString)
	factory.AddFunc("int", ParseInt)
	factory.AddFunc("float", ParseFloat)
	factory.AddFunc("ExcelDate", ParseExcelDate)
	return factory
}

func (factory *ParseDataFuncFactory) Get(dataType string) ParseDataFunc {
	if dataType == "" {
		return factory.ParseDataFuncs["string"]
	}
	return factory.ParseDataFuncs[dataType]
}

func (factory *ParseDataFuncFactory) AddFunc(dataType string, pdFunc ParseDataFunc) {
	factory.ParseDataFuncs[dataType] = pdFunc
}

var ParseDataFuncs *ParseDataFuncFactory

func init() {
	ParseDataFuncs = GetParseDataFuncFactory()
}

type ParseDataFunc func(valueStr string, ops interface{}) (interface{}, error)

func ParseString(valueStr string, ops interface{}) (interface{}, error) {
	return valueStr, nil
}

func ParseInt(valueStr string, ops interface{}) (interface{}, error) {
	intValue, err := strconv.ParseInt(valueStr, 0, 64)
	if err != nil {
		return nil, err
	}
	return intValue, nil
}

func ParseFloat(valueStr string, ops interface{}) (interface{}, error) {
	floatValue, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return nil, err
	}
	return floatValue, nil
}

//func ParseExcelDate(valueStr string, ops interface{}) (interface{}, error) {
//	if valueStr == "" {
//		return "", nil
//	}
//	excelDate, err := strconv.ParseFloat(valueStr, 64)
//	if err != nil {
//		return nil, err
//	}
//	excelTime, err := excelize.ExcelDateToTime(excelDate, false)
//	// excelize.time
//	if err != nil {
//		return nil, err
//	}
//	return excelTime, nil
//}

func ParseExcelDate(valueStr string, ops interface{}) (interface{}, error) {
	if valueStr == "" {
		return "", nil
	}
	excelDate, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		t, err := ParseTime(valueStr)
		if err != nil {
			return nil, err
		} else {
			return t, nil
		}
	}
	excelTime, err := excelize.ExcelDateToTime(excelDate, false)
	// excelize.time
	if err != nil {
		return nil, err
	}
	return excelTime, nil
}

func ParseTime(value string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05Z07:00",
		"Jan 2, 2006 at 3:04pm (MST)",
		"02/01/2006 15:04:05",
		"2/1/06 15:04",
		"2006/2/1 15:04:05",
		"2006-02-01 15:04:05",
	}
	for _, layout := range layouts {
		t, err := time.Parse(layout, value)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("无法解析时间字符串")
}
