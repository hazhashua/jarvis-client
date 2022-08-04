package nodeexporter

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func CpuUsageGet() {
	// 获取cpu相关信息
	if fls, err := cpu.Percent(time.Second, true); err == nil {
		for _, f := range fls {
			fmt.Println("every cpu usage: ", f)
		}
	} else {
		fmt.Println("获取cpu使用率失败 ", err.Error())
	}

	f, _ := cpu.Percent(time.Second, false)
	fmt.Println("cpu usage: ", f)

	//获取cpu配额信息
	if infoStats, err := cpu.Info(); err == nil {
		for _, infoStat := range infoStats {
			fmt.Println("infoStat: ", infoStat)
		}
	}

}

func MemUsageGet() {
	// 获取物理内存使用信息
	if vms, err := mem.VirtualMemory(); err == nil {
		fmt.Println("vms.Used: ", vms.Used)
		fmt.Println("vms.UsedPercent: ", vms.UsedPercent)

	} else {
		fmt.Println("获取物理内存使用率失败 ", err.Error())
	}

	// 获取交换内存使用信息
	if sms, err := mem.SwapMemory(); err == nil {
		fmt.Println("vms.Used: ", sms.Used)
		fmt.Println("vms.UsedPercent: ", sms.UsedPercent)

	} else {
		fmt.Println("获取交换内存使用率失败 ", err.Error())
	}
}

func DiskUsageGet() {
	//获取各个磁盘的信息及使用率
	if ps, err := disk.Partitions(true); err == nil {
		for _, partitionInfo := range ps {
			fmt.Println("partitionInfo.Device: ", partitionInfo.Device)
			fmt.Println("partitionInfo.Fstype: ", partitionInfo.Fstype)
			fmt.Println("partitionInfo.Mountpoint: ", partitionInfo.Mountpoint)
			fmt.Println("partitionInfo.Opts: ", partitionInfo.Opts)

			if usage, err2 := disk.Usage(partitionInfo.Mountpoint); err2 == nil {
				fmt.Println("usage.Used, usage.Free, usage.Total, usage.UsedPercent")
				fmt.Println(usage.Used, usage.Free, usage.Total, usage.UsedPercent)
			}
		}

	}
	// 获取磁盘的io信息
	if ioStatMap, err := disk.IOCounters(); err == nil {
		for key, value := range ioStatMap {
			fmt.Println("key: ", key)
			fmt.Println("value: ", value)
		}
	}
}

func HostInfoGet() {
	if infoStat, err := host.Info(); err == nil {
		fmt.Println(infoStat.Hostname)
		fmt.Println(infoStat.BootTime)
		fmt.Println(infoStat.OS)
		fmt.Println(infoStat.Platform)
		fmt.Println(infoStat.KernelVersion)
	}
}

func NetInfoGet() {
	// 获取网卡信息及读写相关信息
	if connectionStats, err := net.Connections("all"); err == nil {
		fmt.Println("获取数据库的连接信息.....")
		for _, connectionStat := range connectionStats {
			fmt.Println("connectionStat: ", connectionStat)
		}
	}

	ioStats, _ := net.IOCounters(true)
	for _, ioStat := range ioStats {
		fmt.Println("ioStat: ", ioStat)
	}

	fmt.Println("arg is false ......")
	ioStats, _ = net.IOCounters(false)
	for _, ioStat := range ioStats {
		fmt.Println("ioStat: ", ioStat)
	}

}
