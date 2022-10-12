package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"metric_exporter/config"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// mysql:
//   ip: 192.168.10.70
//   port: 3306
//   username: root
//   password: pwd@123

// postgres:
//   ip: 192.168.10.68
//   port: 5432
//   username: postgres
//   password: pwd@123
//   datasourceinfo:
//     schema: cluster

func ParseDbConfig() *config.DbConfigure {
	if bytes, err := ioutil.ReadFile("/etc/config.yaml"); err == nil {
		var dbconfig config.DbConfigure
		// fmt.Print(dbconfig.Cluster.HttpPort)
		yaml.Unmarshal(bytes, &dbconfig)
		return &dbconfig

	} else {
		fmt.Printf("解析配置文件出错! %s\n", err.Error())
		Logger.Printf("解析配置文件出错! %s\n", err.Error())
		return nil
	}
}

type configData struct {
}

type configStruct struct {
	ConfigData map[string]interface{}
	Modes      []string
	Name       string
	Path       string
}

var ConfigStruct configStruct

type configI interface {
	// 初始化配置的基础信息
	init()
	// 加载所有的配置文件信息
	loadAll()
	// 加载特定的配置信息
	load(modelName string) (ite interface{})
}

// 初始化获取所有的模块
func (cf configStruct) init() {
	ConfigStruct = cf
	modes := make([]string, 0)
	allMode := map[string]int{"hadoop": 0, "hbase": 0, "hive": 0, "kafka": 0, "micro_service": 0,
		"mysql": 0, "node_exporter": 0, "redis": 0, "service_alive": 0, "skywalking": 0, "spark": 0, "zookeeper": 0}
	// 扫描项目目录， 加载所有模块的相关配置
	if fss, err := ioutil.ReadDir("./"); err == nil {
		for _, fs := range fss {
			if fs.IsDir() {
				if _, find := allMode[fs.Name()]; find == true {
					modes = append(modes, fs.Name())
				}
			}
		}
	}
	ConfigStruct.Modes = modes
	fmt.Println("ConfigStruct.Modes: ", ConfigStruct.Modes)
	ConfigStruct.ConfigData = make(map[string]interface{})
	// ConfigStruct.ConfigData["hadoop"] = config.HadoopConfigure{}
	// ConfigStruct.ConfigData["hbase"] = config.HbaseConfigure{}
	// ConfigStruct.ConfigData["hive"] = config.HiveConfig{}
	// ConfigStruct.ConfigData["kafka"] = config.KafkConfigure{}
	// ConfigStruct.ConfigData["micro_service"] = config.K8sYamlConfig{}
	// ConfigStruct.ConfigData["mysql"] = config.MysqlConfig{}
	// ConfigStruct.ConfigData["node_exporter"] = config.NodeConfig{}
	// ConfigStruct.ConfigData["redis"] = config.RedisConfig{}
	// ConfigStruct.ConfigData["skywalking"] = config.SkyWalkingConfig{}
	// ConfigStruct.ConfigData["spark"] = config.SparkConfig{}
	// ConfigStruct.ConfigData["zookeeper"] = config.ZookeepeConfig{}
}

// 解析所有的yaml配置文件，公用类
func (cf configStruct) loadAll() {
	ConfigStruct = cf
	for _, model := range cf.Modes {
		// var confid configData
		path := fmt.Sprintf("./%s/config.yaml", model)
		fmt.Println("path: ", path)
		if bytes, err := ioutil.ReadFile(path); err == nil {
			fmt.Println("bytes: ", string(bytes))
			var err2 error
			switch model {
			case "hadoop":
				configR := config.HadoopConfigure{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					Logger.Printf("解析hadoop配置文件成功...")
				}
			case "hbase":
				configR := config.HbaseConfigure{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					Logger.Printf("解析hbase配置文件成功...")
				}
			case "hive":
				configR := config.HiveConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					Logger.Printf("解析hive配置文件成功...")
				}
			case "kafka":
				configR := config.KafkaConfigure{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					Logger.Printf("解析kafka配置文件成功...")
				}
			case "micro_service":
				configR := config.K8sYamlConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					Logger.Printf("解析微服务配置文件成功...")
				}
			case "mysql":
				configR := config.MysqlConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					Logger.Printf("解析mysql配置文件成功...")
				}
			case "node_exporter":
				configR := config.NodeConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					Logger.Printf("解析node配置文件成功...")
				}
			case "redis":
				configR := config.RedisConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					Logger.Printf("解析redis配置文件成功...")
				}
			case "service_alive":
				Logger.Printf("服务存活未有配置文件")
			case "skywalking":
				configR := config.SkyWalkingConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					Logger.Printf("解析skywalking配置文件成功...")
				}
			case "spark":
				configR := config.SparkConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					Logger.Printf("解析spark配置文件成功...")
				}
			case "zookeeper":
				configR := config.ZookeepeConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					Logger.Printf("解析zookeeper配置文件成功...")
				}
			}

		} else {
			Logger.Printf("读配置文件失败: %s ...\n", err.Error())
		}
	}
}

// 加载特定的配置文件
func (cf configStruct) load(model string) (iface interface{}) {
	// 加载特定的配置文件
	if bytes, err := ioutil.ReadFile(fmt.Sprintf("./%s", model)); err == nil {
		switch model {
		case "hadoop":
			configs := config.HadoopConfigure{}
			if err2 := yaml.Unmarshal(bytes, &configs); err2 == nil {
				// 解析配置
				cf.ConfigData[model] = configs
				return configs
			}
		case "hbase":
			configs := config.HbaseConfigure{}
			if err2 := yaml.Unmarshal(bytes, &configs); err2 == nil {
				// 解析配置
				cf.ConfigData[model] = configs
				return configs
			}
		case "spark":
			configs := config.SparkConfig{}
			if err2 := yaml.Unmarshal(bytes, &configs); err2 == nil {
				// 解析配置
				cf.ConfigData[model] = configs
				return configs
			}
		case "hive":
			configs := config.HiveConfig{}
			if err2 := yaml.Unmarshal(bytes, &configs); err2 == nil {
				// 解析配置
				cf.ConfigData[model] = configs
				return configs
			}
		case "kafka":
			configs := config.KafkaConfigure{}
			if err2 := yaml.Unmarshal(bytes, &configs); err2 == nil {
				// 解析配置
				cf.ConfigData[model] = configs
				return configs
			}
		case "micro_service":
			configs := config.K8sYamlConfig{}
			if err2 := yaml.Unmarshal(bytes, &configs); err2 == nil {
				// 解析配置
				cf.ConfigData[model] = configs
				return configs
			}
		case "mysql":
			configs := config.MysqlConfig{}
			if err2 := yaml.Unmarshal(bytes, &configs); err2 == nil {
				// 解析配置
				cf.ConfigData[model] = configs
				return configs
			}
		case "node_exporter":
			configs := config.NodeConfig{}
			if err2 := yaml.Unmarshal(bytes, &configs); err2 == nil {
				// 解析配置
				cf.ConfigData[model] = configs
				return configs
			}
		case "redis":
			configs := config.RedisConfig{}
			if err2 := yaml.Unmarshal(bytes, &configs); err2 == nil {
				// 解析配置
				cf.ConfigData[model] = configs
				return configs
			}
		case "service_alive":
			fmt.Println("没有配置文件...")
		case "skywalking":
			configs := config.SkyWalkingConfig{}
			if err2 := yaml.Unmarshal(bytes, &configs); err2 == nil {
				// 解析配置
				cf.ConfigData[model] = configs
				return configs
			}
		case "zookeeper":
			configs := config.ZookeepeConfig{}
			if err2 := yaml.Unmarshal(bytes, &configs); err2 == nil {
				// 解析配置
				cf.ConfigData[model] = configs
				return configs
			}
		default:
			Logger.Printf("unknown datasource")
		}
	}
	return nil
}

// 全局Logger对象
var Logger *log.Logger

func init() {
	if Logger == nil {
		fmt.Println("日志对象为空，创建日志对象...")
		if logFile, err := os.OpenFile("./exporter.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			Logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
		}
	}
}

// 全局Db对象
var Db *gorm.DB
var DbConfig *config.DbConfigure

func init() {
	config := ParseDbConfig()
	// 赋值全局配置变量
	DbConfig = config
	// config := dbConfig{
	// 	Ip:       "192.168.10.68",
	// 	Port:     5432,
	// 	User:     "postgres",
	// 	Password: "pwd@123",
	// }
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=cluster port=%d sslmode=disable TimeZone=Asia/Shanghai", config.Cluster.Postgres.Ip, config.Cluster.Postgres.Username, config.Cluster.Postgres.Password, config.Cluster.Postgres.Port)
	var err error
	if Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err == nil {
		Logger.Println("*************************connect to db success")
	} else {
		Logger.Println("*************************connect to db error")
	}
}

// 初始化配置
func init() {
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@")
	// ClusterName = "bigdata-dev-cluster"
	ConfigStruct = configStruct{}
	ConfigStruct.init()
	fmt.Println("modes: ", ConfigStruct.Modes)
	if Logger != nil {
		Logger.Println("执行完configure初始化...")
	}

	// 从数据库加载配置
	maps := getSourceAddr()
	if len(maps) != 0 {
		Logger.Printf("从数据库加载的配置: %v\n", maps)
		for model, datas := range maps {
			switch model {
			case config.HADOOP:
				// cluster:
				// name: bigdata-dev-cluster
				// services:
				// 	- 192.168.10.220:8088
				// 	- 192.168.10.220:50070
				// 	- 192.168.10.220:8020
				// 	- 192.168.10.221:8020
				// 	- 192.168.10.222:8020
				// servicenum: 5
				// namenodes:
				// 	- 192.168.10.220
				// 	- 192.168.10.221
				// 	- 192.168.10.222
				// namenodehosts:
				// 	- bigdata-dev01
				// 	- bigdata-dev02
				// 	- bigdata-dev03
				// namenodehttpport: 50070
				// namenoderpcport: 8020
				// datanodes:
				// 	- 192.168.10.220
				// 	- 192.168.10.221
				// 	- 192.168.10.222
				// datanodehosts:
				// 	- bigdata-dev01
				// 	- bigdata-dev02
				// 	- bigdata-dev03
				// datanodehttpport: 9864
				// datanoderpcport: 9867
				// resourcemanagers:
				// 	- 192.168.10.220
				// 	- 192.168.10.221
				// 	- 192.168.10.222
				// resourcemanagerhosts:
				// 	- bigdata-dev01
				// 	- bigdata-dev02
				// 	- bigdata-dev03
				// resourcemanagerurl: http://192.168.10.220:8088/jmx
				// resourcemanagerhost: bigdata-dev01
				// resourcemanagerhttpport: 8088
				hc := config.HadoopConfigure{}
				hc.Cluster.Name = DbConfig.Cluster.Name
				for _, data := range datas {
					if *data.ChildService == "resourcemanager" {
						hc.Cluster.ResourceManagerHosts = append(hc.Cluster.ResourceManagerHosts, "")
						hc.Cluster.ResourceManagers = append(hc.Cluster.ResourceManagers, *data.IP)
						hc.Cluster.ResourcemanagerHttpPort = int(data.Port.Int64)
						hc.Cluster.ResourceManagerUrl = fmt.Sprintf("http://%s:%d/jmx", *data.IP, data.Port.Int64)
					}
					if *data.ChildService == "namenode" {
						hc.Cluster.NamenodeHttpPort = int(data.Port.Int64)
						hc.Cluster.NamenodeHosts = append(hc.Cluster.NamenodeHosts, "")
						hc.Cluster.Namenodes = append(hc.Cluster.Namenodes, *data.IP)
					}
					if *data.ChildService == "datanode" {
						hc.Cluster.DatanodeHosts = append(hc.Cluster.DatanodeHosts, "")
						hc.Cluster.Datanodes = append(hc.Cluster.Datanodes, *data.IP)
						hc.Cluster.DatanodeHttpPort = int(data.Port.Int64)
					}
				}
				ConfigStruct.ConfigData[model] = hc
			case config.HBASE:
				// cluster:
				// 	clustername: bigdata-dev-cluster
				// 	masterjmxport: 16010
				// 	regionserverjmxport: 16030
				// 	hosts:
				// 		- 192.168.10.220
				// 		- 192.168.10.221
				// 		- 192.168.10.222
				// 	names:
				// 		- bigdata-dev01
				// 		- bigdata-dev02
				// 		- bigdata-dev03
				hbaseConf := config.HbaseConfigure{}
				hbaseConf.Cluster.ClusterName = DbConfig.Cluster.Name
				for _, data := range datas {
					hbaseConf.Cluster.Names = append(hbaseConf.Cluster.Names, "")
					if *data.ChildService == "hmaster" {
						hbaseConf.Cluster.MasterJmxPort = fmt.Sprintf("%d", data.Port.Int64)
						hbaseConf.Cluster.Hosts = append(hbaseConf.Cluster.Hosts, *data.IP)
					}
					if *data.ChildService == "regionserver" {
						hbaseConf.Cluster.RegionserverJmxPort = fmt.Sprintf("%d", data.Port.Int64)
					}
				}
				ConfigStruct.ConfigData[config.HBASE] = hbaseConf
			case config.HIVE:
				// cluster:
				// 	name: bigdata-dev-cluster
				// 	hosts:
				// 		- 192.168.10.220
				// 		- 192.168.10.221
				// 		- 192.168.10.222
				// 	rpcport: 10000
				// 	mysql:
				// 		host: 192.168.10.223
				// 		port: 3306
				// 		user: root
				// 		password: pwd@123
				// 	scrapehost: bigdata-dev01
				// 	scrapeip: 192.168.10.220
				hiveConf := config.HiveConfig{}
				hiveConf.Cluster.Name = DbConfig.Cluster.Name
				for _, data := range datas {
					hiveConf.Cluster.Hosts = append(hiveConf.Cluster.Hosts, "")
					if *data.ChildService == "metastore" {
						hiveConf.Cluster.Mysql.Host = *data.IP
						hiveConf.Cluster.Mysql.Port = int(data.Port.Int64)
						hiveConf.Cluster.Mysql.User = *data.Username
						hiveConf.Cluster.Mysql.Password = *data.Password
					}
				}
				ConfigStruct.ConfigData[config.HIVE] = hiveConf
			case config.KAFKA:
				// cluster:
				// 	name: bigdata-dev-cluster
				// 	hosts:
				// 		- 192.168.10.220
				// 		- 192.168.10.221
				// 		- 192.168.10.222
				// 	port: 9092
				// 	env: dev
				kafkaConf := config.KafkaConfigure{}
				kafkaConf.Cluster.Name = DbConfig.Cluster.Name
				for _, data := range datas {
					kafkaConf.Cluster.Hosts = append(kafkaConf.Cluster.Hosts, *data.IP)
					kafkaConf.Cluster.Port = int(data.Port.Int64)
				}
				ConfigStruct.ConfigData[config.KAFKA] = kafkaConf
			case config.MICROSERVICE:
				// cluster:
				// 	name: 测试kubernetes集群
				// 	master:
				// 		- 192.168.10.20
				// 		- 192.168.10.21
				// 		- 192.168.10.22
				// 	nodes:
				// 		- 192.168.10.20
				// 		- 192.168.10.21
				// 		- 192.168.10.22
				// 		- 192.168.10.23
				// 		- 192.168.10.24
				// 		- 192.168.10.32
				// 		- 192.168.10.63
				// 		- 192.168.10.111
				// 	apiserverport: 8080
				k8syamlConf := config.K8sYamlConfig{}
				k8syamlConf.Cluster.Name = DbConfig.Cluster.Name
				for _, data := range datas {
					fmt.Println("micro_service: ", data)
					if *data.ChildService == "apiserver" {
						k8syamlConf.Cluster.Master = append(k8syamlConf.Cluster.Master, *data.IP)
						k8syamlConf.Cluster.ApiServerPort = fmt.Sprintf("%d", data.Port.Int64)
					}
				}
				ConfigStruct.ConfigData[config.MICROSERVICE] = k8syamlConf
			case config.MYSQL:

				// cluster:
				// 	name: bigdata-dev-cluster
				// 	ips:
				// 	- 192.168.10.70
				// 	port: 3306
				// 	username: root
				// 	password: pwd@123
				// 	defaultdb: information_schema
				// 	role:
				// 	- master

				mysqlConf := config.MysqlConfig{}
				mysqlConf.Cluster.Name = DbConfig.Cluster.Name
				for _, data := range datas {
					if *data.ChildService == "mysqld" {
						mysqlConf.Cluster.Ips = append(mysqlConf.Cluster.Ips, *data.IP)
						mysqlConf.Cluster.Port = int(data.Port.Int64)
						mysqlConf.Cluster.Username = *data.Username
						mysqlConf.Cluster.Password = *data.Password
						mysqlConf.Cluster.DefaultDB = "information_schema"
					}
				}
				ConfigStruct.ConfigData[config.MYSQL] = mysqlConf

			case config.REDIS:
				// cluster:
				// 	name: bigdata-dev-cluster
				// 	ips:
				// 		- 192.168.10.107
				// 		- 192.168.10.108
				// 		- 192.168.10.109
				// 	hosts:
				// 		- redis-dev-1
				// 		- redis-dev-2
				// 		- redis-dev-3
				// 	ippwds:
				// 		- rhcloud@123.com
				// 		- rhcloud@123.com
				// 		- rhcloud@123.com
				// 	scrapehost: redis-dev-1
				// 	scrapeip: 192.168.10.107
				// 	redisport: 6379
				redisConf := config.RedisConfig{}
				redisConf.Cluster.Name = DbConfig.Cluster.Name
				for _, data := range datas {
					if *data.ChildService == "redis" {
						redisConf.Cluster.Hosts = append(redisConf.Cluster.Hosts, *data.IP)
						redisConf.Cluster.Ips = append(redisConf.Cluster.Ips, *data.IP)
						redisConf.Cluster.Ippwds = append(redisConf.Cluster.Ippwds, *data.Password)
						redisConf.Cluster.RedisPort = int(data.Port.Int64)
					}
				}
				ConfigStruct.ConfigData[config.REDIS] = redisConf

			case config.SKYWALKING:
				// cluster:
				// 	name: bigdata-dev-cluster
				// 	elasticsearch:
				// 		ips:
				// 		- 192.168.10.65
				// 		port: 9200
				swConf := config.SkyWalkingConfig{}
				swConf.Cluster.Name = DbConfig.Cluster.Name
				for _, data := range datas {
					if *data.ChildService == "elasticsearch" {
						swConf.Cluster.ElasticSearch.Ips = append(swConf.Cluster.ElasticSearch.Ips, *data.IP)
						swConf.Cluster.ElasticSearch.Port = int(data.Port.Int64)
					}
				}
				ConfigStruct.ConfigData[config.SKYWALKING] = swConf

			case config.SPARK:
				// cluster: bigdata-dev-cluster
				// 	masterhttp:
				// 	ips:
				// 		- 192.168.10.220
				// 		- 192.168.10.221
				// 		- 192.168.10.222
				// 	port: 8080
				// 	path: /metrics/prometheus

				// 	workerhttp:
				// 	ips:
				// 		- 192.168.10.220
				// 		- 192.168.10.221
				// 		- 192.168.10.222
				// 	port: 8081
				// 	path: /metrics/prometheus

				// 	applicationhttp:
				// 	ips:
				// 		- 192.168.10.220
				// 		- 192.168.10.221
				// 		- 192.168.10.222
				// 	ports:
				// 		- 4040
				// 		- 4041
				// 		- 4042
				// 	mainpath: /metrics/prometheus
				// 	executorpath: /metrics/executors/prometheus
				sparkConf := config.SparkConfig{}
				sparkConf.Cluster = DbConfig.Cluster.Name
				for _, data := range datas {
					if *data.ChildService == "master" && *data.PortType == "http" {
						sparkConf.Masterhttp.Ips = append(sparkConf.Masterhttp.Ips, *data.IP)
						sparkConf.Masterhttp.Port = int(data.Port.Int64)
						sparkConf.Masterhttp.Path = "/metrics/prometheus"

						sparkConf.Applicationhttp.Ips = append(sparkConf.Applicationhttp.Ips, *data.IP)
						sparkConf.Applicationhttp.Ports = []int{4040, 4041, 4042}
						sparkConf.Applicationhttp.MainPath = "/metrics/prometheus"
					}
					if *data.ChildService == "worker" && *data.PortType == "http" {
						sparkConf.Workerhttp.Ips = append(sparkConf.Workerhttp.Ips, *data.IP)
						sparkConf.Workerhttp.Port = int(data.Port.Int64)
						sparkConf.Workerhttp.Path = "/metrics/prometheus"
					}
				}
				ConfigStruct.ConfigData[config.SPARK] = sparkConf

			case config.ZOOKEEPER:
				// cluster:
				// 	name: bigdata-dev-cluster
				// 	hosts:
				// 		- 192.168.10.220
				// 		- 192.168.10.221
				// 		- 192.168.10.222
				// 	clientport: 2181
				zkConf := config.ZookeepeConfig{}
				zkConf.Cluster.Name = DbConfig.Cluster.Name
				for _, data := range datas {
					if *data.ChildService == "zookeeper" {
						zkConf.Cluster.Hosts = append(zkConf.Cluster.Hosts, *data.IP)
						zkConf.Cluster.ClientPort = fmt.Sprintf("%d", data.Port.Int64)
					}
				}
				ConfigStruct.ConfigData[config.ZOOKEEPER] = zkConf

			}
		}
		// 默认加载node配置
		nodeConf := config.NodeConfig{}
		nodeConf.Cluster.Name = DbConfig.Cluster.Name
		ConfigStruct.ConfigData[config.NODE] = nodeConf
	} else {
		// 从配置文件读取配置
		ConfigStruct.loadAll()
	}

	fmt.Println("ConfigStruct.ConfigData: ", ConfigStruct.ConfigData)
	if Logger != nil {
		Logger.Println("执行完configure配置")
	}
}

// 数据库源数据的加载
// 把ip和端口信息封装成访问地址类型信息
// 配置默认的端口，防止在信息不全情况下数据获取的障碍
func getSourceAddr() map[string][]ServicePort {
	servicePorts := PgServiceQuery(Db)
	Logger.Printf("读取数据库数据的记录数: %d\n", len(servicePorts))
	sps := make(map[string][]ServicePort)
	for _, sp := range servicePorts {
		fmt.Println("servicePort: ", *sp.ServiceName, *sp.ChildService, *sp.IP, *sp.Comment)
		if len(sps[*sp.ServiceName]) == 0 {
			sps[*sp.ServiceName] = make([]ServicePort, 0)
		}
		sps[*sp.ServiceName] = append(sps[*sp.ServiceName], sp)
	}

	// 提前hadoop相关源数据信息
	resourceMap := make(map[string][]ServicePort, 0)
	hadoopRes := make([]ServicePort, 0)
	for _, ele := range sps[config.HADOOP] {
		switch *ele.ChildService {
		case "resourcemanager":
			if *ele.PortType == "http" && strings.Contains(*ele.Comment, "jmx") {
				hadoopRes = append(hadoopRes, ele)
			}
		case "namenode":
			if *ele.PortType == "http" && strings.Contains(*ele.Comment, "jmx") {
				hadoopRes = append(hadoopRes, ele)
			}
		case "datanode":
			if *ele.PortType == "http" && strings.Contains(*ele.Comment, "jmx") {
				hadoopRes = append(hadoopRes, ele)
			}
		}
	}
	resourceMap[config.HADOOP] = hadoopRes

	// 提取hbase相关源数据信息
	hbaseRes := make([]ServicePort, 0)
	for _, ele := range sps[config.HBASE] {
		switch *ele.ChildService {
		case "hmaster":
			if *ele.PortType == "http" && strings.Contains(*ele.Comment, "jmx") {
				hbaseRes = append(hbaseRes, ele)
			}
		case "regionserver":
			if *ele.PortType == "http" && strings.Contains(*ele.Comment, "jmx") {
				hbaseRes = append(hbaseRes, ele)
			}
		}
	}
	resourceMap[config.HBASE] = hbaseRes

	// 读取hive的相关源数据信息
	hiveRes := make([]ServicePort, 0)
	for _, ele := range sps[config.HIVE] {
		switch *ele.ChildService {
		case "metastore":
			hiveRes = append(hiveRes, ele)
		}
	}
	resourceMap[config.HIVE] = hiveRes

	// 读取kafka的相关源数据信息
	kafkaRes := make([]ServicePort, 0)
	for _, ele := range sps[config.KAFKA] {
		kafkaRes = append(kafkaRes, ele)
	}
	resourceMap[config.KAFKA] = kafkaRes

	// 读取micro_service的相关源数据信息
	micSerRes := make([]ServicePort, 0)
	for _, ele := range sps[config.MICROSERVICE] {
		switch *ele.ChildService {
		case "apiserver":
			micSerRes = append(micSerRes, ele)
		}
	}
	resourceMap[config.MICROSERVICE] = micSerRes

	// 读取mysql的相关源数据信息
	mysqlRes := make([]ServicePort, 0)
	for _, ele := range sps[config.MYSQL] {
		mysqlRes = append(mysqlRes, ele)
	}
	resourceMap[config.MYSQL] = mysqlRes

	// 读取redis的相关源数据信息
	redisRes := make([]ServicePort, 0)
	for _, ele := range sps[config.REDIS] {
		redisRes = append(redisRes, ele)
	}
	resourceMap[config.REDIS] = redisRes

	// 读取skywalking相关源数据信息
	skyRes := make([]ServicePort, 0)
	for _, ele := range sps[config.SKYWALKING] {
		if *ele.ChildService == "elasticsearch" {
			skyRes = append(skyRes, ele)
		}
	}
	resourceMap[config.SKYWALKING] = skyRes

	// 读取spark相关的源数据信息
	sparkRes := make([]ServicePort, 0)
	for _, ele := range sps[config.SPARK] {
		switch *ele.ChildService {
		case "master":
			if *ele.PortType == "http" {
				sparkRes = append(sparkRes, ele)
			}
		case "worker":
			if *ele.PortType == "http" {
				sparkRes = append(sparkRes, ele)
			}
		}
	}
	resourceMap[config.SPARK] = sparkRes

	// 读取zookeeper相关的源数据信息
	zookeeperRes := make([]ServicePort, 0)
	for _, ele := range sps[config.ZOOKEEPER] {
		zookeeperRes = append(zookeeperRes, ele)
	}
	resourceMap[config.ZOOKEEPER] = zookeeperRes

	return resourceMap

}
