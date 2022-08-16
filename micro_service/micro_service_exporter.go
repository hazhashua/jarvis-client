package micro_service

import (
	"fmt"
	"metric_exporter/config"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// microservice指标导出器
type MicroServiceExporter struct {

	// type MyK8sNodeInfo struct {
	// 	Name              string           `json:"name"`
	// 	Ip                string           `json:"ip"`
	// 	CreationTimestamp string           `json:"creationTimestamp"`
	// 	NodeCapacityS     *NodeCapacity    `json:"nodeCapacity"`
	// 	NodeAllocatableS  *NodeAllocatable `json:"nodeAllocatable"`
	// 	MemoryPressure    bool             `json:"memoryPressure"`
	// 	DiskPressure      bool             `json:"diskPressure"`
	// 	PidPressure       bool             `json:"pidPressure"`
	// 	IsReady           bool             `json:"isReady"`
	// }

	k8sConfig K8sConfig
	// 微服务资源及节点状态数据
	// 包含多个node节点的资源及状态
	nodeDescs     []K8sNodeDesc
	nodeDatas     []*MyK8sNodeInfo
	nodeInfoDescs []K8sNodeInfoDesc

	// 微服务服务数据
	serviceInfoDescs []K8sServiceDesc
	serviceInfoDatas map[string]ServiceInfo

	// 微服务上应用的pod状态
	podInfoDescs []K8sPodDesc
	podInfoDatas []*MyK8sPodInfo

	// type MyK8sPodInfo struct {
	// 	Name              string               `json:"name"`
	// 	CreationTimestamp string               `json:"creationTimestamp"`
	// 	App               string               `json:"app"`
	// 	Containers        int                  `json:"containers"`
	// 	Status            string               `json:"status"`
	// 	IsInitialized     bool                 `json:"isInitialized"`
	// 	IsReady           bool                 `json:"isReady"`
	// 	IsContainersReady bool                 `json:"isContainersReady"`
	// 	IsPodScheduled    bool                 `json:"isPodScheduled"`
	// 	containersStatus  []*MyContainerStatus `json:""`
	// }

}
type K8sServiceDesc struct {
	ServiceInfoDesc    *prometheus.Desc
	ServiceInfoValType prometheus.ValueType
}

type K8sPodDesc struct {
	PodInfoDesc    *prometheus.Desc
	PodInfoValType prometheus.ValueType
}

type K8sNodeInfoDesc struct {
	NodeInfoDesc    *prometheus.Desc
	NodeInfoValType prometheus.ValueType
}

type K8sNodeDesc struct {
	MaxCpuDesc            *prometheus.Desc
	MaxCpuValType         prometheus.ValueType
	MaxDiskStorageDesc    *prometheus.Desc
	MaxDiskStorageValType prometheus.ValueType
	MaxMemoryDesc         *prometheus.Desc
	MaxMemoryValType      prometheus.ValueType
	MaxPodsDesc           *prometheus.Desc
	MaxPodsValType        prometheus.ValueType

	AllocateCpuDesc            *prometheus.Desc
	AllocateCpuValType         prometheus.ValueType
	AllocateDiskStorageDesc    *prometheus.Desc
	AllocateDiskStorageValType prometheus.ValueType
	AllocateMemoryDesc         *prometheus.Desc
	AllocateMemoryValType      prometheus.ValueType
	AllocatePodsDesc           *prometheus.Desc
	AllocatePodsValType        prometheus.ValueType

	CpuUsedDesc       *prometheus.Desc
	CpuUsedValType    prometheus.ValueType
	MemoryUsedDesc    *prometheus.Desc
	MemoryUsedValType prometheus.ValueType

	// MemoryPressureOk      *prometheus.Desc
	// MemoryPressureValType prometheus.ValueType
	// DiskPressureOk        *prometheus.Desc
	// DiskPressureValType   prometheus.ValueType
	// PidPressureOk         *prometheus.Desc
	// PidPressureValType    prometheus.ValueType
	// IsReady        *prometheus.Desc
	// IsReadyValType prometheus.ValueType
}

func NewMicroServiceExporter() *MicroServiceExporter {

	// 抓取k8s的相关配置
	k8s_config := Parse_k8s_config()
	fmt.Println("k8s_config: ", k8s_config.Cluster.Name)
	master0 := k8s_config.Cluster.Master[0]

	var k8sConfig config.K8sConfig = config.K8sConfig{
		ServiceURL:      fmt.Sprintf("http://%s:%s/api/v1/services", master0, k8s_config.Cluster.ApiServerPort),  //"http://124.65.131.14:38080/api/v1/services",
		EndpointURL:     fmt.Sprintf("http://%s:%s/api/v1/endpoints", master0, k8s_config.Cluster.ApiServerPort), //"http://124.65.131.14:38080/api/v1/endpoints",
		NodeURL:         fmt.Sprintf("http://%s:%s/api/v1/nodes", master0, k8s_config.Cluster.ApiServerPort),
		PodURL:          fmt.Sprintf("http://%s:%s/api/v1/pods", master0, k8s_config.Cluster.ApiServerPort),
		NodeResourceURL: fmt.Sprintf("http://%s:%s/apis/metrics.k8s.io/v1beta1/nodes", master0, k8s_config.Cluster.ApiServerPort),
	}

	myk8sNodeInfos := GetNodeInfo(k8sConfig.NodeURL)
	nodeDescs := make([]K8sNodeDesc, 0)
	nodeInfoDescs := make([]K8sNodeInfoDesc, 0)
	fmt.Println("len(myk8sNodeInfos): ", len(myk8sNodeInfos))

	// 抓取node资源使用情况
	nodeResourceUsedData := GetResourceUsed(k8sConfig.NodeResourceURL)
	for _, nodeInfo := range myk8sNodeInfos {
		// nodeInfo.Ip
		if _, ok := nodeResourceUsedData[nodeInfo.Ip]; ok {
			cpuUsed := nodeResourceUsedData[nodeInfo.Ip].Cpu
			if strings.Contains(*cpuUsed, "n") {
				if cpuUsedint, err := strconv.ParseUint((*cpuUsed)[:len(*cpuUsed)-1], 10, 64); err == nil {
					nodeInfo.CpuUsedN = cpuUsedint
				} else {
					fmt.Println("获取cpu使用数据失败!")
				}
			}
			memoryUsed := nodeResourceUsedData[nodeInfo.Ip].Memory
			if strings.Contains(*memoryUsed, "Ki") {
				if memoryUsedData, err := strconv.ParseUint((*memoryUsed)[:len(*memoryUsed)-2], 10, 64); err == nil {
					nodeInfo.MemoryUsedKB = memoryUsedData
				} else {
					fmt.Println("获取内存使用数据失败!")
				}
			}
		} else {
			fmt.Println("没有抓取到这个ip主机的IP资源使用情况")
		}
	}

	for i := 0; i < len(myk8sNodeInfos); i++ {
		var k8sNodeDesc K8sNodeDesc

		k8sNodeDesc.MaxCpuDesc = prometheus.NewDesc("max_cpu_total", "主机的cpu总数",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		k8sNodeDesc.MaxCpuValType = prometheus.GaugeValue

		k8sNodeDesc.MaxDiskStorageDesc = prometheus.NewDesc("max_disk_total", "主机的磁盘总大小",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		k8sNodeDesc.MaxDiskStorageValType = prometheus.GaugeValue

		k8sNodeDesc.MaxMemoryDesc = prometheus.NewDesc("max_memory_total", "主机的总内存大小",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		k8sNodeDesc.MaxMemoryValType = prometheus.GaugeValue

		k8sNodeDesc.MaxPodsDesc = prometheus.NewDesc("max_pod_total", "主机的总的pod数量",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		k8sNodeDesc.MaxPodsValType = prometheus.GaugeValue

		k8sNodeDesc.AllocateCpuDesc = prometheus.NewDesc("allocate_cpu_total", "分配的cpu总数",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		k8sNodeDesc.AllocateCpuValType = prometheus.GaugeValue

		k8sNodeDesc.AllocateDiskStorageDesc = prometheus.NewDesc("allocate_disk_total", "分配的磁盘总大小",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		k8sNodeDesc.AllocateDiskStorageValType = prometheus.GaugeValue

		k8sNodeDesc.AllocateMemoryDesc = prometheus.NewDesc("allocate_memory_total", "分配的总内存大小",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		k8sNodeDesc.AllocateMemoryValType = prometheus.GaugeValue

		k8sNodeDesc.AllocatePodsDesc = prometheus.NewDesc("allocate_pod_total", "分配的总的pod数量",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		k8sNodeDesc.AllocatePodsValType = prometheus.GaugeValue

		k8sNodeDesc.CpuUsedDesc = prometheus.NewDesc("cpu_used_n", "主机cpu使用量",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		k8sNodeDesc.CpuUsedValType = prometheus.GaugeValue

		k8sNodeDesc.MemoryUsedDesc = prometheus.NewDesc("memory_used_kb", "主机内存使用量",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		k8sNodeDesc.MemoryUsedValType = prometheus.GaugeValue

		nodeDescs = append(nodeDescs, k8sNodeDesc)

		var k8sNodeinfodesc K8sNodeInfoDesc
		k8sNodeinfodesc.NodeInfoDesc = prometheus.NewDesc("k8s_node_ready", "展示k8s每一个节点资源负载情况及是否ready",
			[]string{"cluster", "host", "ip", "memory_pressure_ok", "disk_pressure_ok", "pid_pressure_ok"},
			prometheus.Labels{})
		k8sNodeinfodesc.NodeInfoValType = prometheus.GaugeValue
		nodeInfoDescs = append(nodeInfoDescs, k8sNodeinfodesc)
	}

	serviceinfo := GetServiceInfo(k8sConfig.ServiceURL)
	fmt.Println(" len(serviceinfo):", len(serviceinfo))
	serviceinfoDescs := make([]K8sServiceDesc, 0)
	for i := 0; i < len(serviceinfo); i++ {
		var k8sServiceDesc K8sServiceDesc
		k8sServiceDesc.ServiceInfoDesc = prometheus.NewDesc("service_info", "显示k8s集群中所有的服务信息",
			[]string{"cluster", "service_name", "is_nodeport"},
			prometheus.Labels{})
		k8sServiceDesc.ServiceInfoValType = prometheus.GaugeValue
		serviceinfoDescs = append(serviceinfoDescs, k8sServiceDesc)
	}

	// 微服务上应用的pod状态
	// podInfoDescs []K8sPodDesc
	myk8spodinfo := GetPodInfo(k8sConfig.PodURL)
	podInfoDescs := make([]K8sPodDesc, 0)
	fmt.Println("len(myk8spodinfo): ", len(myk8spodinfo))
	for i := 0; i < len(myk8spodinfo); i++ {
		var k8spodDesc K8sPodDesc
		k8spodDesc.PodInfoDesc = prometheus.NewDesc("pod_info", "显示k8s集群节点的pod信息",
			[]string{"cluster", "pod_name", "app", "phase", "run_host_ip", "restart_count"},
			prometheus.Labels{})
		k8spodDesc.PodInfoValType = prometheus.GaugeValue
		podInfoDescs = append(podInfoDescs, k8spodDesc)
	}

	fmt.Println("", *k8s_config)
	fmt.Println("nodeDescs: ", nodeDescs)
	fmt.Println("nodeInfoDescs: ", nodeInfoDescs)
	fmt.Println("serviceinfoDescs: ", serviceinfoDescs)
	fmt.Println("podInfoDescs: ", podInfoDescs)
	fmt.Println("myk8sNodeInfos: ", myk8sNodeInfos)
	fmt.Println("serviceinfo: ", serviceinfo)
	fmt.Println("myk8spodinfo: ", myk8spodinfo)
	return &MicroServiceExporter{
		k8sConfig:        *k8s_config,
		nodeDescs:        nodeDescs,
		nodeInfoDescs:    nodeInfoDescs,
		serviceInfoDescs: serviceinfoDescs,
		podInfoDescs:     podInfoDescs,
		nodeDatas:        myk8sNodeInfos,
		serviceInfoDatas: serviceinfo,
		podInfoDatas:     myk8spodinfo,
	}

}

// Describe(chan<- *Desc)

// func (Collector).Describe(chan<- *Desc)
func (e *MicroServiceExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, nodeDesc := range e.nodeDescs {
		ch <- nodeDesc.MaxCpuDesc
		ch <- nodeDesc.MaxDiskStorageDesc
		ch <- nodeDesc.MaxMemoryDesc
		ch <- nodeDesc.MaxPodsDesc

		ch <- nodeDesc.AllocateCpuDesc
		ch <- nodeDesc.AllocateDiskStorageDesc
		ch <- nodeDesc.AllocateMemoryDesc
		ch <- nodeDesc.AllocatePodsDesc

		ch <- nodeDesc.CpuUsedDesc
		ch <- nodeDesc.MemoryUsedDesc
	}

	for _, nodeInfoDesc := range e.nodeInfoDescs {
		ch <- nodeInfoDesc.NodeInfoDesc
	}

	for _, serviceInfoDesc := range e.serviceInfoDescs {
		ch <- serviceInfoDesc.ServiceInfoDesc
	}

	for _, podInfoDesc := range e.podInfoDescs {
		ch <- podInfoDesc.PodInfoDesc
	}

}

//
func (e *MicroServiceExporter) Collect(ch chan<- prometheus.Metric) {
	// 基于抓取node, service, pod数据，输出指标
	k8sNodeInfo := e.nodeDatas
	for _, node_info := range k8sNodeInfo {
		fmt.Println(node_info.Name)
		fmt.Println(node_info.CreationTimestamp)
		fmt.Println(node_info.NodeCapacityS.cpuCores)
		fmt.Println(node_info.NodeCapacityS.diskStorage)
		fmt.Println(node_info.NodeCapacityS.memory)
		fmt.Println(node_info.NodeCapacityS.pods)
	}
	for idx, nodeDesc := range e.nodeDescs {
		ch <- prometheus.MustNewConstMetric(nodeDesc.MaxCpuDesc, nodeDesc.MaxCpuValType, float64(e.nodeDatas[idx].NodeCapacityS.cpuCores),
			e.k8sConfig.Cluster.Name, e.nodeDatas[idx].Ip, e.nodeDatas[idx].Ip)
		ch <- prometheus.MustNewConstMetric(nodeDesc.MaxDiskStorageDesc, nodeDesc.MaxDiskStorageValType, float64(e.nodeDatas[idx].NodeCapacityS.diskStorage),
			e.k8sConfig.Cluster.Name, e.nodeDatas[idx].Ip, e.nodeDatas[idx].Ip)
		ch <- prometheus.MustNewConstMetric(nodeDesc.MaxMemoryDesc, nodeDesc.MaxMemoryValType, float64(e.nodeDatas[idx].NodeCapacityS.memory),
			e.k8sConfig.Cluster.Name, e.nodeDatas[idx].Ip, e.nodeDatas[idx].Ip)
		ch <- prometheus.MustNewConstMetric(nodeDesc.MaxPodsDesc, nodeDesc.MaxPodsValType, float64(e.nodeDatas[idx].NodeCapacityS.pods),
			e.k8sConfig.Cluster.Name, e.nodeDatas[idx].Ip, e.nodeDatas[idx].Ip)

		ch <- prometheus.MustNewConstMetric(nodeDesc.AllocateCpuDesc, nodeDesc.AllocateCpuValType, float64(e.nodeDatas[idx].NodeAllocatableS.cpuCores),
			e.k8sConfig.Cluster.Name, e.nodeDatas[idx].Ip, e.nodeDatas[idx].Ip)
		ch <- prometheus.MustNewConstMetric(nodeDesc.AllocateDiskStorageDesc, nodeDesc.AllocateDiskStorageValType, float64(e.nodeDatas[idx].NodeAllocatableS.diskStorage),
			e.k8sConfig.Cluster.Name, e.nodeDatas[idx].Ip, e.nodeDatas[idx].Ip)
		ch <- prometheus.MustNewConstMetric(nodeDesc.AllocateMemoryDesc, nodeDesc.AllocateMemoryValType, float64(e.nodeDatas[idx].NodeAllocatableS.memory),
			e.k8sConfig.Cluster.Name, e.nodeDatas[idx].Ip, e.nodeDatas[idx].Ip)
		ch <- prometheus.MustNewConstMetric(nodeDesc.AllocatePodsDesc, nodeDesc.AllocatePodsValType, float64(e.nodeDatas[idx].NodeAllocatableS.pods),
			e.k8sConfig.Cluster.Name, e.nodeDatas[idx].Ip, e.nodeDatas[idx].Ip)

		ch <- prometheus.MustNewConstMetric(nodeDesc.CpuUsedDesc, nodeDesc.CpuUsedValType, float64(e.nodeDatas[idx].CpuUsedN),
			e.k8sConfig.Cluster.Name, e.nodeDatas[idx].Ip, e.nodeDatas[idx].Ip)
		ch <- prometheus.MustNewConstMetric(nodeDesc.MemoryUsedDesc, nodeDesc.MemoryUsedValType, float64(e.nodeDatas[idx].MemoryUsedKB),
			e.k8sConfig.Cluster.Name, e.nodeDatas[idx].Ip, e.nodeDatas[idx].Ip)

		// 写k8s_node_ready指标数据
		// "cluster", "host", "ip", "memory_pressure_ok", "disk_pressure_ok", "pid_pressure_ok"
		isReady, _ := strconv.ParseFloat(BoolToString(e.nodeDatas[idx].IsReady), 32)
		ch <- prometheus.MustNewConstMetric(e.nodeInfoDescs[idx].NodeInfoDesc, e.nodeInfoDescs[idx].NodeInfoValType, isReady,
			e.k8sConfig.Cluster.Name, e.nodeDatas[idx].Ip, e.nodeDatas[idx].Ip, BoolToString(e.nodeDatas[idx].MemoryPressure),
			BoolToString(e.nodeDatas[idx].DiskPressure), BoolToString(e.nodeDatas[idx].PidPressure))
	}

	// 写k8s service相关的指标
	keys := make([]string, 0, len(e.serviceInfoDatas))
	for k := range e.serviceInfoDatas {
		keys = append(keys, k)
		fmt.Println("serviceinfo key: ", k)
	}
	for idx, service_info := range e.serviceInfoDescs {
		// "cluster", "service_name", "is_nodeport"
		ch <- prometheus.MustNewConstMetric(service_info.ServiceInfoDesc, service_info.ServiceInfoValType, 1,
			e.k8sConfig.Cluster.Name, keys[idx], BoolToString(e.serviceInfoDatas[keys[idx]].IsNodePort))
	}

	// 写k8s pod相关的指标
	// "cluster", "pod_name", "app", "phase", "run_host_ip", "restart_count"
	for idx, pod_info := range e.podInfoDescs {
		var restartCount int
		for _, status := range e.podInfoDatas[idx].containersStatus {
			status.RestartCount += status.RestartCount
		}
		ch <- prometheus.MustNewConstMetric(pod_info.PodInfoDesc, pod_info.PodInfoValType, 1,
			e.k8sConfig.Cluster.Name, e.podInfoDatas[idx].Name, e.podInfoDatas[idx].App, e.podInfoDatas[idx].Status,
			e.podInfoDatas[idx].App, fmt.Sprintf("%d", restartCount))

	}
}

func BoolToString(boolV bool) string {
	if boolV == true {
		return "1"
	} else {
		return "0"
	}
}
