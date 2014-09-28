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

func isItDaemon() error {
	pidFile := filepath.Join(workDir, pidFileName)
	if pid, err := ioutil.ReadFile(pidFile); err == nil {
		if pid, err := strconv.Atoi(string(pid)); err == nil {
			if pid == os.Getpid() {
				println("run server")
			}
			return nil
		} else {
			return err
		}
	}
	return nil
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
