package main

import (
	"github.com/sevlyar/go-daemon"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

func daemonize() {
	// before starting the child process
	log.Println("starting kamino as a daemon...")

	// remove the daemon flag so that the daemon doesn't try to daemonize itself
	args := make([]string, len(os.Args))
	for _, v := range os.Args {
		if v != "-d" || v != "--daemon" {
			args = append(args, v)
		}
	}

	context := &daemon.Context{
		PidFileName: pidFileName,
		PidFilePerm: 0644,
		LogFileName: "kamino_daemon.log",
		LogFilePerm: 0640,
		WorkDir:     workDir,
		Umask:       027,
		Args:        args,
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
