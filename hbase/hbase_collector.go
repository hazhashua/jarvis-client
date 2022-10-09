package hbase

import (
	"fmt"
	"io/ioutil"
	"metric_exporter/config"
	"metric_exporter/utils"
	"net/http"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

type jmxHttpUrl struct {
	masterUrl         *string
	masterBackupUrls  *[]string
	regionserversUrls *[]string
}

type hmasterData struct {
	numRegionServers         int64
	numDeadRegionServers     int64
	ritCount                 int64
	ritCountOverThreshold    int64
	masterNumActiveHandler   int64
	masterReceivedBytes      int64
	masterSentBytes          int64
	masterNumOpenConnections int64
	cluster                  *string
	host                     *string
	ip                       *string
}

type tableData struct {
	namespace         string
	tableName         string
	regionCount       int64
	storefileCount    int64
	readRequestCount  int64
	writeRequestCount int64
	tableSize         int64
	regionServer      string
}

type regionData struct {
	blockCacheCountHitPercent   float32
	blockCacheExpressHitPercent float32
	numActiveHandler            int64
	receivedBytes               int64
	sentBytes                   int64
	numOpenConnections          int64
	authenticationFailures      int64
	authenticationSuccesses     int64
	readRequestCount            int64
	writeRequestCount           int64
	regionCount                 int64
	storeFileCount              int64
	slowGetCount                int64
	slowPutCount                int64
	slowDeleteCount             int64
	slowAppendCount             int64
	slowIncrementCount          int64
	fsReadTimeMax               int64
	fsWriteTimeMax              int64
	tableDatas                  []tableData
	cluster                     string
	host                        string
	ip                          string
}

type hbaseData struct {
	masterData  *hmasterData
	regionDatas []regionData
}

//hbase jmx数据采集转换存入
// hbase jmx URL
// master
// http://192.168.10.220:16010/jmx
// http://192.168.10.221:16010/jmx
// regionserver jmx url
// http://192.168.10.220:16030/jmx
// http://192.168.10.221:16030/jmx
// http://192.168.10.222:16030/jmx

// 解析hbase配置文件
func ParseHbaseConfig() *config.HbaseConfigure {
	hbase_config := new(config.HbaseConfigure)
	bytes, err := ioutil.ReadFile("./hbase/config.yaml")
	if err != nil {
		// fmt.Println("err: ", err.Error())
		utils.Logger.Printf("读取hbase配置文件出错 error:%s\n", err.Error())
	}
	err2 := yaml.Unmarshal(bytes, hbase_config)
	if err2 != nil {
		// fmt.Println("err2: ", err2.Error())
		utils.Logger.Printf("解析hbase配置文件出错 error:%s", err2.Error())
	}
	return hbase_config
}

func initUrl() (int, *jmxHttpUrl) {
	/*
		从配置中获取请求地址，及active master信息
	*/
	// bytes, err := ioutil.ReadFile("./hbase/config.yaml")
	// if err != nil {
	// 	fmt.Println("*****************************")
	// 	fmt.Println("err: ", err.Error())
	// }
	// hbase_config := new(config.HbaseConfigure)
	// err = yaml.Unmarshal(bytes, hbase_config)
	// if err != nil {
	// 	fmt.Println("err: ", err.Error())
	// }

	// hbase_config := ParseHbaseConfig()
	hbase_config := (utils.ConfigStruct.ConfigData[config.HBASE]).(config.HbaseConfigure)

	// 获取active的master
	var jmx_url, master_jmx_url string
	active_master_index := -1
	var master_backup_urls []string
	for idx, host := range hbase_config.Cluster.Hosts {
		jmx_url = fmt.Sprintf("http://%s:%s/jmx", host, hbase_config.Cluster.MasterJmxPort)
		master_url := fmt.Sprintf("http://%s:%s/master-status", host, hbase_config.Cluster.MasterJmxPort)
		utils.Logger.Printf("jmx_url: %s \t master_url: %s \n", jmx_url, master_url)
		response, err2 := http.Get(master_url)
		if err2 != nil {
			// fmt.Println("err2: ", err2.Error())
			utils.Logger.Printf("获取master_url失败 error:%s\n", err2.Error())
			continue
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			// handle error
			utils.Logger.Printf("ioutil.ReadAll(response.Body) error:%s", err.Error())
			panic(err)
		}
		body_str := string(body)
		if strings.Contains(body_str, "<title>Master:") {
			master_jmx_url = jmx_url
			active_master_index = idx
		} else {
			master_backup_urls = append(master_backup_urls, jmx_url)
		}
	}

	// masterUrl := "http://192.168.10.220:16010/jmx"
	regionserver_urls := make([]string, 0)
	for _, host := range hbase_config.Cluster.Hosts {
		regionserver_urls = append(regionserver_urls, fmt.Sprintf("http://%s:%s/jmx", host, hbase_config.Cluster.RegionserverJmxPort))
	}
	// regionserversUrl := [3]string{"http://124.65.131.14:16030/jmx", "http://124.65.131.14:16030/jmx", "http://124.65.131.14:16030/jmx"}
	return active_master_index, &jmxHttpUrl{
		masterUrl:         &master_jmx_url,
		masterBackupUrls:  &master_backup_urls,
		regionserversUrls: &regionserver_urls,
	}
}

func HttpRequest(is_master bool, jmx_http_url *jmxHttpUrl, uri string, region_no int) []byte {
	utils.Logger.Printf("config.masterUrl: %s\t config.masterBackupUrls: %s \t config.regionserversUrls: %s\n", *jmx_http_url.masterUrl, *jmx_http_url.masterBackupUrls, jmx_http_url.regionserversUrls)
	// fmt.Println((*jmx_http_url.masterUrl) + uri)
	var httpErr error
	var response *http.Response
	if is_master {
		// fmt.Println("master url: ", (*jmx_http_url.masterUrl)+uri)
		utils.Logger.Printf("master url: %s\n", (*jmx_http_url.masterUrl)+uri)
		response, httpErr = http.Get((*jmx_http_url.masterUrl) + uri)
	} else {
		// fmt.Println("regionserver url: ", (*jmx_http_url.regionserversUrls)[region_no-1]+uri)
		utils.Logger.Printf("regionserver url: %s\n", (*jmx_http_url.regionserversUrls)[region_no-1]+uri)
		response, httpErr = http.Get((*jmx_http_url.regionserversUrls)[region_no-1] + uri)
	}
	if response != nil {
		defer response.Body.Close()
	}
	if httpErr != nil {
		utils.Logger.Printf("http.Get error:%s \n", httpErr.Error())
		return []byte{}
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// handle error
		utils.Logger.Printf("ioutil.ReadAll(response.Body) error:%s\n", err.Error())
		panic(err)
	}
	//fmt.Println("------response:", string(body))
	return body
}

func QueryMetric() *hbaseData {
	var hmaster_data hmasterData
	// 查询master的特定指标
	query_url := fmt.Sprintf("?qry=%s", "Hadoop:service=HBase,name=Master,sub=Server")
	utils.Logger.Printf("query_url: %s", query_url)

	active_master_index, jmx_http_url := initUrl()
	if active_master_index == -1 {
		utils.Logger.Printf("initUrl() return -1\n")
		// 没有找到active master, 返回
		return &hbaseData{
			masterData:  nil,
			regionDatas: []regionData{},
		}
	}
	// fmt.Println(*jmx_http_url.masterUrl)
	// fmt.Println(*jmx_http_url.masterBackupUrls)
	// fmt.Println(jmx_http_url.regionserversUrls)

	// 获取集群主机的名称IP信息
	// hbase_config := ParseHbaseConfig()
	hbase_config := (utils.ConfigStruct.ConfigData[config.HBASE]).(config.HbaseConfigure)

	host := hbase_config.Cluster.Hosts[active_master_index]
	cluster := hbase_config.Cluster.ClusterName
	ip := hbase_config.Cluster.Names[active_master_index]
	hmaster_data.cluster = &cluster
	hmaster_data.host = &host
	hmaster_data.ip = &ip

	body := HttpRequest(true, jmx_http_url, query_url, 0)
	if mm, unmarshalErr := UnmarshalMasterMain(body); unmarshalErr == nil {
		hmaster_data.numRegionServers = *mm.Beans[0].NumRegionServers
		hmaster_data.numDeadRegionServers = *mm.Beans[0].NumDeadRegionServers
		utils.Logger.Printf("hmaster_data.numRegionServers: %d \t hmaster_data.numDeadRegionServers: %d \n", *mm.Beans[0].NumRegionServers, *mm.Beans[0].NumDeadRegionServers)
	}

	query_url = fmt.Sprintf("?qry=%s", "Hadoop:service=HBase,name=Master,sub=AssignmentManager")
	body = HttpRequest(true, jmx_http_url, query_url, 0)
	if assignment_manager, unmarshalErr := UnmarshalAssignmentManager(body); unmarshalErr == nil {
		hmaster_data.ritCount = *assignment_manager.Beans[0].RitCount
		hmaster_data.ritCountOverThreshold = *assignment_manager.Beans[0].RitCountOverThreshold
		// fmt.Println(*assignment_manager.Beans[0].RitCount)
		// fmt.Println(*assignment_manager.Beans[0].RitCountOverThreshold)
		utils.Logger.Printf("hmaster_data.ritCount: %d \t hmaster_data.ritCountOverThreshold: %d \n", *assignment_manager.Beans[0].RitCount, *assignment_manager.Beans[0].RitCountOverThreshold)
	}

	//Hadoop:service=HBase,name=Master,sub=IPC
	query_url = fmt.Sprintf("?qry=%s", "Hadoop:service=HBase,name=Master,sub=IPC")
	body = HttpRequest(true, jmx_http_url, query_url, 0)
	if master_ipc, unmarshalErr := UnmarshalMasterIPC(body); unmarshalErr == nil {
		hmaster_data.masterNumActiveHandler = *master_ipc.Beans[0].NumActiveHandler
		hmaster_data.masterReceivedBytes = *master_ipc.Beans[0].ReceivedBytes
		hmaster_data.masterSentBytes = *master_ipc.Beans[0].SentBytes
		hmaster_data.masterNumOpenConnections = *master_ipc.Beans[0].NumOpenConnections
		// fmt.Println(*master_ipc.Beans[0].NumActiveHandler)
		// // 接收的数据量
		// fmt.Println(*master_ipc.Beans[0].ReceivedBytes)
		// // 发送的数据量
		// fmt.Println(*master_ipc.Beans[0].SentBytes)
		// // 打开的ipc连接数
		// fmt.Println(*master_ipc.Beans[0].NumOpenConnections)
		utils.Logger.Printf("NumActiveHandler: %d\t ReceivedBytes: %d\t SentBytes: %d \t, NumOpenConnections: %d \n", *master_ipc.Beans[0].NumActiveHandler, *master_ipc.Beans[0].ReceivedBytes, *master_ipc.Beans[0].SentBytes, *master_ipc.Beans[0].NumOpenConnections)
	}

	utils.Logger.Printf("jmx_http_url.regionserversUrls: %s", jmx_http_url.regionserversUrls)
	region_num := len(*jmx_http_url.regionserversUrls)
	region_no := 1
	// var region_datas [3]regionData
	region_datas := make([]regionData, 0)
	// tdL := make([]*tableData, 0)
	for {
		// 查询所有regionserver的特定指标
		//?qry=Hadoop:service=HBase,name=RegionServer,sub=IPC
		var region_data regionData
		query_url = fmt.Sprintf("?qry=%s", "Hadoop:service=HBase,name=RegionServer,sub=IPC")
		utils.Logger.Printf("query url: %s\n", query_url)
		body = HttpRequest(false, jmx_http_url, query_url, region_no)
		// utils.Logger.Printf("response body: %s", string(body))
		if region_ipc, unmarshalErr := UnmarshalRegionserverIPC(body); unmarshalErr == nil {
			region_data.numActiveHandler = *region_ipc.Beans[0].NumActiveHandler
			region_data.receivedBytes = *region_ipc.Beans[0].ReceivedBytes
			region_data.sentBytes = *region_ipc.Beans[0].SentBytes
			region_data.numOpenConnections = *region_ipc.Beans[0].NumOpenConnections
			region_data.authenticationFailures = *region_ipc.Beans[0].AuthenticationFailures
			region_data.authenticationSuccesses = *region_ipc.Beans[0].AuthenticationSuccesses

			// fmt.Println("NumActiveHandler: ", *region_ipc.Beans[0].NumActiveHandler)
			// // 接收的数据量
			// fmt.Println("ReceivedBytes: ", *region_ipc.Beans[0].ReceivedBytes)
			// // 发送的数据量
			// fmt.Println("SentBytes: ", *region_ipc.Beans[0].SentBytes)
			// // 打开的ipc连接数
			// fmt.Println("NumOpenConnections: ", *region_ipc.Beans[0].NumOpenConnections)
			// // rpc认证失败次数
			// fmt.Println("AuthenticationFailures: ", *region_ipc.Beans[0].AuthenticationFailures)
			// // rpc认证成功次数
			// fmt.Println("AuthenticationSuccesses: ", *region_ipc.Beans[0].AuthenticationSuccesses)

			utils.Logger.Printf("NumActiveHandler: %d\t ReceivedBytes: %d\t SentBytes: %d\t NumOpenConnections: %d\t AuthenticationFailures: %d\t AuthenticationSuccesses: %d\n", *region_ipc.Beans[0].NumActiveHandler, *region_ipc.Beans[0].ReceivedBytes, *region_ipc.Beans[0].SentBytes, *region_ipc.Beans[0].NumOpenConnections, *region_ipc.Beans[0].AuthenticationFailures, *region_ipc.Beans[0].AuthenticationSuccesses)
		}

		//Hadoop:service=HBase,name=RegionServer,sub=Server
		query_url = fmt.Sprintf("?qry=%s", "Hadoop:service=HBase,name=RegionServer,sub=Server")
		body = HttpRequest(false, jmx_http_url, query_url, region_no)
		if region_server, unmarshalErr := UnmarshalRegionserverServer(body); unmarshalErr == nil {
			if len(region_server.Beans) != 0 {
				// jmx没有抓取到数据
				region_data.blockCacheCountHitPercent = *region_server.Beans[0].BlockCacheCountHitPercent
				region_data.blockCacheExpressHitPercent = *region_server.Beans[0].BlockCacheExpressHitPercent
				region_data.readRequestCount = *region_server.Beans[0].ReadRequestCount
				region_data.writeRequestCount = *region_server.Beans[0].WriteRequestCount
				region_data.regionCount = *region_server.Beans[0].RegionCount
				region_data.storeFileCount = *region_server.Beans[0].StoreFileCount
				region_data.slowGetCount = *region_server.Beans[0].SlowGetCount
				region_data.slowPutCount = *region_server.Beans[0].SlowPutCount
				region_data.slowDeleteCount = *region_server.Beans[0].SlowDeleteCount
				region_data.slowAppendCount = *region_server.Beans[0].SlowAppendCount
				region_data.slowIncrementCount = *region_server.Beans[0].SlowIncrementCount

				// // server的读请求数
				// fmt.Println(*region_server.Beans[0].ReadRequestCount)
				// // server的写请求数
				// fmt.Println(*region_server.Beans[0].WriteRequestCount)
				// // regionserver的region个数
				// fmt.Println(*region_server.Beans[0].RegionCount)
				// // regionserver的store file个数
				// fmt.Println(*region_server.Beans[0].StoreFileCount)
				// // regionserver的slow get count
				// fmt.Println(*region_server.Beans[0].SlowGetCount)
				// // regionserver的slow put count
				// fmt.Println(*region_server.Beans[0].SlowPutCount)
				// // regionserver的slow delete count
				// fmt.Println(*region_server.Beans[0].SlowDeleteCount)
				// // regionserver的slow delete count
				// fmt.Println(*region_server.Beans[0].SlowAppendCount)
				// // regionserver的slow delete count
				// fmt.Println(*region_server.Beans[0].SlowIncrementCount)

			} else {
				utils.Logger.Printf("jmx中数据为空......")
				region_data.blockCacheCountHitPercent = -1
				region_data.blockCacheExpressHitPercent = -1
				region_data.readRequestCount = -1
				region_data.writeRequestCount = -1
				region_data.regionCount = -1
				region_data.storeFileCount = -1
				region_data.slowGetCount = -1
				region_data.slowPutCount = -1
				region_data.slowDeleteCount = -1
				region_data.slowAppendCount = -1
				region_data.slowIncrementCount = -1

			}
		} else {
			utils.Logger.Printf("解析jmx:%s 数据出错    %s\n", "Hadoop:service=HBase,name=RegionServer,sub=Server", unmarshalErr.Error())
		}

		// Hadoop:service=HBase,name=RegionServer,sub=IO
		query_url = fmt.Sprintf("?qry=%s", "Hadoop:service=HBase,name=RegionServer,sub=IO")
		body = HttpRequest(false, jmx_http_url, query_url, region_no)
		if region_io, unmarshalErr := UnmarshalRegionserverIO(body); unmarshalErr == nil {
			region_data.fsReadTimeMax = *region_io.Beans[0].FSReadTimeMax
			region_data.fsWriteTimeMax = *region_io.Beans[0].FSWriteTimeMax
			// // 文件系统最大读时间
			// fmt.Println(*region_io.Beans[0].FSReadTimeMax)
			// // 文件系统最大写时间
			// fmt.Println(*region_io.Beans[0].FSWriteTimeMax)
		}

		// 解析hbase table相关的数据
		// Hadoop:service=HBase,name=RegionServer,sub=Tables
		query_url = fmt.Sprintf("?qry=%s", "Hadoop:service=HBase,name=RegionServer,sub=Tables")
		body = HttpRequest(false, jmx_http_url, query_url, region_no)
		fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
		if region_tables, unmarshalErr := UnmarshalTables(body); unmarshalErr == nil {
			if len(region_tables.Beans) != 0 {
				tds := make(map[string]*tableData, 0)
				var hostName, tableName string
				for key, value := range region_tables.Beans[0] {
					reg := regexp.MustCompile("Namespace_([^_]+?)_table_(.*)")
					rss := reg.FindSubmatch([]byte(key))
					// var tabled tableData
					if len(rss) == 3 {
						fmt.Printf("key: %s, namespace: %s   table: %s\n", key, rss[1], rss[2])
						splits := strings.Split(string(rss[2]), "_metric_")
						tableName = splits[0]
						fmt.Println("splits: ", splits)
						fmt.Println("tableName: ", tableName)
						if _, ok := tds[tableName]; !ok {
							tds[tableName] = &tableData{namespace: string(rss[1]), tableName: tableName}
							tableName = string(tableName)
						}
					}
					if key == "tag.Hostname" {
						hostName = *value.String
					}
					if idx := strings.Index(key, "_metric_"); idx != -1 {
						metric := key[idx+len("_metric_"):]
						switch metric {
						case "regionCount":
							tds[tableName].regionCount = *value.Integer
							utils.Logger.Printf("tds[%s].regionCount: %d \n", tableName, *value.Integer)
						case "storeFileCount":
							// tabled.storefileCount = *value.Integer
							tds[tableName].storefileCount = *value.Integer
							utils.Logger.Printf("tds[%s].storefileCount: %d \n", tableName, *value.Integer)
						case "readRequestCount":
							// tabled.readRequestCount = *value.Integer
							tds[tableName].readRequestCount = *value.Integer
							utils.Logger.Printf("tds[%s].readRequestCount: %d \n", tableName, *value.Integer)
						case "writeRequestCount":
							// tabled.writeRequestCount = *value.Integer
							tds[tableName].writeRequestCount = *value.Integer
							utils.Logger.Printf("tds[%s].writeRequestCount: %d \n", tableName, *value.Integer)
						case "tableSize":
							// tabled.writeRequestCount = *value.Integer
							tds[tableName].tableSize = *value.Integer
							utils.Logger.Printf("tds[%s].tableSize: %d \n", tableName, *value.Integer)
						default:
							// utils.Logger.Printf("忽略的hbase table相关指标: %s \n", metric)
						}
					}
				}
				if hostName != "" {
					for _, value := range tds {
						value.regionServer = hostName
						utils.Logger.Printf("table info: %v\n", value)
						region_data.tableDatas = append(region_data.tableDatas, *value)
					}
				}
			}
		}

		// cluster, host, ip := "cluster1", fmt.Sprintf("dev%02d", region_no), "192.168.10.220"
		cluster = hbase_config.Cluster.ClusterName
		host = hbase_config.Cluster.Names[region_no-1]
		ip = hbase_config.Cluster.Hosts[region_no-1]
		utils.Logger.Printf("cluster:%s host:%s ip:%s\n", cluster, host, ip)
		region_data.cluster = cluster
		region_data.host = host
		region_data.ip = ip

		region_datas = append(region_datas, region_data)
		region_no += 1
		if region_no > region_num {
			break
		}

	}
	return &hbaseData{
		masterData:  &hmaster_data,
		regionDatas: region_datas,
	}

}

type hbaseCollector struct {
	masterMetrics        hbaseMasterMetric
	regionMetrics        []hbaseRegionMetric
	regionDynamicMetrics [][]hbaseRegionDynamicMetric
	datas                hbaseData
}

type hbaseRegionMetric struct {
	// 下面是regionserver的指标
	BlockCacheCountHitPercent          *prometheus.Desc
	BlockCacheCountHitPercentValType   prometheus.ValueType
	BlockCacheExpressHitPercent        *prometheus.Desc
	BlockCacheExpressHitPercentValType prometheus.ValueType
	NumActiveHandler                   *prometheus.Desc
	NumActiveHandlerValType            prometheus.ValueType
	ReceivedBytes                      *prometheus.Desc
	ReceivedBytesValType               prometheus.ValueType
	SentBytes                          *prometheus.Desc
	SentBytesValType                   prometheus.ValueType
	NumOpenConnections                 *prometheus.Desc
	NumOpenConnectionsValType          prometheus.ValueType
	AuthenticationFailures             *prometheus.Desc
	AuthenticationFailuresValType      prometheus.ValueType
	AuthenticationSuccesses            *prometheus.Desc
	AuthenticationSuccessesValType     prometheus.ValueType
	ReadRequestCount                   *prometheus.Desc
	ReadRequestCountValType            prometheus.ValueType
	WriteRequestCount                  *prometheus.Desc
	WriteRequestCountValType           prometheus.ValueType
	RegionCount                        *prometheus.Desc
	RegionCountValType                 prometheus.ValueType
	StoreFileCount                     *prometheus.Desc
	StoreFileCountValType              prometheus.ValueType
	SlowGetCount                       *prometheus.Desc
	SlowGetCountValType                prometheus.ValueType
	SlowPutCount                       *prometheus.Desc
	SlowPutCountValType                prometheus.ValueType
	SlowDeleteCount                    *prometheus.Desc
	SlowDeleteCountValType             prometheus.ValueType
	SlowAppendCount                    *prometheus.Desc
	SlowAppendCountValType             prometheus.ValueType
	SlowIncrementCount                 *prometheus.Desc
	SlowIncrementCountValType          prometheus.ValueType
	FSReadTimeMax                      *prometheus.Desc
	FSReadTimeMaxValType               prometheus.ValueType
	FSWriteTimeMax                     *prometheus.Desc
	FSWriteTimeMaxValType              prometheus.ValueType
}

type hbaseRegionDynamicMetric struct {
	TableInfo        *prometheus.Desc
	TableInfoValType prometheus.ValueType
}

type hbaseMasterMetric struct {
	// 下面是master的指标
	NumRegionServers                *prometheus.Desc
	NumRegionServersValType         prometheus.ValueType
	NumDeadRegionServers            *prometheus.Desc
	NumDeadRegionServersValueType   prometheus.ValueType
	RitCount                        *prometheus.Desc
	RitCountValType                 prometheus.ValueType
	RitCountOverThreshold           *prometheus.Desc
	RitCountOverThresholdValType    prometheus.ValueType
	MasterNumActiveHandler          *prometheus.Desc
	MasterNumActiveHandlerValType   prometheus.ValueType
	MasterReceivedBytes             *prometheus.Desc
	MasterReceivedBytesValType      prometheus.ValueType
	MasterSentBytes                 *prometheus.Desc
	MasterSentBytesValType          prometheus.ValueType
	MasterNumOpenConnections        *prometheus.Desc
	MasterNumOpenConnectionsValType prometheus.ValueType
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func NewHbaseCollector() *hbaseCollector {
	var master_metrics hbaseMasterMetric
	var regionMetricList []hbaseRegionMetric
	_, jmx_http_url := initUrl()
	fmt.Println("jmx_http_url: ", jmx_http_url)
	regionserver_num := len(*jmx_http_url.regionserversUrls)
	fmt.Println("region_num: ", regionserver_num)

	hbasedatas := QueryMetric()
	region_dynamicss := make([][]hbaseRegionDynamicMetric, 0)

	for length := 0; length < regionserver_num; length++ {
		// var service_alive_collector hbaseRegionMetric
		var region_metrics hbaseRegionMetric
		region_metrics.BlockCacheCountHitPercent = prometheus.NewDesc("blockcache_count_hit_percent", "show the hit percent of the blockcache to all read request",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.BlockCacheCountHitPercentValType = prometheus.GaugeValue

		region_metrics.BlockCacheExpressHitPercent = prometheus.NewDesc("blockcache_express_hit_percent", "show the hit percent of the blockcache to the request to the cache",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.BlockCacheExpressHitPercentValType = prometheus.GaugeValue

		region_metrics.NumActiveHandler = prometheus.NewDesc("num_active_handler", "Show active handler's num",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.NumActiveHandlerValType = prometheus.GaugeValue
		region_metrics.NumOpenConnections = prometheus.NewDesc("num_open_connections", "show open connection's num",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.NumOpenConnectionsValType = prometheus.GaugeValue

		region_metrics.SentBytes = prometheus.NewDesc("sent_bytes", "this regionserver sent's bytes",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.SentBytesValType = prometheus.GaugeValue

		region_metrics.ReceivedBytes = prometheus.NewDesc("receive_bytes", "this regionserver recive's bytes",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.ReceivedBytesValType = prometheus.GaugeValue

		region_metrics.AuthenticationFailures = prometheus.NewDesc("authentication_failures", "authentication failure counts",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.AuthenticationFailuresValType = prometheus.CounterValue

		region_metrics.AuthenticationSuccesses = prometheus.NewDesc("authentication_successes", "authentication success counts",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.AuthenticationSuccessesValType = prometheus.CounterValue

		region_metrics.ReadRequestCount = prometheus.NewDesc("read_request_count", "the regionserver's read request count",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.ReadRequestCountValType = prometheus.CounterValue

		region_metrics.WriteRequestCount = prometheus.NewDesc("write_request_count", "the regionserver's write request count",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.WriteRequestCountValType = prometheus.CounterValue

		//RegionCount
		region_metrics.RegionCount = prometheus.NewDesc("region_count", "the regionserver's region count",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.RegionCountValType = prometheus.GaugeValue

		region_metrics.StoreFileCount = prometheus.NewDesc("storefile_count", "the regionserver's storefile count",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.StoreFileCountValType = prometheus.GaugeValue

		region_metrics.SlowAppendCount = prometheus.NewDesc("slow_append_count", "the regionserver's slow append count",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.StoreFileCountValType = prometheus.CounterValue

		region_metrics.SlowDeleteCount = prometheus.NewDesc("slow_delete_count", "the regionserver's slow delete count",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.SlowDeleteCountValType = prometheus.CounterValue

		region_metrics.SlowGetCount = prometheus.NewDesc("slow_get_count", "the regionserver's slow get count",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.SlowGetCountValType = prometheus.CounterValue

		region_metrics.SlowIncrementCount = prometheus.NewDesc("slow_increment_count", "the regionserver's slow increment count",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.SlowIncrementCountValType = prometheus.CounterValue

		region_metrics.SlowPutCount = prometheus.NewDesc("slow_put_count", "the regionserver's slow put count",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.SlowPutCountValType = prometheus.CounterValue

		region_metrics.FSReadTimeMax = prometheus.NewDesc("fs_readtime_max", "the regionserver's fs read time max",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.FSReadTimeMaxValType = prometheus.GaugeValue

		region_metrics.FSWriteTimeMax = prometheus.NewDesc("fs_writetime_max", "the regionserver's fs write time max",
			[]string{"cluster", "host", "ip"},
			prometheus.Labels{})
		region_metrics.FSWriteTimeMaxValType = prometheus.GaugeValue
		regionMetricList = append(regionMetricList, region_metrics)

		// namespace         string
		// tableName         string
		// regionCount       int64
		// storefileCount    int64
		// readRequestCount  int64
		// writeRequestCount int64
		// tableSize         int64
		// regionServer      string
		region_dynamic_metricL := make([]hbaseRegionDynamicMetric, 0)
		// 集群状态正常时，创建table_info desc
		if len(hbasedatas.regionDatas) > length {
			for i := 0; i < len(hbasedatas.regionDatas[length].tableDatas); i++ {
				var region_dynamic_metrics hbaseRegionDynamicMetric
				region_dynamic_metrics.TableInfo = prometheus.NewDesc("table_info", "the table info on every regionservers",
					[]string{"cluster", "namespace", "table_name", "region_count", "storefile_count", "read_request_count", "write_request_count", "table_size", "regionserver"},
					prometheus.Labels{})
				region_dynamic_metrics.TableInfoValType = prometheus.GaugeValue
				region_dynamic_metricL = append(region_dynamic_metricL, region_dynamic_metrics)
			}
			region_dynamicss = append(region_dynamicss, region_dynamic_metricL)
		}

	}

	fmt.Println("regionserver metric init over......")
	master_metrics.NumRegionServers = prometheus.NewDesc("num_regionservers", "the num of regionserver",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	master_metrics.NumDeadRegionServersValueType = prometheus.GaugeValue

	master_metrics.NumDeadRegionServers = prometheus.NewDesc("num_dead_regionservers", "the num of dead regionservers",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	master_metrics.NumDeadRegionServersValueType = prometheus.GaugeValue

	master_metrics.RitCount = prometheus.NewDesc("rit_count", "the num of regionserver which in  rit status",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	master_metrics.RitCountValType = prometheus.GaugeValue

	master_metrics.RitCountOverThreshold = prometheus.NewDesc("rit_count_over_threshold", "the threshold of regionserver in  rit status",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	master_metrics.RitCountOverThresholdValType = prometheus.GaugeValue

	master_metrics.MasterNumActiveHandler = prometheus.NewDesc("master_num_active_handler", "the num of active handler of master",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	master_metrics.MasterNumActiveHandlerValType = prometheus.GaugeValue

	master_metrics.MasterReceivedBytes = prometheus.NewDesc("master_received_bytes", "the received bytes of master",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	master_metrics.MasterReceivedBytesValType = prometheus.GaugeValue

	master_metrics.MasterSentBytes = prometheus.NewDesc("master_sent_bytes", "the sent bytes of master",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	master_metrics.MasterSentBytesValType = prometheus.GaugeValue

	master_metrics.MasterNumOpenConnections = prometheus.NewDesc("master_num_open_connections", "the num of open connections of master",
		[]string{"cluster", "host", "ip"},
		prometheus.Labels{})
	master_metrics.MasterNumOpenConnectionsValType = prometheus.GaugeValue

	return &hbaseCollector{
		masterMetrics:        master_metrics,
		regionMetrics:        regionMetricList,
		regionDynamicMetrics: region_dynamicss,
		datas:                *hbasedatas,
	}

}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *hbaseCollector) Describe(ch chan<- *prometheus.Desc) {
	for idx, metric := range collector.regionMetrics {
		ch <- metric.BlockCacheCountHitPercent
		ch <- metric.BlockCacheExpressHitPercent
		ch <- metric.NumActiveHandler
		ch <- metric.ReceivedBytes
		ch <- metric.SentBytes
		ch <- metric.NumOpenConnections
		ch <- metric.AuthenticationFailures
		ch <- metric.AuthenticationSuccesses
		ch <- metric.ReadRequestCount
		ch <- metric.WriteRequestCount
		ch <- metric.RegionCount
		ch <- metric.StoreFileCount
		ch <- metric.SlowGetCount
		ch <- metric.SlowDeleteCount
		ch <- metric.SlowAppendCount
		ch <- metric.SlowPutCount
		ch <- metric.SlowIncrementCount
		ch <- metric.FSReadTimeMax
		ch <- metric.FSWriteTimeMax
		if len(collector.regionDynamicMetrics) > idx {
			for _, metric := range collector.regionDynamicMetrics[idx] {
				ch <- metric.TableInfo
			}
		}
	}
	ch <- collector.masterMetrics.NumRegionServers
	ch <- collector.masterMetrics.NumDeadRegionServers
	ch <- collector.masterMetrics.RitCount
	ch <- collector.masterMetrics.RitCountOverThreshold
	ch <- collector.masterMetrics.MasterNumActiveHandler
	ch <- collector.masterMetrics.MasterReceivedBytes
	ch <- collector.masterMetrics.MasterSentBytes
	ch <- collector.masterMetrics.MasterNumOpenConnections

}

//Collect implements required collect function for all promehteus collectors
func (collector *hbaseCollector) Collect(ch chan<- prometheus.Metric) {

	collector = NewHbaseCollector()
	hbase_data := collector.datas

	if hbase_data.masterData == nil || len(hbase_data.regionDatas) == 0 {
		// 集群异常，数据为空，直接返回
		return
	}

	// for _, alive := range da {
	for index, region_info := range collector.regionMetrics {
		// ch <- prometheus.MustNewConstMetric(alive.aliveMetric, prometheus.GaugeValue, float64(da[index].MetricValue), *da[index].ClusterName, *da[index].ServiceName, *da[index].ChildService, *da[index].IP, fmt.Sprintf("%d", da[index].Port), *da[index].PortType)
		//NumActiveHandler
		ch <- prometheus.MustNewConstMetric(region_info.BlockCacheCountHitPercent, prometheus.GaugeValue, float64(hbase_data.regionDatas[index].blockCacheCountHitPercent), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.BlockCacheExpressHitPercent, prometheus.GaugeValue, float64(hbase_data.regionDatas[index].blockCacheExpressHitPercent), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.NumActiveHandler, prometheus.GaugeValue, float64(hbase_data.regionDatas[index].numActiveHandler), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.ReceivedBytes, prometheus.GaugeValue, float64(hbase_data.regionDatas[index].receivedBytes), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.SentBytes, prometheus.GaugeValue, float64(hbase_data.regionDatas[index].sentBytes), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.NumOpenConnections, prometheus.GaugeValue, float64(hbase_data.regionDatas[index].numOpenConnections), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.AuthenticationFailures, prometheus.CounterValue, float64(hbase_data.regionDatas[index].authenticationFailures), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.AuthenticationSuccesses, prometheus.CounterValue, float64(hbase_data.regionDatas[index].authenticationSuccesses), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.ReadRequestCount, prometheus.CounterValue, float64(hbase_data.regionDatas[index].readRequestCount), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.WriteRequestCount, prometheus.CounterValue, float64(hbase_data.regionDatas[index].writeRequestCount), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.RegionCount, prometheus.GaugeValue, float64(hbase_data.regionDatas[index].regionCount), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.StoreFileCount, prometheus.GaugeValue, float64(hbase_data.regionDatas[index].storeFileCount), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.SlowGetCount, prometheus.CounterValue, float64(hbase_data.regionDatas[index].slowGetCount), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.SlowPutCount, prometheus.CounterValue, float64(hbase_data.regionDatas[index].slowPutCount), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.SlowDeleteCount, prometheus.CounterValue, float64(hbase_data.regionDatas[index].slowDeleteCount), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.SlowAppendCount, prometheus.CounterValue, float64(hbase_data.regionDatas[index].slowAppendCount), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.SlowIncrementCount, prometheus.CounterValue, float64(hbase_data.regionDatas[index].slowIncrementCount), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.FSWriteTimeMax, prometheus.GaugeValue, float64(hbase_data.regionDatas[index].fsWriteTimeMax), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)
		ch <- prometheus.MustNewConstMetric(region_info.FSReadTimeMax, prometheus.GaugeValue, float64(hbase_data.regionDatas[index].fsReadTimeMax), hbase_data.regionDatas[index].cluster, hbase_data.regionDatas[index].host, hbase_data.regionDatas[index].ip)

		for idx, dynamicMetric := range collector.regionDynamicMetrics[index] {
			// "namespace", "table_name", "region_count", "storefile_count", "read_request_count", "write_request_count", "table_size", "regionserver"
			ch <- prometheus.MustNewConstMetric(dynamicMetric.TableInfo, dynamicMetric.TableInfoValType, 1,
				hbase_data.regionDatas[index].cluster,
				hbase_data.regionDatas[index].tableDatas[idx].namespace, hbase_data.regionDatas[index].tableDatas[idx].tableName,
				fmt.Sprintf("%d", hbase_data.regionDatas[index].tableDatas[idx].regionCount),
				fmt.Sprintf("%d", hbase_data.regionDatas[index].tableDatas[idx].storefileCount),
				fmt.Sprintf("%d", hbase_data.regionDatas[index].tableDatas[idx].readRequestCount),
				fmt.Sprintf("%d", hbase_data.regionDatas[index].tableDatas[idx].writeRequestCount),
				fmt.Sprintf("%d", hbase_data.regionDatas[index].tableDatas[idx].tableSize),
				hbase_data.regionDatas[index].tableDatas[idx].regionServer)
		}
	}
	// NumRegionServers                *prometheus.Desc
	// NumDeadRegionServers            *prometheus.Desc
	// RitCount                        *prometheus.Desc
	// RitCountOverThreshold           *prometheus.Desc
	// MasterNumActiveHandler          *prometheus.Desc
	// MasterReceivedBytes             *prometheus.Desc
	// MasterSentBytes                 *prometheus.Desc
	// MasterNumOpenConnections        *prometheus.Desc

	ch <- prometheus.MustNewConstMetric(collector.masterMetrics.NumRegionServers, prometheus.GaugeValue, float64(hbase_data.masterData.numRegionServers), *hbase_data.masterData.cluster, *hbase_data.masterData.host, *hbase_data.masterData.ip)
	ch <- prometheus.MustNewConstMetric(collector.masterMetrics.NumDeadRegionServers, prometheus.GaugeValue, float64(hbase_data.masterData.numDeadRegionServers), *hbase_data.masterData.cluster, *hbase_data.masterData.host, *hbase_data.masterData.ip)
	ch <- prometheus.MustNewConstMetric(collector.masterMetrics.RitCount, prometheus.GaugeValue, float64(hbase_data.masterData.ritCount), *hbase_data.masterData.cluster, *hbase_data.masterData.host, *hbase_data.masterData.ip)
	ch <- prometheus.MustNewConstMetric(collector.masterMetrics.RitCountOverThreshold, prometheus.GaugeValue, float64(hbase_data.masterData.ritCountOverThreshold), *hbase_data.masterData.cluster, *hbase_data.masterData.host, *hbase_data.masterData.ip)
	ch <- prometheus.MustNewConstMetric(collector.masterMetrics.MasterNumActiveHandler, prometheus.GaugeValue, float64(hbase_data.masterData.masterNumActiveHandler), *hbase_data.masterData.cluster, *hbase_data.masterData.host, *hbase_data.masterData.ip)
	ch <- prometheus.MustNewConstMetric(collector.masterMetrics.MasterNumOpenConnections, prometheus.GaugeValue, float64(hbase_data.masterData.masterNumOpenConnections), *hbase_data.masterData.cluster, *hbase_data.masterData.host, *hbase_data.masterData.ip)
	ch <- prometheus.MustNewConstMetric(collector.masterMetrics.MasterReceivedBytes, prometheus.GaugeValue, float64(hbase_data.masterData.masterReceivedBytes), *hbase_data.masterData.cluster, *hbase_data.masterData.host, *hbase_data.masterData.ip)
	ch <- prometheus.MustNewConstMetric(collector.masterMetrics.MasterSentBytes, prometheus.GaugeValue, float64(hbase_data.masterData.masterSentBytes), *hbase_data.masterData.cluster, *hbase_data.masterData.host, *hbase_data.masterData.ip)

}
