package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/sevlyar/go-daemon"
)

const pidFileName = "kamino.pid"
const workDir = "./"

var isDaemon bool

func isItDaemon() (bool, error) {
	pidFile := filepath.Join(workDir, pidFileName)

	fileContent, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return false, err
	}

	pid, err := strconv.Atoi(string(fileContent))
	if err != nil {
		return false, err
	}

	if pid == os.Getpid() {
		return true, nil
	}
	return false, nil
}

// Checks whether the process is already daemonized:
// If not, daemonizes it.
// If yes, starts the server
func daemonizeIfNeeded(c *cli.Context) {
	isDaemon, err := isItDaemon()
	if err != nil {
		log.Fatal(err)
	}

	if isDaemon {
		startServer(c)
	} else {
		daemonize()
	}
}

func daemonize() {
	// before starting the child process
	log.Println("starting kamino as a daemon...")

	context := &daemon.Context{
		PidFileName: pidFileName,
		PidFilePerm: 0644,
		LogFileName: "kamino_daemon.log",
		LogFilePerm: 0640,
		WorkDir:     workDir,
		Umask:       027,
		Args:        os.Args,
	}

	child, err := context.Reborn()
	if err != nil {
		panic(err)
	}

	if child != nil {
		// succeeded to start a child process. here the parent process writes to
		// a .pid file so that the child could know it's daemonized
		log.Println("kamino started as a daemon on pid ", child.Pid)
		pid := strconv.AppendInt(make([]byte, 0), int64(child.Pid), 10)
		ioutil.WriteFile(context.PidFileName, pid, context.PidFilePerm)
	} else {
		defer context.Release()
		//postchild()
	}
}
