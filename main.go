package main

import (
	"flag"
	"fmt"
	"metric_exporter/config"
	"metric_exporter/hadoop"
	"metric_exporter/hbase"
	"metric_exporter/hive"
	"metric_exporter/kafka"
	"metric_exporter/micro_service"
	"metric_exporter/mysql"
	nodeexporter "metric_exporter/node_exporter"
	"metric_exporter/redis"
	"metric_exporter/service_alive"
	"metric_exporter/skywalking"
	"metric_exporter/spark"
	"metric_exporter/utils"
	_ "metric_exporter/utils"
	"metric_exporter/zookeeper"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// func comineServiceInfo() (map[string]map[string]string, []*micro_service.MyK8sNodeInfo) {
// 	k8s_config := micro_service.Parse_k8s_config()
// 	fmt.Println("k8s_config: ", k8s_config.Cluster.Name)
// 	master0 := k8s_config.Cluster.Master[0]

// 	var k8sConfig config.K8sConfig = config.K8sConfig{
// 		ServiceURL:  fmt.Sprintf("http://%s:%s/api/v1/services", master0, k8s_config.Cluster.ApiServerPort),  //"http://124.65.131.14:38080/api/v1/services",
// 		EndpointURL: fmt.Sprintf("http://%s:%s/api/v1/endpoints", master0, k8s_config.Cluster.ApiServerPort), //"http://124.65.131.14:38080/api/v1/endpoints",
// 		NodeURL:     fmt.Sprintf("http://%s:%s/api/v1/nodes", master0, k8s_config.Cluster.ApiServerPort),
// 	}

// 	serviceinfo := micro_service.GetServiceInfo(k8sConfig.ServiceURL)
// 	endpointinfo := micro_service.GetEndpointInfo(k8sConfig.EndpointURL)

// 	var service_all_info map[string]map[string]string
// 	service_all_info = make(map[string]map[string]string)

// 	for key, _ := range serviceinfo {
// 		data := make(map[string]string)

// 		if value, ok := endpointinfo[key]; ok {
// 			data["ip"] = value.IP
// 		}
// 		if serviceinfo[key].IsNodePort == true {
// 			data["is_node_port"] = "true"
// 		} else {
// 			data["is_node_port"] = "false"
// 		}
// 		data["service_name"] = key
// 		data["port"] = fmt.Sprintf("%d", serviceinfo[key].Port)
// 		fmt.Println("port: ", serviceinfo[key].Port)
// 		service_all_info[key] = data
// 	}
// 	fmt.Println("service_all_info: ", service_all_info)
// 	myK8SNodeInfo := micro_service.GetNodeInfo(k8sConfig.NodeURL)
// 	return service_all_info, myK8SNodeInfo
// }

// 向主程序发布要执行的采集模块
func publish(model string, ch chan interface{}) {
	ch <- model
}

// 读取发布的待采集模块，并实现启动
func subscribe(ch chan interface{}) {
	model := <-ch
	// 启动对应的exporter
	fmt.Printf("启动%s ...", model)
}

func parseArgs() string {

	modelPtr := flag.String("model", "all", "the model to export")
	flag.String("help", "", "please input the model below: \n\thadoop hbase hive kafka micro_service mysql node redis alive skywalking spark zookeeper \n use , split")
	// fmt.Printf("*modelPtr: %s\n", *modelPtr)
	// fmt.Printf("*helpPtr: %s\n", *helpPtr)
	// 解析命令行参数
	flag.Parse()
	return *modelPtr
}

// 暴露所有的服务指标数据
func exportAll() {
	// 激活微服务exporter
	microServiceExporter := micro_service.NewMicroServiceExporter()
	microServiceR := prometheus.NewRegistry()
	microServiceR.MustRegister(microServiceExporter)
	microServiceHandler := promhttp.HandlerFor(microServiceR, promhttp.HandlerOpts{})
	http.Handle(config.MICROSERVICE_METRICPATH, microServiceHandler)

	// 激活服务存活exporter
	serviceCollector := service_alive.NewServiceAliveCollector()
	r := prometheus.NewRegistry()
	r.MustRegister(serviceCollector)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	http.Handle(config.ALIVE_METRICPATH, handler)

	// 激活hbase exporter
	hbaseCollector := hbase.NewHbaseCollector()
	hbaseR := prometheus.NewRegistry()
	hbaseR.MustRegister(hbaseCollector)
	hbaseHandler := promhttp.HandlerFor(hbaseR, promhttp.HandlerOpts{})
	http.Handle(config.HBASE_METRICPATH, hbaseHandler)

	// 激活spark exporter
	// 数组传入所有的master和standby地址
	// 查询spark的metric信息，默认为查询测试集群
	print_metrics := spark.GetMetrics()
	sparkHandler := spark.SparkHandler{Metrics: print_metrics}
	http.Handle(config.SPARK_METRICPATH, sparkHandler)
	fmt.Println("命令行的参数有", len(os.Args))

	// 激活kafka exporter
	kafka_collector := kafka.NewKafkaCollector()
	kafka_r := prometheus.NewRegistry()
	kafka_r.MustRegister(kafka_collector)
	kafka_handler := promhttp.HandlerFor(kafka_r, promhttp.HandlerOpts{})
	http.Handle(config.KAFKA_METRICPATH, kafka_handler)

	// 激活hadoop exporter
	hadoop_exporter := hadoop.NewHadoopCollector()
	hadoop_r := prometheus.NewRegistry()
	hadoop_r.MustRegister(hadoop_exporter)
	hadoop_handler := promhttp.HandlerFor(hadoop_r, promhttp.HandlerOpts{})
	http.Handle(config.HADOOP_METRICPATH, hadoop_handler)

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
	http.Handle(config.HIVE_METRICPATH, hive_handler)

	// 激活mysql exporter
	mysql_exporter := mysql.NewMysqlExporter()
	mysql_r := prometheus.NewRegistry()
	mysql_r.MustRegister(mysql_exporter)
	mysql_handler := promhttp.HandlerFor(mysql_r, promhttp.HandlerOpts{})
	http.Handle(config.MYSQL_METRICPATH, mysql_handler)

	// 激活物理机指标采集脚本
	node_exporter := nodeexporter.NewNodeExporter()
	node_r := prometheus.NewRegistry()
	node_r.MustRegister(node_exporter)
	node_handler := promhttp.HandlerFor(node_r, promhttp.HandlerOpts{})
	http.Handle(config.NODE_METRICPATH, node_handler)

	// 激活skywalking exporter
	skywalking_exporter := skywalking.NewSkywalkingExporter()
	skywalking_r := prometheus.NewRegistry()
	skywalking_r.MustRegister(skywalking_exporter)
	skywalking_handler := promhttp.HandlerFor(skywalking_r, promhttp.HandlerOpts{})
	http.Handle(config.SKYWALKING_METRICPATH, skywalking_handler)
}

func main() {

	modelStart := make(map[string]bool, 0)

	modelV := parseArgs()
	if modelV == "all" {
		exportAll()
	} else {
		//只导出关心指标的数据
		models := strings.Split(modelV, ",")
		for _, model := range models {
			switch model {
			case "hadoop":
				if modelStart[model] == false {
					// 激活hadoop exporter
					hadoop_exporter := hadoop.NewHadoopCollector()
					hadoop_r := prometheus.NewRegistry()
					hadoop_r.MustRegister(hadoop_exporter)
					hadoop_handler := promhttp.HandlerFor(hadoop_r, promhttp.HandlerOpts{})
					http.Handle("/hadoop/metrics", hadoop_handler)
					modelStart[model] = true
				}
			case "hbase":
				if modelStart[model] == false {
					// 激活hbase exporter
					hbaseCollector := hbase.NewHbaseCollector()
					hbaseR := prometheus.NewRegistry()
					hbaseR.MustRegister(hbaseCollector)
					hbaseHandler := promhttp.HandlerFor(hbaseR, promhttp.HandlerOpts{})
					http.Handle("/hbase/metrics", hbaseHandler)
					modelStart[model] = true
				}

			case "hive":
				if modelStart[model] == false {
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
					modelStart[model] = true
				}
			case "kafka":
				if modelStart[model] == false {
					// 激活kafka exporter
					kafka_collector := kafka.NewKafkaCollector()
					kafka_r := prometheus.NewRegistry()
					kafka_r.MustRegister(kafka_collector)
					kafka_handler := promhttp.HandlerFor(kafka_r, promhttp.HandlerOpts{})
					http.Handle("/kafka/metrics", kafka_handler)
					modelStart[model] = true
				}
			case "micro_service":
				if modelStart[model] == false {
					// 激活微服务exporter
					microServiceExporter := micro_service.NewMicroServiceExporter()
					microServiceR := prometheus.NewRegistry()
					microServiceR.MustRegister(microServiceExporter)
					microServiceHandler := promhttp.HandlerFor(microServiceR, promhttp.HandlerOpts{})
					http.Handle("/micro_service/metrics", microServiceHandler)
					modelStart[model] = true
				}
			case "mysql":
				if modelStart[model] == false {
					// 激活mysql exporter
					mysql_exporter := mysql.NewMysqlExporter()
					mysql_r := prometheus.NewRegistry()
					mysql_r.MustRegister(mysql_exporter)
					mysql_handler := promhttp.HandlerFor(mysql_r, promhttp.HandlerOpts{})
					http.Handle("/mysql/metrics", mysql_handler)
					modelStart[model] = true
				}
			case "node":
				if modelStart[model] == false {
					// 激活物理机指标采集脚本
					node_exporter := nodeexporter.NewNodeExporter()
					node_r := prometheus.NewRegistry()
					node_r.MustRegister(node_exporter)
					node_handler := promhttp.HandlerFor(node_r, promhttp.HandlerOpts{})
					http.Handle("/node/metrics", node_handler)
					modelStart[model] = true
				}
			case "redis":
				if modelStart[model] == false {
					// 激活redis exporter
					redis.RedisExporter()
					modelStart[model] = true
				}
			case "alive":
				if modelStart[model] == false {
					// 激活服务存活exporter
					serviceCollector := service_alive.NewServiceAliveCollector()
					r := prometheus.NewRegistry()
					r.MustRegister(serviceCollector)
					handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
					http.Handle("/alive/metrics", handler)
					modelStart[model] = true
				}
			case "skywalking":
				if modelStart[model] == false {
					// 激活skywalking exporter
					skywalking_exporter := skywalking.NewSkywalkingExporter()
					skywalking_r := prometheus.NewRegistry()
					skywalking_r.MustRegister(skywalking_exporter)
					skywalking_handler := promhttp.HandlerFor(skywalking_r, promhttp.HandlerOpts{})
					http.Handle("/skywalking/metrics", skywalking_handler)
					modelStart[model] = true
				}
			case "spark":
				if modelStart[model] == false {
					// 激活spark exporter
					// 数组传入所有的master和standby地址
					// 查询spark的metric信息，默认为查询测试集群
					print_metrics := spark.GetMetrics()
					sparkHandler := spark.SparkHandler{Metrics: print_metrics}
					http.Handle("/spark/metrics", sparkHandler)
					fmt.Println("命令行的参数有", len(os.Args))
					modelStart[model] = true
				}
			case "zookeeper":
				if modelStart[model] == false {
					// 激活zookeeper exporter
					zookeeper.ZookeeperExporter()
					modelStart[model] = true
				}
			case "config":
				http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
					fmt.Println("执行配置函数")
					pyaml := utils.LoadYaml()
					configs := pyaml.ScrapeConfigs
					// 读取数据库配置
					dss := utils.PgDataStoreQuery(utils.Db)
					for _, ds := range dss {
						name := ds.DataName
						path := ds.Path
						ip := ds.Ip
						configs = append(configs, struct {
							JobName       string "yaml:\"job_name,omitempty\" mapstructure:\"job_name\""
							MetricsPath   string "yaml:\"metrics_path,omitempty\" mapstructure:\"metrics_path\""
							StaticConfigs []struct {
								Targets []string "yaml:\"targets,omitempty\""
							} "yaml:\"static_configs\" mapstructure:\"static_configs\""
						}{
							JobName:     name,
							MetricsPath: path,
							StaticConfigs: []struct {
								Targets []string "yaml:\"targets,omitempty\""
							}{
								{
									Targets: []string{ip},
								},
							},
						})

					}
					pyaml.Global.ScrapeInterval = "15s"
					pyaml.ScrapeConfigs = configs
					// 基于数据库配置数据,生成新的yaml文件
					yamlBytes := utils.GenerateYamlFile(pyaml, "./prometheus_auto.yml")
					w.Write(yamlBytes)
				})
			default:
				fmt.Println("unknown model...")
			}
		}
	}

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

	// http://bigdata-dev01:8088/jmx?qry=Hadoop:service=ResourceManager,name=QueueMetrics,q0=root,q1=default

	// escape := url.QueryEscape("redis_keyspace_hits_total/(redis_keyspace_misses_total+redis_keyspace_hits_total)")
	// urlstr := fmt.Sprintf("http://192.168.10.221:9090/api/v1/query?query=%s", escape)

	// if r2, err := http.Get(urlstr); err == nil {
	// 	var body []byte
	// 	body, err = ioutil.ReadAll(r2.Body)
	// 	fmt.Println("response: ", string(body))
	// } else {
	// 	fmt.Println("request error!")
	// }

	// utils.Migirate()

	// // 测试文件传输
	// var sc *utils.SftpClient
	// var err error
	// if sc, err = utils.NewSessionWithPassword("192.168.10.220", 22, "root", "pwd@123"); err != nil {
	// 	fmt.Println("ssh 创建连接失败! ")
	// }
	// if err = sc.ScopyRmoteFile("/root/collector", "collector_test"); err != nil {
	// 	fmt.Println("拷贝远程文件到本地, 失败!")
	// }

	log.Info("Beginning to serve on port :38080")
	log.Fatal(http.ListenAndServe(":38080", nil))

}
