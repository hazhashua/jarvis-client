package spark

import (
	"fmt"
	"io/ioutil"
	"metric_exporter/config"
	"metric_exporter/utils"
	"net/http"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

/*
*  spark获取指标数据时：
*	1. spark-submit 添加参数，打开metric http地址 参数包括: --conf "spark.ui.prmetheus.enabled=true"
*   2. spark集群需要添加配置:
*        metrics.properties
*          master.source.jvm.class=org.apache.spark.metrics.source.JvmSource
*          worker.source.jvm.class=org.apache.spark.metrics.source.JvmSource
*          driver.source.jvm.class=org.apache.spark.metrics.source.JvmSource
*          executor.source.jvm.class=org.apache.spark.metrics.source.JvmSource
*          applications.source.jvm.class=org.apache.spark.metrics.source.JvmSource
*
*          *.sink.prometheusServlet.class=org.apache.spark.metrics.sink.PrometheusServlet
*          *.sink.prometheusServlet.path=/metrics/prometheus
 */

type SparkHandler struct {
	Metrics []string
}

func (handler SparkHandler) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	handler.Metrics = GetMetrics()
	switch r.URL.Path {
	case "/spark/metrics":
		for _, value := range handler.Metrics {
			fmt.Fprintf(writer, "%s", value)
		}
	default:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "no such page: %s\n", r.URL)
	}
}

// spark 分为http的接口获取数据的方式
// :8080/metrics/master/prometheus master汇总的相关信息地址
func GetMetrics() []string {
	// url_array := []string{"http://192.168.10.220", "http://192.168.10.221", "http://192.168.10.222"}
	arrs := make([]string, 0)
	sparkConfig, err := (utils.ConfigStruct.ConfigData[config.SPARK]).(config.SparkConfig)

	if err != true {
		fmt.Println("load spark configure failed!")
		return []string{}
	}

	url_array := make([]string, 0)
	for _, ip := range sparkConfig.Masterhttp.Ips {
		url_array = append(url_array, fmt.Sprintf("http://%s", ip))
	}
	url_s := sparkConfig.Masterhttp.Ips
	fmt.Println("url_s:", url_s)
	ports := sparkConfig.Applicationhttp.Ports
	ports = append(ports, sparkConfig.Workerhttp.Port)
	fmt.Println("ports: ", ports)
	// 抓取master的网页地址，获取active的地址
	// 添加is_active_node指标
	var active_node_index int
	for idx, url := range url_array {

		for _, port := range ports {
			//获取driver进程的堆栈内存使用率
			// master_metric_url := yamlConfig.Masterhttp.Ips[idx] + ":" + fmt.Sprintf("%d", yamlConfig.Masterhttp.Port) + yamlConfig.Masterhttp.Path
			// fmt.Println("master_metric_url", master_metric_url)
			metric_url := sparkConfig.Applicationhttp.Ips[idx] + ":" + fmt.Sprintf("%d", port) + sparkConfig.Applicationhttp.MainPath
			fmt.Println("metric_url: ", metric_url)

			metric_url = fmt.Sprintf(url+":%d/metrics/prometheus", port)
			fmt.Println("metric_url: ", metric_url)

			metric_response := utils.GetUrl(metric_url)
			if metric_response == "" {
				fmt.Println(fmt.Sprintf("机器:%s 上没有运行的程序", url))
			}

			// fmt.Println("driver_response: ", driver_response)
			for _, line := range strings.Split(metric_response, "\n") {
				if strings.Contains(line, "_driver_jvm_heap_usage_Value") {
					fmt.Println("contains: ", line)
					reg := regexp.MustCompile("metrics_(.*)_driver_jvm_heap_usage_Value.*")
					app_name := reg.FindStringSubmatch(line)[1]
					fmt.Println("app name: ", app_name)
					arrs = append(arrs, "driver_jvm_heap_usage{type=\"gauges\", application_name=\""+app_name+"\", host=\""+url+"\" }\n")
				}

				if strings.Contains(line, "metrics_jvm_heap_usage_Value") {
					line = strings.ReplaceAll(line, "metrics_jvm_heap_usage_Value", "worker_jvm_heap_usage")
					regexp := regexp.MustCompile("[^{]*{(.*)}.*")
					ss := regexp.FindStringSubmatch(line)
					line = strings.ReplaceAll(line, ss[1], ss[1]+","+"host=\""+url+"\" ")
					arrs = append(arrs, line+"\n")
				}

				// if strings.Contains(line, "metrics_jvm_heap_usage_Value") {
				// 	line = strings.ReplaceAll(line, "metrics_jvm_heap_usage_Value", "master_jvm_heap_usage")
				// 	ss := regexp.FindStringSubmatch(line)
				// 	line = strings.ReplaceAll(line, ss[1], ss[1]+","+"host=\""+url_array[active_node_index]+"\" ")
				// 	arrs = append(arrs, line+"\n")
				// }
			}
		}

		//查询所有maser的状态
		// res := utils.GetUrl(url + ":28080")
		res := utils.GetUrl(url + ":8080")
		host_list := strings.Split(url, "/")
		fmt.Println("host_list: ", host_list)
		host := host_list[2]
		is_active_node := strings.Contains(res, "<strong>Status:</strong> ALIVE")
		fmt.Println("is_active_node: ", is_active_node)
		is_standby_node := strings.Contains(res, "<strong>Status:</strong> STANDBY")
		fmt.Println("is_standby_node: ", is_standby_node)

		if is_active_node {
			arrs = append(arrs, "is_active_master"+"{type=\"gauges\", host=\""+host+"\"} 1\n")
			active_node_index = idx

			// 获取完成的app数量
			fmt.Println("active_url: ", url+fmt.Sprintf(":%d", sparkConfig.Masterhttp.Port))
			// response := utils.GetUrl(urls[active_node_index] + ":28080")
			response := utils.GetUrl(url_array[active_node_index] + ":8080")
			reg := regexp.MustCompile("(\\d+) <a href=\"#completed-app\">Completed</a>")
			match_strings := reg.FindStringSubmatch(response)
			fmt.Println("match strings: ", match_strings[1])
			// strconv.Itoa(match_strings[1])    int to string...
			arrs = append(arrs, "master_finished_apps{type=\"gauges\", host=\""+host+"\"} "+match_strings[1]+"\n")

			reg = regexp.MustCompile("(\\d+) <a href=\"#running-app\">Running</a>")
			match_strings = reg.FindStringSubmatch(response)
			fmt.Println("match strings: ", match_strings[1])
			// strconv.Itoa(match_strings[1])    int to string...
			arrs = append(arrs, "master_running_apps{type=\"gauges\", host=\""+host+"\"} "+match_strings[1]+"\n")
		} else if is_standby_node {
			arrs = append(arrs, "is_active_master"+"{type=\"gauges\", host=\""+host+"\"} 0\n")
		} else {
			continue
		}
		fmt.Println("url: ", url)
	}
	fmt.Println("active node index: ", active_node_index)

	// 0 <a href="#completed-app">Completed</a>
	// fmt.Println("match strings: ", match_strings)

	// 查询active http metric数据
	response := utils.GetUrl(url_array[active_node_index] + ":8080/metrics/prometheus/")
	if response == "" {
		fmt.Println("active master的指标数据为空！！！！")
	}

	regexp := regexp.MustCompile("[^{]*{(.*)}.*")
	for _, line := range strings.Split(response, "\n") {
		// fmt.Println("line: ", line)
		if strings.Contains(line, "metrics_master_aliveWorkers_Value") {
			line = strings.ReplaceAll(line, "metrics_master_aliveWorkers_Value", "master_alive_workers")
			arrs = append(arrs, line+"\n")
		}
		if strings.Contains(line, "metrics_master_apps_Value") {
			line = strings.ReplaceAll(line, "metrics_master_apps_Value", "master_apps")
			arrs = append(arrs, line+"\n")
		}
		if strings.Contains(line, "metrics_master_waitingApps_Value") {
			line = strings.ReplaceAll(line, "metrics_master_waitingApps_Value", "master_waiting_apps")
			arrs = append(arrs, line+"\n")
		}
		if strings.Contains(line, "metrics_master_workers_Value") {
			line = strings.ReplaceAll(line, "metrics_master_workers_Value", "master_workers")
			arrs = append(arrs, line+"\n")
		}
		// master 进程堆内存使用率
		if strings.Contains(line, "metrics_jvm_heap_usage_Value") {
			line = strings.ReplaceAll(line, "metrics_jvm_heap_usage_Value", "master_jvm_heap_usage")
			ss := regexp.FindStringSubmatch(line)
			line = strings.ReplaceAll(line, ss[1], ss[1]+","+"host=\""+url_array[active_node_index]+"\" ")
			arrs = append(arrs, line+"\n")
		}
	}
	cluster := fmt.Sprintf("cluster=\"%s\"", sparkConfig.Cluster)
	print_metrics := []string{}
	for _, line := range arrs {
		ss := regexp.FindStringSubmatch(line)
		//匹配到的话 0为全串，1，2...为()内的串
		// fmt.Println("find string: ", ss[1])
		s := strings.ReplaceAll(line, ss[1], ss[1]+","+cluster)
		print_metrics = append(print_metrics, s)
	}
	fmt.Println("prit_metrics: ", print_metrics)

	// 解析active地址中的有用的metric信息
	return print_metrics
}

func ParseSparkConf() (*config.SparkConfig, error) {
	config := new(config.SparkConfig)
	// var yamlConfig YamlConfig
	dir, _ := os.Getwd()
	confPath := dir + "/spark/config.yaml"
	fmt.Println("confPath: ", confPath)

	data, _ := ioutil.ReadFile(confPath)
	err := yaml.Unmarshal(data, config)
	if err != nil {
		fmt.Println("err: ", err)
	}
	fmt.Println("configStruct: ", config)

	fmt.Println("yamlConfig data: ", config.Masterhttp.Ips, config.Masterhttp.Port, config.Masterhttp.Path)
	fmt.Println("configStruct.IpList: ", config.Applicationhttp.ExecutorPath, config.Applicationhttp.Ips, config.Applicationhttp.Ports, config.Applicationhttp.MainPath)
	return config, nil
}
