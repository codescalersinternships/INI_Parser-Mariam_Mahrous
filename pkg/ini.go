package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatalf("Can't Read file %s\n", e)
	}
}

type IniParser struct {
	section map[string]map[string]string
}

func (parser *IniParser) Add(section, key, value string) {
	if parser.section == nil {
		parser.section = make(map[string]map[string]string)
	}
	if parser.section[section] == nil {
		parser.section[section] = make(map[string]string)
	}
	parser.section[section][key] = value
}

func (parser *IniParser) LoadFromString(content string) *IniParser {
	var currentSection string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]
		} else if strings.Contains(line, "=") {
			values := strings.Split(line, " = ")
			if parser.section == nil {
				parser.section = make(map[string]map[string]string)
			}
			if parser.section[currentSection] == nil {
				parser.section[currentSection] = make(map[string]string)
			}
			parser.section[currentSection][values[0]] = values[1]
		}
	}
	return parser
}

func (parser *IniParser) LoadFromFile(fileName string) *IniParser {
	dat, err := os.ReadFile(fileName)
	check(err)
	return parser.LoadFromString(string(dat))
}

func (parser *IniParser) GetSectionNames() []string {
	var sections []string
	for s := range parser.section {
		sections = append(sections, s)
	}
	return sections
}

func (parser *IniParser) GetSection() map[string]map[string]string {
	return parser.section
}

// Check lw i tried to get a section msh mawgod yet
func (parser *IniParser) GetValue(sectionName, key string) string {
	return parser.section[sectionName][key]
}

// Check lw i tried to set a section msh mawgod yet
func (parser *IniParser) SetValue(section, key, value string) {
	parser.section[section][key] = value
}

func (parser *IniParser) ToString() string {
	sectionNames := parser.GetSectionNames()
	fmt.Println(sectionNames)
	var content string
	for _, section := range sectionNames {
		content += "\n" + "[" + section + "]\n"
		smallMap := parser.section[section]
		for k, v := range smallMap {
			content += k + " = " + v + "\n"
		}
	}
	return content
}

func (parser *IniParser) SaveToFile(path string) {
	dat := parser.ToString()
	f, err := os.Create(path + ".txt")
	check(err)
	defer f.Close()
	_, err = f.WriteString(dat)
	check(err)
}
