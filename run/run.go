package run

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/splitscriptjs/cli/config"
	"github.com/splitscriptjs/cli/utils"
)

var previous int = -1
var closedProcess int = -1

func Run(conf config.Config) {
	if previous != -1 {
		p, err := os.FindProcess(previous)
		if err != nil {
			fmt.Println(utils.Error.Render(fmt.Sprintf("Failed to kill %d: %s", previous, err.Error())))
			os.Exit(1)
		}
		p.Signal(syscall.SIGTERM)
		fmt.Printf("Killed %d\n", previous)
		closedProcess = previous
	}
	fileToRun, err := utils.GenerateDevFileName(conf, conf.Main)
	fmt.Println(utils.Info.Render("Running `" + conf.Main + "`"))
	if err != nil {
		fmt.Println(utils.Error.Render(fmt.Sprintf("Failed to run: %s", err.Error())))
		os.Exit(1)
	}
	node := exec.Command("node", fileToRun)
	node.Env = append(node.Env, "ROOT"+conf.Out, "CONFIG_LOCATION=./")

	stdout, err := node.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
		}
		node.Wait()
		if closedProcess != node.Process.Pid {
			fmt.Println("Finished running, waiting for updates")
		}
	}()
	stderr, err := node.StderrPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(utils.Error.Render(m))
		}
		node.Wait()
	}()
	node.Start()
	previous = node.Process.Pid
}
