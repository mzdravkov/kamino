package main

import (
	"github.com/msbranco/goconfig"
	"os"
	"path"
	"path/filepath"
)

var Config map[string]string

func ConfigToMap(conf *goconfig.ConfigFile) (m map[string]string) {
	m = make(map[string]string)
	options, err := conf.GetOptions("default")
	if err != nil {
		panic("Can't convert configs to map")
	}
	for _, k := range options {
		m[k], _ = conf.GetString("default", k)
	}
	return
}

// convert configs to map[string]string
func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("Can't find current directory for file config.go")
	}

	rawConfig, err := goconfig.ReadConfigFile(path.Join(dir, "config.cfg"))
	if err != nil {
		panic("Can't read config file")
	}

	Config = ConfigToMap(rawConfig)
}
