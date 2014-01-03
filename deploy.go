package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func isPortFree(port uint16) bool {
	ip, _ := Config["server_ip"]
	ipWithPort := strings.Join([]string{ip, ":", strconv.Itoa(int(port))}, "")
	command := strings.Join([]string{"netstat -na | grep tcp | grep ", ipWithPort}, "")
	output, _ := exec.Command("sh", "-c", command).Output()
	if len(output) > 0 {
		return false
	}
	return true
}

//returns random free port as uint16
func findFreePort() uint16 {
	port := randUint16(1025, 64000)
	for !isPortFree(port) {
		port = randUint16(1025, 64000)
	}
	return port
}

//returns pseudo random uint16 between min and max
func randUint16(min, max uint16) uint16 {
	rand.Seed(time.Now().UnixNano())
	rInt := rand.Int31n(int32(max - min))
	return uint16(rInt) + min
}

func dockerRunOptions(name string, port uint16) []string {
	portStr := strconv.Itoa(int(port))
	confSrc := Config["tenants_configs_dir"]
	confDest := Config["tenants_config_path"]
	configVolumeOption := strings.Join([]string{"-v ", confSrc, ":", confDest, ":ro"}, "")
	pluginsSrc := filepath.Join(Config["tenants_plugins_dir"], name)
	pluginsDest := Config["tenants_plugins_path"]
	pluginsVolumeOption := strings.Join([]string{"-v ", pluginsSrc, ":", pluginsDest, ":ro"}, "")
	options := []string{
		"-d",
		"-name " + name,
		strings.Join([]string{"-p ", portStr, ":", Config["tenants_port"]}, ""),
		configVolumeOption,
		pluginsVolumeOption}

	return options
}

func dockerRunArguments() []string {
	runArguments := []string{Config["docker_image"]}
	if Config["container_entry_command"] != "" {
		return append(runArguments, Config["container_entry_command"])
	}
	return runArguments
}

func locationOptions(port uint16) map[string]string {
	return map[string]string{
		"proxy_pass": "http://" + Config["server_ip"] + ":" + strconv.Itoa(int(port)) + "/$1",
		"resolver":   Config["nginx_location_resolver"]}
}

//copies the default config file to tenant's specific config file
func makeTenantConfig(name string) (err error) {
	data, err := ioutil.ReadFile(Config["tenants_default_config"])
	filename := filepath.Join(Config["tenants_configs_dir"], name+".yml")
	ioutil.WriteFile(filename, data, os.ModePerm)
	return
}

func makePluginsDir(name string) (err error) {
	dir := filepath.Join(Config["tenants_plugins_dir"], name)
	err = os.MkdirAll(dir, os.ModePerm)
	return
}

func Deploy(name string, port uint16) (err error) {
	opts := strings.Join(dockerRunOptions(name, port), " ")
	args := strings.Join(dockerRunArguments(), " ")
	if err = makeTenantConfig(name); err != nil {
		return
	}
	if err = makePluginsDir(name); err != nil {
		return
	}
	cmd := exec.Command("sh", "-c", "docker run "+opts+" "+args)
	if err = cmd.Run(); err != nil {
		return
	}
	locationOpts := locationOptions(port)
	addLocation(name, locationOpts)
	nginxReload := exec.Command("sh", "-c", "sudo "+Config["nginx_bin"]+" -s reload")
	err = nginxReload.Run()
	return
}
