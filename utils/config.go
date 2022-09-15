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
					fmt.Println("解析配置文件成功...")
				}
			case "hbase":
				configR := config.HbaseConfigure{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					fmt.Println("解析配置文件成功...")
				}
			case "hive":
				configR := config.HiveConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					fmt.Println("解析配置文件成功...")
				}
			case "kafka":
				configR := config.KafkaConfigure{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					fmt.Println("解析配置文件成功...")
				}
			case "micro_service":
				configR := config.K8sYamlConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					fmt.Println("解析配置文件成功...")
				}
			case "mysql":
				configR := config.MysqlConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					fmt.Println("解析配置文件成功...")
				}
			case "node_exporter":
				configR := config.NodeConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					fmt.Println("解析配置文件成功...")
				}
			case "redis":
				configR := config.RedisConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					fmt.Println("解析配置文件成功...")
				}
			case "service_alive":
				fmt.Println("...")
			case "skywalking":
				configR := config.SkyWalkingConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					fmt.Println("解析配置文件成功...")
				}
			case "spark":
				configR := config.SparkConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					fmt.Println("解析配置文件成功...")
				}
			case "zookeeper":
				configR := config.ZookeepeConfig{}
				if err2 = yaml.Unmarshal(bytes, &configR); err2 == nil {
					ConfigStruct.ConfigData[model] = configR
					fmt.Println("解析配置文件成功...")
				}
			}

		} else {
			fmt.Println("读文件失败: ", err.Error())
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
			fmt.Println("unknown datasource...")
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

func init() {
	config := dbConfig{
		Ip:       "192.168.10.68",
		Port:     5432,
		User:     "postgres",
		Password: "pwd@123",
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=cluster port=%d sslmode=disable TimeZone=Asia/Shanghai", config.Ip, config.User, config.Password, config.Port)
	var err error
	if Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err == nil {
		fmt.Println("*************************connect to db success")
		Logger.Println("*************************connect to db success")
	} else {
		fmt.Println("*************************connet to db error!")
		Logger.Println("*************************connect to db error")
	}
}

// 初始化配置
func init() {
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@")
	ConfigStruct = configStruct{}
	ConfigStruct.init()
	fmt.Println("modes: ", ConfigStruct.Modes)
	if Logger != nil {
		Logger.Println("执行完configure初始化...")
	}

	maps := getSourceAddr()
	fmt.Println("从数据库加载的配置: ", maps)

	ConfigStruct.loadAll()
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
	sps := make(map[string][]ServicePort)
	for _, sp := range servicePorts {
		fmt.Println("servicePort: ", sp)
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
