package service_alive

import (
	"alive_exporter/micro_service"
	"alive_exporter/utils"
	"fmt"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
// type serviceAliveCollector struct {
// 	aliveMetric *prometheus.Desc
// 	// labelInfo   *prometheus.Labels
// 	channel chan []DatsourceAlive //uint64
// }

type serviceCollector struct {
	serviceAliveCollector []serviceAlive2Collector
}

type serviceAlive2Collector struct {
	aliveMetric *prometheus.Desc
	valType     prometheus.ValueType
}

type DatsourceAlive struct {
	ServiceName  *string
	ChildService *string
	ClusterName  *string
	IP           *string
	Port         int
	PortType     *string
	MetricValue  float32
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func NewServiceAliveCollector() *serviceCollector {
	var serviceAliveList []serviceAlive2Collector
	datasource_count := utils.ValueQuery("")
	for length := 0; length < datasource_count; length++ {
		var service_alive_collector serviceAlive2Collector
		service_alive_collector.aliveMetric = prometheus.NewDesc("alive_metric", "Show whether the ip:port is alive",
			[]string{"cluster", "service_name", "child_service", "ip", "port", "port_type"},
			prometheus.Labels{})
		service_alive_collector.valType = prometheus.GaugeValue
		serviceAliveList = append(serviceAliveList, service_alive_collector)

	}
	return &serviceCollector{serviceAliveCollector: serviceAliveList}

}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *serviceCollector) Describe(ch chan<- *prometheus.Desc) {
	//Update this section with the each metric you create for a given collector
	for _, metric := range collector.serviceAliveCollector {
		ch <- metric.aliveMetric
	}

}

// func generateaAliveValue(channel chan []DatsourceAlive) {
// 	// // var channel chan int = make(chan int)
// 	// var value uint64 = 1
// 	// for {
// 	// 	value += 1
// 	// 	channel <- DatsourceAlive{}

// 	// 	fmt.Println("put value:", value)
// 	// }

// 	var datasource_alives []DatsourceAlive
// 	for {
// 		da := GetPortInfos()
// 		for _, datasourceAlive := range da {
// 			datasource_alives = append(datasource_alives, datasourceAlive)
// 		}
// 		channel <- datasource_alives
// 	}
// }

// func getAliveValueLoop(channel chan []DatsourceAlive) {
// 	for {
// 		time.Sleep(1 * time.Second)
// 		<-channel
// 	}
// }

// func getAliveValue(channel chan []DatsourceAlive) (value []DatsourceAlive) {
// 	return <-channel
// }

//Collect implements required collect function for all promehteus collectors
func (collector *serviceCollector) Collect(ch chan<- prometheus.Metric) {
	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.

	da := GetPortInfos()
	// for _, alive := range da {
	for index, alive := range collector.serviceAliveCollector {
		ch <- prometheus.MustNewConstMetric(alive.aliveMetric, prometheus.GaugeValue, float64(da[index].MetricValue), *da[index].ClusterName, *da[index].ServiceName, *da[index].ChildService, *da[index].IP, fmt.Sprintf("%d", da[index].Port), *da[index].PortType)
		// break
	}

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	// ch <- prometheus.MustNewConstMetric(collector.aliveMetric, prometheus.CounterValue, metricValue, "cluster1", "hbase", "regionserver", "127.0.0.1", "10000", "tcp")
	// ch <- prometheus.MustNewConstMetric(collector.barMetric, prometheus.CounterValue, metricValue)

}

func getServices(bytes []byte) {
	// 获取k8s service的返回内容
	micro_service.UnmarshalAPIV1Services([]byte("hello go"))

}

func GetPortInfos() []DatsourceAlive {

	//执行数据序列化
	// utils.Serilize()

	//数据反序列化
	sp := utils.ReSerialize()

	dataSources := make([]DatsourceAlive, 0)
	for _, servicePort := range sp {
		// fmt.Println(string(*servicePort.IP) + string(servicePort.Port))
		// fmt.Sprintf("%s:%d", *servicePort.IP, servicePort.Port)
		var datasourceAlive DatsourceAlive
		datasourceAlive.ServiceName = servicePort.ServiceName
		datasourceAlive.ChildService = servicePort.ChildService
		datasourceAlive.ClusterName = servicePort.ClusterName
		datasourceAlive.IP = servicePort.IP
		datasourceAlive.Port = servicePort.Port
		datasourceAlive.PortType = servicePort.PortType
		fmt.Println("***: ", *servicePort.IP, servicePort.Port)
		datasourceAlive.MetricValue = float32(CheckPorts(fmt.Sprintf("%s:%d", *servicePort.IP, servicePort.Port), "tcp"))
		fmt.Println("datsourceAlive: ", datasourceAlive)

		dataSources = append(dataSources, datasourceAlive)
		// CheckPorts(fmt.Sprintf("%s:%d", *servicePort.IP, servicePort.Port), "tcp")
	}
	// CheckPorts("localhost:13306", "tcp")
	// CheckPorts("localhost:2379", "tcp")
	// CheckPorts("localhost:9115", "tcp")
	return dataSources
}

// 检测端口
func CheckPorts(ip_port string, port_type string) int {
	check := 0
	now := time.Now().Format("2006-01-02 15:04:05")
	// 检测端口
	conn, err := net.DialTimeout(port_type, ip_port, 1*time.Second)
	if err != nil {
		fmt.Println("err: ", err)
		fmt.Println("["+now+"]", ip_port, "端口未开启(fail)!")
	} else {
		if conn != nil {
			check = 1
			fmt.Println("["+now+"]", ip_port, "端口已开启(success)!")
			conn.Close()
		} else {
			fmt.Println("["+now+"]", ip_port, "端口未开启(fail)!")
		}
	}
	return check

}
