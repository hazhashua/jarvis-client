package utils

import (
	"strings"

	"github.com/shirou/gopsutil/net"
)

type FlowInfo struct {
	SentBytes      uint64
	ReceiveBytes   uint64
	PackageSent    uint64
	PackageReceive uint64
	Errin          uint64
	Errout         uint64
	Dropin         uint64
	Dropout        uint64
}

type NetInfo struct {
	DeviceIds map[string]FlowInfo
	EthInfo   map[string]string
	Ip        string
}

// 获取网卡网络信息
func NetInfoGet() *NetInfo {
	// 获取网卡信息及读写相关信息
	// //网络连接相关信息
	// if connectionStats, err := net.Connections("all"); err == nil {
	// 	fmt.Println("获取网络的连接信息.....")
	// 	for _, connectionStat := range connectionStats {
	// 		// fmt.Println("connectionStat: ", connectionStat)
	// 		fmt.Println("localAddr: ", connectionStat.Laddr.IP, "   destAddr: ", connectionStat.Raddr.IP)
	// 	}
	// }

	netInfo := NetInfo{}
	// var interfaceName, ip string
	interfaceInfo := make(map[string]string, 0)
	interfaces, err := net.Interfaces()
	if err != nil {
		return &netInfo
	}
	for _, interfaceStat := range interfaces {
		for _, v := range interfaceStat.Addrs {
			// fmt.Println("interfaceStat.Name: ", interfaceStat.Name, " net.InterfaceAddr: ", v.String(), v.Addr)
			Logger.Printf("interfaceStat.Name: %s, net.InterfaceAddr: %s\n", interfaceStat.Name, v.String())
			ips := strings.Split(v.Addr, "/")
			Logger.Printf("interface名称: %s\n", interfaceStat.Name)
			if ips[0] != "127.0.0.1" && len(strings.Split(ips[0], ".")) == 4 {
				interfaceInfo[interfaceStat.Name] = ips[0]
				netInfo.Ip = ips[0]
			}
		}
	}

	deviceFlows := make(map[string]FlowInfo, 0)
	ioStats, _ := net.IOCounters(true)
	for _, ioStat := range ioStats {
		deviceFlows[ioStat.Name] = FlowInfo{
			SentBytes:      ioStat.BytesSent,
			ReceiveBytes:   ioStat.BytesRecv,
			PackageSent:    ioStat.PacketsSent,
			PackageReceive: ioStat.PacketsRecv,
			Errin:          ioStat.Errin,
			Errout:         ioStat.Errout,
			Dropin:         ioStat.Dropin,
			Dropout:        ioStat.Dropout,
		}
	}
	netInfo.DeviceIds = deviceFlows
	netInfo.EthInfo = interfaceInfo
	return &netInfo
}
