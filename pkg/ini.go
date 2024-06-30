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
			err := parser.SetValue(currentSection, strings.TrimSpace(values[0]), strings.TrimSpace(values[1]))
			if err != nil {
				return err
			}
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
	var sections []string
	for s := range parser.section {
		sections = append(sections, s)
	}
	return sections
}

// Serialize and convert the INI content into a map with the format {section_name: {key1: val1, key2: val2}, ...}.
func (parser *IniParser) GetSection() map[string]map[string]string {
	m2 := make(map[string]map[string]string, len(parser.section))
	maps.Copy(m2, parser.section)
	return m2
}

// Retrieve the value of a key in a specific section.
// It could return an error if the section/key not found.
func (parser *IniParser) GetValue(sectionName, key string) (string, error) {
	sectionMap, ok := parser.section[sectionName]
	if !ok {
		return sectionMap[key], fmt.Errorf("can't get: %s section not found", sectionName)
	} else if sectionMap[key] == "" {
		return sectionMap[key], fmt.Errorf("can't get: %s key not found", key)
	}
	return sectionMap[key], nil
}

// Set the value of a key in a specific section.
// It can also be used to add a new key-value pair in a new section.
func (parser *IniParser) SetValue(section, key, value string) error {
	if section == "" {
		return errors.New("trying to add key and value for an empty section")
	} else if key == "" {
		return errors.New("can't set/add missing key")
	} else if value == "" {
		return errors.New("can't set/add missing value")
	}
	if parser.section == nil {
		parser.section = make(map[string]map[string]string)
	}
	if parser.section[section] == nil {
		parser.section[section] = make(map[string]string)
	}
	parser.section[section][key] = value
	return nil
}

// Convert the INI content back to a string format.
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

// Save the INI content to a file.
// Takes the path of the filename as an argument.
// An error could occur if it can't create/write to the file.
func (parser *IniParser) SaveToFile(path string) error {
	dat := parser.ToString()
	f, err := os.Create(path + ".ini")
	if err != nil {
		return fmt.Errorf("can't create %s file", path+".ini")
	}
	defer f.Close()
	_, err = f.WriteString(dat)
	if err != nil {
		return fmt.Errorf("can't write to %s file", path+".ini")
	}
	return err
}
