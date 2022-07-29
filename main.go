package main

import (
	"metric_exporter/config"
	"metric_exporter/hadoop"
	"metric_exporter/hive"
	"metric_exporter/kafka"
	"metric_exporter/micro_service"
	"metric_exporter/redis"
	"metric_exporter/service_alive"
	"metric_exporter/spark"
	"metric_exporter/utils"
	"metric_exporter/zookeeper"
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
	k8s_config := micro_service.Parse_k8s_config()
	fmt.Println("k8s_config: ", k8s_config.Cluster.Name)
	master0 := k8s_config.Cluster.Master[0]

	var k8sConfig config.K8sConfig = config.K8sConfig{
		ServiceURL:  fmt.Sprintf("http://%s:%s/api/v1/services", master0, k8s_config.Cluster.ApiServerPort),  //"http://124.65.131.14:38080/api/v1/services",
		EndpointURL: fmt.Sprintf("http://%s:%s/api/v1/endpoints", master0, k8s_config.Cluster.ApiServerPort), //"http://124.65.131.14:38080/api/v1/endpoints",
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
	handler.metrics = spark.GetMetrics()
	switch r.URL.Path {
	case "/spark/metrics":
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

	// 激活服务存活exporter
	fmt.Println("*&&&&&&&&&&&&&&&&&&", utils.ValueQuery(""))
	serviceCollector := service_alive.NewServiceAliveCollector()
	r := prometheus.NewRegistry()
	r.MustRegister(serviceCollector)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	http.Handle("/alive/metrics", handler)

	// 激活hbase exporter
	hbaseCollector := newHbaseCollector()
	hbaseR := prometheus.NewRegistry()
	hbaseR.MustRegister(hbaseCollector)
	hbaseHandler := promhttp.HandlerFor(hbaseR, promhttp.HandlerOpts{})
	http.Handle("/hbase/metrics", hbaseHandler)

	// 激活spark exporter
	// 数组传入所有的master和standby地址
	// 查询spark的metric信息，默认为查询测试集群
	print_metrics := spark.GetMetrics()
	sparkHandler := SparkHandler{metrics: print_metrics}
	http.Handle("/spark/metrics", sparkHandler)
	fmt.Println("命令行的参数有", len(os.Args))

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

	// 激活kafka exporter
	kafka_collector := kafka.NewKafkaCollector()
	kafka_r := prometheus.NewRegistry()
	kafka_r.MustRegister(kafka_collector)
	kafka_handler := promhttp.HandlerFor(kafka_r, promhttp.HandlerOpts{})
	http.Handle("/kafka/metrics", kafka_handler)

	// http://bigdata-dev01:8088/jmx?qry=Hadoop:service=ResourceManager,name=QueueMetrics,q0=root,q1=default

	// serviceCollector := micro_service.newServiceAliveCollector()
	// service_r := prometheus.NewRegistry()
	// r.MustRegister(serviceCollector)
	// handler := promhttp.HandlerFor(service_r, promhttp.HandlerOpts{})
	// http.Handle("/alive/metrics", handler)

	// 激活hadoop exporter
	hadoop_exporter := hadoop.NewHadoopCollector()
	hadoop_r := prometheus.NewRegistry()
	hadoop_r.MustRegister(hadoop_exporter)
	hadoop_handler := promhttp.HandlerFor(hadoop_r, promhttp.HandlerOpts{})
	http.Handle("/hadoop/metrics", hadoop_handler)

	// 激活redis exporter
	redis.RedisExporter()

	// 激活zookeeper exporter
	zookeeper.ZookeeperExporter()
	// zookeeper.Watch()

	hive_exporter := hive.NewHiveExporter()
	if hive_exporter == nil {
		fmt.Println("hive_exporter is nil")
	}
	fmt.Printf("hive_exporter: %v \n", hive_exporter)
	hive_r := prometheus.NewRegistry()
	fmt.Println("hive_exporter is nil ", hive_exporter == nil)
	hive_r.MustRegister(hive_exporter)
	hive_handler := promhttp.HandlerFor(hive_r, promhttp.HandlerOpts{})
	http.Handle("/hive/metrics", hive_handler)

	log.Info("Beginning to serve on port :38080")
	log.Fatal(http.ListenAndServe(":38080", nil))

	// time.Sleep(100)
	// kafka_collector = kafka.NewKafkaCollector()
	// prometheus.MustRegister(kafka_collector)
	// http.Handle("/kafka/metrics", promhttp.Handler())

}
