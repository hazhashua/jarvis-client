package main

import (
	"database/sql"
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
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/minms/shutdown"
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

func parseArgs() (string, string) {

	modelPtr := flag.String("model", "all", "the model to export")
	modelExcudePtr := flag.String("exclude_model", "", "the model unexport")
	flag.String("help", "", "please input the model below: \n\thadoop hbase hive kafka micro_service mysql node redis alive skywalking spark zookeeper \n use , split")
	// fmt.Printf("*modelPtr: %s\n", *modelPtr)
	// fmt.Printf("*helpPtr: %s\n", *helpPtr)
	// 解析命令行参数
	flag.Parse()
	return *modelPtr, *modelExcudePtr
}

// 注册 exporter地址数据 到数据库
func registerEndpoint(dataName string, port int, metricPath string) {
	ds := new(utils.Data_store_configure)
	ds.DataName = dataName
	ni := utils.NetInfoGet()
	ds.Ip = fmt.Sprintf("%s:%d", ni.Ip, utils.DbConfig.Cluster.HttpPort)
	ds.CreateTime = time.Now()
	ds.UpdateTime = time.Now()
	ds.Path = fmt.Sprintf("http://%s%s", ds.Ip, config.MetricPathMap[dataName])
	serviceDataType := make(map[string]int, 0)
	if dataName == "spark" || dataName == "alive" || dataName == "mysql" || dataName == "zookeeper" || dataName == "hadoop" || dataName == "hbase" || dataName == "hive" || dataName == "kafka" || dataName == "redis" {
		serviceDataType[dataName] = 3
	} else if dataName == "apisix" {
		serviceDataType[dataName] = 2
	} else if dataName == "skywalking" {
		serviceDataType[dataName] = 4
	} else if dataName == "micro_service" {
		serviceDataType[dataName] = 5
	} else if dataName == "node_exporter" {
		serviceDataType[dataName] = 1
	}
	ds.DataType = fmt.Sprintf("%d", serviceDataType[dataName])

	var dss []utils.Data_store_configure
	if dataName == "apisix" {
		// 如果是网关，只更新一次
		utils.Db.Where("data_name=?", dataName).Find(&dss)
		if len(dss) == 0 {
			utils.PgDataStoreInsert(utils.Db, ds)
		} else {
			utils.Logger.Printf("网关配置后,不再更新！")
		}
		// apisix endpoint信息退出不清理
		//config.MetricIpMap[strings.ToUpper(dataName)] = ds.Ip
		return
	}

	// var dss []utils.Data_store_configure
	utils.Db.Where("data_name=?", dataName).Where("ip=?", ds.Ip).Find(&dss)
	if len(dss) == 0 {
		// 数据库没有数据插入
		var id []sql.NullInt32
		utils.Db.Raw(fmt.Sprintf("select max(id) as id from public.%s", utils.DbConfig.Cluster.Postgres.ExportTable)).Pluck("id", &id)
		if id[0].Valid {
			ds.Id = int(id[0].Int32) + 1
		} else {
			ds.Id = 1
		}
		utils.PgDataStoreInsert(utils.Db, ds)
	} else {
		// 数据库有数据执行更新
		utils.Db.Select("id").Where("data_name=?", dataName).Where("ip=?", ds.Ip).Take(&dss)
		ds.Id = dss[0].Id
		utils.Logger.Printf("更新数据: %v\n", ds)
		utils.Db.Save(ds)
	}
	// 存入数据库后 更新exporter的配置信息，方便异常退出时清理
	config.MetricIpMap[strings.ToUpper(dataName)] = ds.Ip

}

// 注册config endpoint
func registerConfigEndpoint() {
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		pyaml := utils.LoadYaml()
		configs := pyaml.ScrapeConfigs
		// 读取数据库配置
		dss := utils.PgDataStoreQuery(utils.Db, utils.DbConfig.Cluster.Postgres.ExportTable)
		sort.Slice(dss, func(i, j int) bool {
			if dss[i].DataName >= dss[j].DataName {
				return true
			} else {
				return false
			}
		})
		var name string
		num := 0
		for _, ds := range dss {
			// 如果有相同job_name, 则追加数字后缀
			if ds.DataName == name && name != "" {
				num += 1
				name = fmt.Sprintf("%s_%d", ds.DataName, num)
			} else {
				name = ds.DataName
				num = 0
			}
			var ip, path string
			reg, err := regexp.Compile("http://([^/]*)(.*)")
			if err == nil {
				pathInfo := reg.FindStringSubmatch(ds.Path)
				if len(pathInfo) > 2 {
					path = pathInfo[2]
					ip = pathInfo[1]
					fmt.Printf("path: %s  ip: %s", path, ip)
				} else {
					fmt.Println("解析路径错误! ", ds.Path, pathInfo)
				}
			}
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
		pyaml.Global.EvaluationInterval = "30s"
		pyaml.ScrapeConfigs = configs
		// 基于数据库配置数据,生成新的yaml文件
		yamlBytes := utils.GenerateYamlFile(pyaml, "./prometheus_auto.yml")
		w.Write(yamlBytes)
	})
}

// 暴露所有的服务指标数据
func exportAll(allModels map[string]string) {
	if _, ok := allModels[config.MICROSERVICE]; ok {
		// 存在micro_service, 注册微服务endpoint
		// 激活微服务exporter
		microServiceExporter := micro_service.NewMicroServiceExporter()
		microServiceR := prometheus.NewRegistry()
		microServiceR.MustRegister(microServiceExporter)
		microServiceHandler := promhttp.HandlerFor(microServiceR, promhttp.HandlerOpts{})
		http.Handle(config.MICROSERVICE_METRICPATH, microServiceHandler)
	}

	if _, ok := allModels[config.ALIVE]; ok {
		// 激活服务存活exporter
		serviceCollector := service_alive.NewServiceAliveCollector()
		r := prometheus.NewRegistry()
		r.MustRegister(serviceCollector)
		handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
		http.Handle(config.ALIVE_METRICPATH, handler)
	}

	if _, ok := allModels[config.HBASE]; ok {
		// 激活hbase exporter
		hbaseCollector := hbase.NewHbaseCollector()
		hbaseR := prometheus.NewRegistry()
		hbaseR.MustRegister(hbaseCollector)
		hbaseHandler := promhttp.HandlerFor(hbaseR, promhttp.HandlerOpts{})
		http.Handle(config.HBASE_METRICPATH, hbaseHandler)
	}

	if _, ok := allModels[config.SPARK]; ok {
		// 激活spark exporter
		// 数组传入所有的master和standby地址
		// 查询spark的metric信息，默认为查询测试集群
		print_metrics := spark.GetMetrics()
		sparkHandler := spark.SparkHandler{Metrics: print_metrics}
		http.Handle(config.SPARK_METRICPATH, sparkHandler)
		fmt.Println("命令行的参数有", len(os.Args))
	}

	if _, ok := allModels[config.KAFKA]; ok {
		// 激活kafka exporter
		kafka_collector := kafka.NewKafkaCollector()
		kafka_r := prometheus.NewRegistry()
		kafka_r.MustRegister(kafka_collector)
		kafka_handler := promhttp.HandlerFor(kafka_r, promhttp.HandlerOpts{})
		http.Handle(config.KAFKA_METRICPATH, kafka_handler)
	}

	if _, ok := allModels[config.HADOOP]; ok {
		// 激活hadoop exporter
		hadoop_exporter := hadoop.NewHadoopCollector()
		hadoop_r := prometheus.NewRegistry()
		hadoop_r.MustRegister(hadoop_exporter)
		hadoop_handler := promhttp.HandlerFor(hadoop_r, promhttp.HandlerOpts{})
		http.Handle(config.HADOOP_METRICPATH, hadoop_handler)

	}

	if _, ok := allModels[config.REDIS]; ok {
		// 激活redis exporter
		redis.RedisExporter()
	}

	if _, ok := allModels[config.ZOOKEEPER]; ok {
		// 激活zookeeper exporter
		zookeeper.ZookeeperExporter()
		// zookeeper.Watch()
	}

	if _, ok := allModels[config.HIVE]; ok {
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
	}

	if _, ok := allModels[config.MYSQL]; ok {
		// 激活mysql exporter
		mysql_exporter := mysql.NewMysqlExporter()
		mysql_r := prometheus.NewRegistry()
		mysql_r.MustRegister(mysql_exporter)
		mysql_handler := promhttp.HandlerFor(mysql_r, promhttp.HandlerOpts{})
		http.Handle(config.MYSQL_METRICPATH, mysql_handler)
	}

	if _, ok := allModels[config.NODE]; ok {
		// 激活物理机指标采集脚本
		node_exporter := nodeexporter.NewNodeExporter()
		node_r := prometheus.NewRegistry()
		node_r.MustRegister(node_exporter)
		node_handler := promhttp.HandlerFor(node_r, promhttp.HandlerOpts{})
		http.Handle(config.NODE_METRICPATH, node_handler)
	}

	if _, ok := allModels[config.SKYWALKING]; ok {
		// 激活skywalking exporter
		skywalking_exporter := skywalking.NewSkywalkingExporter()
		skywalking_r := prometheus.NewRegistry()
		skywalking_r.MustRegister(skywalking_exporter)
		skywalking_handler := promhttp.HandlerFor(skywalking_r, promhttp.HandlerOpts{})
		http.Handle(config.SKYWALKING_METRICPATH, skywalking_handler)
	}

	// 默认注册config endpoint
	registerConfigEndpoint()
}

func main() {

	modelStart := make(map[string]bool, 0)

	modelV, excludeModelV := parseArgs()

	if excludeModelV == "all" {
		utils.Logger.Printf("所有的模块exporter都不启动,程序退出！\n")
		return
	} else if excludeModelV != "" {
		utils.Logger.Printf("启动排除之外的所有exporter！\n")
		// 启动exclude指定之外的所有exporter
		excludeModels := strings.Split(excludeModelV, ",")
		models := make(map[string]string)

		for model, _ := range config.MetricPathMap {
			models[model] = ""
		}
		for _, model := range excludeModels {
			// 删除exclude的元素
			delete(models, model)
		}
		utils.Logger.Printf("export models : %v", models)
		exportAll(models)
		for dataName, path := range config.MetricPathMap {
			if _, ok := models[dataName]; ok {
				utils.Logger.Printf("注册endpoint数据: %s %s\n", dataName, path)
				registerEndpoint(dataName, utils.DbConfig.Cluster.HttpPort, path)
			}
		}
	}

	if modelV == "all" && excludeModelV == "" {
		utils.Logger.Printf("启动全部exporter!\n")
		exportAll(config.MetricPathMap)
		for dataName, path := range config.MetricPathMap {
			registerEndpoint(dataName, utils.DbConfig.Cluster.HttpPort, path)
		}
	}

	if modelV != "" && excludeModelV == "" {
		utils.Logger.Printf("启动指定设置的exporter！\n")
		//只导出关心指标的数据
		models := strings.Split(modelV, ",")
		for _, model := range models {
			switch model {
			case config.HADOOP:
				if modelStart[model] == false {
					// 激活hadoop exporter
					hadoop_exporter := hadoop.NewHadoopCollector()
					hadoop_r := prometheus.NewRegistry()
					hadoop_r.MustRegister(hadoop_exporter)
					hadoop_handler := promhttp.HandlerFor(hadoop_r, promhttp.HandlerOpts{})
					http.Handle(config.HADOOP_METRICPATH, hadoop_handler)
					modelStart[model] = true
					registerEndpoint(config.HADOOP, utils.DbConfig.Cluster.HttpPort, config.MetricPathMap[config.HADOOP])
				}
			case config.HBASE:
				if modelStart[model] == false {
					// 激活hbase exporter
					hbaseCollector := hbase.NewHbaseCollector()
					hbaseR := prometheus.NewRegistry()
					hbaseR.MustRegister(hbaseCollector)
					hbaseHandler := promhttp.HandlerFor(hbaseR, promhttp.HandlerOpts{})
					http.Handle(config.HBASE_METRICPATH, hbaseHandler)
					modelStart[model] = true
					registerEndpoint(config.HBASE, utils.DbConfig.Cluster.HttpPort, config.MetricPathMap[config.HBASE])

				}

			case config.HIVE:
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
					http.Handle(config.HIVE_METRICPATH, hive_handler)
					modelStart[model] = true
					registerEndpoint(config.HIVE, utils.DbConfig.Cluster.HttpPort, config.MetricPathMap[config.HIVE])

				}
			case config.KAFKA:
				if modelStart[model] == false {
					// 激活kafka exporter
					kafka_collector := kafka.NewKafkaCollector()
					kafka_r := prometheus.NewRegistry()
					kafka_r.MustRegister(kafka_collector)
					kafka_handler := promhttp.HandlerFor(kafka_r, promhttp.HandlerOpts{})
					http.Handle(config.KAFKA_METRICPATH, kafka_handler)
					modelStart[model] = true
					registerEndpoint(config.KAFKA, utils.DbConfig.Cluster.HttpPort, config.MetricPathMap[config.KAFKA])

				}
			case config.MICROSERVICE:
				if modelStart[model] == false {
					// 激活微服务exporter
					microServiceExporter := micro_service.NewMicroServiceExporter()
					microServiceR := prometheus.NewRegistry()
					microServiceR.MustRegister(microServiceExporter)
					microServiceHandler := promhttp.HandlerFor(microServiceR, promhttp.HandlerOpts{})
					http.Handle(config.MICROSERVICE_METRICPATH, microServiceHandler)
					modelStart[model] = true
					registerEndpoint(config.MICROSERVICE, utils.DbConfig.Cluster.HttpPort, config.MetricPathMap[config.MICROSERVICE])

				}
			case config.MYSQL:
				if modelStart[model] == false {
					// 激活mysql exporter
					mysql_exporter := mysql.NewMysqlExporter()
					mysql_r := prometheus.NewRegistry()
					mysql_r.MustRegister(mysql_exporter)
					mysql_handler := promhttp.HandlerFor(mysql_r, promhttp.HandlerOpts{})
					http.Handle(config.MYSQL_METRICPATH, mysql_handler)
					modelStart[model] = true
					registerEndpoint(config.MYSQL, utils.DbConfig.Cluster.HttpPort, config.MetricPathMap[config.MYSQL])

				}
			case config.NODE:
				if modelStart[model] == false {
					// 激活物理机指标采集脚本
					node_exporter := nodeexporter.NewNodeExporter()
					node_r := prometheus.NewRegistry()
					node_r.MustRegister(node_exporter)
					node_handler := promhttp.HandlerFor(node_r, promhttp.HandlerOpts{})
					http.Handle(config.NODE_METRICPATH, node_handler)
					modelStart[model] = true
					registerEndpoint(config.NODE, utils.DbConfig.Cluster.HttpPort, config.MetricPathMap[config.NODE])

				}
			case config.REDIS:
				if modelStart[model] == false {
					// 激活redis exporter
					redis.RedisExporter()
					modelStart[model] = true
					registerEndpoint(config.REDIS, utils.DbConfig.Cluster.HttpPort, config.MetricPathMap[config.REDIS])

				}
			case config.ALIVE:
				if modelStart[model] == false {
					// 激活服务存活exporter
					serviceCollector := service_alive.NewServiceAliveCollector()
					r := prometheus.NewRegistry()
					r.MustRegister(serviceCollector)
					handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
					http.Handle(config.ALIVE_METRICPATH, handler)
					modelStart[model] = true
					utils.Logger.Printf("注册alive endpoint到数据库！")
					registerEndpoint(config.ALIVE, utils.DbConfig.Cluster.HttpPort, config.MetricPathMap[config.ALIVE])

				}
			case config.SKYWALKING:
				if modelStart[model] == false {
					// 激活skywalking exporter
					skywalking_exporter := skywalking.NewSkywalkingExporter()
					skywalking_r := prometheus.NewRegistry()
					skywalking_r.MustRegister(skywalking_exporter)
					skywalking_handler := promhttp.HandlerFor(skywalking_r, promhttp.HandlerOpts{})
					http.Handle(config.SKYWALKING_METRICPATH, skywalking_handler)
					modelStart[model] = true
					registerEndpoint(config.SKYWALKING, utils.DbConfig.Cluster.HttpPort, config.MetricPathMap[config.SKYWALKING])

				}
			case config.SPARK:
				if modelStart[model] == false {
					// 激活spark exporter
					// 数组传入所有的master和standby地址
					// 查询spark的metric信息，默认为查询测试集群
					print_metrics := spark.GetMetrics()
					sparkHandler := spark.SparkHandler{Metrics: print_metrics}
					http.Handle(config.SPARK_METRICPATH, sparkHandler)
					fmt.Println("命令行的参数有", len(os.Args))
					modelStart[model] = true
					registerEndpoint(config.SPARK, utils.DbConfig.Cluster.HttpPort, config.MetricPathMap[config.SPARK])

				}
			case config.ZOOKEEPER:
				if modelStart[model] == false {
					// 激活zookeeper exporter
					zookeeper.ZookeeperExporter()
					modelStart[model] = true
					registerEndpoint(config.ZOOKEEPER, utils.DbConfig.Cluster.HttpPort, config.MetricPathMap[config.ZOOKEEPER])

				}
			case "config":
				// 注册config endpoint
				registerConfigEndpoint()
			default:
				fmt.Println("unknown model...")
			}
		}
	}

	//默认写表中apisix配置信息
	registerEndpoint("apisix", utils.DbConfig.Cluster.HttpPort, "")

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

	go graceExit()

	log.Info(fmt.Sprintf("Beginning to serve on port :%d", utils.DbConfig.Cluster.HttpPort))
	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", utils.DbConfig.Cluster.HttpPort), nil))
	http.ListenAndServe(fmt.Sprintf(":%d", utils.DbConfig.Cluster.HttpPort), nil)

	// 	// handler
	//     handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//         time.Sleep(2 * time.Second)
	//         fmt.Fprintln(w, "hello")
	//     })

	//     // server
	//     srv := http.Server{
	//         Addr:    *addr,
	//         Handler: handler,
	//     }

	//     // make sure idle connections returned
	//     processed := make(chan struct{})
	//     go func() {
	//         c := make(chan os.Signal, 1)
	//         signal.Notify(c, os.Interrupt)
	//         <-c

	//         ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//         defer cancel()
	//         if err := srv.Shutdown(ctx); nil != err {
	//             log.Fatalf("server shutdown failed, err: %v
	// ", err)
	//         }
	//         log.Println("server gracefully shutdown")

	//         close(processed)
	//     }()

	//     // serve
	//     err := srv.ListenAndServe()
	//     if http.ErrServerClosed != err {
	//         log.Fatalf("server not gracefully shutdown, err :%v
	// ", err)
	//     }

	//     // waiting for goroutine above processed
	//     <-processed

}

func graceExit() {
	shutdown.WaitTerminationSignal(func() {
		//异常终止, 删除exporter注册地址
		utils.PgDataStoreRemove(utils.Db)
		utils.Logger.Printf("程序异常退出，清除exporter暴露地址！")
		os.Exit(1)
	})
}

// 定时刷新配置的源路径
func autoRefreshSourceConfig() {

}
