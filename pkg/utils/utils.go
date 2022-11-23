// Copyright (c) 2022 CrowdStrike, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package utils

import (
	"os"
	"path/filepath"
	"regexp"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func ValidateRegExp(regex string, value string) bool {
	if ok, _ := regexp.MatchString(regex, value); !ok {
		return false
	}

	return true
}

func ConfigExists(path string) {
	path = filepath.Clean(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(path), 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(filepath.Clean(path), flags, 0600)
	if err != nil {
		log.Fatalf("Failed to open %s for writing: %s", path, err)
	}
	if err = file.Close(); err != nil {
		log.Fatalf("Failed to close %s: %s", path, err)
	}
}

func ReadYAML(yamlContent interface{}, fileName string) any {
	yamlFile, err := os.ReadFile(filepath.Clean(fileName))
	if err != nil {
		log.Fatalf("Unable to read data from the file: %v, %s", err, fileName)
	}
	err = yaml.Unmarshal(yamlFile, &yamlContent)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return yamlContent
}

// WriteYAML writes a YAML file to disk.
func WriteYAML(yamlContent interface{}, fileName string) {
	cfg, err := yaml.Marshal(&yamlContent)
	if err != nil {
		log.Errorf("Error while Marshaling. %v", err)
	}

	err = os.WriteFile(filepath.Clean(fileName), cfg, 0600)
	if err != nil {
		log.Fatalf("Unable to write data into the file: %v, %s", err, fileName)
	}
}
