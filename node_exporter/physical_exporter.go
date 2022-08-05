package nodeexporter

import "github.com/prometheus/client_golang/prometheus"

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
	cpuCoresDesc        *prometheus.Desc
	cpuCoresValType     prometheus.ValueType
	cpuUsageDesc        *prometheus.Desc
	cpuUsageValType     prometheus.ValueType
	memTotalDesc        *prometheus.Desc
	memTotalValType     prometheus.ValueType
	memUsageDesc        *prometheus.Desc
	memUsageValType     prometheus.ValueType
	diskTotalDesc       *prometheus.Desc
	diskTotalValType    prometheus.ValueType
	diskUsageDesc       *prometheus.Desc
	diskUsageValType    prometheus.ValueType
	diskReadDesc        *prometheus.Desc
	diskReadValType     prometheus.ValueType
	diskWriteDesc       *prometheus.Desc
	diskWriteValType    prometheus.ValueType
	networkReadDesc     *prometheus.Desc
	networkReadValType  prometheus.ValueType
	networkWriteDesc    *prometheus.Desc
	networkWriteValType prometheus.ValueType
}

type VirtualDesc struct {
}

type MachineExporter struct {
	physicalMetrics PhysicalDesc
	virtualMetrics  VirtualDesc
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
		[]string{},
		prometheus.Labels{})
	physicalMetrics.memUsageValType = prometheus.GaugeValue
	physicalMetrics.diskTotalDesc = prometheus.NewDesc("disk_total", "机器的总磁盘大小",
		[]string{},
		prometheus.Labels{})
	physicalMetrics.diskTotalValType = prometheus.GaugeValue
	physicalMetrics.diskUsageDesc = prometheus.NewDesc("disk_usage", "机器的磁盘使用率",
		[]string{},
		prometheus.Labels{})
	physicalMetrics.diskUsageValType = prometheus.GaugeValue
	physicalMetrics.diskReadDesc = prometheus.NewDesc("disk_read_bytes", "磁盘每秒读速率",
		[]string{},
		prometheus.Labels{})
	physicalMetrics.diskReadValType = prometheus.GaugeValue
	physicalMetrics.diskWriteDesc = prometheus.NewDesc("disk_write_bytes", "磁盘每秒写速率",
		[]string{},
		prometheus.Labels{})
	physicalMetrics.diskWriteValType = prometheus.GaugeValue
	physicalMetrics.networkReadDesc = prometheus.NewDesc("network_receive_bytes", "网络设备每秒接收到的字节数",
		[]string{},
		prometheus.Labels{})
	physicalMetrics.networkReadValType = prometheus.GaugeValue
	physicalMetrics.networkWriteDesc = prometheus.NewDesc("network_transmit_bytes", "网络设备每秒发送的字节数",
		[]string{},
		prometheus.Labels{})
	physicalMetrics.networkWriteValType = prometheus.GaugeValue
	var virtualMetrics VirtualDesc
	return &MachineExporter{
		physicalMetrics: physicalMetrics,
		virtualMetrics:  virtualMetrics,
	}

}

func (e *MachineExporter) Describe(ch chan<- *prometheus.Desc) {

	ch <- e.physicalMetrics.cpuCoresDesc
	ch <- e.physicalMetrics.cpuUsageDesc
	ch <- e.physicalMetrics.memTotalDesc
	ch <- e.physicalMetrics.memUsageDesc
	ch <- e.physicalMetrics.diskTotalDesc
	ch <- e.physicalMetrics.diskUsageDesc
	ch <- e.physicalMetrics.diskReadDesc
	ch <- e.physicalMetrics.diskWriteDesc
	ch <- e.physicalMetrics.networkReadDesc
	ch <- e.physicalMetrics.networkWriteDesc

}

func (e *MachineExporter) collect(ch chan<- prometheus.Metric) {
	// CpuUsageGet()
	// MemUsageGet()
	// DiskUsageGet()
	// NetInfoGet()
	// HostInfoGet()
	// ProcessnfoGet()

}
