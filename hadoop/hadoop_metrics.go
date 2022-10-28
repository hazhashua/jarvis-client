package hadoop

import (
	"metric_exporter/utils"
)

// 8020为namenode rpc端口
// 50070为namenode http端口
// 9864为datanode http端口

// dfs存储信息：
// url:
//	namenode
//		// hdfs文件系统的使用率
//		addr:  192.168.10.220:50070/jmx?qry=Hadoop:service=NameNode,name=FSNamesystem
//
//  datanode
//		addr:  192.168.10.220:9864/jmx?qry=

//  resourcemanager
//      // 获取namenode存活信息
//		addr:  192.168.10.222:8088/jmx?qry=Hadoop:service=ResourceManager,name=ClusterMetrics

//  yarn rest api
//  curl --compressed -H "Accept: application/json" -X  GET "http://host.domain.com:8088/ws/v1/cluster/apps/application_1326821518301_0010"

func GetAppInfo(yarnUrl string) (apps_submitted *int64, apps_running *int64, apps_pending *int64,
	apps_killed *int64, apps_failed *int64, apps_completed *int64, running_0 *int64,
	running_60 *int64, running_300 *int64, running_1440 *int64) {
	//抓取进程存活信息

	url := yarnUrl + "?qry=Hadoop:service=ResourceManager,name=QueueMetrics,q0=root,q1=default"
	// fmt.Println("url: ", url)
	utils.Logger.Printf("request url: %s", url)
	res := utils.GetUrl(url)
	// fmt.Println("url response: ", res)
	// utils.Logger.Printf("url:%s  response: %s\n", url, res)
	app_info := []byte(res)
	rma, err := UnmarshalResourceManagerApp(app_info)
	if err != nil {
		// fmt.Println("parse data error!", err.Error())
		utils.Logger.Printf("UnmarshalResourceManagerApp error error:%s\n", err.Error())
	}
	if len(rma.Beans) > 0 {
		apps_submitted = rma.Beans[0].AppsSubmitted
		apps_running = rma.Beans[0].AppsRunning
		apps_pending = rma.Beans[0].AppsPending
		apps_killed = rma.Beans[0].AppsKilled
		apps_failed = rma.Beans[0].AppsFailed
		apps_completed = rma.Beans[0].AppsCompleted
		running_0 = rma.Beans[0].Running0
		running_60 = rma.Beans[0].Running60
		running_300 = rma.Beans[0].Running300
		running_1440 = rma.Beans[0].Running1440
	}

	// utils.Logger.Printf("rma.Beans[0].AppsSubmitted: %d\n", *rma.Beans[0].AppsSubmitted)
	// utils.Logger.Printf("rma.Beans[0].AppsRunning: %d\n", *rma.Beans[0].AppsRunning)
	// utils.Logger.Printf("rma.Beans[0].AppsPending: %d\n", *rma.Beans[0].AppsPending)
	// utils.Logger.Printf("rma.Beans[0].AppsKilled: %d\n", *rma.Beans[0].AppsKilled)
	// utils.Logger.Printf("rma.Beans[0].AppsFailed: %d\n", *rma.Beans[0].AppsFailed)
	// utils.Logger.Printf("rma.Beans[0].AppsCompleted: %d\n", *rma.Beans[0].AppsCompleted)
	// utils.Logger.Printf("rma.Beans[0].running_0: %d\n", *rma.Beans[0].Running0)
	// utils.Logger.Printf("rma.Beans[0].running_60: %d\n", *rma.Beans[0].Running60)
	// utils.Logger.Printf("rma.Beans[0].running_300: %d\n", *rma.Beans[0].Running300)
	// utils.Logger.Printf("rma.Beans[0].running_1440: %d\n", *rma.Beans[0].Running1440)
	return
}

func GetJvmMetricsInfo(http_url string) (mem_non_heap_usedm float64, mem_non_heap_committedm float64, mem_heap_usedm float64, mem_heap_committedm float64) {
	// fmt.Println("jvm_url: ", http_url)
	res := utils.GetUrl(http_url)
	body_byte := []byte(res)
	jvm_metrics, err := UnmarshalJVMMetrics(body_byte)
	if err == nil {
		utils.Logger.Printf("UnmarshalJVMMetrics(body_byte) error: %s\n", err.Error())
	}
	// fmt.Println(jvm_metrics.Beans[0].Name)
	// fmt.Println(jvm_metrics.Beans[0].MemNonHeapUsedM)
	// fmt.Println(jvm_metrics.Beans[0].MemNonHeapCommittedM)
	// fmt.Println(jvm_metrics.Beans[0].MemHeapUsedM)
	// fmt.Println(jvm_metrics.Beans[0].MemHeapCommittedM)
	if len(jvm_metrics.Beans) > 0 {
		mem_non_heap_usedm = *jvm_metrics.Beans[0].MemNonHeapUsedM
		mem_non_heap_committedm = *jvm_metrics.Beans[0].MemNonHeapCommittedM
		mem_heap_usedm = *jvm_metrics.Beans[0].MemHeapUsedM
		mem_heap_committedm = *jvm_metrics.Beans[0].MemHeapCommittedM
	}

	return
}

func GetDFSInfo(namenoeUrl string) (nondfs_gb *float64, capacity_total_gb *float64, capacity_remaining_gb *float64,
	capacity_used_gb *float64, blocks_total *int64, corrupt_blocks *int64,
	pending_deletion_blocks *int64, pending_replication_blocks *int64,
	files_total *int64, tag_ha_state *string) {

	url := namenoeUrl + "?qry=Hadoop:service=NameNode,name=FSNamesystem"
	utils.Logger.Printf("request url: %s\n", url)
	response := utils.GetUrl(url)
	fs_namesystem_bytes := []byte(response)
	fs_namesystem, err := UnmarshalFSNamesystem(fs_namesystem_bytes)
	if err != nil {
		utils.Logger.Printf("UnmarshalFSNamesystem(fs_namesystem_bytes) error:%s\n", err.Error())
	}

	nondfs_gbv := float64(*fs_namesystem.Beans[0].CapacityUsedNonDFS) / float64(1024.0) / 1024 / 1024
	nondfs_gb = &nondfs_gbv
	capacity_total_gb = fs_namesystem.Beans[0].CapacityTotalGB
	capacity_remaining_gb = fs_namesystem.Beans[0].CapacityRemainingGB
	capacity_used_gb = fs_namesystem.Beans[0].CapacityUsedGB

	// fmt.Println("fs_namesystem.Beans[0].CapacityTotalGB: ", fs_namesystem.Beans[0].CapacityTotalGB)
	// fmt.Println("fs_namesystem.Beans[0].CapacityRemainingGB: ", fs_namesystem.Beans[0].CapacityRemainingGB)
	// fmt.Println("fs_namesystem.Beans[0].CapacityUsedGB: ", fs_namesystem.Beans[0].CapacityUsedGB)
	// fmt.Println("fs_namesystem.Beans[0].CapacityUsedNonDFS", fs_namesystem.Beans[0].CapacityUsedNonDFS)

	// 总的block数量
	blocks_total = fs_namesystem.Beans[0].BlocksTotal
	// 已损坏的block数量
	corrupt_blocks = fs_namesystem.Beans[0].CorruptBlocks
	//未被验证的block个数
	pending_deletion_blocks = fs_namesystem.Beans[0].PendingDeletionBlocks
	//等待被备份的block个数
	pending_replication_blocks = fs_namesystem.Beans[0].PendingReplicationBlocks

	// 总文件的数量
	files_total = fs_namesystem.Beans[0].FilesTotal
	// fmt.Println("fs_namesystem.Beans[0].FilesTotal: ", fs_namesystem.Beans[0].FilesTotal)
	// namenode的角色 active或者standby
	tag_ha_state = fs_namesystem.Beans[0].TagHAState
	// fmt.Println("fs_namesystem.Beans[0].TagHAState: ", fs_namesystem.Beans[0].TagHAState)
	return
}

// 抓取存活数据
func GetAliveInfo(yarnUrls []string, namenodeUrls []string) (ok bool, num_active_nms *int64, num_lost_nms *int64,
	num_shutdown_nms *int64, num_unhealthy_nms *int64, num_live_datanodes *int64,
	num_dead_datanodes *int64, num_decom_livedatanodes *int64, num_decom_missioningdatanodes *int64,
	num_decommissioning_datanodes *int64, blocks_total *int64, files_total *int64) {

	//获取namenode存活信息
	var response string
	for idx, url := range yarnUrls {
		request_url := url + "?qry=Hadoop:service=ResourceManager,name=ClusterMetrics"
		response = utils.GetUrl(request_url)
		if response != "" {
			break
		}
		utils.Logger.Printf("访问地址失败: %s\n", request_url)
		if idx == len(yarnUrls)-1 {
			utils.Logger.Println("yarn url 访问全失败")
			return
		}
	}
	// utils.Logger.Println("response: ", response)
	cluster_metrics_bytes := []byte(response)
	cm, err := UnmarshalClusterMetrics(cluster_metrics_bytes)
	if err != nil {
		utils.Logger.Printf("UnmarshalClusterMetrics(cluster_metrics_bytes) error:%s\n", err.Error())
	}
	if len(cm.Beans) == 0 {
		num_active_nms = new(int64)
		*num_active_nms = -1
		num_lost_nms = new(int64)
		*num_lost_nms = -1
		num_shutdown_nms = new(int64)
		*num_shutdown_nms = -1
		num_unhealthy_nms = new(int64)
		*num_unhealthy_nms = -1
	} else {
		num_active_nms = cm.Beans[0].NumActiveNMS
		// utils.Logger.Printf("cm.Beans[0].NumActiveNMS: %d\n", *cm.Beans[0].NumActiveNMS)
		num_lost_nms = cm.Beans[0].NumLostNMS
		// utils.Logger.Printf("cm.Beans[0].NumLostNMS: %d\n", *cm.Beans[0].NumLostNMS)
		num_shutdown_nms = cm.Beans[0].NumShutdownNMS
		// utils.Logger.Printf("cm.Beans[0].NumShutdownNMS: %d\n", *cm.Beans[0].NumShutdownNMS)
		num_unhealthy_nms = cm.Beans[0].NumUnhealthyNMS
		// utils.Logger.Printf("cm.Beans[0].NumUnhealthyNMS: %d\n", *cm.Beans[0].NumUnhealthyNMS)

	}

	//获取datanode存活数据
	for idx, url := range namenodeUrls {
		request_url := url + "?qry=Hadoop:service=NameNode,name=FSNamesystemState"
		response = utils.GetUrl(request_url)
		if response != "" {
			break
		}
		if idx == len(namenodeUrls)-1 {
			utils.Logger.Println("len(namenodeUrls)-1: 访问地址失败!")
			panic("访问地址失败！")
		}
	}
	fsnamesystem_state_bytes := []byte(response)
	fs, err2 := UnmarshalFSNamesystemState(fsnamesystem_state_bytes)
	if err2 != nil {
		utils.Logger.Printf("UnmarshalFSNamesystemState(fsnamesystem_state_bytes) error: %s", err2.Error())
		// fmt.Println("error: ", err2.Error())
	}
	num_live_datanodes = fs.Beans[0].NumLiveDataNodes
	// utils.Logger.Printf("fs.Beans[0].NumLiveDataNodes: %d\n", *fs.Beans[0].NumLiveDataNodes)
	num_dead_datanodes = fs.Beans[0].NumDeadDataNodes
	// utils.Logger.Printf("fs.Beans[0].NumDeadDataNodes: %d\n", *fs.Beans[0].NumDeadDataNodes)
	num_decom_livedatanodes = fs.Beans[0].NumDecomLiveDataNodes
	// utils.Logger.Printf("fs.Beans[0].NumDecomLiveDataNodes: %d\n", *fs.Beans[0].NumDecomLiveDataNodes)
	num_decom_missioningdatanodes = fs.Beans[0].NumDecommissioningDataNodes
	// utils.Logger.Printf("fs.Beans[0].NumDecommissioningDataNodes: %d\n", *fs.Beans[0].NumDecommissioningDataNodes)
	//hdfs 总块数
	blocks_total = fs.Beans[0].BlocksTotal
	// utils.Logger.Printf("fs.Beans[0].BlocksTotal: %d", *fs.Beans[0].BlocksTotal)
	files_total = fs.Beans[0].FilesTotal
	// utils.Logger.Printf("fs.Beans[0].FilesTotal: %d\n", *fs.Beans[0].FilesTotal)

	return

}

func GetNameNodeRPCInfo(namenode_url string) (call_queue_length *int64, rpc_slow_calls *int64,
	num_open_connections *int64, num_dropped_connections *int64, rpc_authentication_successes *int64,
	rpc_authentication_failures *int64, sent_bytes *int64, received_bytes *int64,
	call_queuetime_avgtime *float64, tag_hostname *string, tag_port *string) {
	url := namenode_url + "?qry=Hadoop:service=NameNode,name=RpcActivityForPort8020"
	response := utils.GetUrl(url)
	rpc_activity_bytes := []byte(response)

	rfp, err := UnmarshalRPCActivityForPort8020(rpc_activity_bytes)
	if err != nil {
		// fmt.Println("error: ", err.Error())
		utils.Logger.Printf("UnmarshalRPCActivityForPort8020(rpc_activity_bytes): %s", err.Error())
	}

	call_queue_length = rfp.Beans[0].CallQueueLength
	rpc_slow_calls = rfp.Beans[0].RPCSlowCalls
	num_open_connections = rfp.Beans[0].NumOpenConnections
	num_dropped_connections = rfp.Beans[0].NumDroppedConnections
	rpc_authentication_successes = rfp.Beans[0].RPCAuthenticationSuccesses
	rpc_authentication_failures = rfp.Beans[0].RPCAuthenticationFailures
	sent_bytes = rfp.Beans[0].SentBytes
	received_bytes = rfp.Beans[0].ReceivedBytes
	call_queuetime_avgtime = rfp.Beans[0].RPCQueueTimeAvgTime
	tag_hostname = rfp.Beans[0].TagHostname
	tag_port = rfp.Beans[0].TagPort

	// utils.Logger.Printf("rfp.Beans[0].CallQueueLength: %d\t rfp.Beans[0].RPCSlowCalls: %d\t rfp.Beans[0].NumOpenConnections: %d\n", *rfp.Beans[0].CallQueueLength, *rfp.Beans[0].RPCSlowCalls, *rfp.Beans[0].NumOpenConnections)
	// utils.Logger.Printf("rfp.Beans[0].NumDroppedConnections: %d\t rfp.Beans[0].RPCAuthenticationSuccesses: %d\t rfp.Beans[0].RPCAuthenticationFailures: %d\n", *rfp.Beans[0].NumDroppedConnections, *rfp.Beans[0].RPCAuthenticationSuccesses, *rfp.Beans[0].RPCAuthenticationFailures)
	// utils.Logger.Printf("rfp.Beans[0].SentBytes: %d\t rfp.Beans[0].ReceivedBytes: %d\t rfp.Beans[0].RPCQueueTimeAvgTime: %f\t\n", *rfp.Beans[0].SentBytes, *rfp.Beans[0].ReceivedBytes, *rfp.Beans[0].RPCQueueTimeAvgTime)

	// 主机名和密码
	utils.Logger.Printf("rfp.Beans[0].TagHostname: %s\t rfp.Beans[0].TagPort: %s\n", *rfp.Beans[0].TagHostname, *rfp.Beans[0].TagPort)
	return

}

// 获取datanode相关的指标
func GetDataNodeRPCInfo(datanode_url string) {
	url := datanode_url + "?qry=Hadoop:service=DataNode,name=RpcActivityForPort9867"
	response := utils.GetUrl(url)
	datanode_response_byte := []byte(response)
	rfp, err := UnmarshalRPCActivityForPort9867(datanode_response_byte)
	if err != nil {
		// fmt.Println("error: ", err.Error())
		utils.Logger.Printf("UnmarshalRPCActivityForPort9867(datanode_response_byte): %s", err.Error())
	}

	utils.Logger.Printf("rfp.Beans[0].RPCAuthenticationFailures: %d\n", *rfp.Beans[0].RPCAuthenticationFailures)
	// fmt.Println("rfp.Beans[0].RPCAuthenticationFailures: ", rfp.Beans[0].RPCAuthenticationFailures)
	// fmt.Println("rfp.Beans[0].RPCAuthenticationSuccesses: ", rfp.Beans[0].RPCAuthenticationSuccesses)
	// fmt.Println("rfp.Beans[0].RPCSlowCalls: ", rfp.Beans[0].RPCSlowCalls)
	// fmt.Println("rfp.Beans[0].RPCQueueTimeAvgTime: ", rfp.Beans[0].RPCQueueTimeAvgTime)
	// fmt.Println("rfp.Beans[0].CallQueueLength: ", rfp.Beans[0].CallQueueLength)
	// fmt.Println("rfp.Beans[0].ReceivedBytes: ", rfp.Beans[0].ReceivedBytes)
	// fmt.Println("rfp.Beans[0].SentBytes: ", rfp.Beans[0].SentBytes)
	// fmt.Println("rfp.Beans[0].TagHostname: ", rfp.Beans[0].TagHostname)
	// fmt.Println("rfp.Beans[0].TagPort: ", rfp.Beans[0].TagPort)

}

// 获取namenode上的各种操作数
func GetNameNodeOps(namenode_url string) (fileinfo_ops *int64, createfile_ops *int64, getlisting_ops *int64, deletefile_ops *int64) {
	url := namenode_url + "?qry=Hadoop:service=NameNode,name=NameNodeActivity"
	response := utils.GetUrl(url)
	namenode_activity_byte := []byte(response)
	nna, err := UnmarshalNameNodeActivity(namenode_activity_byte)
	if err != nil {
		// fmt.Println("error: ", err.Error())
		utils.Logger.Printf("UnmarshalNameNodeActivity(namenode_activity_byte) error:%s\n", err.Error())
	}
	fileinfo_ops = nna.Beans[0].CreateFileOps
	createfile_ops = nna.Beans[0].CreateFileOps
	getlisting_ops = nna.Beans[0].GetListingOps
	deletefile_ops = nna.Beans[0].DeleteFileOps

	// fmt.Println("nna.Beans[0].CreateFileOps: ", *nna.Beans[0].CreateFileOps)
	// fmt.Println("nna.Beans[0].GetListingOps: ", *nna.Beans[0].GetListingOps)
	// fmt.Println("nna.Beans[0].DeleteFileOps: ", *nna.Beans[0].DeleteFileOps)
	// fmt.Println("nna.Beans[0].FileInfoOps: ", *nna.Beans[0].FileInfoOps)

	utils.Logger.Printf("nna.Beans[0].CreateFileOps: %d \t, nna.Beans[0].GetListingOps: %d \t, nna.Beans[0].DeleteFileOps: %d \t nna.Beans[0].FileInfoOps: %d\n", *nna.Beans[0].CreateFileOps, *nna.Beans[0].GetListingOps, *nna.Beans[0].DeleteFileOps, *nna.Beans[0].FileInfoOps)

	return
}

// 获取datanode上的数据读写信息
func GetDataNodeInfo(datanode_url string) (bytes_read *int64, bytes_written *int64,
	remote_bytes_read *int64, remote_bytes_written *int64,
	heartbeats_numops *int64, heartbeats_avgtime *float64, tag_hostname *string) {
	url := datanode_url + "?qry=Hadoop:service=DataNode,name=DataNodeActivity-*"
	response := utils.GetUrl(url)
	datanode_activity_byte := []byte(response)
	dna, err := UnmarshalDataNodeActivity(datanode_activity_byte)
	if err != nil {
		// fmt.Println("error: ", err.Error())
		utils.Logger.Printf("UnmarshalDataNodeActivity(datanode_activity_byte): %s\n", err.Error())
	}

	bytes_read = dna.Beans[0].BytesRead
	bytes_written = dna.Beans[0].BytesWritten
	heartbeats_numops = dna.Beans[0].HeartbeatsNumOps
	heartbeats_avgtime = dna.Beans[0].HeartbeatsAvgTime
	remote_bytes_read = dna.Beans[0].RemoteBytesRead
	remote_bytes_written = dna.Beans[0].RemoteBytesWritten
	tag_hostname = dna.Beans[0].TagHostname

	// // datanode读取的数据量
	// fmt.Println("dna.Beans[0].BytesRead: ", dna.Beans[0].BytesRead)
	// utils.Logger.Printf("dna.Beans[0].BytesRead: %d\n", *dna.Beans[0].BytesRead)
	// // datanode写入的数据量
	// fmt.Println("dna.Beans[0].BytesWritten: ", dna.Beans[0].BytesWritten)
	// utils.Logger.Printf("dna.Beans[0].BytesWritten: %d\n", *dna.Beans[0].BytesWritten)

	// // 💗数量
	// fmt.Println("dna.Beans[0].HeartbeatsNumOps: ", dna.Beans[0].HeartbeatsNumOps)
	// utils.Logger.Printf("dna.Beans[0].HeartbeatsNumOps: %d\n", *dna.Beans[0].HeartbeatsNumOps)

	// // 💗的平均时间
	// fmt.Println("dna.Beans[0].HeartbeatsAvgTime: ", dna.Beans[0].HeartbeatsAvgTime)
	// utils.Logger.Printf("dna.Beans[0].HeartbeatsAvgTime: %f\n", *dna.Beans[0].HeartbeatsAvgTime)

	// // 远端读取数据量
	// // fmt.Println("dna.Beans[0].RemoteBytesRead: ", dna.Beans[0].RemoteBytesRead)
	// utils.Logger.Printf("dna.Beans[0].RemoteBytesRead: %d\n", *dna.Beans[0].RemoteBytesRead)

	// // 远端写入数据量
	// // fmt.Println("dna.Beans[0].RemoteBytesWritten: ", dna.Beans[0].RemoteBytesWritten)
	// utils.Logger.Printf("dna.Beans[0].RemoteBytesWritten: %d\n", *dna.Beans[0].RemoteBytesWritten)

	// // 主机名
	// // fmt.Println("dna.Beans[0].TagHostname: ", dna.Beans[0].TagHostname)
	// utils.Logger.Printf("dna.Beans[0].TagHostname: %s\n", *dna.Beans[0].TagHostname)

	return
}
