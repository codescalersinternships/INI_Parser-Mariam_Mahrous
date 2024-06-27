package main

import (
	"os"
	"reflect"
	"slices"
	"testing"
)

func TestParser_Add(t *testing.T) {
	parser := &IniParser{}
	parser.Add("section1", "hello", "world")

	if parser.section["section1"]["hello"] != "world" {
		t.Errorf("got %q want world", parser.section["section1"]["hello"])
	}
}

func TestParser_LoadFromString(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	testContent := `[forge.example]
	User = hg
	[topsecret.server.example]
	Port = 50022
	ForwardX11 = no`
	parser.LoadFromString(testContent)
	t.Run("Checking first section", func(t *testing.T) {
		got := parser.section["forge.example"]["User"]
		want := "hg"
		if got != want {
			t.Errorf("got %s want %s given", got, want)
		}
	})
	t.Run("Checking second section", func(t *testing.T) {
		got := parser.section["topsecret.server.example"]["Port"]
		want := "50022"
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
	parser.LoadFromFile("golden_file.txt")
	t.Run("Checking first section", func(t *testing.T) {
		got := parser.section["forge.example"]["User"]
		want := "hg"
		if got != want {
			t.Errorf("got %s want %s given", got, want)
		}
	})
	t.Run("Checking second section", func(t *testing.T) {
		got := parser.section["topsecret.server.example"]["Port"]
		want := "50022"
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

func TestParser_GetSectionNames(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	parser.LoadFromFile("golden_file.txt")
	got := parser.GetSectionNames()
	want := []string{"forge.example", "topsecret.server.example"}
	if !slices.Contains(got, "forge.example") || !slices.Contains(got, "topsecret.server.example") {
		t.Errorf("got %v want %v given", got, want)
	}
}

func TestParser_GetSection(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	parser.LoadFromFile("golden_file.txt")
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
}

func TestParser_GetValue(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	parser.LoadFromFile("golden_file.txt")
	got := parser.GetValue("forge.example", "User")
	want := "hg"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestParser_SetValue(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	parser.LoadFromFile("golden_file.txt")
	parser.SetValue("forge.example", "user", "mariam")
	got := parser.section["forge.example"]["user"]

	if "mariam" != got {
		t.Errorf("got %s want mariam", got)
	}
}

func TestParser_ToString(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	parser.LoadFromFile("golden_file.txt")
	got := parser.ToString()
	parser.LoadFromString(got)
	want := `
[forge.example]
User = hg

[topsecret.server.example]
Port = 50022
ForwardX11 = no
`
	if parser.section["forge.example"]["User"] != "hg" && parser.section["topsecret.server.example"]["ForwardX11"] != "no" && parser.section["topsecret.server.example"]["Port"] != "50022" {
		t.Errorf("want:\n%s\nGot:\n%s", want, got)
	}
}

func TestParser_SaveToFile(t *testing.T) {
	parser := &IniParser{}
	parser.LoadFromFile("golden_file.txt")
	parser.SaveToFile("test_output")
	content, err := os.ReadFile("test_output.txt")
	if err != nil {
		t.Fatalf("Failed to read the file: %s", err)
	}
	parser.LoadFromString(string(content))
	want := `
[forge.example]
User = hg

[topsecret.server.example]
Port = 50022
ForwardX11 = no
`

	if parser.section["forge.example"]["User"] != "hg" && parser.section["topsecret.server.example"]["ForwardX11"] != "no" && parser.section["topsecret.server.example"]["Port"] != "50022" {
		t.Errorf("want:\n%s\nGot:\n%s", want, content)
	}

	os.Remove("test_output.txt")
}
