package main

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"strings"
)

// IniParser is a struct for parsing INI files/strings.
type IniParser struct {
	section map[string]map[string]string
}

// Intializing the parser
func initialize() *IniParser {
	parser := &IniParser{}
	if parser.section == nil {
		parser.section = make(map[string]map[string]string)
	}
	return parser
}

// Load INI content from a string.
// It could return an error if the provided string isn't valid.
func (parser *IniParser) LoadFromString(content string) error {
	var currentSection string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]
		} else if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, ";") {
			values := strings.Split(line, "=")
			if len(values) == 2 {
				err := parser.Set(currentSection, strings.TrimSpace(values[0]), strings.TrimSpace(values[1]))
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("invalid ini syntax %s", line)
			}
		} else if len(line) > 1 && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, ";") {
			return fmt.Errorf("invalid ini syntax %s", line)
		}
	}
	return nil
}

// Load INI content from a file.
// It could return an error if the file isn't valid or if the file isn't found.
func (parser *IniParser) LoadFromFile(fileName string) error {
	dat, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("file %s not found", fileName)
	}
	err = parser.LoadFromString(string(dat))
	return err
}

// Retrieve a list of all section names.
func (parser *IniParser) GetSectionNames() []string {
	sections := make([]string, len(parser.section))

	for s := range parser.section {
		sections = append(sections, s)
	}
	return sections
}

// Serialize and convert the INI content into a map with the format {section_name: {key1: val1, key2: val2}, ...}.
func (parser *IniParser) GetSections() map[string]map[string]string {
	m2 := make(map[string]map[string]string, len(parser.section))
	maps.Copy(m2, parser.section)
	return m2
}

// Retrieve the value of a key in a specific section.
// It could return an error if the section/key not found.
func (parser *IniParser) Get(sectionName, key string) (string, error) {
	sectionMap, ok := parser.section[sectionName]
	if !ok {
		return "", fmt.Errorf("can't get: %s section not found", sectionName)
	}
	return sectionMap[key], nil
}

// Set the value of a key in a specific section.
// It can also be used to add a new key-value pair in a new section.
func (parser *IniParser) Set(section, key, value string) error {
	if section == "" {
		return errors.New("trying to add key and value for an empty section")
	}
	if parser.section[section] == nil {
		parser.section[section] = make(map[string]string)
	}
	parser.section[section][key] = value
	return nil
}

// Convert the INI content back to a string format.
func (parser *IniParser) String() string {
	sectionNames := parser.GetSectionNames()
	var content string
	for _, section := range sectionNames {
		content += fmt.Sprintf("[%v]\n", section)
		keyValueMap := parser.section[section]
		for k, v := range keyValueMap {
			content += fmt.Sprintf("%v = %v \n", k, v)
		}
	}
	return content
}

// Save the INI content to a file.
// Takes the path of the filename as an argument.
// An error could occur if it can't create/write to the file.
func (parser *IniParser) SaveToFile(path string) error {
	dat := parser.String()
	err := os.WriteFile(path, []byte(dat), 0644)
	if err != nil {
		return errors.New("couldn't save to new file")
	}
	return err
}
