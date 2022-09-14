package service_alive

import (
	"fmt"
	"os/exec"
	"strings"
)

func IsProcessRunning(processName string) bool {
	//进程探活，探测进程是否正在运行

	// byteList := bytes.Buffer{}
	// cmd := exec.Command("ps", "-C", processName)
	// cmd.Stdout = &byteList
	// cmd.Run()

	find := false
	cmd := exec.Command("ps", "-C", processName)
	processes := make([]string, 0)
	if bytes, err := cmd.Output(); err == nil {
		processes = strings.Fields(string(bytes))
		fmt.Println("processes: ", processes)
	}
	for _, process := range processes {
		if process == processName || strings.Contains(processName, process) {
			find = true
			fmt.Println("发现进程: ", processName)
		}
	}
	return find
}

type Host struct {
	domainAddr *string
	host       *string
	ip         *string
}

func PingAddr(address string) bool {
	pingOK := false
	cmd := exec.Command("ping", "-c5", address)
	err := cmd.Run()
	if err != nil {
		fmt.Println("地址ping失败 ", err.Error())
	} else {
		fmt.Println("地址ping成功 ")
		pingOK = true
	}
	return pingOK
}

func PingOk(host Host) bool {
	// 探测主机或者ip是否网络联通
	connected := false
	if host.domainAddr != nil {
		connected = PingAddr(*host.domainAddr)
	}
	if host.host != nil {
		//基于主机名探测
		connected = PingAddr(*host.host)
	}
	if host.ip != nil {
		connected = PingAddr(*host.ip)
	}
	return connected
}
