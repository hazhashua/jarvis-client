package nodeexporter

import (
	"fmt"
	"metric_exporter/config"
	"metric_exporter/utils"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/load"
)

type PhysicalMetricsValue struct {
	cpuCores       int
	cpuUsage       float32
	memTotal       int
	memUsage       float32
	diskTotal      int
	diskUsage      float32
	diskRead       int32
	diskWrite      int32
	diskReadCount  int32
	diskWriteCount int32
	networkRead    int32
	networkWrite   int32
}

type PhysicalDesc struct {
	upTimeDesc    *prometheus.Desc
	upTimeValType prometheus.ValueType

	cpuCoresDesc    *prometheus.Desc
	cpuCoresValType prometheus.ValueType
	cpuUsageDesc    *prometheus.Desc
	cpuUsageValType prometheus.ValueType

	cpuSysUsageDesc     *prometheus.Desc
	cpuSysUsageValType  prometheus.ValueType
	cpuUserUsageDesc    *prometheus.Desc
	cpuUserUsageValType prometheus.ValueType
	cpuIoUsageDesc      *prometheus.Desc
	cpuIoUsageValType   prometheus.ValueType

	memTotalDesc          *prometheus.Desc
	memTotalValType       prometheus.ValueType
	memUsageDesc          *prometheus.Desc
	memUsageValType       prometheus.ValueType
	diskTotalDesc         []*prometheus.Desc
	diskTotalValType      []prometheus.ValueType
	diskUsedDesc          []*prometheus.Desc
	diskUsedValType       []prometheus.ValueType
	diskReadDesc          []*prometheus.Desc
	diskReadValType       []prometheus.ValueType
	diskWriteDesc         []*prometheus.Desc
	diskWriteValType      []prometheus.ValueType
	diskReadCountDesc     []*prometheus.Desc
	diskReadCountValType  []prometheus.ValueType
	diskWriteCountDesc    []*prometheus.Desc
	diskWriteCountValType []prometheus.ValueType
	load1Desc             *prometheus.Desc
	load1ValueType        prometheus.ValueType
	load5Desc             *prometheus.Desc
	load5ValueType        prometheus.ValueType
	load15Desc            *prometheus.Desc
	load15ValueType       prometheus.ValueType
	networkReceiveDesc    []*prometheus.Desc
	networkReceiveValType []prometheus.ValueType
	networkSentDesc       []*prometheus.Desc
	networkSentValType    []prometheus.ValueType
	processInfoDesc       []*prometheus.Desc
	processInfoValType    []prometheus.ValueType
}

type VirtualDesc struct {
}

type MachineExporter struct {
	physicalDiskNum      int
	physicalIoDiskNum    int
	physicalNetDeviceNum int
	processNum           int
	physicalMetrics      PhysicalDesc
	physicalMetricsData  *ProcessInfo
	netInfoData          *utils.NetInfo
	cpuData              *CpuInfo
	hostInfoData         *HostInfo
	memoryData           *Memory
	diskData             *Disk
	load1                float64
	load5                float64
	load15               float64
	virtualMetrics       VirtualDesc
}

func NewNodeExporter() *MachineExporter {

	// 构建MachineExporter对象
	var physicalMetrics PhysicalDesc

	physicalMetrics.upTimeDesc = prometheus.NewDesc("uptime", "系统运行时长",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	physicalMetrics.upTimeValType = prometheus.CounterValue

	physicalMetrics.cpuCoresDesc = prometheus.NewDesc("cpu_cores_total", "cpu总核心数",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	physicalMetrics.cpuCoresValType = prometheus.GaugeValue
	physicalMetrics.cpuUsageDesc = prometheus.NewDesc("cpu_usage", "cpu当前使用率",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	physicalMetrics.cpuUsageValType = prometheus.GaugeValue

	physicalMetrics.cpuSysUsageDesc = prometheus.NewDesc("cpu_sys_usage", "cpu当前系统态使用率",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	physicalMetrics.cpuSysUsageValType = prometheus.GaugeValue

	physicalMetrics.cpuUserUsageDesc = prometheus.NewDesc("cpu_user_usage", "cpu当前用户态使用率",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	physicalMetrics.cpuUserUsageValType = prometheus.GaugeValue

	physicalMetrics.cpuIoUsageDesc = prometheus.NewDesc("cpu_io_usage", "cpu当前io堵塞使用率",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	physicalMetrics.cpuIoUsageValType = prometheus.GaugeValue

	physicalMetrics.memTotalDesc = prometheus.NewDesc("memory_total", "总内存大小",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	physicalMetrics.memTotalValType = prometheus.GaugeValue
	physicalMetrics.memUsageDesc = prometheus.NewDesc("memory_usage", "内存使用率",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	physicalMetrics.memUsageValType = prometheus.GaugeValue
	deviceNum := DiskDeviceNum()
	deviceIoNum := DiskIoDeviceNum()
	physicalMetrics.diskTotalDesc = make([]*prometheus.Desc, deviceNum)
	physicalMetrics.diskTotalValType = make([]prometheus.ValueType, deviceNum)
	physicalMetrics.diskUsedDesc = make([]*prometheus.Desc, deviceNum)
	physicalMetrics.diskUsedValType = make([]prometheus.ValueType, deviceNum)
	for i := 0; i < deviceNum; i++ {
		physicalMetrics.diskTotalDesc[i] = prometheus.NewDesc("disk_total", "机器的总磁盘大小",
			[]string{"cluster", "host", "ip", "device_id", "filesystem_type", "mount_path"},
			prometheus.Labels{})
		physicalMetrics.diskTotalValType[i] = prometheus.GaugeValue

		physicalMetrics.diskUsedDesc[i] = prometheus.NewDesc("disk_used", "机器使用的磁盘大小",
			[]string{"cluster", "host", "ip", "device_id", "filesystem_type", "mount_path"},
			prometheus.Labels{})
		physicalMetrics.diskUsedValType[i] = prometheus.GaugeValue

	}
	physicalMetrics.diskReadDesc = make([]*prometheus.Desc, deviceIoNum)
	physicalMetrics.diskReadValType = make([]prometheus.ValueType, deviceIoNum)
	physicalMetrics.diskWriteDesc = make([]*prometheus.Desc, deviceIoNum)
	physicalMetrics.diskWriteValType = make([]prometheus.ValueType, deviceIoNum)
	physicalMetrics.diskReadCountDesc = make([]*prometheus.Desc, deviceIoNum)
	physicalMetrics.diskReadCountValType = make([]prometheus.ValueType, deviceIoNum)
	physicalMetrics.diskWriteCountDesc = make([]*prometheus.Desc, deviceIoNum)
	physicalMetrics.diskWriteCountValType = make([]prometheus.ValueType, deviceIoNum)

	for i := 0; i < deviceIoNum; i++ {
		physicalMetrics.diskReadDesc[i] = prometheus.NewDesc("disk_read_bytes", "磁盘每秒读速率",
			[]string{"cluster", "host", "ip", "device"},
			prometheus.Labels{})
		physicalMetrics.diskReadValType[i] = prometheus.GaugeValue

		physicalMetrics.diskWriteDesc[i] = prometheus.NewDesc("disk_write_bytes", "磁盘每秒写速率",
			[]string{"cluster", "host", "ip", "device"},
			prometheus.Labels{})
		physicalMetrics.diskWriteValType[i] = prometheus.GaugeValue

		physicalMetrics.diskReadCountDesc[i] = prometheus.NewDesc("disk_read_count", "磁盘read次数",
			[]string{"cluster", "host", "ip", "device"},
			prometheus.Labels{})
		physicalMetrics.diskReadCountValType[i] = prometheus.CounterValue

		physicalMetrics.diskWriteCountDesc[i] = prometheus.NewDesc("disk_write_count", "磁盘write次数",
			[]string{"cluster", "host", "ip", "device"},
			prometheus.Labels{})
		physicalMetrics.diskWriteCountValType[i] = prometheus.CounterValue

	}

	// 声明最近1 5 15 分钟的系统负载
	physicalMetrics.load1Desc = prometheus.NewDesc("os_load1", "最近1分钟内的负载",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	physicalMetrics.load1ValueType = prometheus.GaugeValue

	physicalMetrics.load5Desc = prometheus.NewDesc("os_load5", "最近5分钟内的负载",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	physicalMetrics.load5ValueType = prometheus.GaugeValue

	physicalMetrics.load15Desc = prometheus.NewDesc("os_load15", "最近15分钟内的负载",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	physicalMetrics.load15ValueType = prometheus.GaugeValue

	netDeviceNum := NetDeviceNum()
	physicalMetrics.networkSentDesc = make([]*prometheus.Desc, netDeviceNum)
	physicalMetrics.networkSentValType = make([]prometheus.ValueType, netDeviceNum)
	physicalMetrics.networkReceiveDesc = make([]*prometheus.Desc, netDeviceNum)
	physicalMetrics.networkReceiveValType = make([]prometheus.ValueType, netDeviceNum)
	for i := 0; i < netDeviceNum; i++ {
		physicalMetrics.networkReceiveDesc[i] = prometheus.NewDesc("network_receive_bytes", "网络设备每秒接收到的字节数",
			[]string{"cluster", "host", "ip", "net_name"},
			prometheus.Labels{})
		physicalMetrics.networkReceiveValType[i] = prometheus.GaugeValue
		physicalMetrics.networkSentDesc[i] = prometheus.NewDesc("network_transmit_bytes", "网络设备每秒发送的字节数",
			[]string{"cluster", "host", "ip", "net_name"},
			prometheus.Labels{})
		physicalMetrics.networkSentValType[i] = prometheus.GaugeValue
	}
	diskIoDeviceNum := DiskIoDeviceNum()

	processNum, processInfos := ProcessInfoGet()
	physicalMetrics.processInfoDesc = make([]*prometheus.Desc, processNum)
	physicalMetrics.processInfoValType = make([]prometheus.ValueType, processNum)
	for i := 0; i < processNum; i++ {
		physicalMetrics.processInfoDesc[i] = prometheus.NewDesc("process_info", "主机运行进程的相关信息",
			[]string{"cluster", "host", "ip", "id", "read_bytes", "write_bytes"},
			prometheus.Labels{})
		physicalMetrics.processInfoValType[i] = prometheus.GaugeValue
	}

	netInfo := utils.NetInfoGet()
	hostInfo := HostInfoGet()
	memoryInfo := MemUsageGet()
	diskInfo := DiskUsageGet()
	cpuInfo := cpuUsageDetailGet()

	var load1, load5, load15 float64
	if avgSta, err := load.Avg(); err != nil {
		load1 = avgSta.Load1
		load5 = avgSta.Load5
		load15 = avgSta.Load15
	}

	var virtualMetrics VirtualDesc
	return &MachineExporter{
		physicalDiskNum:      deviceNum,
		physicalIoDiskNum:    diskIoDeviceNum,
		processNum:           processNum,
		physicalMetrics:      physicalMetrics,
		physicalMetricsData:  processInfos,
		physicalNetDeviceNum: netDeviceNum,
		cpuData:              cpuInfo,
		netInfoData:          netInfo,
		hostInfoData:         hostInfo,
		memoryData:           memoryInfo,
		diskData:             diskInfo,
		load1:                load1,
		load5:                load5,
		load15:               load15,
		virtualMetrics:       virtualMetrics,
	}
}

func (e *MachineExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.physicalMetrics.upTimeDesc
	ch <- e.physicalMetrics.cpuCoresDesc
	ch <- e.physicalMetrics.cpuUsageDesc
	ch <- e.physicalMetrics.cpuSysUsageDesc
	ch <- e.physicalMetrics.cpuUserUsageDesc
	ch <- e.physicalMetrics.cpuIoUsageDesc
	ch <- e.physicalMetrics.memTotalDesc
	ch <- e.physicalMetrics.memUsageDesc
	for i := 0; i < e.physicalDiskNum; i++ {
		ch <- e.physicalMetrics.diskTotalDesc[i]
		ch <- e.physicalMetrics.diskUsedDesc[i]
	}
	for i := 0; i < e.physicalIoDiskNum; i++ {
		ch <- e.physicalMetrics.diskReadDesc[i]
		ch <- e.physicalMetrics.diskWriteDesc[i]
		ch <- e.physicalMetrics.diskReadCountDesc[i]
		ch <- e.physicalMetrics.diskWriteCountDesc[i]
	}
	for i := 0; i < e.physicalNetDeviceNum; i++ {
		ch <- e.physicalMetrics.networkReceiveDesc[i]
		ch <- e.physicalMetrics.networkSentDesc[i]
	}
	for i := 0; i < e.processNum; i++ {
		ch <- e.physicalMetrics.processInfoDesc[i]
	}

}

func (e *MachineExporter) Collect(ch chan<- prometheus.Metric) {
	e = NewNodeExporter()
	// nodeConfig := parseNodeConfig()
	nodeConfig, _ := (utils.ConfigStruct.ConfigData[config.NODE]).(config.NodeConfig)
	utils.Logger.Println("nodeConfig: ", nodeConfig)
	hostInfo := e.hostInfoData
	utils.Logger.Printf("获取到的主机名称: %s\n", hostInfo.hostName)
	cpuInfo := e.cpuData
	netInfo := e.netInfoData

	//输出系统运行时长
	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.upTimeDesc, e.physicalMetrics.upTimeValType, hostInfo.upTime, nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip)

	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.cpuCoresDesc, e.physicalMetrics.cpuCoresValType,
		float64(cpuInfo.cores), nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip)
	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.cpuUsageDesc, e.physicalMetrics.cpuUsageValType,
		cpuInfo.usage, nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip)

	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.cpuSysUsageDesc, e.physicalMetrics.cpuSysUsageValType, cpuInfo.sys, nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip)
	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.cpuUserUsageDesc, e.physicalMetrics.cpuUserUsageValType, cpuInfo.user, nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip)
	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.cpuIoUsageDesc, e.physicalMetrics.cpuIoUsageValType, cpuInfo.io, nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip)

	memory := e.memoryData
	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.memTotalDesc, e.physicalMetrics.memTotalValType,
		float64(memory.total), nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip)

	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.memUsageDesc, e.physicalMetrics.memTotalValType,
		float64(memory.usedPercent), nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip)

	disk := e.diskData
	for i := 0; i < e.physicalDiskNum; i++ {
		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.diskTotalDesc[i], e.physicalMetrics.diskTotalValType[i],
			float64(disk.total[i]), nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip, disk.deviceIds[i], disk.filesystemType[i], disk.mountPoint[i])

		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.diskUsedDesc[i], e.physicalMetrics.memTotalValType,
			float64(disk.used[i]), nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip, disk.deviceIds[i], disk.filesystemType[i], disk.mountPoint[i])

	}
	// devices := make([]string, 0)
	// readBytess := make([]uint64, 0)
	// writeBytess := make([]uint64, 0)
	// for device, bytes := range disk.readBytes {
	// 	devices = append(devices, device)
	// 	readBytess = append(readBytess, bytes)
	// }
	// for _, bytes := range disk.writeBytes {
	// 	writeBytess = append(writeBytess, bytes)
	// }

	// for i := 0; i < e.physicalIoDiskNum; i++ {
	// 	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.diskReadDesc[i], e.physicalMetrics.diskReadValType[i],
	// 		float64(readBytess[i]), nodeConfig.Cluster.name, hostInfo.hostName, "", disk.deviceIds[i], disk.mountPoint[i])

	// 	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.diskWriteDesc[i], e.physicalMetrics.diskReadValType[i],
	// 		float64(writeBytess[i]), nodeConfig.Cluster.name, hostInfo.hostName, "", disk.deviceIds[i], disk.mountPoint[i])
	// }
	i := 0
	for key, value := range disk.readBytes {
		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.diskReadDesc[i], e.physicalMetrics.diskReadValType[i],
			float64(value), nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip, key)
		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.diskReadCountDesc[i], e.physicalMetrics.diskReadCountValType[i],
			float64(disk.readCount[key]), nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip, key)
		i += 1
	}
	i = 0
	for key, value := range disk.writeBytes {
		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.diskWriteDesc[i], e.physicalMetrics.diskWriteValType[i],
			float64(value), nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip, key)
		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.diskWriteCountDesc[i], e.physicalMetrics.diskWriteCountValType[i],
			float64(disk.writeCount[key]), nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip, key)
		i += 1
	}

	// 写入负载的相关信息
	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.load1Desc, e.physicalMetrics.load1ValueType,
		e.load1, nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip)
	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.load5Desc, e.physicalMetrics.load5ValueType,
		e.load5, nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip)
	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.load15Desc, e.physicalMetrics.load15ValueType,
		e.load15, nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip)

	deviceNames := make([]string, 0)
	flowInfos := make([]utils.FlowInfo, 0)
	for deviceName, flowInfo := range netInfo.DeviceIds {
		utils.Logger.Println("deviceName: ", deviceName, " flowInfo: ", flowInfo)
		deviceNames = append(deviceNames, deviceName)
		flowInfos = append(flowInfos, flowInfo)
	}
	for idx, deviceName := range deviceNames {
		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.networkReceiveDesc[idx], e.physicalMetrics.networkReceiveValType[idx],
			float64(flowInfos[idx].ReceiveBytes), nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip, deviceName)

		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.networkSentDesc[idx], e.physicalMetrics.networkSentValType[idx],
			float64(flowInfos[idx].SentBytes), nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip, deviceName)
	}

	processInfos := e.physicalMetricsData

	idx := 0
	for key, value := range processInfos.processIoMap {

		// {"cluster", "host", "ip", "id", "read_bytes", "write_bytes"}
		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.processInfoDesc[idx], e.physicalMetrics.processInfoValType[idx],
			1, nodeConfig.Cluster.Name, hostInfo.hostName, netInfo.Ip, fmt.Sprintf("%d", key), fmt.Sprintf("%d", value.readBytes), fmt.Sprintf("%d", value.writeBytes))
		idx += 1
	}

}
