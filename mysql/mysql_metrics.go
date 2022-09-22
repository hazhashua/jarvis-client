package mysql

import (
	"io/ioutil"
	"metric_exporter/config"
	"metric_exporter/utils"

	"gopkg.in/yaml.v2"
)

// cluster:
//   name: bigdata-dev-cluster
//   hosts:
//     - 192.168.10.70
//   port: 3306
//   roles:
//     - master

func Parse_mysql_config() *config.MysqlConfig {
	if bytes, err := ioutil.ReadFile("./mysql/config.yaml"); err == nil {
		var mysql_config config.MysqlConfig
		yaml.Unmarshal(bytes, &mysql_config)
		return &mysql_config

	} else {
		// fmt.Println("解析配置文件出错! ", err.Error())
		utils.Logger.Println("解析配置文件出错! ", err.Error())
		return nil
	}
}
