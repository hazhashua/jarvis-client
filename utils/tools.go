package utils

import (
	"fmt"
	"net"
	"time"
)

func CheckPorts(ip_port string, port_type string) bool {
	check := false
	now := time.Now().Format("2006-01-02 15:04:05")
	// 检测端口
	conn, err := net.DialTimeout(port_type, ip_port, 1*time.Second)
	if err != nil {
		fmt.Printf("检测%s超时, [%v], ip_port, 端口未开启(fail), error: %s\n", ip_port, now, err.Error())
	} else {
		if conn != nil {
			check = true
			Logger.Println("["+now+"]", ip_port, "端口已开启(success)!")
			conn.Close()
		} else {
			Logger.Println("["+now+"]", ip_port, "端口未开启(fail)!")
		}
	}
	return check

}

func Errorrecover() {
	if r := recover(); r != nil {
		Logger.Println("recover()返回异常，err: ", r)
	} else {
		Logger.Println("recover()返回正常")
	}
	DataStoreRemove(SourceMysqlDb)
	Logger.Printf("服务退出，清空数据库注册端口数据！")

}
