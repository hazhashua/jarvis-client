package micro_service

import (
	// utils "alive_exporter/utils"
	"fmt"
	"io/ioutil"
	"net/http"
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
	// fmt.Println("endpoint_data: ", endpoint_data)
	aPIV1Endpoints, _ := UnmarshalAPIV1Endpoints([]byte(endpoint_data))
	// fmt.Println(*aPIV1Endpoints.APIVersion, *aPIV1Endpoints.Kind, aPIV1Endpoints.Items, *aPIV1Endpoints.Metadata)

	var endpointInfoMap map[string]EndpointInfo
	endpointInfoMap = make(map[string]EndpointInfo)

	for _, data := range aPIV1Endpoints.Items {
		// fmt.Println("*data.Metadata.Name: ", *data.Metadata.Name)
		// 如果数据中subsets长度大于0
		if len(data.Subsets) > 0 {
			var endpointInfo EndpointInfo = EndpointInfo{
				EndpointName: *data.Metadata.SelfLink,
				ClusterIP:    *data.Subsets[0].Addresses[0].IP,
			}
			if data.Subsets[0].Addresses[0].NodeName != nil {
				endpointInfo.IP = *data.Subsets[0].Addresses[0].NodeName
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
