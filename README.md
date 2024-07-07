# INI_Parser-Mariam_Mahrous

A Go-compatible library for reading/writing INI data from files or strings.

## How to Install and Run the Project

1. To use the IniParser package, simply import it into your Go project:
    ```go
    import "github.com/codescalersinternships/INI_Parser-Mariam_Mahrous"
    ```

2. To use the IniParser package, you will have to create an object first:
    ```go
    iniparser := &IniParser{}
    ```

## Overview

This library consists of 8 main functionalities:

1. **LoadFromString**

   Load INI content from a string.
   ```go
   content := `
   [forge.example]
   User = hg

   [topsecret.server.example]
   Port = 50022
   ForwardX11 = no
    `
   err := iniparser.LoadFromString(content)
   if err != nil {
        //handle error
   }
   ```

2. **LoadFromFile**

   Load INI content from a file.
   ```go
   err := iniparser.LoadFromFile("filename.txt")
   if err != nil {
       //handle error
   }
   ```

3. **GetSectionNames**

   Retrieve a list of all section names.
   ```go
   sections := iniparser.GetSectionNames()
   ```

4. **GetSections**

   Serialize and convert the INI content into a map with the format `{section_name: {key1: val1, key2: val2}, ...}`.
   ```go
   allSections := iniparser.GetSection()
   ```

5. **Get**

   Retrieve the value of a key in a specific section.
   ```go
   value, err := iniparser.GetValue("forge.example", "User")
   if err != nil {
       //handle error
   }
   ```

6. **Set**

   Set the value of a key in a specific section. 
   <br>
   It can also be used to add a new key-value pair in a new section.
   ```go
   err := iniparser.SetValue("forge.example", "User", "mariam")
   if err != nil {
        //handle error
   }
   ```

7. **ToString**

   Convert the INI content back to a string format.
   ```go
   contentString := iniparser.ToString()
   ```

8. **SaveToFile**

   Save the INI content to a file.
   ```go
   err := iniparser.SaveToFile("output")
   if err != nil {
        //handle error
   }
   ```

## Running Tests 
```sh
cd pkg
go test
```
    

## Important Remarks
1. There are no global keys, every key needs to be part of a section
2. Only key-value separator is "="
3. Comments aren't supported