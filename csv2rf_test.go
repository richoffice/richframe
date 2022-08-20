package richframe

import (
	"fmt"
	"testing"
)

func TestLoadCSV(t *testing.T) {
	rf, err := LoadCSV("./testfiles/msgs.csv", []string{"from", "to", "date"})

	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

	fmt.Println(rf)

	err = SaveCSV("./testfiles/msgs-out.csv", rf, nil, true)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

	err = SaveCSV("./testfiles/msgs-out.csv", rf, []string{"from", "to"}, true)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

}

func TestLoadCSVWithTitles(t *testing.T) {
	rf, err := LoadCSV("./testfiles/msgs-withtitles.csv", nil)

	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

	fmt.Println(rf)
}
