package skywalking

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

type EventInfo struct {
	EventInfoDesc    *prometheus.Desc
	EventInfoValType prometheus.ValueType
}

// type SkyWalkingInfo struct{
// 	eventInfos []EventInfo
// }

// cluster:
//   name: bigdata-dev-cluster
//   elasticsearch:
//     ips:
//       - 192.168.10.65
//     port: 9200

type SkyWalkingConfig struct {
	Cluster struct {
		Name          string `json:"name"`
		ElasticSearch struct {
			Ips  []string `json:"ips"`
			Port int      `json:port`
		}
	}
}

type SkyWalkingExporter struct {
	EventInfos []EventInfo
}

// 读取skywalking的相关配置
func ParseSkyWalkingConfig() *SkyWalkingConfig {
	var skywalkingConfig SkyWalkingConfig
	if bytes, err := ioutil.ReadFile("./skywalking/config.yaml"); err == nil {
		yaml.Unmarshal(bytes, &skywalkingConfig)
	} else {
		fmt.Println("解析本地skywalking配置文件失败!")
	}
	return &skywalkingConfig
}

// 创建skywalking exporter对象
func NewSkywalkingExporter() *SkyWalkingExporter {
	event_info_list := make([]EventInfo, 0)
	// skywalkingConfig := ParseSkyWalkingConfig()
	year, month, day := time.Now().Date()
	eventIndex := fmt.Sprintf("sw_events-%04d%02d%02d", year, month, day)
	var typ SkwEvent
	events := GetAll(eventIndex, "_doc", typ)
	for _, event := range events {
		switch ret := event.(type) {
		case string:
			fmt.Println("event.(type): ", ret)
		case SkwEvent:
			typ := SkwEvent(ret)
			fmt.Printf("event: %v \n", typ)
		default:
		}
		eventInfo := prometheus.NewDesc("event_info", "描述事件的详细信息",
			[]string{"name", "type", "service_name", "start_time", "end_time", "message"},
			prometheus.Labels{})
		evnetInfoValType := prometheus.GaugeValue
		event_info_list = append(event_info_list, EventInfo{
			EventInfoDesc:    eventInfo,
			EventInfoValType: evnetInfoValType,
		})
	}
	return &SkyWalkingExporter{
		EventInfos: event_info_list,
	}
}

func (e *SkyWalkingExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, event := range e.EventInfos {
		ch <- event.EventInfoDesc
	}
}

// 收集skywalking事件方法
func (e *SkyWalkingExporter) Collect(ch chan<- prometheus.Metric) {

	println("e address: %p", e, "**************************************")
	// 抓取当天的event索引数据
	year, month, day := time.Now().Date()
	eventIndex := fmt.Sprintf("sw_events-%04d%02d%02d", year, month, day)
	var typ SkwEvent
	events := GetAll(eventIndex, "_doc", typ)
	convertedObjs := make([]SkwEvent, 0)
	for _, event := range events {
		switch ret := event.(type) {
		case string:
			fmt.Println("event.(type): ", ret)
		case SkwEvent:
			typ := SkwEvent(ret)
			fmt.Printf("event: %v \n", typ)
			fmt.Println("typ.Service: ", typ.Service)
			fmt.Println("typ.Name: ", typ.Name)
			fmt.Println("typ.Message: ", typ.Message)
			convertedObjs = append(convertedObjs, typ)
		default:
		}
	}

	for idx, eventInfo := range e.EventInfos {
		// "name", "type", "service_name", "start_time", "end_time", "message"
		// se := SkyEvent(events[idx])
		ch <- prometheus.MustNewConstMetric(eventInfo.EventInfoDesc, eventInfo.EventInfoValType, 1,
			convertedObjs[idx].Name, convertedObjs[idx].Type, convertedObjs[idx].Service, fmt.Sprintf("%d", convertedObjs[idx].StartTime), fmt.Sprintf("%d,", convertedObjs[idx].EndTime), convertedObjs[idx].Message)
	}
}
