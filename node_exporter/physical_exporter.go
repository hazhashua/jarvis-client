package nodeexporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type PhysicalMetricsValue struct {
	cpuCores     int
	cpuUsage     float32
	memTotal     int
	memUsage     float32
	diskTotal    int
	diskUsage    float32
	diskRead     int32
	diskWrite    int32
	networkRead  int32
	networkWrite int32
}

type PhysicalDesc struct {
	cpuCoresDesc          *prometheus.Desc
	cpuCoresValType       prometheus.ValueType
	cpuUsageDesc          *prometheus.Desc
	cpuUsageValType       prometheus.ValueType
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
	networkReceiveDesc    []*prometheus.Desc
	networkReceiveValType []prometheus.ValueType
	networkSentDesc       []*prometheus.Desc
	networkSentValType    []prometheus.ValueType
}

type VirtualDesc struct {
}

type MachineExporter struct {
	physicalDiskNum      int
	physicalIoDiskNum    int
	physicalNetDeviceNum int
	physicalMetrics      PhysicalDesc
	virtualMetrics       VirtualDesc
}

func NewNodeExporter() *MachineExporter {
	// 构建MachineExporter对象
	var physicalMetrics PhysicalDesc
	physicalMetrics.cpuCoresDesc = prometheus.NewDesc("cpu_cores_total", "cpu总核心数",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	physicalMetrics.cpuCoresValType = prometheus.GaugeValue
	physicalMetrics.cpuUsageDesc = prometheus.NewDesc("cpu_usage", "cpu当前使用率",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	physicalMetrics.cpuUsageValType = prometheus.GaugeValue
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
	for i := 0; i < deviceNum; i++ {
		physicalMetrics.diskTotalDesc[i] = prometheus.NewDesc("disk_total", "机器的总磁盘大小",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		physicalMetrics.diskTotalValType[i] = prometheus.GaugeValue

		physicalMetrics.diskUsedDesc[i] = prometheus.NewDesc("disk_usage", "机器的磁盘使用率",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		physicalMetrics.diskUsedValType[i] = prometheus.GaugeValue
	}

	for i := 0; i < deviceIoNum; i++ {
		physicalMetrics.diskReadDesc[i] = prometheus.NewDesc("disk_read_bytes", "磁盘每秒读速率",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		physicalMetrics.diskReadValType[i] = prometheus.GaugeValue

		physicalMetrics.diskWriteDesc[i] = prometheus.NewDesc("disk_write_bytes", "磁盘每秒写速率",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		physicalMetrics.diskWriteValType[i] = prometheus.GaugeValue
	}

	netDeviceNum := NetDeviceNum()
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

	var virtualMetrics VirtualDesc
	return &MachineExporter{
		physicalDiskNum:      deviceNum,
		physicalMetrics:      physicalMetrics,
		physicalNetDeviceNum: netDeviceNum,
		virtualMetrics:       virtualMetrics,
	}
}

func (e *MachineExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.physicalMetrics.cpuCoresDesc
	ch <- e.physicalMetrics.cpuUsageDesc
	ch <- e.physicalMetrics.memTotalDesc
	ch <- e.physicalMetrics.memUsageDesc
	for i := 0; i < e.physicalDiskNum; i++ {
		ch <- e.physicalMetrics.diskTotalDesc[i]
		ch <- e.physicalMetrics.diskUsedDesc[i]
	}
	for i := 0; i < e.physicalIoDiskNum; i++ {
		ch <- e.physicalMetrics.diskReadDesc[i]
		ch <- e.physicalMetrics.diskWriteDesc[i]
	}
	for i := 0; i < e.physicalNetDeviceNum; i++ {
		ch <- e.physicalMetrics.networkReceiveDesc[i]
		ch <- e.physicalMetrics.networkSentDesc[i]
	}
}

func (e *MachineExporter) collect(ch chan<- prometheus.Metric) {

	// ch <- e.physicalMetrics.cpuCoresDesc
	// ch <- e.physicalMetrics.cpuUsageDesc
	// ch <- e.physicalMetrics.memTotalDesc
	// ch <- e.physicalMetrics.memUsageDesc
	// ch <- e.physicalMetrics.diskTotalDesc
	// ch <- e.physicalMetrics.diskUsageDesc
	// ch <- e.physicalMetrics.diskReadDesc
	// ch <- e.physicalMetrics.diskWriteDesc
	// ch <- e.physicalMetrics.networkReadDesc
	// ch <- e.physicalMetrics.networkWriteDesc

	nodeConfig := parseNodeConfig()
	hostInfo := HostInfoGet()
	fmt.Println(hostInfo.hostName)

	cpuInfo := CpuUsageGet()
	// {"cluster", "host", "ip"}

	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.cpuCoresDesc, e.physicalMetrics.cpuCoresValType,
		float64(cpuInfo.cores), nodeConfig.Cluster.name, hostInfo.hostName, "")
	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.cpuUsageDesc, e.physicalMetrics.cpuUsageValType,
		cpuInfo.usage, nodeConfig.Cluster.name, hostInfo.hostName, "")

	memory := MemUsageGet()
	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.memTotalDesc, e.physicalMetrics.memTotalValType,
		float64(memory.total), nodeConfig.Cluster.name, hostInfo.hostName, "")

	ch <- prometheus.MustNewConstMetric(e.physicalMetrics.memUsageDesc, e.physicalMetrics.memTotalValType,
		float64(memory.usedPercent), nodeConfig.Cluster.name, hostInfo.hostName, "")

	disk := DiskUsageGet()
	for i := 0; i < e.physicalDiskNum; i++ {
		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.diskTotalDesc[i], e.physicalMetrics.diskTotalValType[i],
			float64(disk.total[i]), nodeConfig.Cluster.name, hostInfo.hostName, "")

		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.diskUsedDesc[i], e.physicalMetrics.memTotalValType,
			float64(disk.used[i]), nodeConfig.Cluster.name, hostInfo.hostName, "")
	}
	devices := make([]string, 0)
	readBytess := make([]uint64, 0)
	writeBytess := make([]uint64, 0)
	for device, bytes := range disk.readBytes {
		devices = append(devices, device)
		readBytess = append(readBytess, bytes)
	}
	for _, bytes := range disk.writeBytes {
		writeBytess = append(writeBytess, bytes)
	}

	for i := 0; i < e.physicalIoDiskNum; i++ {
		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.diskReadDesc[i], e.physicalMetrics.diskReadValType[i],
			float64(readBytess[i]), nodeConfig.Cluster.name, hostInfo.hostName, "")

		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.diskWriteDesc[i], e.physicalMetrics.diskReadValType[i],
			float64(writeBytess[i]), nodeConfig.Cluster.name, hostInfo.hostName, "")
	}
	netInfo := NetInfoGet()
	deviceNames := make([]string, 0)
	flowInfos := make([]FlowInfo, 0)
	for deviceName, flowInfo := range netInfo.deviceIds {
		fmt.Println("deviceName: ", deviceName)
		fmt.Println("flowInfo: ", flowInfo)
		deviceNames = append(deviceNames, deviceName)
		flowInfos = append(flowInfos, flowInfo)
	}
	for idx, deviceName := range deviceNames {
		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.networkReceiveDesc[idx], e.physicalMetrics.networkReceiveValType[idx],
			float64(flowInfos[idx].receiveBytes), nodeConfig.Cluster.name, hostInfo.hostName, "", deviceName)

		ch <- prometheus.MustNewConstMetric(e.physicalMetrics.networkSentDesc[idx], e.physicalMetrics.networkSentValType[idx],
			float64(flowInfos[idx].sentBytes), nodeConfig.Cluster.name, hostInfo.hostName, "", deviceName)
	}

}
