package micro_service

import (
	// utils "alive_exporter/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// 服务的地址信息，包括服务的名称和端口类型及端口
type ServiceInfo struct {
	ServiceName string `json:"service_name,omitempty"`
	IsNodePort  bool   `json:"is_nodeport,omitempty"`
	IP          string `json:"ip,omitempty"`
	Port        int    `json:"port,omitempty"`
}

//服务的名称及服务的IP信息
type EndpointInfo struct {
	EndpointName string `json:"endpoint_name,omitempty"`
	ClusterIP    string `json:"cluster_ip,omitempty"`
	IP           string `json:"ip,omitempty"`
}

type NodeCapacity struct {
	cpuCores    float32
	diskStorage uint64
	memory      uint64
	pods        int64
}

type NodeAllocatable struct {
	cpuCores    float32
	diskStorage uint64
	memory      uint64
	pods        int64
}

type MyK8sNodeInfo struct {
	Name              string           `json:"name"`
	Ip                string           `json:"ip"`
	CreationTimestamp string           `json:"creationTimestamp"`
	NodeCapacityS     *NodeCapacity    `json:"nodeCapacity"`
	NodeAllocatableS  *NodeAllocatable `json:"nodeAllocatable"`
	MemoryPressure    bool             `json:"memoryPressure"`
	DiskPressure      bool             `json:"diskPressure"`
	PidPressure       bool             `json:"pidPressure"`
	IsReady           bool             `json:"isReady"`
}

type MyContainerStatus struct {
	ContainerName string `json:"containerName"`
	IsReady       bool   `json:"isReady"`
	IsStart       bool   `json:"isStart"`
	RestartCount  int    `json:"restartCount"`
}

type MyK8sPodInfo struct {
	Name              string               `json:"name"`
	CreationTimestamp string               `json:"creationTimestamp"`
	App               string               `json:"app"`
	Containers        int                  `json:"containers"`
	Status            string               `json:"status"`
	RunHostIP         string               `json:"hostip"`
	IsInitialized     bool                 `json:"isInitialized"`
	IsReady           bool                 `json:"isReady"`
	IsContainersReady bool                 `json:"isContainersReady"`
	IsPodScheduled    bool                 `json:"isPodScheduled"`
	containersStatus  []*MyContainerStatus `json:""`
}

// cluster:
//   name: 测试kubernetes集群
//   master:
//     - 192.168.10.20
//     - 192.168.10.21
//     - 192.168.10.22
//   nodes:
//     - 192.168.10.20
//     - 192.168.10.21
//     - 192.168.10.22
//     - 192.168.10.23
//     - 192.168.10.24
//     - 192.168.10.32
//     - 192.168.10.63
//     - 192.168.10.111

type K8sConfig struct {
	Cluster struct {
		Name          string   `yaml:"name"`
		Master        []string `yaml:"master"`
		Nodes         []string `yaml:"nodes"`
		ApiServerPort string   `yaml:"apiserverport"`
	}
}

func Parse_k8s_config() *K8sConfig {
	bytes, err := ioutil.ReadFile("./micro_service/config.yaml")
	if err != nil {
		fmt.Println("读文件出错！")
		return nil
	}
	k8sConfig := new(K8sConfig)
	err2 := yaml.Unmarshal(bytes, &k8sConfig)
	if err2 != nil {
		fmt.Println("Unmarshal失败！")
	}
	return k8sConfig
}

func Get(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	return string(body)
}

func GetServiceInfo(url string) map[string]ServiceInfo {
	/*
		基于k8sapi 获取所有service的名称及IP元信息
	*/
	//url 为apiservice service的路径地址，后面改成配置化
	data := Get(url)
	fmt.Println("&&&&&&&", data)
	aPIV1Services, _ := UnmarshalAPIV1Services([]byte(data))
	fmt.Println(*aPIV1Services.APIVersion, *aPIV1Services.Kind)
	// 存储所有serviceinfo信息
	var serviceInfoMap map[string]ServiceInfo
	serviceInfoMap = make(map[string]ServiceInfo)
	for _, item := range aPIV1Services.Items {
		// 取出所有service的名称
		// fmt.Println("service_name:", *(*item.Metadata).Name)
		var service_name = *(*item.Metadata).Name
		serviceInfo := ServiceInfo{ServiceName: service_name}
		//获得所有service的地址和端口
		for _, port := range (*item.Spec).Ports {
			if port.Name != nil {
				fmt.Println("port.name: ", *port.Name)
			}
			if port.NodePort != nil {
				serviceInfo.IsNodePort = true
				serviceInfo.Port = int(*port.NodePort)
			} else {
				serviceInfo.IsNodePort = false
			}
			serviceInfoMap[service_name] = serviceInfo
		}
	}
	fmt.Println("serviceinfoMap: ", serviceInfoMap)
	return serviceInfoMap
}

func GetEndpointInfo(url string) map[string]EndpointInfo {
	/**
	解析endpoint api内容
	*/
	endpoint_data := Get(url)
	fmt.Println("endpoint_data: ", endpoint_data)
	aPIV1Endpoints, _ := UnmarshalAPIV1Endpoints([]byte(endpoint_data))
	// fmt.Println(*aPIV1Endpoints.APIVersion, *aPIV1Endpoints.Kind, aPIV1Endpoints.Items, *aPIV1Endpoints.Metadata)
	var endpointInfoMap map[string]EndpointInfo
	endpointInfoMap = make(map[string]EndpointInfo)

	for _, data := range aPIV1Endpoints.Items {
		fmt.Println("*data.Metadata.Name: ", *data.Metadata.Name)
		// 如果数据中subsets长度大于0
		if len(data.Subsets) > 0 {
			var clusterIp string
			if len(data.Subsets[0].Addresses) > 0 {
				clusterIp = *data.Subsets[0].Addresses[0].IP
			} else {
				clusterIp = ""
			}
			var endpointInfo EndpointInfo = EndpointInfo{
				EndpointName: *data.Metadata.SelfLink,
				ClusterIP:    clusterIp,
			}
			if len(data.Subsets[0].Addresses) > 0 {
				if data.Subsets[0].Addresses[0].NodeName != nil {
					endpointInfo.IP = *data.Subsets[0].Addresses[0].NodeName
				}
			} else {
				endpointInfo.IP = ""
			}
			endpointInfoMap[*data.Metadata.Name] = endpointInfo
		} else {
			endpointInfoMap[*data.Metadata.Name] = EndpointInfo{
				EndpointName: *data.Metadata.SelfLink,
			}
		}
	}
	fmt.Println("endpointInfoMap:      ", endpointInfoMap)
	return endpointInfoMap
}

func GetNodeInfo(url string) []*MyK8sNodeInfo {
	/*
		基于k8sapi 获取所有所有node的节点信息
	*/
	node_data := Get(url)
	fmt.Println("node_data: ", node_data)
	k8sNodeInfo, _ := UnmarshalK8sNodeInfo([]byte(node_data))
	myNodeInfos := make([]*MyK8sNodeInfo, 0)
	for _, data := range k8sNodeInfo.Items {
		var myNodeInfo MyK8sNodeInfo
		fmt.Println("*k8sNodeInfo: ", *data.Metadata.Name)
		// 获得node的主机信息
		myNodeInfo.Name = *data.Metadata.Name
		myNodeInfo.CreationTimestamp = *data.Metadata.CreationTimestamp
		if len(data.Status.Addresses) > 0 {
			for _, address := range data.Status.Addresses {
				if *address.Type == InternalIP {
					myNodeInfo.Ip = *address.Address
					break
				}
			}
		}
		var nodeCapacity NodeCapacity
		if data.Status.Capacity != nil {
			cpu := *data.Status.Capacity.CPU
			if strings.Contains(cpu, "m") {
				cpucores, _ := strconv.ParseFloat(cpu[:len(cpu)-1], 32)
				nodeCapacity.cpuCores = float32(cpucores) / 1000
			} else {
				cpucores, _ := strconv.ParseFloat(cpu, 32)
				nodeCapacity.cpuCores = float32(cpucores)
			}

			if strings.Contains(*data.Status.Capacity.EphemeralStorage, "Ki") {
				storageK := (*data.Status.Capacity.EphemeralStorage)[:len(*data.Status.Capacity.EphemeralStorage)-2]
				storageKv, _ := strconv.ParseUint(storageK, 10, 32)
				nodeCapacity.diskStorage = storageKv * 1024
			} else {
				storagev, _ := strconv.ParseUint(*data.Status.Capacity.EphemeralStorage, 10, 32)
				nodeCapacity.diskStorage = storagev
			}

			if strings.Contains(*data.Status.Capacity.Memory, "Ki") {
				memoryK := (*data.Status.Capacity.Memory)[:len(*data.Status.Capacity.Memory)-2]
				memoryKv, _ := strconv.ParseUint(memoryK, 10, 32)
				nodeCapacity.memory = memoryKv * 1024
			}
			nodeCapacity.pods, _ = strconv.ParseInt(*data.Status.Capacity.Pods, 10, 32)
		}
		myNodeInfo.NodeCapacityS = &nodeCapacity

		var allocatable NodeAllocatable
		if data.Status.Allocatable != nil {
			cpu := *data.Status.Allocatable.CPU
			if strings.Contains(cpu, "m") {
				cpucores, _ := strconv.ParseFloat(cpu[:len(cpu)-1], 32)
				allocatable.cpuCores = float32(cpucores) / 1000
			} else {
				cpucores, _ := strconv.ParseFloat(cpu, 32)
				allocatable.cpuCores = float32(cpucores)
			}

			if strings.Contains(*data.Status.Allocatable.EphemeralStorage, "Ki") {
				storageK := (*data.Status.Allocatable.EphemeralStorage)[:len(*data.Status.Allocatable.EphemeralStorage)-2]
				storageKv, _ := strconv.ParseUint(storageK, 10, 32)
				allocatable.diskStorage = storageKv * 1024
			} else {
				storagev, _ := strconv.ParseUint(*data.Status.Allocatable.EphemeralStorage, 10, 32)
				allocatable.diskStorage = storagev
			}

			if strings.Contains(*data.Status.Allocatable.Memory, "Ki") {
				memoryK := (*data.Status.Allocatable.Memory)[:len(*data.Status.Allocatable.Memory)-2]
				memoryKv, _ := strconv.ParseUint(memoryK, 10, 32)
				allocatable.memory = memoryKv * 1024
			}
			allocatable.pods, _ = strconv.ParseInt(*data.Status.Allocatable.Pods, 10, 32)
		}
		myNodeInfo.NodeAllocatableS = &allocatable

		if len(data.Status.Conditions) > 0 {
			for _, item := range data.Status.Conditions {
				switch *item.Type {
				case MemoryPressure:
					if *item.Status == False {
						myNodeInfo.MemoryPressure = false
					} else {
						myNodeInfo.MemoryPressure = true
					}
				case DiskPressure:
					if *item.Status == False {
						myNodeInfo.DiskPressure = false
					} else {
						myNodeInfo.DiskPressure = true
					}
				case PIDPressure:
					if *item.Status == False {
						myNodeInfo.PidPressure = false
					} else {
						myNodeInfo.PidPressure = true
					}
				case Ready:
					if *item.Status == True {
						myNodeInfo.IsReady = true
					} else {
						myNodeInfo.IsReady = false
					}
				default:
					break
				}
			}
		}
		myNodeInfos = append(myNodeInfos, &myNodeInfo)
	}
	return myNodeInfos
}

func GetPodInfo(podUrl string) []*MyK8sPodInfo {
	// 基于K8S api获取所有pod的信息
	pod_data := Get(podUrl)
	fmt.Println("pod data: ", pod_data)
	k8sPodInfo, _ := UnmarshalPodInfo([]byte(pod_data))
	myPodInfos := make([]*MyK8sPodInfo, 0)

	// myPodInfo.Name
	// myPodInfo.App
	// myPodInfo.CreationTimestamp
	// myPodInfo.Containers
	// myPodInfo.Status
	// myPodInfo.IsInitialized
	// myPodInfo.IsReady
	// myPodInfo.IsContainersReady
	// myPodInfo.IsPodScheduled

	// 解析api数据并写入MyK8sPodInfo对象中
	for _, data := range k8sPodInfo.Items {
		var myPodInfo MyK8sPodInfo
		myPodInfo.Name = *data.Metadata.Name
		myPodInfo.CreationTimestamp = *data.Metadata.CreationTimestamp
		fmt.Println("", data.Metadata.Name)
		fmt.Println("", data.Metadata.CreationTimestamp)
		if data.Metadata.Labels != nil {
			myPodInfo.App = *data.Metadata.Labels.App
		}
		if data.Spec != nil {
			myPodInfo.Containers = len(data.Spec.Containers)
		}
		if data.Status != nil {
			myPodInfo.Status = string(*data.Status.Phase)
			myPodInfo.RunHostIP = string(*data.Status.HostIP)
		}
		if data.Status.Conditions != nil {
			for _, condition := range data.Status.Conditions {
				switch *condition.Type {
				case "Initialized":
					if *condition.Status == "True" {
						myPodInfo.IsInitialized = true
					} else {
						myPodInfo.IsInitialized = false
					}
				case "Ready":
					if *condition.Status == "True" {
						myPodInfo.IsReady = true
					} else {
						myPodInfo.IsReady = false
					}
				case "ContainersReady":
					if *condition.Status == "True" {
						myPodInfo.IsContainersReady = true
					} else {
						myPodInfo.IsContainersReady = false
					}
				case "PodScheduled":
					if *condition.Status == "True" {
						myPodInfo.IsPodScheduled = true
					} else {
						myPodInfo.IsPodScheduled = false
					}
				default:
					fmt.Println("unknown condition.type: ", *condition.Type)
				}
			}
		}
		if data.Status.ContainerStatuses != nil && len(data.Status.ContainerStatuses) > 0 {
			// 抓取pod下各个container的数据
			containerStatusList := make([]*MyContainerStatus, 0)
			for _, containerStatus := range data.Status.ContainerStatuses {
				var myContainerStatus MyContainerStatus
				myContainerStatus.ContainerName = *containerStatus.Name
				myContainerStatus.IsReady = *containerStatus.Ready
				myContainerStatus.IsStart = *containerStatus.Started
				myContainerStatus.RestartCount = int(*containerStatus.RestartCount)
				containerStatusList = append(containerStatusList, &myContainerStatus)
			}
			myPodInfo.containersStatus = containerStatusList
		}

		myPodInfos = append(myPodInfos, &myPodInfo)
	}
	return myPodInfos
}
