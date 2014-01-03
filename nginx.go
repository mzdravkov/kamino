package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func addLocation(path string, options map[string]string) (err error) {
	conf, err := ioutil.ReadFile(Config["nginx_config_file"])
	if err != nil {
		return err
	}
	regex := regexp.MustCompile("server\\s*{")
	match := regex.FindIndex(conf)
	importToConf(conf, match, locationBlock(path, options))
	return
}

//pass string after the opening curly and you will receive the index of the closing or nil
func findMatchingCurly(str []byte) int {
	balance := 1
	for i, c := range str {
		if c == '{' {
			balance++
		} else if c == '}' {
			balance--
		}
		if balance == 0 {
			return i
		}
	}
	return -1
}

func importToConf(conf []byte, match []int, location string) error {
	closingCurlyIndex := findMatchingCurly(conf[match[1]+1:]) + match[1] + 1
	if closingCurlyIndex == -1 {
		panic("Can't find closing curly bracket for nginx conf")
	}
	block := make([]byte, len(location))
	copy(block[:], location)
	config := make([]byte, 0, len(conf)+len(location))
	config = append(config, conf[:match[1]]...)
	config = append(config, conf[match[1]:closingCurlyIndex]...)
	config = append(config, block...)
	config = append(config, conf[closingCurlyIndex:]...)

	err := ioutil.WriteFile(Config["nginx_config_file"], config, os.ModePerm)
	return err
}

func locationBlock(path string, options map[string]string) string {
	opts := make([]string, len(options))
	i := 0
	for k, v := range options {
		opts[i] = k + " " + v + ";"
		i++
	}
	return "\n\t\tlocation ~ ^/" + path + "(/.*|$) {\n\t\t\t" +
		strings.Join(opts, "\n\t\t\t") +
		"\n\t\t}\n"
}
