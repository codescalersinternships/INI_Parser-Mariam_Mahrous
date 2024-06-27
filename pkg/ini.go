package main

import (
	"fmt"
	"maps"
	"os"
	"strings"
)

type IniParser struct {
	section map[string]map[string]string
}

func (parser *IniParser) LoadFromString(content string) {
	var currentSection string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]
		} else if strings.Contains(line, "=") {
			values := strings.Split(line, " = ")
			parser.SetValue(currentSection, values[0], values[1])
		}
	}
}

func (parser *IniParser) LoadFromFile(fileName string) error {
	dat, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Errorf("file %s not found", fileName)
	}
	parser.LoadFromString(string(dat))
	return err
}

func (parser *IniParser) GetSectionNames() []string {
		var sections []string
	for s := range parser.section {
		sections = append(sections, s)
	}
	return sections
}

func (parser *IniParser) GetSection() map[string]map[string]string {
	m2 := make(map[string]map[string]string, len(parser.section))
	maps.Copy(m2, parser.section)
	return m2
}

func (parser *IniParser) GetValue(sectionName, key string) string {
	return parser.section[sectionName][key]
}

func (parser *IniParser) SetValue(section, key, value string) {
	if parser.section == nil {
		parser.section = make(map[string]map[string]string)
	}
	if parser.section[section] == nil {
		parser.section[section] = make(map[string]string)
	}
	parser.section[section][key] = value
}

func (parser *IniParser) ToString() string {
	sectionNames := parser.GetSectionNames()
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

func (parser *IniParser) SaveToFile(path string) error {
	dat := parser.ToString()
	f, err := os.Create(path + ".txt")
	if err != nil {
		fmt.Errorf("Can't Create  %q file", path+".txt")
		return err
	}
	defer f.Close()
	_, err = f.WriteString(dat)
	if err != nil {
		fmt.Errorf("Can't write in   %q file", path+".txt")
	}
	return err
}
