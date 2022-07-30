package richframe

import (
	"fmt"
	"reflect"
	"testing"
)

func getData1() RichFrame {
	return RichFrame{
		RichMap{
			"key1": "abc1",
		},
		RichMap{
			"key1": "abc2",
		},
		RichMap{
			"key1": "abc3",
		},
	}
}

func TestRichFrame_ToString(t *testing.T) {

	expected := `key1:abc1,
key1:abc2,
key1:abc3,
`
	rf := getData1()

	if expected != rf.String() {
		t.Errorf("Expected %v but got %v", expected, rf.String())
	}

	fmt.Println(rf)

}

func TestRichFrame_Apply(t *testing.T) {
	rf := getData1()

	expected := RichFrame{
		RichMap{
			"key1": "abc1",
			"key2": "abc1 abc1",
		},
		RichMap{
			"key1": "abc2",
			"key2": "abc2 abc2",
		},
		RichMap{
			"key1": "abc3",
			"key2": "abc3 abc3",
		},
	}
	rf.Apply(func(rm RichMap) {
		rm["key2"] = fmt.Sprintf("%v %v", rm["key1"], rm["key1"])
	})

	if !reflect.DeepEqual(rf, expected) {
		t.Errorf("expected %v, but got %v", expected, rf)
	}
}

func TestRichFrame_Filter(t *testing.T) {
	rf := getData1()
	expected := RichFrame{
		RichMap{
			"key1": "abc1",
		},
		RichMap{
			"key1": "abc2",
		},
	}

	rf.Filter(func(rm RichMap) bool {
		return rm["key1"] != "abc3"
	})

	if !reflect.DeepEqual(rf, expected) {
		t.Errorf("expected %v, but got %v", expected, rf)
	}

}