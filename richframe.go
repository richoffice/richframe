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

func (rf *RichFrame) Rename(old string, new string) *RichFrame {
	for _, row := range rf.RichMaps {
		row[new] = row[old]
		delete(row, old)
	}
	return rf
}

func (rf *RichFrame) Distinct(col string) []interface{} {
	values := make(map[interface{}]interface{}, 0)
	for _, row := range rf.Rows() {
		values[row[col]] = 0
	}
	keys := make([]interface{}, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}

	return keys
}

func (rf *RichFrame) Col(col string) []interface{} {
	values := make([]interface{}, 0, len(rf.Rows()))
	for _, row := range rf.Rows() {
		values = append(values, row[col])
	}
	return values
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
		outMap := GetGroup(out.Rows(), row, groupBy)
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

func GetGroup(out []RichMap, origin RichMap, groupBy []string) RichMap {
	for _, r := range out {
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

func (rf *RichFrame) FilterNew(f FilterFunc) *RichFrame {

	tmpRM := []RichMap{}

	for _, row := range rf.RichMaps {
		if f(row) {
			tmpRM = append(tmpRM, row)
		}
	}

	return &RichFrame{tmpRM}
}

func (rf *RichFrame) Join(rights *RichFrame, leftBys []string, rightBys []string, defaults map[string]interface{}) {
	for i := 0; i < len(rf.Rows()); i++ {
		left := rf.Rows()[i]
		matched := GetMatchRichMap(left, rights.Rows(), leftBys, rightBys)
		if matched == nil {
			for k, v := range defaults {
				left[k] = v
			}
		} else {
			for k, v := range matched {
				if _, ok := left[k]; !ok {
					left[k] = v
				}
			}
		}
	}
}

func GetMatchRichMap(left RichMap, rights []RichMap, leftBys []string, rightBys []string) RichMap {
	for _, right := range rights {
		if MatchGroup(left, right, leftBys, rightBys) {
			return right
		}
	}
	return nil
}

func MatchGroup(left RichMap, right RichMap, leftBys []string, rightBys []string) bool {
	for byIndex, leftBy := range leftBys {
		if left[leftBy] != right[rightBys[byIndex]] {
			return false
		}
	}

	return true
}

// func contains(s []interface{}, str interface{}) bool {
// 	for _, v := range s {
// 		if v == str {
// 			return true
// 		}
// 	}

// 	return false
// }
