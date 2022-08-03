package richframe

import (
	"bytes"
	"fmt"
)

type RichMap map[string]interface{}

type RichFrame struct {
	RichMaps []RichMap
}

func (rf *RichFrame) Append(rm RichMap) {
	rf.RichMaps = append(rf.RichMaps, rm)
}

// type RichFrames map[string]*RichFrame

// func (rfs RichFrames) Get(key string) interface{} {
// 	return rfs[key]
// }

type ApplyFunc func(RichMap)

type MutateFunc func(RichMap) interface{}

type AggregateFunc func(interface{}, RichMap) interface{}

type FilterFunc func(RichMap) bool

func (rf *RichFrame) String() string {
	var buffer bytes.Buffer
	for _, row := range rf.RichMaps {
		buffer.WriteString(fmt.Sprintf("%v \n", row))
		break
	}

	buffer.WriteString(fmt.Sprintf("len(%v) \n", len(rf.RichMaps)))

	return buffer.String()
}

func (rf *RichFrame) Rows() []RichMap {
	return rf.RichMaps
}

func (rf *RichFrame) Apply(f ApplyFunc) *RichFrame {
	for _, row := range rf.RichMaps {
		f(row)
	}
	return rf
}

func (rf *RichFrame) Mutate(title string, f MutateFunc) *RichFrame {
	for _, row := range rf.RichMaps {
		row[title] = f(row)
	}
	return rf
}

func (rf *RichFrame) Aggregate(groupBy []string, cols []string, funcs []AggregateFunc) *RichFrame {
	out := &RichFrame{}

	for _, row := range rf.Rows() {
		outMap := GetGroup(out, row, groupBy)
		// fmt.Println(outMap)
		if outMap == nil {
			// fmt.Println(outMap)
			outMap = RichMap{}
			for _, by := range groupBy {
				outMap[by] = row[by]
			}

			for i := 0; i < len(cols); i++ {
				outMap[cols[i]] = funcs[i](nil, row)
			}

			out.Append(outMap)
		} else {
			for i := 0; i < len(cols); i++ {
				outMap[cols[i]] = funcs[i](outMap[cols[i]], row)
			}
		}

	}

	return out
}

func GetGroup(out *RichFrame, origin RichMap, groupBy []string) RichMap {
	for _, r := range out.Rows() {
		for byIndex, by := range groupBy {
			if r[by] != origin[by] {
				break
			}
			if byIndex == len(groupBy)-1 {
				return r
			}
		}
	}

	return nil
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
