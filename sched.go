package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sj20082663/evloop"
)

var mainloop *evloop.EventLoop

func main() {
	mainloop = evloop.NewEventLoop()
	go thread()
	mainloop.Run()
}

func thread() {
	var cfg = newConfigFromFile(kConfigPathDefault)

	for _, schedTask := range cfg.SchedTasks {
		var d = schedTask.Duration * int(time.Second) / 10
		mainloop.RepeatFunc(func(stop *bool) {
			execSchedCmd(schedTask.CmdString)
			schedTask.Count++
			if model.Count >= schedTask.ExecNum {
				return
			}
		}, d)
	}
}

func execSchedCmd(s string) {
	var parts = strings.Split(s, " ")
	if len(parts) <= 0 {
		log.Printf("Invalid cmd %v", s)
		return
	}

	var cmd = exec.Command(parts[0], parts[1:]...)
	var err = cmd.Run()
	if err != nil {
		log.Print(err)
	}
}

func saveSession(task *tSchedTask) {

}

type tSchedTask struct {
	CmdString string
	Duration  int
	Id        string
	Count     int
	ExecNum   int
}

const kConfigPathDefault = "$HOME/.sched/config.json"

type tSchedConfig struct {
	SchedTasks []tSchedTask
}

func newConfigFromFile(path string) *tSchedConfig {
	path = os.ExpandEnv(path)
	var data, err = ioutil.ReadFile(path)
	if err != nil {
		log.Printf("open config.json error %v", err)
		return nil
	}

	var cfg = new(tSchedConfig)
	err = json.Unmarshal(data, cfg)
	if err != nil {
		log.Printf("config.json parse error %v", err)
		return nil
	}

	return cfg
}
