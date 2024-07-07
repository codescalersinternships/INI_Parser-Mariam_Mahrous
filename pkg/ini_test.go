package main

import (
	"os"
	"reflect"
	"slices"
	"testing"
)

const testParserString = `[forge.example]
User = hg

[topsecret.server.example]
Port = 50022
ForwardX11 = no
`

func handleTestsErrors(t *testing.T, err error) {
	if err != nil {
		t.Errorf("An error occured %v", err)
		t.Fail()
	}
}

func TestParser_LoadFromString(t *testing.T) {

	parser := initialize()
	t.Run("Checking empty string", func(t *testing.T) {
		err := parser.LoadFromString("")
		handleTestsErrors(t, err)
		got := len(parser.section)
		want := 0
		if got != want {
			t.Errorf("got %d want %d given", got, want)
		}
	})
	t.Run("valid string", func(t *testing.T) {
		err := parser.LoadFromString(testParserString)
		handleTestsErrors(t, err)
		got := parser.section["forge.example"]["User"]
		want := "hg"
		if got != want {
			t.Errorf("got %s want %s given", got, want)
		}
		got = parser.section["topsecret.server.example"]["Port"]
		want = "50022"
		if got != want {
			t.Errorf("got %s want %s given", got, want)
		}
		got = parser.section["topsecret.server.example"]["ForwardX11"]
		want = "no"
		if got != want {
			t.Errorf("got %s want %s given", got, want)
		}
	})
}

func TestParser_LoadFromFile(t *testing.T) {

	parser := initialize()
	t.Run("valid file", func(t *testing.T) {
		err := parser.LoadFromFile("./testdata/input_file.ini")
		handleTestsErrors(t, err)
		got := parser.section["forge.example"]["User"]
		want := "hg"
		if got != want {
			t.Errorf("got %s want %s given", got, want)
		}
		got = parser.section["topsecret.server.example"]["Port"]
		want = "50022"
		if got != want {
			t.Errorf("got %s want %s given", got, want)
		}
		got = parser.section["topsecret.server.example"]["ForwardX11"]
		want = "no"
		if got != want {
			t.Errorf("got %s want %s given", got, want)
		}
	})
	t.Run("invalid file", func(t *testing.T) {
		err := parser.LoadFromFile("goldefile.ini")
		if err == nil {
			t.Fail()
		}
	})
}

func TestParser_GetSectionNames(t *testing.T) {

	parser := initialize()
	t.Run("Get section name from empty parser", func(t *testing.T) {
		err := parser.LoadFromString("")
		handleTestsErrors(t, err)
		got := parser.GetSectionNames()
		want := []string{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("Get section name from populated parser", func(t *testing.T) {
		err := parser.LoadFromFile("./testdata/input_file.ini")
		handleTestsErrors(t, err)
		got := parser.GetSectionNames()
		want := []string{"forge.example", "topsecret.server.example"}
		if !slices.Contains(got, "forge.example") || !slices.Contains(got, "topsecret.server.example") {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestParser_GetSections(t *testing.T) {

	parser := initialize()
	t.Run("Get section from empty parser", func(t *testing.T) {
		err := parser.LoadFromString("")
		handleTestsErrors(t, err)
		got := parser.GetSections()
		var want = make(map[string]map[string]string)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("Get section from populated parser", func(t *testing.T) {
		err := parser.LoadFromFile("./testdata/input_file.ini")
		handleTestsErrors(t, err)
		got := parser.GetSections()
		want := make(map[string]map[string]string)
		want["forge.example"] = make(map[string]string)
		want["forge.example"]["User"] = "hg"
		want["topsecret.server.example"] = make(map[string]string)
		want["topsecret.server.example"]["Port"] = "50022"
		want["topsecret.server.example"]["ForwardX11"] = "no"

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

}

func TestParser_Get(t *testing.T) {

	parser := initialize()
	err := parser.LoadFromFile("./testdata/input_file.ini")
	handleTestsErrors(t, err)
	setTests := []struct {
		test    string
		section string
		key     string
		want    string
	}{
		{"Gettting a value that exists", "forge.example", "User", "hg"},
		{"Gettting a value from a key that doesn't exitsts", "forge.example", "user", ""},
		{"Gettting a value from a section that doesn't exitsts", "default", "user", ""},
	}
	for _, tt := range setTests {
		got, _ := parser.Get(tt.section, tt.key)
		if got != tt.want {
			t.Errorf("test %s got %s want %s", tt.test, got, tt.want)
		}
	}
}

func TestParser_Set(t *testing.T) {

	parser := initialize()
	err := parser.LoadFromFile("./testdata/input_file.ini")
	handleTestsErrors(t, err)
	setTests := []struct {
		test    string
		section string
		key     string
		value   string
		want    string
	}{
		{"Setting a value", "forge.example", "User", "mariam", "mariam"},
		{"Adding a value in an existing section", "forge.example", "IP_address", "80:60:244:32", "80:60:244:32"},
		{"Adding a value in a new section", "Default", "IP_address", "80:60:244:32", "80:60:244:32"},
	}
	for _, tt := range setTests {
		err := parser.Set(tt.section, tt.key, tt.value)
		handleTestsErrors(t, err)
		got := parser.section[tt.section][tt.key]
		if got != tt.want {
			t.Errorf("test %s got %s want %s", tt.test, got, tt.want)
		}
	}
}

func TestParser_String(t *testing.T) {

	parser := initialize()
	err := parser.LoadFromFile("./testdata/input_file.ini")
	handleTestsErrors(t, err)
	got := parser.String()
	newparser := initialize()
	err = newparser.LoadFromString(got)
	handleTestsErrors(t, err)
	want := testParserString
	if newparser.section["forge.example"]["User"] != "hg" && newparser.section["topsecret.server.example"]["ForwardX11"] != "no" && newparser.section["topsecret.server.example"]["Port"] != "50022" {
		t.Errorf("want:\n%s\nGot:\n%s", want, got)
	}
}

func TestParser_SaveToFile(t *testing.T) {
	parser := initialize()
	err := parser.LoadFromFile("./testdata/input_file.ini")
	handleTestsErrors(t, err)
	err = parser.SaveToFile("test_output.ini")
	handleTestsErrors(t, err)
	content, e := os.ReadFile("test_output.ini")
	handleTestsErrors(t, e)
	newparser := initialize()
	err = newparser.LoadFromString(string(content))
	handleTestsErrors(t, err)
	want := testParserString
	if newparser.section["forge.example"]["User"] != "hg" && newparser.section["topsecret.server.example"]["ForwardX11"] != "no" && newparser.section["topsecret.server.example"]["Port"] != "50022" {
		t.Errorf("want:\n%s\nGot:\n%s", want, content)
	}
	os.Remove("test_output.ini")
}
