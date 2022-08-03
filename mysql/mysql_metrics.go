package mysql

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// cluster:
//   name: bigdata-dev-cluster
//   hosts:
//     - 192.168.10.70
//   port: 3306
//   roles:
//     - master

// 存储mysql配置文件信息
type MysqlConfig struct {
	Cluster struct {
		Name string   `yaml:"name"`
		Ips  []string `yaml:"ips"`
		Port int      `yaml:"port"`
		Role []string `yaml:"role"`
	}
}

func Parse_mysql_config() *MysqlConfig {
	if bytes, err := ioutil.ReadFile("./mysql/config.yaml"); err == nil {
		var mysql_config MysqlConfig
		yaml.Unmarshal(bytes, &mysql_config)
		return &mysql_config

	} else {
		fmt.Println("解析配置文件出错! ", err.Error())
		return nil
	}
}
