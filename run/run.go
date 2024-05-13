package run

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"splitscript/config"
	"splitscript/utils"
)

var previous int = -1

func Run(conf config.Config) {
	if previous != -1 {
		err := syscall.Kill(previous, syscall.SIGKILL)
		if err != nil {
			fmt.Println(utils.Error.Render(fmt.Sprintf("Failed to kill %d: %s", previous, err.Error())))
			os.Exit(1)
		}
		fmt.Printf("Killed %d\n", previous)
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
		fmt.Println("Finished running, waiting for updates")
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
