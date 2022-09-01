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
	configData map[string]interface{}
	modes      []string
	name       string
	path       string
}

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
	cf.modes = modes
}

// 解析所有的yaml配置文件，公用类
func (cf configStruct) loadAll() {
	for _, model := range cf.modes {
		var confid configData
		path := fmt.Sprintf("./%s", model)
		if bytes, err := ioutil.ReadFile(path); err == nil {
			fmt.Println("bytes: ", string(bytes))
			if err2 := yaml.Unmarshal(bytes, &confid); err2 == nil {
				fmt.Println("解析配置文件成功...")
				cf.configData[model] = confid
			}
		}
	}
}

// 加载特定的配置文件
func (cf configStruct) load(model string) {
	// 加载特定的配置文件
	var configd interface{}
	if bytes, err := ioutil.ReadFile(fmt.Sprintf("./%s", model)); err == nil {
		switch model {
		case "hadoop":
			configd = config.HadoopConfigure{}
		case "hbase":
			configd = config.HbaseConfigure{}
		case "spark":
			// configd := config.
		case "hive":
		case "kafka":
		case "micro_service":
		case "mysql":
		case "node_exporter":
		case "redis":
		case "service_alive":
		case "skywalking":
		case "zookeeper":

		default:
			fmt.Println("unknown datasource...")
		}
		if err2 := yaml.Unmarshal(bytes, &configd); err2 == nil {
			// 解析配置

		}

	}
}
