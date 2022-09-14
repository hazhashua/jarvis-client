package hadoop

import (
	"fmt"
	"io/ioutil"
	"metric_exporter/config"
	"metric_exporter/utils"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

type HadoopMetric struct {
	// 下面是master的指标
	// ServiceStatus        []*prometheus.Desc
	// ServiceStatusValType prometheus.ValueType
	// hdfs块大小
	BlockSize        *prometheus.Desc
	BlockSizeValType prometheus.ValueType
	// hdfs副本数
	ReplicationNum        *prometheus.Desc
	ReplicationNumValType prometheus.ValueType

	NumLiveDataNodes        *prometheus.Desc
	NumLiveDataNodesValType prometheus.ValueType
	NumDeadDataNodes        *prometheus.Desc
	NumDeadDataNodesValType prometheus.ValueType
	NumLiveNameNodes        *prometheus.Desc
	NumLiveNameNodesValType prometheus.ValueType
	NumDeadNameNodes        *prometheus.Desc
	NumDeadNameNodesValType prometheus.ValueType
	// jvm内存的相关指标
	MemNonHeapUsedM             []*prometheus.Desc
	MemNonHeapUsedMValType      prometheus.ValueType
	MemNonHeapCommittedM        []*prometheus.Desc
	MemNonHeapCommittedMValType prometheus.ValueType
	MemHeapUsedM                []*prometheus.Desc
	MemHeapUsedMValType         prometheus.ValueType
	MemHeapCommittedM           []*prometheus.Desc
	MemHeapCommittedMValType    prometheus.ValueType
	// hadoop应用提交的相关信息
	AppsSubmitted        []*prometheus.Desc
	AppsSubmittedValType prometheus.ValueType
	AppsRunning          []*prometheus.Desc
	AppsRunningValType   prometheus.ValueType
	AppsCompleted        []*prometheus.Desc
	AppsCompletedValType prometheus.ValueType

	AppsPending       []*prometheus.Desc
	AppsPendinValType prometheus.ValueType
	AppsKilled        []*prometheus.Desc
	AppsKilledValType prometheus.ValueType
	AppsFailed        []*prometheus.Desc
	AppsFailedValType prometheus.ValueType

	// 运行时长0～60分钟的任务数
	Running0        []*prometheus.Desc
	Running0ValType prometheus.ValueType
	// 运行时长60～300分钟的任务数
	Running60        []*prometheus.Desc
	Running60ValType prometheus.ValueType
	// 运行时长300～1440分钟的任务数
	Running300        []*prometheus.Desc
	Running300ValType prometheus.ValueType

	//hdfs相关的指标
	CapacityTotalGB            *prometheus.Desc
	CapacityTotalGBValType     prometheus.ValueType
	CapacityUsedGB             *prometheus.Desc
	CapacityUsedGBValType      prometheus.ValueType
	CapacityRemainingGB        *prometheus.Desc
	CapacityRemainingGBValType prometheus.ValueType
	BlocksTotal                *prometheus.Desc
	BlocksTotalValType         prometheus.ValueType
	FilesTotal                 *prometheus.Desc
	FilesTotalValType          prometheus.ValueType
	FileInfoOps                []*prometheus.Desc
	FileInfoOpsValType         prometheus.ValueType
	CreateFileOps              []*prometheus.Desc
	CreateFileOpsValType       prometheus.ValueType
	GetlistingOps              []*prometheus.Desc
	GetlistingOpsValType       prometheus.ValueType
	DeleteFileOps              []*prometheus.Desc
	DeleteFileOpsValType       prometheus.ValueType
	// 请求连接相关的指标
	RpcQueueTimeAvgTime        []*prometheus.Desc
	RpcQueueTimeAvgTimeValType prometheus.ValueType
	NumOpenConnections         []*prometheus.Desc
	NumOpenConnectionsValType  prometheus.ValueType
	RpcSlowCalls               []*prometheus.Desc
	RpcSlowCallsValType        prometheus.ValueType
	CallQueueLength            []*prometheus.Desc
	CallQueueLengthValType     prometheus.ValueType

	RpcAuthorizationFailures         []*prometheus.Desc
	RpcAuthorizationFailuresValType  prometheus.ValueType
	RpcAuthorizationSuccesses        []*prometheus.Desc
	RpcAuthorizationSuccessesValType prometheus.ValueType

	// datanode相关的指标
	BytesWritten              []*prometheus.Desc
	BytesWrittenValType       prometheus.ValueType
	BytesRead                 []*prometheus.Desc
	BytesReadValType          prometheus.ValueType
	RemoteBytesWritten        []*prometheus.Desc
	RemoteBytesWrittenValType prometheus.ValueType
	RemoteBytesRead           []*prometheus.Desc
	RemoteBytesReadValType    prometheus.ValueType
	HeartbeatsNum             []*prometheus.Desc
	HeartbeatsNumValType      prometheus.ValueType
	HeartbeatsAvgTime         []*prometheus.Desc
	HeartbeatsAvgTimeValType  prometheus.ValueType
}

type HadoopCollector struct {
	hadoopMetrics HadoopMetric
}

func Parse_hadoop_config() *config.HadoopConfigure {
	bytes, _ := ioutil.ReadFile("./hadoop/config.yaml")
	hadoopConfig := new(config.HadoopConfigure)
	err := yaml.Unmarshal(bytes, hadoopConfig)
	if err != nil {
		fmt.Println("Unmarshal failed: ", err)
	}
	utils.Logger.Printf("hadoopConfig.Cluster.ServiceNum: %d\n", hadoopConfig.Cluster.ServiceNum)
	utils.Logger.Printf("hadoopConfig.Cluster.Services:   %s\n", hadoopConfig.Cluster.Services)
	// fmt.Println("hadoopConfig.Cluster.ServiceNum: ", hadoopConfig.Cluster.ServiceNum)
	// fmt.Println("hadoopConfig.Cluster.Services: ", hadoopConfig.Cluster.Services)
	return hadoopConfig
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func NewHadoopCollector() *HadoopCollector {
	var hadoop_metrics HadoopMetric

	// 使用全局配置文件
	hadoop_config := (utils.ConfigStruct.ConfigData["hadoop"]).(config.HadoopConfigure)
	// hadoop_config := Parse_hadoop_config()

	// for i := 0; i < hadoop_config.Cluster.ServiceNum; i++ {
	// 	hadoop_metrics.ServiceStatus = append(hadoop_metrics.ServiceStatus, prometheus.NewDesc("service_status", "show service status of hadoop cluster",
	// 		[]string{"cluster", "host", "ip", "port", "service_name"},
	// 		prometheus.Labels{}))
	// }
	// hadoop_metrics.ServiceStatusValType = prometheus.GaugeValue
	hadoop_metrics.BlockSize = prometheus.NewDesc("hdfs_block_size", "hdfs块大小,单位MB",
		[]string{"cluster"},
		prometheus.Labels{})
	hadoop_metrics.BlockSizeValType = prometheus.GaugeValue

	hadoop_metrics.ReplicationNum = prometheus.NewDesc("hdfs_replication_num", "每个hdfs块的副本个数",
		[]string{"cluster"},
		prometheus.Labels{})
	hadoop_metrics.ReplicationNumValType = prometheus.GaugeValue

	hadoop_metrics.NumLiveDataNodes = prometheus.NewDesc("num_live_datanodes", "存活的datanode个数",
		[]string{"cluster"},
		prometheus.Labels{})
	hadoop_metrics.NumLiveDataNodesValType = prometheus.GaugeValue
	hadoop_metrics.NumDeadDataNodes = prometheus.NewDesc("num_dead_datanodes", "离线的datanode个数",
		[]string{"cluster"},
		prometheus.Labels{})
	hadoop_metrics.NumDeadDataNodesValType = prometheus.GaugeValue
	hadoop_metrics.NumLiveNameNodes = prometheus.NewDesc("num_live_namenodes", "存活的namenode个数",
		[]string{"cluster"},
		prometheus.Labels{})
	hadoop_metrics.NumLiveNameNodesValType = prometheus.GaugeValue
	hadoop_metrics.NumDeadNameNodes = prometheus.NewDesc("num_dead_namenodes", "离线的namenode个数",
		[]string{"cluster"},
		prometheus.Labels{})
	hadoop_metrics.NumDeadNameNodesValType = prometheus.GaugeValue

	// datanode_http_port := hadoop_config.Cluster.DatanodeHttpPort
	// namenode_http_port := hadoop_config.Cluster.NamenodeHttpPort
	// datanode_rpc_port := hadoop_config.Cluster.DatanodeRpcPort
	// namenode_rpc_port := hadoop_config.Cluster.NamenodeRpcPort
	datanode_list := hadoop_config.Cluster.Datanodes
	namenode_list := hadoop_config.Cluster.Namenodes

	// 添加resource manager的jvm指标
	hadoop_metrics.MemHeapUsedM = append(hadoop_metrics.MemHeapUsedM, prometheus.NewDesc("mem_heap_used_memory", "当前已使用的堆栈内存大小",
		[]string{"cluster", "host", "ip", "port", "service_name"},
		prometheus.Labels{}))
	hadoop_metrics.MemHeapCommittedM = append(hadoop_metrics.MemHeapCommittedM, prometheus.NewDesc("mem_heap_committed_memory", "当前已提交的堆栈内存大小",
		[]string{"cluster", "host", "ip", "port", "service_name"},
		prometheus.Labels{}))
	hadoop_metrics.MemNonHeapUsedM = append(hadoop_metrics.MemNonHeapUsedM, prometheus.NewDesc("mem_non_heap_used_momory", "当前已使用非堆栈内存大小",
		[]string{"cluster", "host", "ip", "port", "service_name"},
		prometheus.Labels{}))
	hadoop_metrics.MemNonHeapCommittedM = append(hadoop_metrics.MemNonHeapCommittedM, prometheus.NewDesc("mem_non_heap_committed_momory", "当前已提交非堆栈内存大小",
		[]string{"cluster", "host", "ip", "port", "service_name"},
		prometheus.Labels{}))

	// 添加datanode jvm指标数据
	for i := 0; i < len(datanode_list); i++ {
		hadoop_metrics.MemHeapUsedM = append(hadoop_metrics.MemHeapUsedM, prometheus.NewDesc("mem_heap_used_memory", "当前已使用的堆栈内存大小",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.MemHeapUsedMValType = prometheus.GaugeValue
		hadoop_metrics.MemHeapCommittedM = append(hadoop_metrics.MemHeapCommittedM, prometheus.NewDesc("mem_heap_committed_memory", "当前已提交的堆栈内存大小",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.MemHeapCommittedMValType = prometheus.GaugeValue
		hadoop_metrics.MemNonHeapUsedM = append(hadoop_metrics.MemNonHeapUsedM, prometheus.NewDesc("mem_non_heap_used_momory", "当前已使用非堆栈内存大小",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.MemNonHeapUsedMValType = prometheus.GaugeValue
		hadoop_metrics.MemNonHeapCommittedM = append(hadoop_metrics.MemNonHeapCommittedM, prometheus.NewDesc("mem_non_heap_committed_momory", "当前已提交非堆栈内存大小",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.MemNonHeapCommittedMValType = prometheus.GaugeValue
	}

	// 添加namenode jvm指标数据
	for i := 0; i < len(namenode_list); i++ {
		hadoop_metrics.MemHeapUsedM = append(hadoop_metrics.MemHeapUsedM, prometheus.NewDesc("mem_heap_used_memory", "当前已使用的堆栈内存大小",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.MemHeapUsedMValType = prometheus.GaugeValue
		hadoop_metrics.MemHeapCommittedM = append(hadoop_metrics.MemHeapCommittedM, prometheus.NewDesc("mem_heap_committed_memory", "当前已提交的堆栈内存大小",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.MemHeapCommittedMValType = prometheus.GaugeValue
		hadoop_metrics.MemNonHeapUsedM = append(hadoop_metrics.MemNonHeapUsedM, prometheus.NewDesc("mem_non_heap_used_momory", "当前已使用非堆栈内存大小",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.MemNonHeapUsedMValType = prometheus.GaugeValue
		hadoop_metrics.MemNonHeapCommittedM = append(hadoop_metrics.MemNonHeapCommittedM, prometheus.NewDesc("mem_non_heap_committed_momory", "当前已提交非堆栈内存大小",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.MemNonHeapCommittedMValType = prometheus.GaugeValue
	}

	// 每个队列都有提交的应用信息
	queue_num := 1
	for i := 0; i < queue_num; i++ {
		hadoop_metrics.AppsSubmitted = append(hadoop_metrics.AppsSubmitted, prometheus.NewDesc("apps_submitted", "提交应用个数",
			[]string{"cluster", "queue"},
			prometheus.Labels{}))
		hadoop_metrics.AppsSubmittedValType = prometheus.GaugeValue
		hadoop_metrics.AppsCompleted = append(hadoop_metrics.AppsCompleted, prometheus.NewDesc("apps_completed", "完成应用个数",
			[]string{"cluster", "queue"},
			prometheus.Labels{}))
		hadoop_metrics.AppsCompletedValType = prometheus.GaugeValue
		hadoop_metrics.AppsRunning = append(hadoop_metrics.AppsRunning, prometheus.NewDesc("apps_running", "正在运行的任务个数",
			[]string{"cluster", "queue"},
			prometheus.Labels{}))
		hadoop_metrics.AppsRunningValType = prometheus.GaugeValue
		hadoop_metrics.AppsPending = append(hadoop_metrics.AppsPending, prometheus.NewDesc("apps_pending", "被挂载的任务个数",
			[]string{"cluster", "queue"},
			prometheus.Labels{}))
		hadoop_metrics.AppsPendinValType = prometheus.GaugeValue
		hadoop_metrics.AppsKilled = append(hadoop_metrics.AppsKilled, prometheus.NewDesc("apps_killed", "被终止的任务个数",
			[]string{"cluster", "queue"},
			prometheus.Labels{}))
		hadoop_metrics.AppsKilledValType = prometheus.GaugeValue
		hadoop_metrics.AppsFailed = append(hadoop_metrics.AppsFailed, prometheus.NewDesc("apps_failed", "失败的任务个数",
			[]string{"cluster", "queue"},
			prometheus.Labels{}))
		hadoop_metrics.AppsFailedValType = prometheus.GaugeValue
		hadoop_metrics.Running0 = append(hadoop_metrics.Running0, prometheus.NewDesc("running_60", "当前运行任务中，运行时长在60分钟内的任务个数",
			[]string{"cluster", "queue"},
			prometheus.Labels{}))
		hadoop_metrics.Running0ValType = prometheus.GaugeValue
		hadoop_metrics.Running60 = append(hadoop_metrics.Running60, prometheus.NewDesc("running_300", "当前运行任务中，运行时长在60-300分钟内的任务个数",
			[]string{"cluster", "queue"},
			prometheus.Labels{}))
		hadoop_metrics.Running60ValType = prometheus.GaugeValue
		hadoop_metrics.Running300 = append(hadoop_metrics.Running300, prometheus.NewDesc("running_1440", "当前运行任务中，运行时长在300-1440分钟内的任务个数",
			[]string{"cluster", "queue"},
			prometheus.Labels{}))
		hadoop_metrics.Running300ValType = prometheus.GaugeValue
	}

	hadoop_metrics.CapacityTotalGB = prometheus.NewDesc("capacity_total_gb", "hadoop集群hdfs的总容量",
		[]string{"cluster"},
		prometheus.Labels{})
	hadoop_metrics.CapacityTotalGBValType = prometheus.GaugeValue
	hadoop_metrics.CapacityUsedGB = prometheus.NewDesc("capacity_used_gb", "hadoop集群hdfs已经使用的容量",
		[]string{"cluster"},
		prometheus.Labels{})
	hadoop_metrics.CapacityUsedGBValType = prometheus.GaugeValue
	hadoop_metrics.CapacityRemainingGB = prometheus.NewDesc("capacity_remaining_gb", "hadoop集群剩余的容量",
		[]string{"cluster"},
		prometheus.Labels{})
	hadoop_metrics.CapacityRemainingGBValType = prometheus.GaugeValue
	hadoop_metrics.BlocksTotal = prometheus.NewDesc("blocks_total", "hdfs的block个数",
		[]string{"cluster"},
		prometheus.Labels{})
	hadoop_metrics.BlocksTotalValType = prometheus.GaugeValue
	hadoop_metrics.FilesTotal = prometheus.NewDesc("files_total", "hdfs的总文件个数",
		[]string{"cluster"},
		prometheus.Labels{})
	hadoop_metrics.FilesTotalValType = prometheus.GaugeValue

	for i := 0; i < len(namenode_list); i++ {
		hadoop_metrics.FileInfoOps = append(hadoop_metrics.FileInfoOps, prometheus.NewDesc("fileinfo_ops", "fileinfo操作次数",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.FileInfoOpsValType = prometheus.GaugeValue
		hadoop_metrics.CreateFileOps = append(hadoop_metrics.CreateFileOps, prometheus.NewDesc("createfile_ops", "createfile操作次数",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.CreateFileOpsValType = prometheus.GaugeValue
		hadoop_metrics.GetlistingOps = append(hadoop_metrics.GetlistingOps, prometheus.NewDesc("getlisting_ops", "getlisting操作次数",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.GetlistingOpsValType = prometheus.GaugeValue
		hadoop_metrics.DeleteFileOps = append(hadoop_metrics.DeleteFileOps, prometheus.NewDesc("deletefile_ops", "deletefile操作次数",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.DeleteFileOpsValType = prometheus.GaugeValue
		// rpc队列中的平均等待时间
		hadoop_metrics.RpcQueueTimeAvgTime = append(hadoop_metrics.RpcQueueTimeAvgTime, prometheus.NewDesc("rpcqueuetime_avg_time", "rpc在交互中的平均等待时间",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.RpcQueueTimeAvgTimeValType = prometheus.GaugeValue
		// rpc打开的连接数目
		hadoop_metrics.NumOpenConnections = append(hadoop_metrics.NumOpenConnections, prometheus.NewDesc("num_open_connections", "rpc打开的连接数",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.NumOpenConnectionsValType = prometheus.GaugeValue
		// rpc slowcall的次数
		hadoop_metrics.RpcSlowCalls = append(hadoop_metrics.RpcSlowCalls, prometheus.NewDesc("rpc_slow_calls", "rpc slowcall的次数",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.RpcSlowCallsValType = prometheus.GaugeValue
		// rpc队列长度
		hadoop_metrics.CallQueueLength = append(hadoop_metrics.CallQueueLength, prometheus.NewDesc("rpc_callqueue_length", "rpc调用队列的长度",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.CallQueueLengthValType = prometheus.GaugeValue
		// rpc授权成功次数
		hadoop_metrics.RpcAuthorizationSuccesses = append(hadoop_metrics.RpcAuthorizationSuccesses, prometheus.NewDesc("rpc_authorization_successes", "rpc授权成功次数",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.RpcAuthorizationSuccessesValType = prometheus.GaugeValue
		// rpc授权失败次数
		hadoop_metrics.RpcAuthorizationFailures = append(hadoop_metrics.RpcAuthorizationFailures, prometheus.NewDesc("rpc_authorization_failures", "rpc授权失败的次数",
			[]string{"cluster", "host", "ip", "port", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.RpcAuthorizationFailuresValType = prometheus.GaugeValue
	}
	for i := 0; i < len(datanode_list); i++ {
		hadoop_metrics.BytesRead = append(hadoop_metrics.BytesRead, prometheus.NewDesc("bytes_read", "datanode读取的字节数",
			[]string{"cluster", "host", "ip", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.BytesReadValType = prometheus.GaugeValue
		hadoop_metrics.BytesWritten = append(hadoop_metrics.BytesWritten, prometheus.NewDesc("bytes_written", "datanode写入字节数",
			[]string{"cluster", "host", "ip", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.BytesWrittenValType = prometheus.GaugeValue
		hadoop_metrics.RemoteBytesRead = append(hadoop_metrics.RemoteBytesRead, prometheus.NewDesc("remote_bytes_read", "datanode远端读取字节数",
			[]string{"cluster", "host", "ip", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.RemoteBytesReadValType = prometheus.GaugeValue
		hadoop_metrics.RemoteBytesWritten = append(hadoop_metrics.RemoteBytesWritten, prometheus.NewDesc("remote_bytes_written", "datanode远端写入字节数",
			[]string{"cluster", "host", "ip", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.RemoteBytesWrittenValType = prometheus.GaugeValue
		hadoop_metrics.HeartbeatsNum = append(hadoop_metrics.HeartbeatsNum, prometheus.NewDesc("heartbeat_num", "心跳💗的次数",
			[]string{"cluster", "host", "ip", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.HeartbeatsNumValType = prometheus.CounterValue
		hadoop_metrics.HeartbeatsAvgTime = append(hadoop_metrics.HeartbeatsAvgTime, prometheus.NewDesc("heartbeat_avg_time", "心跳💗的平均时间",
			[]string{"cluster", "host", "ip", "service_name"},
			prometheus.Labels{}))
		hadoop_metrics.HeartbeatsAvgTimeValType = prometheus.GaugeValue
	}
	return &HadoopCollector{hadoopMetrics: hadoop_metrics}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *HadoopCollector) Describe(ch chan<- *prometheus.Desc) {
	// for _, service_status_desc := range collector.hadoopMetrics.ServiceStatus {
	// 	ch <- service_status_desc
	// }
	ch <- collector.hadoopMetrics.NumLiveDataNodes
	ch <- collector.hadoopMetrics.NumDeadDataNodes
	ch <- collector.hadoopMetrics.NumLiveNameNodes
	ch <- collector.hadoopMetrics.NumDeadNameNodes
	for _, mem_non_heap_usedm := range collector.hadoopMetrics.MemNonHeapUsedM {
		ch <- mem_non_heap_usedm
	}
	for _, mem_non_committed_usedm := range collector.hadoopMetrics.MemNonHeapCommittedM {
		ch <- mem_non_committed_usedm
	}
	for _, mem_heap_usedm := range collector.hadoopMetrics.MemHeapUsedM {
		ch <- mem_heap_usedm
	}
	for _, mem_heap_committedm := range collector.hadoopMetrics.MemHeapCommittedM {
		ch <- mem_heap_committedm
	}
	for _, apps_submitted := range collector.hadoopMetrics.AppsSubmitted {
		ch <- apps_submitted
	}
	for _, apps_running := range collector.hadoopMetrics.AppsRunning {
		ch <- apps_running
	}
	for _, apps_completed := range collector.hadoopMetrics.AppsCompleted {
		ch <- apps_completed
	}
	for _, apps_pending := range collector.hadoopMetrics.AppsPending {
		ch <- apps_pending
	}
	for _, apps_killed := range collector.hadoopMetrics.AppsKilled {
		ch <- apps_killed
	}
	for _, apps_failed := range collector.hadoopMetrics.AppsFailed {
		ch <- apps_failed
	}
	for _, running0 := range collector.hadoopMetrics.Running0 {
		ch <- running0
	}
	for _, running60 := range collector.hadoopMetrics.Running60 {
		ch <- running60
	}
	for _, running300 := range collector.hadoopMetrics.Running300 {
		ch <- running300
	}
	ch <- collector.hadoopMetrics.CapacityTotalGB
	ch <- collector.hadoopMetrics.CapacityUsedGB
	ch <- collector.hadoopMetrics.CapacityRemainingGB
	ch <- collector.hadoopMetrics.BlocksTotal
	ch <- collector.hadoopMetrics.FilesTotal
	for _, fileinfo_ops := range collector.hadoopMetrics.FileInfoOps {
		ch <- fileinfo_ops
	}
	for _, createfile_ops := range collector.hadoopMetrics.CreateFileOps {
		ch <- createfile_ops
	}
	for _, getlisting_ops := range collector.hadoopMetrics.GetlistingOps {
		ch <- getlisting_ops
	}
	for _, deletefile_ops := range collector.hadoopMetrics.DeleteFileOps {
		ch <- deletefile_ops
	}
	for _, rpc_queuetime_avgtime := range collector.hadoopMetrics.RpcQueueTimeAvgTime {
		ch <- rpc_queuetime_avgtime
	}
	for _, num_openconnections := range collector.hadoopMetrics.NumOpenConnections {
		ch <- num_openconnections
	}
	for _, rpc_slow_calls := range collector.hadoopMetrics.RpcSlowCalls {
		ch <- rpc_slow_calls
	}
	for _, callqueue_length := range collector.hadoopMetrics.CallQueueLength {
		ch <- callqueue_length
	}
	for _, rpc_authorization_failures := range collector.hadoopMetrics.RpcAuthorizationFailures {
		ch <- rpc_authorization_failures
	}
	for _, rpc_authorization_success := range collector.hadoopMetrics.RpcAuthorizationSuccesses {
		ch <- rpc_authorization_success
	}
	for _, bytes_written := range collector.hadoopMetrics.BytesWritten {
		ch <- bytes_written
	}
	for _, bytes_read := range collector.hadoopMetrics.BytesRead {
		ch <- bytes_read
	}
	for _, remote_bytes_written := range collector.hadoopMetrics.RemoteBytesWritten {
		ch <- remote_bytes_written
	}
	for _, remote_bytes_read := range collector.hadoopMetrics.RemoteBytesRead {
		ch <- remote_bytes_read
	}
	for _, heartbeats_num := range collector.hadoopMetrics.HeartbeatsNum {
		ch <- heartbeats_num
	}
	for _, haertbeats_avgtime := range collector.hadoopMetrics.HeartbeatsAvgTime {
		ch <- haertbeats_avgtime
	}
}

//Collect implements required collect function for all promehteus collectors
func (collector *HadoopCollector) Collect(ch chan<- prometheus.Metric) {

	// collector = NewHadoopCollector()
	// hadoop_config := Parse_hadoop_config()
	hadoop_config, _ := (utils.ConfigStruct.ConfigData[config.HADOOP]).(config.HadoopConfigure)
	fmt.Println("hadoop_config: ", hadoop_config)
	yarn_urls := make([]string, 0)
	namenode_urls := make([]string, 0)
	for _, yarn_ip := range hadoop_config.Cluster.ResourceManagers {
		// 获取hadoop的metric数据
		yarn_urls = append(yarn_urls, fmt.Sprintf("http://%s:%d/jmx", yarn_ip, hadoop_config.Cluster.ResourcemanagerHttpPort))
	}

	for _, namenode_ip := range hadoop_config.Cluster.Namenodes {
		namenode_urls = append(namenode_urls, fmt.Sprintf("http://%s:%d/jmx", namenode_ip, hadoop_config.Cluster.NamenodeHttpPort))
	}

	num_active_nms, num_lost_nms, num_shutdown_nms, num_unhealthy_nms, num_live_datanodes, num_dead_datanodes, num_decom_livedatanodes, num_decom_missioningdatanodes, num_decommissioning_datanodes, blocks_total, files_total := GetAliveInfo(yarn_urls, namenode_urls)
	fmt.Println(num_active_nms, num_lost_nms, num_shutdown_nms, num_unhealthy_nms, num_live_datanodes, num_dead_datanodes, num_decom_livedatanodes, num_decom_missioningdatanodes, num_decommissioning_datanodes, blocks_total, files_total)
	utils.Logger.Printf("num_active_nms:%d  num_lost_nms:%d  num_shutdown_nms:%d  num_unhealthy_nms:%d  num_live_datanodes:%d  num_dead_datanodes:%d  num_decom_livedatanodes:%d  num_decom_missioningdatanodes:%d  num_decommissioning_datanodes:%d  blocks_total:%d  files_total:%d \n", num_active_nms, num_lost_nms, num_shutdown_nms, num_unhealthy_nms, num_live_datanodes, num_dead_datanodes, num_decom_livedatanodes, num_decom_missioningdatanodes, num_decommissioning_datanodes, blocks_total, files_total)

	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.BlockSize, collector.hadoopMetrics.BlockSizeValType, 128, hadoop_config.Cluster.Name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.ReplicationNum, collector.hadoopMetrics.ReplicationNumValType, 3, hadoop_config.Cluster.Name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.NumLiveDataNodes, collector.hadoopMetrics.NumLiveDataNodesValType, float64(*num_live_datanodes), hadoop_config.Cluster.Name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.NumDeadDataNodes, collector.hadoopMetrics.NumDeadDataNodesValType, float64(*num_dead_datanodes), hadoop_config.Cluster.Name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.NumLiveNameNodes, collector.hadoopMetrics.NumDeadNameNodesValType, float64(*num_active_nms), hadoop_config.Cluster.Name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.NumDeadNameNodes, collector.hadoopMetrics.NumDeadNameNodesValType, float64(*num_lost_nms), hadoop_config.Cluster.Name)

	cluster_list := make([]string, 0)
	host_list := make([]string, 0)
	ip_list := make([]string, 0)
	port_list := make([]string, 0)
	service_name := make([]string, 0)

	urls := make([]string, 0)
	urls = append(urls, fmt.Sprintf("%s?qry=Hadoop:service=ResourceManager,name=JvmMetrics", hadoop_config.Cluster.ResourceManagerUrl))
	host_list = append(host_list, hadoop_config.Cluster.ResourceManagerHosts[0])
	ip_list = append(ip_list, hadoop_config.Cluster.ResourceManagers[0])
	port_list = append(port_list, fmt.Sprintf("%d", hadoop_config.Cluster.ResourcemanagerHttpPort))
	service_name = append(service_name, "ResourceManager")
	cluster_list = append(cluster_list, hadoop_config.Cluster.Name)

	for idx, namenode := range hadoop_config.Cluster.Namenodes {
		urls = append(urls, fmt.Sprintf("http://%s:%d/jmx?qry=Hadoop:service=NameNode,name=JvmMetrics", namenode, hadoop_config.Cluster.NamenodeHttpPort))
		host_list = append(host_list, hadoop_config.Cluster.NamenodeHosts[idx])
		ip_list = append(ip_list, namenode)
		port_list = append(port_list, fmt.Sprintf("%d", hadoop_config.Cluster.NamenodeHttpPort))
		service_name = append(service_name, "NameNode")
		cluster_list = append(cluster_list, hadoop_config.Cluster.Name)
	}
	for idx, datanode := range hadoop_config.Cluster.Datanodes {
		urls = append(urls, fmt.Sprintf("http://%s:%d/jmx?qry=Hadoop:service=DataNode,name=JvmMetrics", datanode, hadoop_config.Cluster.DatanodeHttpPort))
		host_list = append(host_list, hadoop_config.Cluster.DatanodeHosts[idx])
		ip_list = append(ip_list, datanode)
		port_list = append(port_list, fmt.Sprintf("%d", hadoop_config.Cluster.DatanodeHttpPort))
		service_name = append(service_name, "DataNode")
		cluster_list = append(cluster_list, hadoop_config.Cluster.Name)
	}
	mem_non_heap_usedm_list := make([]float64, 0)
	mem_non_heap_committedm_list := make([]float64, 0)
	mem_heap_usedm_list := make([]float64, 0)
	mem_heap_committedm_list := make([]float64, 0)
	fmt.Println("hadoop urls: ", urls)
	utils.Logger.Printf("num_active_nms:%d  num_lost_nms:%d  num_shutdown_nms:%d  num_unhealthy_nms:%d  num_live_datanodes:%d  num_dead_datanodes:%d  num_decom_livedatanodes:%d  num_decom_missioningdatanodes:%d  num_decommissioning_datanodes:%d  blocks_total:%d  files_total:%d \n", num_active_nms, num_lost_nms, num_shutdown_nms, num_unhealthy_nms, num_live_datanodes, num_dead_datanodes, num_decom_livedatanodes, num_decom_missioningdatanodes, num_decommissioning_datanodes, blocks_total, files_total)
	// resourcemanager的jvm信息
	for _, url := range urls {
		mem_non_heap_usedm, mem_non_heap_committedm, mem_heap_usedm, mem_heap_committedm := GetJvmMetricsInfo(url)
		mem_non_heap_usedm_list = append(mem_non_heap_usedm_list, mem_non_heap_usedm)
		mem_non_heap_committedm_list = append(mem_non_heap_committedm_list, mem_non_heap_committedm)
		mem_heap_usedm_list = append(mem_heap_usedm_list, mem_heap_usedm)
		mem_heap_committedm_list = append(mem_heap_committedm_list, mem_heap_committedm)
	}

	//"cluster", "host", "ip", "port", "service_name"
	// 写堆栈信息
	for idx, desc := range collector.hadoopMetrics.MemNonHeapUsedM {
		//"cluster", "host", "ip", "port", "service_name"
		ch <- prometheus.MustNewConstMetric(desc, collector.hadoopMetrics.MemNonHeapUsedMValType, float64(mem_non_heap_usedm_list[idx]), cluster_list[idx], host_list[idx], ip_list[idx], port_list[idx], service_name[idx])
	}
	for idx, desc := range collector.hadoopMetrics.MemNonHeapCommittedM {
		ch <- prometheus.MustNewConstMetric(desc, collector.hadoopMetrics.MemNonHeapCommittedMValType, float64(mem_non_heap_committedm_list[idx]), cluster_list[idx], host_list[idx], ip_list[idx], port_list[idx], service_name[idx])
	}
	for idx, desc := range collector.hadoopMetrics.MemHeapUsedM {
		ch <- prometheus.MustNewConstMetric(desc, collector.hadoopMetrics.MemHeapUsedMValType, float64(mem_heap_usedm_list[idx]), cluster_list[idx], host_list[idx], ip_list[idx], port_list[idx], service_name[idx])
	}
	for idx, desc := range collector.hadoopMetrics.MemHeapCommittedM {
		ch <- prometheus.MustNewConstMetric(desc, collector.hadoopMetrics.MemHeapCommittedMValType, float64(mem_heap_committedm_list[idx]), cluster_list[idx], host_list[idx], ip_list[idx], port_list[idx], service_name[idx])
	}

	queue_name := "q0=root,q1=default"
	yarn_url := fmt.Sprintf("http://%s:%d/jmx", hadoop_config.Cluster.ResourceManagers[0], hadoop_config.Cluster.ResourcemanagerHttpPort)
	apps_submitted, apps_running, apps_pending, apps_killed, apps_failed, apps_completed, running_0, running_60, running_300, _ := GetAppInfo(yarn_url)

	// 基于队列的个数给空间赋值
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.AppsSubmitted[0], collector.hadoopMetrics.AppsSubmittedValType, float64(*apps_submitted), hadoop_config.Cluster.Name, queue_name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.AppsRunning[0], collector.hadoopMetrics.AppsRunningValType, float64(*apps_running), hadoop_config.Cluster.Name, queue_name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.AppsCompleted[0], collector.hadoopMetrics.AppsCompletedValType, float64(*apps_completed), hadoop_config.Cluster.Name, queue_name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.AppsPending[0], collector.hadoopMetrics.AppsPendinValType, float64(*apps_pending), hadoop_config.Cluster.Name, queue_name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.AppsKilled[0], collector.hadoopMetrics.AppsKilledValType, float64(*apps_killed), hadoop_config.Cluster.Name, queue_name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.AppsFailed[0], collector.hadoopMetrics.AppsFailedValType, float64(*apps_failed), hadoop_config.Cluster.Name, queue_name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.Running0[0], collector.hadoopMetrics.Running0ValType, float64(*running_0), hadoop_config.Cluster.Name, queue_name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.Running60[0], collector.hadoopMetrics.Running60ValType, float64(*running_60), hadoop_config.Cluster.Name, queue_name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.Running300[0], collector.hadoopMetrics.Running300ValType, float64(*running_300), hadoop_config.Cluster.Name, queue_name)

	namenode_url := fmt.Sprintf("http://%s:%d/jmx", hadoop_config.Cluster.Namenodes[0], hadoop_config.Cluster.NamenodeHttpPort)
	nondfs_gb, capacity_total_gb, capacity_remaining_gb, capacity_used_gb, blocks_total, corrupt_blocks, pending_deletion_blocks, pending_replication_blocks, files_total, tag_ha_state := GetDFSInfo(namenode_url)
	fmt.Println("nondfs_gb: ", *nondfs_gb)
	fmt.Println("capacity_total_gb: ", *capacity_total_gb)
	fmt.Println("capacity_remaining_gb: ", *capacity_remaining_gb)
	fmt.Println("capacity_used_gb: ", *capacity_used_gb)
	fmt.Println("blocks_total: ", *blocks_total)
	fmt.Println("corrupt_blocks: ", *corrupt_blocks)
	fmt.Println("pending_deletion_blocks: ", *pending_deletion_blocks)
	fmt.Println("pending_replication_blocks: ", *pending_replication_blocks)
	fmt.Println("files_total: ", *files_total)
	fmt.Println("tag_ha_state: ", *tag_ha_state)
	utils.Logger.Printf("capacity_total_gb:%f, capacity_remaining_gb:%f, capacity_used_gb:%f, blocks_total:%d, corrupt_blocks:%d, pending_deletion_blocks:%d, pending_replication_blocks:%d files_total:%d  tag_ha_state:%s\n",
		*capacity_total_gb, *capacity_remaining_gb, *capacity_used_gb,
		*blocks_total, *corrupt_blocks, *pending_deletion_blocks, *pending_replication_blocks,
		*files_total, *tag_ha_state)

	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.CapacityTotalGB, collector.hadoopMetrics.CapacityTotalGBValType, float64(*capacity_total_gb), hadoop_config.Cluster.Name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.CapacityUsedGB, collector.hadoopMetrics.CapacityUsedGBValType, float64(*capacity_used_gb), hadoop_config.Cluster.Name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.CapacityRemainingGB, collector.hadoopMetrics.CapacityRemainingGBValType, float64(*capacity_remaining_gb), hadoop_config.Cluster.Name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.BlocksTotal, collector.hadoopMetrics.BlocksTotalValType, float64(*blocks_total), hadoop_config.Cluster.Name)
	ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.FilesTotal, collector.hadoopMetrics.FilesTotalValType, float64(*files_total), hadoop_config.Cluster.Name)

	// "cluster", "host", "ip", "port", "service_name"
	for idx, namenode := range hadoop_config.Cluster.Namenodes {
		namenode_url := fmt.Sprintf("http://%s:%d/jmx", namenode, hadoop_config.Cluster.NamenodeHttpPort)
		fileinfo_ops, createfile_ops, getlisting_ops, deletefile_ops := GetNameNodeOps(namenode_url)
		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.FileInfoOps[idx],
			collector.hadoopMetrics.FileInfoOpsValType, float64(*fileinfo_ops),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.NamenodeHosts[idx],
			namenode, fmt.Sprintf("%d", hadoop_config.Cluster.NamenodeRpcPort), "NameNode")

		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.CreateFileOps[idx],
			collector.hadoopMetrics.CreateFileOpsValType, float64(*createfile_ops),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.NamenodeHosts[idx],
			namenode, fmt.Sprintf("%d", hadoop_config.Cluster.NamenodeRpcPort), "NameNode")

		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.GetlistingOps[idx],
			collector.hadoopMetrics.GetlistingOpsValType, float64(*getlisting_ops),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.NamenodeHosts[idx],
			namenode, fmt.Sprintf("%d", hadoop_config.Cluster.NamenodeRpcPort), "NameNode")

		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.DeleteFileOps[idx],
			collector.hadoopMetrics.DeleteFileOpsValType, float64(*deletefile_ops),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.NamenodeHosts[idx],
			namenode, fmt.Sprintf("%d", hadoop_config.Cluster.NamenodeRpcPort), "NameNode")

		call_queue_length, rpc_slow_calls, num_open_connections, num_dropped_connections, rpc_authentication_successes, rpc_authentication_failures, sent_bytes, received_bytes, call_queuetime_avgtime, tag_hostname, tag_port := GetNameNodeRPCInfo(namenode_url)
		fmt.Printf("call_queue_length:%d\n  rpc_slow_calls:%d\n  num_open_connections:%d\n  num_dropped_connections:%d\n rpc_authentication_successes:%d\n  rpc_authentication_failures:%d\n  sent_bytes:%d\n  received_bytes:%d\n  call_queuetime_avgtime:%f\n  tag_hostname:%s\n  tag_port:%s \n",
			*call_queue_length, *rpc_slow_calls, *num_open_connections, *num_dropped_connections,
			*rpc_authentication_successes, *rpc_authentication_failures, *sent_bytes, *received_bytes,
			*call_queuetime_avgtime, *tag_hostname, *tag_port)
		// RpcQueueTimeAvgTime        []*prometheus.Desc
		// NumOpenConnections         []*prometheus.Desc
		// RpcSlowCalls               []*prometheus.Desc
		// CallQueueLength            []*prometheus.Desc
		// RpcAuthorizationFailures         []*prometheus.Desc
		// RpcAuthorizationSuccesses        []*prometheus.Desc
		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.RpcQueueTimeAvgTime[idx],
			collector.hadoopMetrics.RpcQueueTimeAvgTimeValType, float64(*call_queuetime_avgtime),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.NamenodeHosts[idx],
			namenode, fmt.Sprintf("%d", hadoop_config.Cluster.NamenodeRpcPort), "NameNode")

		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.NumOpenConnections[idx],
			collector.hadoopMetrics.NumOpenConnectionsValType, float64(*num_open_connections),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.NamenodeHosts[idx],
			namenode, fmt.Sprintf("%d", hadoop_config.Cluster.NamenodeRpcPort), "NameNode")

		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.RpcSlowCalls[idx],
			collector.hadoopMetrics.RpcSlowCallsValType, float64(*rpc_slow_calls),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.NamenodeHosts[idx],
			namenode, fmt.Sprintf("%d", hadoop_config.Cluster.NamenodeRpcPort), "NameNode")

		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.CallQueueLength[idx],
			collector.hadoopMetrics.CallQueueLengthValType, float64(*call_queue_length),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.NamenodeHosts[idx],
			namenode, fmt.Sprintf("%d", hadoop_config.Cluster.NamenodeRpcPort), "NameNode")

		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.RpcAuthorizationFailures[idx],
			collector.hadoopMetrics.RpcAuthorizationFailuresValType, float64(*rpc_authentication_failures),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.NamenodeHosts[idx],
			namenode, fmt.Sprintf("%d", hadoop_config.Cluster.NamenodeRpcPort), "NameNode")

		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.RpcAuthorizationSuccesses[idx],
			collector.hadoopMetrics.RpcAuthorizationSuccessesValType, float64(*rpc_authentication_successes),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.NamenodeHosts[idx],
			namenode, fmt.Sprintf("%d", hadoop_config.Cluster.NamenodeRpcPort), "NameNode")
	}
	// datanode相关的指标
	// BytesWritten              []*prometheus.Desc
	// BytesRead                 []*prometheus.Desc
	// RemoteBytesWritten        []*prometheus.Desc
	// RemoteBytesRead           []*prometheus.Desc
	// HeartbeatsNum             []*prometheus.Desc
	// HeartbeatsAvgTime         []*prometheus.Desc

	for idx, datanode := range hadoop_config.Cluster.Datanodes {
		datanode_url := fmt.Sprintf("http://%s:%d/jmx", datanode, hadoop_config.Cluster.DatanodeHttpPort)
		bytes_read, bytes_written, remote_bytes_read, remote_bytes_written, heartbeats_numops, heartbeats_avgtime, tag_hostname := GetDataNodeInfo(datanode_url)

		fmt.Println(bytes_read, bytes_written, remote_bytes_read, remote_bytes_written, heartbeats_numops, heartbeats_avgtime, tag_hostname)
		//"cluster", "host", "ip", "service_name"
		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.BytesWritten[idx],
			collector.hadoopMetrics.BytesWrittenValType, float64(*bytes_written),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.DatanodeHosts[idx], datanode, "DataNode")

		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.BytesRead[idx],
			collector.hadoopMetrics.BytesReadValType, float64(*bytes_read),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.DatanodeHosts[idx], datanode, "DataNode")

		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.RemoteBytesWritten[idx],
			collector.hadoopMetrics.RemoteBytesWrittenValType, float64(*remote_bytes_written),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.DatanodeHosts[idx], datanode, "DataNode")

		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.RemoteBytesRead[idx],
			collector.hadoopMetrics.RemoteBytesReadValType, float64(*remote_bytes_read),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.DatanodeHosts[idx], datanode, "DataNode")

		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.HeartbeatsNum[idx],
			collector.hadoopMetrics.HeartbeatsNumValType, float64(*heartbeats_numops),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.DatanodeHosts[idx], datanode, "DataNode")

		ch <- prometheus.MustNewConstMetric(collector.hadoopMetrics.HeartbeatsAvgTime[idx],
			collector.hadoopMetrics.HeartbeatsAvgTimeValType, float64(*heartbeats_avgtime),
			hadoop_config.Cluster.Name, hadoop_config.Cluster.DatanodeHosts[idx], datanode, "DataNode")
	}

}
