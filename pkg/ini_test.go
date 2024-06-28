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
ForwardX11 = no`

func checkError(t *testing.T , err error) {
	if err !=nil {
		t.Errorf("failed")
		t.Fail()
	}
}

func TestParser_LoadFromString(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	t.Run("Checking empty string", func(t *testing.T) {
		err := parser.LoadFromString("")
		checkError(t,err)
		got := len(parser.section)
		want := 0
		if got != want {
			t.Errorf("got %d want %d given", got, want)
		}
	})
	t.Run("valid string", func(t *testing.T) {
		err:=parser.LoadFromString(testParserString)
		got := parser.section["forge.example"]["User"]
		want := "hg"
		checkError(t,err)
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
	t.Helper()
	parser := &IniParser{}
	t.Run("valid file", func(t *testing.T) {
		err := parser.LoadFromFile("golden_file.txt")
		checkError(t,err)
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
		err := parser.LoadFromFile("goldefile.txt")
		if err == nil {
			//if not null check that it's print
			t.Fail()
		}
	})
}

func TestParser_GetSectionNames(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	t.Run("Get section name from empty parser", func(t *testing.T) {
		err:=parser.LoadFromString("")
		checkError(t,err)
		got := parser.GetSectionNames()
		want := []string{}
		if reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v given", got, want)
		}
	})
	t.Run("Get section name from populated parser", func(t *testing.T) {
		err:=parser.LoadFromFile("golden_file.txt")
		checkError(t,err)
		got := parser.GetSectionNames()
		want := []string{"forge.example", "topsecret.server.example"}
		if !slices.Contains(got, "forge.example") || !slices.Contains(got, "topsecret.server.example") {
			t.Errorf("got %v want %v given", got, want)
		}
	})

}

func TestParser_GetSection(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	t.Run("Get section from empty parser", func(t *testing.T) {
		err:=parser.LoadFromString("")
		checkError(t,err)
		got := parser.GetSection()
		var want map[string]map[string]string
		if reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v given", got, want)
		}
	})
	t.Run("Get section from populated parser", func(t *testing.T) {
		err:=parser.LoadFromFile("golden_file.txt")
		checkError(t,err)
		got := parser.GetSection()
		want := make(map[string]map[string]string)
		want["forge.example"] = make(map[string]string)
		want["forge.example"]["User"] = "hg"
		want["topsecret.server.example"] = make(map[string]string)
		want["topsecret.server.example"]["Port"] = "50022"
		want["topsecret.server.example"]["ForwardX11"] = "no"

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v given", got, want)
		}
	})

}

func TestParser_GetValue(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	err:= parser.LoadFromFile("golden_file.txt")
	checkError(t,err)
	t.Run("Gettting a value that exists", func(t *testing.T) {
		got, err := parser.GetValue("forge.example", "User")
		checkError(t,err)
		want := "hg"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
	t.Run("Gettting a value from a key that doesn't exitsts", func(t *testing.T) {
		got, _ := parser.GetValue("forge.example", "user")
		want := ""
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
	t.Run("Gettting a value from a section that doesn't exitsts", func(t *testing.T) {
		got, _ := parser.GetValue("default", "user")
		want := ""
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

}

func TestParser_SetValue(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	err:=parser.LoadFromFile("golden_file.txt")
	checkError(t,err)
	t.Run("Setting a value that exists", func(t *testing.T) {
		err:=parser.SetValue("forge.example", "User", "mariam")
		checkError(t,err)
		got := parser.section["forge.example"]["User"]

		if "mariam" != got {
			t.Errorf("got %s want mariam", got)
		}
	})
	t.Run("Trying to set a value for a key that doesn't exist", func(t *testing.T) {
		want := "80:60:244:32"
		err:=parser.SetValue("forge.example", "IP_address", want)
		checkError(t,err)
		got := parser.section["forge.example"]["IP_address"]
		if got != want {
			t.Errorf("got %s want %s given", got, want)
		}
	})
	t.Run("Trying to set a value for a section & key that don't exist", func(t *testing.T) {
		want := "80:60:244:32"
		err:= parser.SetValue("Default", "IP_address", want)
		checkError(t,err)
		got := parser.section["Default"]["IP_address"]
		if got != want {
			t.Errorf("got %s want %s given", got, want)
		}
	})
}

func TestParser_ToString(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	err:= parser.LoadFromFile("golden_file.txt")
	checkError(t,err)
	got := parser.ToString()
	err = parser.LoadFromString(got)
	checkError(t,err)
	want := testParserString
	if parser.section["forge.example"]["User"] != "hg" && parser.section["topsecret.server.example"]["ForwardX11"] != "no" && parser.section["topsecret.server.example"]["Port"] != "50022" {
		t.Errorf("want:\n%s\nGot:\n%s", want, got)
	}
}

func TestParser_SaveToFile(t *testing.T) {
	parser := &IniParser{}
	err := parser.LoadFromFile("golden_file.txt")
	checkError(t,err)
	err = parser.SaveToFile("test_output")
	checkError(t,err)
	content, e := os.ReadFile("test_output.txt")
	checkError(t,e)
	err=parser.LoadFromString(string(content))
	checkError(t,err)
	want := testParserString
	if parser.section["forge.example"]["User"] != "hg" && parser.section["topsecret.server.example"]["ForwardX11"] != "no" && parser.section["topsecret.server.example"]["Port"] != "50022" {
		t.Errorf("want:\n%s\nGot:\n%s", want, content)
	}

	os.Remove("test_output.txt")
}
