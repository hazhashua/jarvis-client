package utils

import (
	"fmt"
	"io/ioutil"
	"metric_exporter/config"

	"gopkg.in/yaml.v2"
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

// 初始化配置
func init() {
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@")
	ConfigStruct = configStruct{}
	ConfigStruct.init()
	fmt.Println("modes: ", ConfigStruct.Modes)
	if Logger != nil {
		Logger.Println("执行完configure初始化...")
	}
	ConfigStruct.loadAll()
	fmt.Println("ConfigStruct.ConfigData: ", ConfigStruct.ConfigData)
	if Logger != nil {
		Logger.Println("执行完configure配置")
	}
}

// 数据库源数据的加载
// 把ip和端口信息封装成访问地址类型信息
// 配置默认的端口，防止在信息不全情况下数据获取的障碍
