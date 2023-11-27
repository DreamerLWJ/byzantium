package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

func main() {
	//pids, err := GetPortPID(9999)
	//if err != nil {
	//	panic(err)
	//}

	pids := []int{79422}

	for _, pid := range pids {
		process, err := os.FindProcess(pid)
		if err != nil {
			panic(err)
		}
		err = process.Signal(syscall.Signal(0))
		if err != nil {
			fmt.Println(err)
		}
		err = process.Signal(syscall.SIGTERM)
		if err != nil {
			panic(err)
		}
		err = WaitProcessDone(pid)
		fmt.Printf("process.Wait, err:%s\n", err)
	}
}

func WaitProcessDone(pid int) error {
	for {
		process, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		err = process.Signal(syscall.Signal(0))
		if err != nil {
			if errors.Is(err, os.ErrProcessDone) {
				fmt.Println("process done")
			} else {
				fmt.Printf("WaitProcessDone|err:%s\n", err)
			}
			return nil
		}
	}
}

// GetPortPID 通过 Linux 命令获取占用指定端口的进程 PID
func GetPortPID(port int) (pids []int, err error) {
	cmd := exec.Command("lsof", "-t", fmt.Sprintf("-i:%d", port))
	output, err := cmd.Output()
	if err != nil {
		return pids, err
	}

	outputStr := string(output)

	outputStr = strings.Trim(outputStr, "\n")
	rawPids := strings.Split(outputStr, "\n")
	pidMap := make(map[int]struct{})
	for _, pidStr := range rawPids {
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return pids, errors.Errorf("GetPortPID|err:%s", err)
		}
		pidMap[pid] = struct{}{}
	}

	for k, _ := range pidMap {
		pids = append(pids, k)
	}
	return pids, nil
}

// KillPid 杀死进程
func KillPid(pid string) error {
	killCmd := exec.Command("kill", "-15", pid)
	killCmd.Stdout = os.Stdout
	killCmd.Stderr = os.Stderr
	if err := killCmd.Run(); err != nil {
		return errors.Errorf("KillPid|err:%s", err)
	}
	return nil
}

// WaitPid 等待进程结束
func WaitPid(pid string) error {
	waitCmd := exec.Command("wait", pid)
	waitCmd.Stdout = os.Stdout
	waitCmd.Stderr = os.Stderr
	if err := waitCmd.Run(); err != nil {
		return errors.Errorf("WaitPid|err:%s", err)
	}
	return nil
}
