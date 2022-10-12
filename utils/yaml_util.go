package utils

import (
	"fmt"
	"io/fs"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// type ZookeeperConfig struct {
// 	Cluster struct {
// 		Name       string   `yaml:"name"`
// 		Hosts      []string `yaml:"hosts"`
// 		ClientPort string   `yaml:"clientport"`
// 	}
// }

type prometheusYaml struct {
	Global struct {
		ScrapeInterval     string `yaml:"scrape_interval" mapstructure:"scrape_interval" `
		EvaluationInterval string `yaml:"evaluation_interval" mapstructure:"evaluation_interval"`
	}

	RuleFiles []string `yaml:"rule_files,omitempty" mapstructure:"rule_files"`

	ScrapeConfigs []struct {
		JobName       string `yaml:"job_name,omitempty" mapstructure:"job_name"`
		MetricsPath   string `yaml:"metrics_path,omitempty" mapstructure:"metrics_path"`
		StaticConfigs []struct {
			// Ts struct {
			Targets []string `yaml:"targets,omitempty"`
			// }
		} `yaml:"static_configs" mapstructure:"static_configs"`
	} `yaml:"scrape_configs" mapstructure:"scrape_configs"`
}

type ExporterConfig struct {
	JobName       string `yaml:"job_name,omitempty" mapstructure:"job_name"`
	MetricsPath   string `yaml:"metrics_path,omitempty" mapstructure:"metrics_path"`
	StaticConfigs []struct {
		// Ts struct {
		Targets []string `yaml:"targets,omitempty"`
		// }
	} `yaml:"static_configs" mapstructure:"static_configs"`
}

// 生成prometheus配置文件
func LoadYaml() prometheusYaml {
	var pyaml prometheusYaml
	bytes, _ := ioutil.ReadFile("./prometheus.yml")
	var err2 error
	if err2 = yaml.Unmarshal(bytes, &pyaml); err2 != nil {
		Logger.Println("解析yaml文件失败!", err2.Error())
	}
	return pyaml
}

// yaml对象生成yaml文件
func GenerateYamlFile(pyaml prometheusYaml, absoluteFile string) []byte {
	fmt.Println("pyaml: ", pyaml, pyaml.Global)
	var bytes []byte
	var err error
	if bytes, err = yaml.Marshal(&pyaml); err != nil {
		Logger.Printf("解析yaml对象失败")
	}
	if err = ioutil.WriteFile(absoluteFile, bytes, fs.ModeAppend); err != nil {
		return nil
	}

	Logger.Println("yaml写入文件中成功！")
	return bytes

}
