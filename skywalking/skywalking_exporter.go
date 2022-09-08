package skywalking

import (
	"fmt"
	"io/ioutil"
	"metric_exporter/config"
	"metric_exporter/utils"
	"sort"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

type EventInfo struct {
	EventInfoDesc    *prometheus.Desc
	EventInfoValType prometheus.ValueType
}

type ServiceCpm struct {
	ServiceCpmDesc    *prometheus.Desc
	ServiceCpmValType prometheus.ValueType
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

type SkyWalkingExporter struct {
	EventInfos    []EventInfo
	ServiceCpms   []ServiceCpm
	SkyEventDatas []SkwEvent
	CpmDatas      []MyCpmInfo
}

// 读取skywalking的相关配置
func ParseSkyWalkingConfig() *config.SkyWalkingConfig {
	var skywalkingConfig config.SkyWalkingConfig
	if bytes, err := ioutil.ReadFile("./skywalking/config.yaml"); err == nil {
		yaml.Unmarshal(bytes, &skywalkingConfig)
	} else {
		fmt.Println("解析本地skywalking配置文件失败!")
	}
	return &skywalkingConfig
}

// 创建skywalking exporter对象
func NewSkywalkingExporter() *SkyWalkingExporter {
	eventInfos := make([]EventInfo, 0)
	// skywalkingConfig := ParseSkyWalkingConfig()
	now := time.Now()
	beforeOneM := now.Add(time.Duration(-1000000000 * 60 * 480))
	year, month, day := beforeOneM.Date()
	eventIndex := fmt.Sprintf("sw_events-%04d%02d%02d", year, month, day)
	fmt.Printf("skywalking 采集当前时间: %04d-%02d-%02d %02d:%02d:%02d\n", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	skyEventDatas := make([]SkwEvent, 0)
	var typ SkwEvent
	events := GetAll(eventIndex, "_doc", typ)
	if events == nil {
		fmt.Println("没有查询到es数据...")
		utils.Logger.Printf("没有查询到es数据...\n")
		return nil
	}
	for _, event := range events {
		switch ret := event.(type) {
		case string:
			fmt.Println("event.(type): ", ret)
		case SkwEvent:
			typ := SkwEvent(ret)
			skyEventDatas = append(skyEventDatas, typ)
			fmt.Printf("event: %v \n", typ)
		default:
		}
		eventInfo := prometheus.NewDesc("event_info", "描述事件的详细信息",
			[]string{"name", "type", "service_name", "start_time", "end_time", "message", "start", "end", "time_bucket"},
			prometheus.Labels{})
		evnetInfoValType := prometheus.GaugeValue
		eventInfos = append(eventInfos, EventInfo{
			EventInfoDesc:    eventInfo,
			EventInfoValType: evnetInfoValType,
		})
	}

	serviceCPMs := make([]ServiceCpm, 0)
	cpmInfoDatas := GetCpmInfo("service_instance_cpm")
	for i := 0; i < len(cpmInfoDatas); i++ {
		cpminfo := prometheus.NewDesc("service_cpm", "服务的cpm",
			[]string{"service_name", "entity", "time_bucket", "cluster", "ip", "export_time_bucket"},
			prometheus.Labels{})
		cpminfoValType := prometheus.GaugeValue
		serviceCPMs = append(serviceCPMs, ServiceCpm{
			ServiceCpmDesc:    cpminfo,
			ServiceCpmValType: cpminfoValType,
		})

	}

	return &SkyWalkingExporter{
		EventInfos:    eventInfos,
		ServiceCpms:   serviceCPMs,
		SkyEventDatas: skyEventDatas,
		CpmDatas:      cpmInfoDatas,
	}
}

func (e *SkyWalkingExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, event := range e.EventInfos {
		ch <- event.EventInfoDesc
	}
}

// 收集skywalking事件方法
func (e *SkyWalkingExporter) Collect(ch chan<- prometheus.Metric) {
	if e = NewSkywalkingExporter(); e == nil {
		utils.Logger.Println("es数据为空")
		return
	}

	sort.Slice(e.CpmDatas, func(i, j int) bool {
		if e.CpmDatas[i].TimeBucket > e.CpmDatas[j].TimeBucket {
			return true
		} else {
			return false
		}
	})

	for _, serviceData := range e.CpmDatas {
		fmt.Println("......", serviceData.ServiceId, "*********", serviceData.EntityId, "*******", serviceData.TimeBucket)
	}

	skywalkingConfig := ParseSkyWalkingConfig()
	println("e address: %p", e, "**************************************")
	// 获取event数据
	eventInfoDatas := e.SkyEventDatas
	for idx, eventInfo := range e.EventInfos {
		// "name", "type", "service_name", "start_time", "end_time", "message"
		// se := SkyEvent(events[idx])
		start := time.Unix(int64(eventInfoDatas[idx].StartTime/1000), 0).Format("2006-01-02 15:04:05")
		end := time.Unix(int64(eventInfoDatas[idx].EndTime/1000), 0).Format("2006-01-02 15:04:05")
		time_bucket := time.Unix(int64(eventInfoDatas[idx].StartTime/1000), 0).Format("200601021504")
		ch <- prometheus.MustNewConstMetric(eventInfo.EventInfoDesc, eventInfo.EventInfoValType, 1,
			eventInfoDatas[idx].Name, eventInfoDatas[idx].Type, eventInfoDatas[idx].Service, fmt.Sprintf("%d", eventInfoDatas[idx].StartTime), fmt.Sprintf("%d,", eventInfoDatas[idx].EndTime), eventInfoDatas[idx].Message, start, end, time_bucket)
	}

	now := time.Now()
	nowTimeBucket := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	serviceDatas := e.CpmDatas
	for idx, cpmInfo := range e.ServiceCpms {
		// "service_id", "service_name", "entity_id", "entity", "time_bucket", "cluster", "ip"
		ch <- prometheus.MustNewConstMetric(cpmInfo.ServiceCpmDesc, cpmInfo.ServiceCpmValType,
			float64(serviceDatas[idx].Value), serviceDatas[idx].ServiceName,
			serviceDatas[idx].Entity, fmt.Sprintf("%d", serviceDatas[idx].TimeBucket),
			skywalkingConfig.Cluster.Name, "", nowTimeBucket)
	}
}
