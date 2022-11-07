package service_alive

import (
	"database/sql"
	"fmt"
	"metric_exporter/micro_service"
	"metric_exporter/utils"
	"net"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
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
	datas                 []DatsourceAlive
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
	Port         sql.NullInt64
	PortType     *string
	MetricValue  float32
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func NewServiceAliveCollector() *serviceCollector {
	var serviceAliveList []serviceAlive2Collector
	// db := utils.DbOpen(nil)
	var db *gorm.DB
	if utils.Db == nil {
		utils.Db = utils.DbOpen(nil)
	}
	db = utils.Db
	datasource_count := utils.PgCountQuery(db, "")

	fmt.Println("查询到的service_port表记录数: ", datasource_count)
	for length := 0; length < datasource_count; length++ {
		var service_alive_collector serviceAlive2Collector
		service_alive_collector.aliveMetric = prometheus.NewDesc("alive_metric", "Show whether the ip:port is alive",
			[]string{"cluster", "service_name", "child_service", "ip", "port", "port_type"},
			prometheus.Labels{})
		service_alive_collector.valType = prometheus.GaugeValue
		serviceAliveList = append(serviceAliveList, service_alive_collector)

	}

	da := GetAliveInfos()

	return &serviceCollector{
		serviceAliveCollector: serviceAliveList,
		datas:                 da,
	}

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
	collector = NewServiceAliveCollector()
	datas := collector.datas
	fmt.Println("collector datas: ", datas)

	// for _, alive := range da {
	for index, alive := range collector.serviceAliveCollector {
		if index >= len(datas) {
			//查询数据已经遍历完，退出
			break
		}
		var portValue string
		if datas[index].Port.Valid == true {
			portValue = fmt.Sprintf("%d", datas[index].Port.Int64)
		} else {
			portValue = ""
		}
		ch <- prometheus.MustNewConstMetric(alive.aliveMetric, prometheus.GaugeValue, float64(datas[index].MetricValue), *datas[index].ClusterName, *datas[index].ServiceName, *datas[index].ChildService, *datas[index].IP, portValue, *datas[index].PortType)
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

func GetAliveInfos() []DatsourceAlive {

	//执行数据序列化
	utils.Serilize()

	//数据反序列化
	sp := utils.ReSerialize()

	dataSources := make([]DatsourceAlive, 0)

	var localIp string
	netInfo := utils.NetInfoGet()
	for ethName, ip := range netInfo.EthInfo {
		if strings.Contains(ethName, "eth") || strings.Contains(ethName, "en") {
			utils.Logger.Printf("网络设备: %s  ip地址: %s\n", ethName, ip)
			localIp = ip
		}
	}
	for _, servicePort := range sp {
		var datasourceAlive DatsourceAlive
		datasourceAlive.ServiceName = servicePort.ServiceName
		datasourceAlive.ChildService = servicePort.ChildService
		datasourceAlive.ClusterName = &utils.DbConfig.Cluster.Name
		datasourceAlive.IP = servicePort.IP
		datasourceAlive.Port = servicePort.Port
		datasourceAlive.PortType = servicePort.PortType
		utils.Logger.Println("***: ", *servicePort.IP, servicePort.Port)
		if *servicePort.ServiceName == "micro_service" {
			// 如果是k8s服务，则使用进程探活
			// 如果不是本地的进程探测数据，则跳过
			if *servicePort.IP == localIp {
				utils.Logger.Printf("ip是本地地址...")
			} else {
				utils.Logger.Printf("待检测的ip地址: %s", *servicePort.IP)
				utils.Logger.Printf("本地ip地址: %s", localIp)
				// 测试时, 不跳过  开发或运行环境 跳过
				continue
			}
			utils.Logger.Printf("执行服务进程检测: %s ...\n", *servicePort.ChildService)
			alive := IsProcessRunning(*servicePort.ChildService)
			if alive == true {
				datasourceAlive.MetricValue = float32(1)
			} else {
				datasourceAlive.MetricValue = float32(0)
			}
		} else {

			if servicePort.Port.Valid == false {
				datasourceAlive.MetricValue = 0
			} else {
				if CheckPorts(fmt.Sprintf("%s:%d", *servicePort.IP, servicePort.Port.Int64), "tcp") {
					datasourceAlive.MetricValue = float32(1)
				} else {
					datasourceAlive.MetricValue = float32(0)
				}
			}
		}
		utils.Logger.Println("datsourceAlive: ", datasourceAlive)
		dataSources = append(dataSources, datasourceAlive)
		// CheckPorts(fmt.Sprintf("%s:%d", *servicePort.IP, servicePort.Port), "tcp")
	}
	// CheckPorts("localhost:13306", "tcp")
	// CheckPorts("localhost:2379", "tcp")
	// CheckPorts("localhost:9115", "tcp")
	return dataSources
}

// 检测端口
func CheckPorts(ip_port string, port_type string) bool {
	check := false
	now := time.Now().Format("2006-01-02 15:04:05")
	// 检测端口
	conn, err := net.DialTimeout(port_type, ip_port, 1*time.Second)
	if err != nil {
		fmt.Printf("检测%s超时, [%v], ip_port, 端口未开启(fail), error: %s\n", ip_port, now, err.Error())
	} else {
		if conn != nil {
			check = true
			utils.Logger.Println("["+now+"]", ip_port, "端口已开启(success)!")
			conn.Close()
		} else {
			utils.Logger.Println("["+now+"]", ip_port, "端口未开启(fail)!")
		}
	}
	return check

}
