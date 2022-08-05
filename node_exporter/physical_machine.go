package nodeexporter

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

type CpuInfo struct {
	cores int
	usage float64
}

func CpuUsageGet() *CpuInfo {
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

	cores := 0
	//获取cpu配额信息
	if infoStats, err := cpu.Info(); err == nil {
		for _, infoStat := range infoStats {
			fmt.Println("infoStat: ", infoStat)
			cores += int(infoStat.Cores)
		}
	}

	return &CpuInfo{
		cores: cores,
		usage: f[0],
	}

}

type Memory struct {
	total       uint64
	used        uint64
	available   uint64
	cached      uint64
	free        uint64
	usedPercent float64
}

func MemUsageGet() *Memory {
	// 获取物理内存使用信息
	var vms *mem.VirtualMemoryStat
	var err error
	if vms, err = mem.VirtualMemory(); err == nil {
		fmt.Println("vms.Used: ", vms.Used)
		fmt.Println("vms.Total: ", vms.Total)
		fmt.Println("vms.Available: ", vms.Available)
		fmt.Println("vms.Cached: ", vms.Cached)
		fmt.Println("vms.Buffers: ", vms.Buffers)
		fmt.Println("vms.Free", vms.Free)
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

	return &Memory{
		total:       vms.Total,
		used:        vms.Used,
		available:   vms.Available,
		cached:      vms.Cached,
		free:        vms.Free,
		usedPercent: vms.UsedPercent,
	}
}

type Disk struct {
	// 磁盘设备编号
	deviceIds  []string
	mountPoint []string
	total      []uint64
	used       []uint64
	free       []uint64
	// 磁盘读写速率
	readBytes  map[string]uint64
	writeBytes map[string]uint64
}

func DiskUsageGet() *Disk {
	//获取各个磁盘的信息及使用率

	deviceIds := make([]string, 0)
	mountPoint := make([]string, 0)
	total := make([]uint64, 0)
	used := make([]uint64, 0)
	free := make([]uint64, 0)

	readBytes := make(map[string]uint64)
	writeBytes := make(map[string]uint64)

	if ps, err := disk.Partitions(false); err == nil {
		for _, partitionInfo := range ps {
			fmt.Println("partitionInfo.Device: ", partitionInfo.Device)
			fmt.Println("partitionInfo.Fstype: ", partitionInfo.Fstype)
			fmt.Println("partitionInfo.Mountpoint: ", partitionInfo.Mountpoint)
			fmt.Println("partitionInfo.Opts: ", partitionInfo.Opts)
			var usage *disk.UsageStat
			var err2 error
			if usage, err2 = disk.Usage(partitionInfo.Mountpoint); err2 == nil {
				fmt.Println("usage.Used, usage.Free, usage.Total, usage.UsedPercent")
				fmt.Println(usage.Used, usage.Free, usage.Total, usage.UsedPercent)
			}

			deviceIds = append(deviceIds, partitionInfo.Device)
			mountPoint = append(mountPoint, partitionInfo.Mountpoint)
			total = append(total, usage.Total)
			used = append(used, usage.Used)
			free = append(free, usage.Free)
		}
	}
	// 获取磁盘的io信息
	if ioStatMap, err := disk.IOCounters(); err == nil {
		for key, value := range ioStatMap {
			fmt.Println("key: ", key)
			fmt.Println("value: ", value)

			readBytes[key] = value.ReadBytes
			writeBytes[key] = value.WriteBytes
		}
	}

	return &Disk{
		deviceIds:  deviceIds,
		mountPoint: mountPoint,
		total:      total,
		used:       used,
		free:       free,

		readBytes:  readBytes,
		writeBytes: writeBytes,
	}
}

type HostInfo struct {
	// 主机名
	hostName string
	os       string
	bootTime uint64
	plaform  string
	// 主机的主要ip
	ip string
}

func HostInfoGet() *HostInfo {
	// 返回主机信息
	var infoStat *host.InfoStat
	var err error
	if infoStat, err = host.Info(); err == nil {
		fmt.Println(infoStat.Hostname)
		fmt.Println(infoStat.BootTime)
		fmt.Println(infoStat.OS)
		fmt.Println(infoStat.Platform)
		fmt.Println(infoStat.KernelVersion)
	}
	return &HostInfo{
		hostName: infoStat.Hostname,
		os:       infoStat.OS,
		bootTime: infoStat.BootTime,
		plaform:  infoStat.Platform,
		ip:       infoStat.HostID,
	}
}

type FlowInfo struct {
	sentBytes      uint64
	receiveBytes   uint64
	packageSent    uint64
	packageReceive uint64
	errin          uint64
	errout         uint64
	dropin         uint64
	dropout        uint64
}

type NetInfo struct {
	deviceIds map[string]FlowInfo
}

// 获取网卡网络信息
func NetInfoGet() {
	// 获取网卡信息及读写相关信息

	//网络连接相关信息
	if connectionStats, err := net.Connections("all"); err == nil {
		fmt.Println("获取网络的连接信息.....")
		for _, connectionStat := range connectionStats {
			fmt.Println("connectionStat: ", connectionStat)
			fmt.Println(connectionStat.Laddr.IP)
			fmt.Println(connectionStat.Raddr.IP)
		}
	}

	netInfo := NetInfo{}
	deviceFlows := make(map[string]FlowInfo, 0)
	ioStats, _ := net.IOCounters(true)
	for _, ioStat := range ioStats {
		fmt.Println("ioStat: ", ioStat)
		deviceFlows[ioStat.Name] = FlowInfo{
			sentBytes:      ioStat.BytesSent,
			receiveBytes:   ioStat.BytesRecv,
			packageSent:    ioStat.PacketsSent,
			packageReceive: ioStat.PacketsRecv,

			errin:   ioStat.Errin,
			errout:  ioStat.Errout,
			dropin:  ioStat.Dropin,
			dropout: ioStat.Dropout,
		}
	}

	netInfo.deviceIds = deviceFlows

}

type ProcessIO struct {
	processId  int32
	readBytes  uint64
	writeBytes uint64
}

// 存储进程相关信息
type ProcessInfo struct {
	processIoMap map[int32]ProcessIO
}

func ProcessnfoGet() *ProcessInfo {

	processInfo := ProcessInfo{}
	IoMap := make(map[int32]ProcessIO)
	if processes, err := process.Processes(); err == nil {
		for _, process := range processes {
			fmt.Println("process.Pid: ", process.Pid)
			ioCounterStat, _ := process.IOCounters()
			fmt.Println("进程读字节数: ", ioCounterStat.ReadBytes)
			fmt.Println("进程写字节数: ", ioCounterStat.WriteBytes)
			IoMap[process.Pid] = ProcessIO{
				processId:  process.Pid,
				readBytes:  ioCounterStat.ReadBytes,
				writeBytes: ioCounterStat.WriteBytes,
			}

		}
	}
	processInfo.processIoMap = IoMap
	return &processInfo
}
