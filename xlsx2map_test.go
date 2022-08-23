package richframe

import (
	"fmt"
	"testing"
)

// func TestUnmarshal(t *testing.T) {

// 	def := &XlsxFileDef{}
// 	file, err := os.Open("sample_def.json")

// 	if err != nil {
// 		t.Errorf("expected no err, but got %v", err)
// 	}

// 	defer file.Close()

// 	loadErr := LoadXlsxFileDef(file, def)
// 	if loadErr != nil {
// 		t.Errorf("expected no err, but got %v", loadErr)
// 	}

// 	xlsxFile := "sample_file.xlsx"

// 	xlsxMaps := make(map[string][]map[string]interface{})

// 	err = Unmarshal(xlsxFile, xlsxMaps, def, nil)
// 	if err != nil {
// 		t.Errorf("Unmarshal() error = %v, wantErr %v", err, nil)
// 	}

// 	if xlsxMaps["visitors"][1]["name"] != "Tom Hanks" {
// 		t.Errorf("Expected 'Tom Hanks', but got %v", xlsxMaps["visitors"][1]["name"])

// 	}

// 	fmt.Println(xlsxMaps)
// 	// if !reflect.DeepEqual(got, tt.want) {
// 	// 	t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
// 	// }

// 	outErr := Marshal("test_out.xlsx", xlsxMaps, def)
// 	if outErr != nil {
// 		t.Errorf("Expected no output err, but got %v", outErr)
// 	}

// }

// func TestLoadFromFile(t *testing.T) {
// 	data, err := LoadFromFile("sample_file.xlsx", "sample_def.json", nil)
// 	if err != nil {
// 		t.Errorf("expected no error, but got %v", err)
// 	}
// 	// fmt.Println(data)

// 	err = ExportToFile(data, "test_out.xlsx", "sample_def.json", nil)
// 	if err != nil {
// 		t.Errorf("expected no error, but got %v", err)
// 	}
// }

func TestLoadRichFrames(t *testing.T) {
	data, err := LoadRichFrames("testfiles/sample_file.xlsx", "testfiles/sample_def.json", nil)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	fmt.Println(data)

	err = ExportRichFrames(data, "testfiles/test_out.xlsx", "testfiles/sample_def.json", nil)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}

func TestExportRichFramesByTemp(t *testing.T) {
	data, err := LoadRichFrames("testfiles/sample_file.xlsx", "testfiles/sample_def.json", nil)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	fmt.Println(data)

	err = ExportRichFramesByTemp(data, "testfiles/test_tmp_out.xlsx", "testfiles/test_tmp.xlsx", "testfiles/sample_def.json", nil)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
}
