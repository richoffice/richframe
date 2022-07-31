package richframe

import (
	"bytes"
	"fmt"
)

type RichMap map[string]interface{}

type RichFrame struct {
	RichMaps []RichMap
}

type RichFrames map[string]*RichFrame

func (rfs RichFrames) Get(key string) interface{} {
	return rfs[key]
}

type ApplyFunc func(RichMap)

type AddFunc func(RichMap) interface{}

type FilterFunc func(RichMap) bool

func (rf *RichFrame) String() string {
	var buffer bytes.Buffer
	for _, row := range rf.RichMaps {
		buffer.WriteString(fmt.Sprintf("%v \n", row))
	}
	return buffer.String()
}

func (rf *RichFrame) Apply(f ApplyFunc) *RichFrame {
	for _, row := range rf.RichMaps {
		f(row)
	}
	return rf
}

func (rf *RichFrame) Add(title string, f AddFunc) *RichFrame {
	for _, row := range rf.RichMaps {
		row[title] = f(row)
	}
	return rf
}

func (rf *RichFrame) Filter(f FilterFunc) *RichFrame {

	tmpRM := []RichMap{}

	for _, row := range rf.RichMaps {
		if f(row) {
			tmpRM = append(tmpRM, row)
		}
	}

	rf.RichMaps = tmpRM
	return rf
}
