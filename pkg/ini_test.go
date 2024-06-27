package main

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	parser := &IniParser{}
	parser.Add("section1", "hello", "world")

	if parser.section["section1"]["hello"] != "world" {
		t.Errorf("got %q want world", parser.section["section1"]["hello"])
	}
}

func TestLoadFromString(t *testing.T) {
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

func TestloadFromFile(t *testing.T) {
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

func TestGetSectionNames(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	parser.LoadFromFile("golden_file.txt")
	t.Run("Get sections names", func(t *testing.T) {
		got := parser.GetSectionNames()
		want := []string{"forge.example", "topsecret.server.example"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v given", got, want)
		}
	})
}

func TestGetSection(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	parser.LoadFromFile("golden_file.txt")
	t.Run("Get sections", func(t *testing.T) {
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

func TestGetValue(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	parser.LoadFromFile("golden_file.txt")
	t.Run("Get value", func(t *testing.T) {
		got := parser.GetValue("forge.example", "User")
		want := "hg"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}

func TestSetValue(t *testing.T) {
	t.Helper()
	parser := &IniParser{}
	parser.LoadFromFile("golden_file.txt")
	t.Run("set value", func(t *testing.T) {
		parser.SetValue("forge.example", "user", "mariam")
		got := parser.section["forge.example"]["user"]

		if "mariam" != got {
			t.Errorf("got %s want mariam", got)
		}
	})
}
