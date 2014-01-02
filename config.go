package main

import (
	"github.com/msbranco/goconfig"
)

var (
	rawConfig, _ = goconfig.ReadConfigFile("config.cfg")
	Config       map[string]string
)

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

func init() {
	Config = ConfigToMap(rawConfig)
}
