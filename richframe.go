package richframe

import (
	"bytes"
	"fmt"
)

type RichMap map[string]interface{}

func (rm RichMap) String() string {
	var buffer bytes.Buffer
	for k, v := range rm {
		buffer.WriteString(fmt.Sprintf("%v:%v,", k, v))
	}
	return buffer.String()
}

type RichFrame []RichMap

type RichFrames map[string]RichFrame

type ApplyFunc func(RichMap)

type AddFunc func(RichMap) interface{}

type FilterFunc func(RichMap) bool

func (rf *RichFrame) String() string {
	var buffer bytes.Buffer
	for _, row := range *rf {
		buffer.WriteString(row.String() + "\n")
	}
	return buffer.String()
}

func (rf *RichFrame) Apply(f ApplyFunc) RichFrame {
	for _, row := range *rf {
		f(row)
	}
	return *rf
}

func (rf *RichFrame) Add(title string, f AddFunc) RichFrame {
	for _, row := range *rf {
		row[title] = f(row)
	}
	return *rf
}

func (rf *RichFrame) Filter(f FilterFunc) RichFrame {

	tmpRf := *rf
	*rf = RichFrame{}

	for _, row := range tmpRf {
		if f(row) {
			*rf = append(*rf, row)
		}
	}
	return *rf
}
