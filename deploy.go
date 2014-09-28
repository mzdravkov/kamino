package main

import (
	"github.com/codegangsta/cli"
	"github.com/dynport/gossh"
	"io/ioutil"
	"log"
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

	if memLim := Config["container_memory_limit"]; memLim != "" {
		options = append(options, "-m "+memLim)
	}

	if runOpts := Config["docker_run_options"]; runOpts != "" {
		options = append(options, runOpts)
	}

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

func serverOptions(subdomain string) map[string]string {
	return map[string]string{
		"listen":      "80",
		"server_name": subdomain + "." + Config["server_ip"]}
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

func deploy(c *cli.Context) {
	if c.Args().First() != "" {
		port := uint16(c.Int("port"))
		if port == 0 {
			port = findFreePort()
		}
		if err := deployRemotely(c.Args().First(), port); err != nil {
			log.Fatal(err)
			return
		}
	} else {
		log.Println("You have to pass a name for the tenant as the first argument to deploy. Use 'kamino help deploy' for more info")
	}
}

func deployRemotely(name string, port uint16) (err error) {
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
	if Config["nginx_use_locations"] == "true" {
		addLocation(name, locationOpts)
	} else {
		serverOpts := serverOptions(name)
		addServer(locationOpts, serverOpts)
	}
	nginxReload := exec.Command("sh", "-c", "sudo "+Config["nginx_bin"]+" -s reload")
	err = nginxReload.Run()
	return
}

// returns a function of type gossh.Writer func(...interface{})
// MakeLogger just adds a prefix (DEBUG, INFO, ERROR)
func makeLogger(prefix string) gossh.Writer {
	return func(args ...interface{}) {
		log.Println((append([]interface{}{prefix}, args...))...)
	}
}

//func deployRemotely(name string, port uint16) (err error) {
//	client := gossh.New("some.host", "user")
//	// my default agent authentication is used. use
//	// client.SetPassword("<secret>")
//	// for password authentication
//	client.DebugWriter = MakeLogger("DEBUG")
//	client.InfoWriter = MakeLogger("INFO ")
//	client.ErrorWriter = MakeLogger("ERROR")
//
//	defer client.Close()
//	rsp, e := client.Execute("uptime")
//	if e != nil {
//		client.ErrorWriter(e.Error())
//	}
//	client.InfoWriter(rsp.String())
//}
