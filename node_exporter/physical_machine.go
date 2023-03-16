package nodeexporter

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	// "strings"
	"metric_exporter/config"
	"metric_exporter/utils"
	"time"

	"github.com/cedarwu/lsblk"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"gopkg.in/yaml.v2"
)

func parseNodeConfig() *config.NodeConfig {
	var nodeConfig config.NodeConfig
	if bytes, err := ioutil.ReadFile("./node_exporter/config.yaml"); err == nil {
		err2 := yaml.Unmarshal(bytes, &nodeConfig)
		if err2 != nil {
			utils.Logger.Printf("解析node配置文件失败 error: %s\n", err2.Error())
		}
	}
	utils.Logger.Println("nodeConfig.Cluster.name: ", nodeConfig.Cluster.Name)
	return &nodeConfig
}

type CpuInfo struct {
	cores int
	usage float64
	sys   float64
	user  float64
	io    float64
}

func cpuUsageDetailGet() *CpuInfo {
	// 获取系统核心数
	f, _ := cpu.Percent(time.Second, false)
	utils.Logger.Println("主机cpu usage: ", f)
	cores := 0
	//获取cpu配额信息
	if infoStats, err := cpu.Info(); err == nil {
		for _, infoStat := range infoStats {
			// fmt.Println("infoStat: ", infoStat)
			cores += int(infoStat.Cores)
		}
	}

	// 获取CPU使用率
	s := fmt.Sprintf(`top -bn1 | fgrep 'Cpu(s)' | awk -F" " '{print $2" "$4" "$8" "$10}'`)
	utils.Logger.Printf("top = %s", s)
	cmd := exec.Command("bash", "-c", s)
	bs, e := cmd.Output()
	if e != nil {
		utils.Logger.Printf("failed due to :%v", e)
		panic(e)
	}
	utils.Logger.Printf("%v", string(bs))
	raw := string(bs)
	percents := strings.Split(raw, " ")
	//用户态使用率 系统态使用率 cpu空闲率 io使用占比
	fmt.Printf("%v, %v, %v, %v", percents[0], percents[1], percents[2], percents[3])

	var usage, sys, user, io float64
	if idle, err := strconv.ParseFloat(percents[2], 64); err != nil {
		usage = 1 - idle/100
	}
	user, _ = strconv.ParseFloat(percents[0], 64)
	sys, _ = strconv.ParseFloat(percents[1], 64)
	io, _ = strconv.ParseFloat(percents[3], 64)

	return &CpuInfo{
		cores: cores,
		usage: usage,
		sys:   sys,
		user:  user,
		io:    io,
	}
}

func CpuUsageGet() *CpuInfo {
	// 获取cpu相关信息
	f, _ := cpu.Percent(time.Second, false)
	utils.Logger.Println("主机cpu usage: ", f)

	cores := 0
	//获取cpu配额信息
	if infoStats, err := cpu.Info(); err == nil {
		for _, infoStat := range infoStats {
			// fmt.Println("infoStat: ", infoStat)
			cores += int(infoStat.Cores)
		}
	}

	return &CpuInfo{
		cores: cores,
		usage: f[0] / 100.0,
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
		utils.Logger.Println("主机虚拟内存信息: ", vms)
		// fmt.Println("vms.Used: ", vms.Used)
		// fmt.Println("vms.Total: ", vms.Total)
		// fmt.Println("vms.Available: ", vms.Available)
		// fmt.Println("vms.Cached: ", vms.Cached)
		// fmt.Println("vms.Buffers: ", vms.Buffers)
		// fmt.Println("vms.Free", vms.Free)
		// fmt.Println("vms.UsedPercent: ", vms.UsedPercent)
	} else {
		utils.Logger.Println("获取主机物理内存使用率失败   error: ", err.Error())
	}

	// 获取交换内存使用信息
	if sms, err := mem.SwapMemory(); err == nil {
		utils.Logger.Println("交换主机内存信息: ", sms)

	} else {
		utils.Logger.Println("获取主机交换内存使用率失败 ", err.Error())
	}

	return &Memory{
		total:       vms.Total,
		used:        vms.Used,
		available:   vms.Available,
		cached:      vms.Cached,
		free:        vms.Free,
		usedPercent: vms.UsedPercent / 100.0,
	}
}

type Disk struct {
	// 磁盘设备编号
	deviceNum      int
	deviceIds      []string
	mountPoint     []string
	filesystemType []string
	total          []uint64
	used           []uint64
	free           []uint64
	// 磁盘读写速率
	ioDeviceNum int
	readBytes   map[string]uint64
	writeBytes  map[string]uint64
	readCount   map[string]uint64
	writeCount  map[string]uint64
}

func DiskDeviceNum() int {
	partitionStats, _ := disk.Partitions(false)
	return len(partitionStats)
}

func Uptime() uint64 {
	if uptime, err := host.Uptime(); err != nil {
		return uptime
	}
	return 0

}

func DiskIoDeviceNum() int {
	ioStatMap, _ := disk.IOCounters()
	utils.Logger.Println("主机io操作的磁盘数: ", len(ioStatMap))
	return len(ioStatMap)
}

func DiskUsageGet() *Disk {
	//获取各个磁盘的信息及使用率
	deviceIds := make([]string, 0)
	mountPoint := make([]string, 0)
	filesystemType := make([]string, 0)
	total := make([]uint64, 0)
	used := make([]uint64, 0)
	free := make([]uint64, 0)

	readBytes := make(map[string]uint64)
	writeBytes := make(map[string]uint64)
	readCount := make(map[string]uint64)
	writeCount := make(map[string]uint64)

	if devices, err3 := lsblk.ListDevices(); err3 != nil {
		for deviceName, deviceObj := range devices {
			utils.Logger.Printf("deviceName: %s fstype: %s", deviceName, deviceObj.Fstype)
		}

	}

	if ps, err := disk.Partitions(false); err == nil {
		for _, partitionInfo := range ps {

			// fmt.Println("partitionInfo.Device: ", partitionInfo.Device)
			// fmt.Println("partitionInfo.Fstype: ", partitionInfo.Fstype)
			// fmt.Println("partitionInfo.Mountpoint: ", partitionInfo.Mountpoint)
			// fmt.Println("partitionInfo.Opts: ", partitionInfo.Opts)
			var usage *disk.UsageStat
			var err2 error
			if usage, err2 = disk.Usage(partitionInfo.Mountpoint); err2 != nil {
				utils.Logger.Printf("获取主机disk.Usage失败   error: %s", err2.Error())
			}

			deviceIds = append(deviceIds, partitionInfo.Device)
			mountPoint = append(mountPoint, partitionInfo.Mountpoint)
			filesystemType = append(filesystemType, partitionInfo.Fstype)
			total = append(total, usage.Total)
			used = append(used, usage.Used)
			free = append(free, usage.Free)
		}
	}
	// 获取磁盘的io信息
	ioDeviceNum := 0
	if ioStatMap, err := disk.IOCounters(); err == nil {
		for key, value := range ioStatMap {
			ioDeviceNum += 1
			readBytes[key] = value.ReadBytes
			writeBytes[key] = value.WriteBytes

			readCount[key] = value.ReadCount
			writeCount[key] = value.WriteCount
		}
	}

	return &Disk{
		deviceNum:      len(deviceIds),
		deviceIds:      deviceIds,
		mountPoint:     mountPoint,
		filesystemType: filesystemType,
		total:          total,
		used:           used,
		free:           free,
		ioDeviceNum:    ioDeviceNum,
		readBytes:      readBytes,
		writeBytes:     writeBytes,
		readCount:      readCount,
		writeCount:     writeCount,
	}
}

type HostInfo struct {
	// 主机名
	hostName string
	os       string
	bootTime uint64
	upTime   float64
	plaform  string
	// 主机的主要ip
	id string
}

func HostInfoGet() *HostInfo {
	// 返回主机信息
	var infoStat *host.InfoStat
	var err error
	var upTime uint64
	if infoStat, err = host.Info(); err != nil {
		utils.Logger.Println("获取的主机信息: ", infoStat)
	}

	if upTime, err = host.Uptime(); err != nil {
		utils.Logger.Printf("获取设备运行时长报错：%v", err)
	}

	return &HostInfo{
		hostName: infoStat.Hostname,
		os:       infoStat.OS,
		bootTime: infoStat.BootTime,
		upTime:   float64(upTime / 60 / 60 / 24),
		plaform:  infoStat.Platform,
		id:       infoStat.HostID,
	}
}

func NetDeviceNum() int {
	ioStats, _ := net.IOCounters(true)
	utils.Logger.Println("获取的主机io设备信息: ", ioStats)
	return len(ioStats)
}

// func getIpFromAddr(addr net.Addr) net.IP {
// 	var ip net.IP
// 	switch v := addr.(type) {
// 	case *net.IPNet:
// 		ip = v.IP
// 	case *net.IPAddr:
// 		ip = v.IP
// 	}
// 	if ip == nil || ip.IsLoopback() {
// 		return nil
// 	}
// 	ip = ip.To4()
// 	if ip == nil {
// 		return nil // not an ipv4 address
// 	}

// 	return ip
// }

type ProcessIO struct {
	processId  int32
	readBytes  uint64
	writeBytes uint64
}

// 存储进程相关信息
type ProcessInfo struct {
	processIoMap map[int32]ProcessIO
}

func ProcessInfoGet() (int, *ProcessInfo) {
	processInfo := ProcessInfo{}
	IoMap := make(map[int32]ProcessIO)
	var processes []*process.Process
	var err error
	if processes, err = process.Processes(); err == nil {
		for _, process := range processes {
			if ioCounterStat, err := process.IOCounters(); err != nil {
				utils.Logger.Printf("process.IOCounters() error:%s\n", err.Error())
			} else {
				// utils.Logger.Printf("process.Pid: %d  进程读字节数: %d  进程写字节数: %d \n", process.Pid, ioCounterStat.ReadBytes, ioCounterStat.WriteBytes)
				IoMap[process.Pid] = ProcessIO{
					processId:  process.Pid,
					readBytes:  ioCounterStat.ReadBytes,
					writeBytes: ioCounterStat.WriteBytes,
				}
			}
		}
	}
	processInfo.processIoMap = IoMap
	return len(processes), &processInfo
}
