package main

import (
	"alive_exporter/config"
	"alive_exporter/hadoop"
	"alive_exporter/kafka"
	"alive_exporter/micro_service"
	"alive_exporter/zookeeper"
	"os"

	// "alive_exporter/utils"
	"fmt"
	"net/http"

	// "github.com/hazhashua/alive_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func comineServiceInfo() map[string]map[string]string {

	var k8sConfig config.K8sConfig = config.K8sConfig{
		ServiceURL:  "http://124.65.131.14:38080/api/v1/services",
		EndpointURL: "http://124.65.131.14:38080/api/v1/endpoints",
	}

	serviceinfo := micro_service.GetServiceInfo(k8sConfig.ServiceURL)
	endpointinfo := micro_service.GetEndpointInfo(k8sConfig.EndpointURL)

	var service_all_info map[string]map[string]string
	service_all_info = make(map[string]map[string]string)

	for key, _ := range serviceinfo {
		data := make(map[string]string)

		if value, ok := endpointinfo[key]; ok {
			data["ip"] = value.IP
		}
		if serviceinfo[key].IsNodePort == true {
			data["is_node_port"] = "true"
		} else {
			data["is_node_port"] = "false"
		}
		data["service_name"] = key
		data["port"] = fmt.Sprintf("%d", serviceinfo[key].Port)
		fmt.Println("port: ", serviceinfo[key].Port)
		service_all_info[key] = data
	}
	fmt.Println("service_all_info: ", service_all_info)
	return service_all_info

}

type SparkHandler struct {
	metrics []string
}

func (handler SparkHandler) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/spark/base/metrics":
		for _, value := range handler.metrics {
			fmt.Fprintf(writer, "%s", value)
		}
	default:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "no such page: %s\n", r.URL)
	}
}

func main() {

	// 抓取微服务的数据信息
	// comineServiceInfo()

	// 查询数据库数据
	// utils.Query("")

	// 服务存活collector
	// fmt.Println("*&&&&&&&&&&&&&&&&&&", utils.ValueQuery(""))
	// serviceCollector := newServiceAliveCollector()
	// // prometheus.MustRegister(serviceCollector)
	// r := prometheus.NewRegistry()
	// r.MustRegister(serviceCollector)
	// handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	hbaseCollector := newHbaseCollector()
	hbaseR := prometheus.NewRegistry()
	hbaseR.MustRegister(hbaseCollector)
	hbaseHandler := promhttp.HandlerFor(hbaseR, promhttp.HandlerOpts{})

	// QueryMetric()
	// http.Handle("/metrics", handler)
	http.Handle("/hbase/metrics", hbaseHandler)

	// // 数组传入所有的master和standby地址
	// url_array := []string{"http://124.65.131.14"}
	// // 查询spark的metric信息，默认为查询测试集群
	// print_metrics := spark.GetMetrics(url_array)
	// sparkHandler := SparkHandler{metrics: print_metrics}
	// http.Handle("/spark/base/metrics", sparkHandler)

	fmt.Println("命令行的参数有", len(os.Args))
	mode := "normal"
	produce := true
	// 遍历 os.Args 切片，就可以得到所有的命令行输入参数值
	for idx, value := range os.Args {
		fmt.Printf("args[%v]=%v\n", idx, value)
		if value == "kafka" {
			mode = "kafka"
		}
		if value == "produce" || value == "producer" {
			produce = true
		}
		if value == "consume" || value == "consumer" {
			produce = false
		}
	}
	if mode == "kafka" {
		if produce == true {
			kafka.AsyncProducer()
		} else {
			kafka.ConsumeTest()
		}

	} else {
		fmt.Println("kafka.Parse_kafka_config()..........................")
		// kafka.Parse_kafka_config()
		// kafka.GetKafkaMetrics()

		// go generateaAliveValue(serviceAliveCollector.channel)
		// go getAliveValueLoop(serviceAliveCollector.channel)

		//Create a new instance of the foocollector and
		//register it with the prometheus client.
		// foo := newFooCollector()
		// prometheus.MustRegister(foo)

		// go generateValue(foo.channel)
		// go getValueLoop(foo.channel)

		// ch := chan <- prometheus.Metric
		// foo.Collectx(make(chan<- prometheus.Metric), 100)

		//This section will start the HTTP server and expose
		//any metrics on the /metrics endpoint.

		// 带全部参数 注册句柄
		// serviceCollector := newServiceAliveCollector()
		// prometheus.MustRegister(serviceCollector)
		// http.Handle("/metrics", promhttp.Handler())

		// kafka_collector := kafka.NewKafkaCollector()
		// prometheus.MustRegister(kafka_collector)
		// http.Handle("/kafka/metrics", promhttp.Handler())

		// http://bigdata-dev01:8088/jmx?qry=Hadoop:service=ResourceManager,name=QueueMetrics,q0=root,q1=default

		hadoop_collector := hadoop.NewHadoopCollector()
		prometheus.MustRegister(hadoop_collector)
		http.Handle("/hadoop/metrics", promhttp.Handler())

		zookeeper.ZookeeperExporter()
		zookeeper.Watch()

		log.Info("Beginning to serve on port :38080")
		log.Fatal(http.ListenAndServe(":38080", nil))

		// time.Sleep(100)
		// kafka_collector = kafka.NewKafkaCollector()
		// prometheus.MustRegister(kafka_collector)
		// http.Handle("/kafka/metrics", promhttp.Handler())

	}

}
