package kafka

import (
	"fmt"
	"io/ioutil"
	"syscall"

	"github.com/Shopify/sarama"
	"gopkg.in/yaml.v2"
)

type KafkConfigure struct {
	Cluster struct {
		Hosts []string `yaml:"hosts"`
		Port  int      `yaml:"port"`
		Env   string   `yaml:"env"`
	}
}

type DiskStatus struct {
	Path string `json:"path"`
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

func Parse_kafka_config() *KafkConfigure {
	bytes, _ := ioutil.ReadFile("./kafka/config.yaml")
	kafkaConfig := new(KafkConfigure)
	err := yaml.Unmarshal(bytes, kafkaConfig)
	if err != nil {
		fmt.Println("Unmarshal failed: ", err)
	}
	fmt.Println("bytes: \n", string(bytes))
	fmt.Println("kafkaConfig.Cluster.env: ", kafkaConfig.Cluster.Env)
	fmt.Println("kafkaConfig.Cluster.Hosts", kafkaConfig.Cluster.Hosts)
	fmt.Println("kafkaConfig.Cluster.Port: ", kafkaConfig.Cluster.Port)
	return kafkaConfig
}

func getBrokerNums(client sarama.Client) (total int, alive int) {
	// 获取活的总的broker个数
	alive_brokers := client.Brokers()
	fmt.Println("活的broker个数: ", len(alive_brokers))

	//总的broker个数，读取配置文件
	kafkaConfig := Parse_kafka_config()
	broker_total := len(kafkaConfig.Cluster.Hosts)
	fmt.Println("总的broker个数: ", broker_total)
	return broker_total, len(alive_brokers)
}

func getDiskUsage(client sarama.Client, config *sarama.Config) []*DiskStatus {
	// 获取磁盘使用率
	// get broker logDir
	brokers := client.Brokers()
	fmt.Printf("活的broker个数: %d:\n", len(brokers))
	disk_status := make([]*DiskStatus, 0)
	for _, broker := range brokers {
		fmt.Println("broker 地址: ", broker.Addr())
		err3 := broker.Open(config)
		if err3 != nil {
			fmt.Println("connect broker error: ", err3.Error())
		}
		conneted, _ := broker.Connected()
		if conneted {
			var request = new(sarama.DescribeLogDirsRequest)
			dldr, err2 := broker.DescribeLogDirs(request)

			if err2 != nil {
				fmt.Println("error: ", err2.Error())
			} else {
				fmt.Println("logDir: ", dldr.LogDirs[0].Path)
			}
			fmt.Printf("%s\n", broker.Addr())

			// disk usage of path/disk
			fs := syscall.Statfs_t{}
			err := syscall.Statfs(dldr.LogDirs[0].Path, &fs)
			if err != nil {
				return nil
			}
			var disk DiskStatus
			disk.Path = broker.Addr()
			disk.All = fs.Blocks * uint64(fs.Bsize)
			disk.Free = fs.Bfree * uint64(fs.Bsize)
			disk.Used = disk.All - disk.Free
			fmt.Println("diskUsage: ", disk)
			disk_status = append(disk_status, &disk)
		}
	}
	return disk_status
}

// type number struct{
// }

func listDistinct(listx []int32) (num int) {
	// 列表的未重复元素个数
	mapObj := make(map[int32]int)
	for _, ele := range listx {
		mapObj[ele] = 0
	}
	return len(mapObj)
}

func listStringDistinct(listx []string) (num int) {
	mapObj := make(map[string]int, 0)
	for _, ele := range listx {
		mapObj[ele] = 0
	}
	return len(mapObj)
}

func getTopicInfo(client sarama.Client, config *sarama.Config) (topic_num_metric int, topic_partition_metric map[string]int, topic_partition_brokers map[string][]string,
	topic_partition_offsets_metric map[string]map[int]int64, topic_partition_replication_metric map[string]map[int32]int, replication_distribution_balanced_rate_metric map[string]float32,
	consumer_group_num_metric int, topic_partition_consumer_group_offsets map[string]map[string]int64, topic_partition_balance_rate map[string]float32) {
	// 获取topic的相关信息
	// metric:  topic个数   		partition个数
	// 			replication个数     topic partition下的偏移量
	topics, err := client.Topics()
	if err != nil {
		fmt.Printf("try get topics err %s\n", err.Error())
		return 0, nil, nil, nil, nil, nil, 0, nil, nil
	}

	topic_num_metric = len(topics)
	// fmt.Printf("topic个数: %d:\n", topic_num)
	topic_partitions := make(map[string][]int32)
	// 记录topic的partition个数
	topic_partition_metric = make(map[string]int, 0)
	// 记录topic partition的副本信息
	topic_partition_replication_metric = make(map[string]map[int32]int, 0)
	// 当前topic的偏移量
	topic_partition_offsets_metric = make(map[string]map[int]int64, 0)

	all, _ := getBrokerNums(client)
	// 副本的列表
	var replication_ids []int32
	// 存放topic下的副本个数
	replication_distribution_balanced_rate_metric = make(map[string]float32)
	topic_partition_brokers = make(map[string][]string, 0)
	topic_partition_balance_rate = make(map[string]float32)
	for _, topic := range topics {
		fmt.Println("topic: ", topic)
		partitions, _ := client.Partitions(topic)

		fmt.Println("topic的分区数: ", partitions)
		topic_partition_metric[topic] = len(partitions)
		// topic_partitions := len(partitions)
		topic_partition_replication_metric[topic] = make(map[int32]int)
		topic_partition_offsets_metric[topic] = make(map[int]int64, 0)
		topic_partition_brokers[topic] = make([]string, 0)
		for _, partition := range partitions {
			// fmt.Print("partition: ", partition, "\t")
			replication_ids, _ = client.Replicas(topic, partition)
			brokerObj, _ := client.Leader(topic, partition)
			fmt.Println("id: ", brokerObj.ID(), " addr: ", brokerObj.Addr())
			topic_partition_brokers[topic] = append(topic_partition_brokers[topic], fmt.Sprintf("%d", brokerObj.ID()))
			fmt.Println("副本的ids: ", replication_ids, "\t")
			// fmt.Println("副本的个数: ", len(replication_ids))
			topic_partition_replication_metric[topic][partition] = len(replication_ids)
			offsets, err5 := client.GetOffset(topic, partition, sarama.OffsetNewest)
			if err5 != nil {
				fmt.Println("GetOffset error: ", err5.Error())
			} else {
				offset_info := fmt.Sprintf("topic:%s  partition:%d  offsets:%d", topic, partition, offsets)
				fmt.Println(offset_info)
				topic_partition_offsets_metric[topic][int(partition)] = offsets
			}
		}

		fmt.Println("topic 分区的分布: ", topic_partition_brokers[topic])
		// 分区的均衡率
		topic_partition_balance_rate[topic] = float32(listStringDistinct(topic_partition_brokers[topic])) / float32(all)

		topic_partitions[topic] = make([]int32, 0)
		topic_partitions[topic] = partitions

		distinct_num := listDistinct(replication_ids)
		replication_distribution_balanced_rate_metric[topic] = float32(distinct_num) / float32(all)
	}

	s := config.Admin
	fmt.Println(s.Retry, s.Timeout)

	cluster_admin, err2 := sarama.NewClusterAdminFromClient(client)
	if err2 != nil {
		fmt.Println("err2: ", err2.Error())
	}
	// {consumer_group: {topic_partition: -1}}
	topic_partition_consumer_group_offsets = make(map[string]map[string]int64, 0)
	// 获取consumer的消费便宜量
	consumer_groups := make([]string, 0)
	m, err3 := cluster_admin.ListConsumerGroups()
	if err3 != nil {
		fmt.Println("err3: ", err3.Error())
	} else {
		for group, _ := range m {
			topic_partition_consumer_group_offsets[group] = make(map[string]int64)
			consumer_groups = append(consumer_groups, group)
			ofr, err5 := cluster_admin.ListConsumerGroupOffsets(group, topic_partitions)
			if err5 != nil {
				fmt.Println("err5: ", err5.Error())
			}
			for topic, offset_infos := range ofr.Blocks {
				// 消费组下partition的消费情况
				for partition, offset_info := range offset_infos {
					topic_partition := fmt.Sprintf("%s_%d", topic, partition)
					topic_partition_consumer_group_offsets[group][topic_partition] = offset_info.Offset
					// fmt.Println("partition: ", partition, "offset_info: ", offset_info.Offset, offset_info.Metadata)
				}
			}
			// fmt.Println()
			// fmt.Println()
		}
	}
	consumer_group_num_metric = len(consumer_groups)
	fmt.Println("topic个数: ", topic_num_metric)
	fmt.Println("topic分区数: ", topic_partition_metric)
	fmt.Println("topic分区的分布: ", topic_partition_brokers)
	fmt.Println("topic 分区的偏移量: ", topic_partition_offsets_metric)
	fmt.Println("topic 分区的副本数: ", topic_partition_replication_metric)
	fmt.Println("topic 副本的分布均衡率: ", replication_distribution_balanced_rate_metric)
	fmt.Println("消费组个数: ", consumer_group_num_metric)
	fmt.Println("消费组的消费偏移量: ", topic_partition_consumer_group_offsets)

	return

	// return topic_num_metric, topic_partition_metric, topic_partition_brokers, topic_partition_offsets_metric,
	// 	topic_partition_replication_metric, replication_distribution_balanced_rate_metric,
	// 	consumer_group_num_metric, topic_partition_consumer_group_offsets, topic_partition_balance_rate
}

func GetClient() (sarama.Client, sarama.Config) {
	kafka_config := Parse_kafka_config()
	config := sarama.NewConfig()
	config.Version = sarama.V2_7_0_0
	var kafka_host string
	if len(kafka_config.Cluster.Hosts) > 0 {
		kafka_host = kafka_config.Cluster.Hosts[0]
	}
	fmt.Println("kafka_host: ", kafka_host)
	client, err := sarama.NewClient([]string{kafka_host + ":9092"}, config)
	if err != nil {
		fmt.Printf("metadata_test try create client err :%s\n", err.Error())
		return nil, *config
	}

	return client, *config

}

func GetKafkaMetrics() (diskStatus []*DiskStatus, total_brokers int, alive_brokers int, topic_num_metric int, topic_partition_metric map[string]int,
	topic_partition_brokers map[string][]string, topic_partition_offsets_metric map[string]map[int]int64,
	topic_partition_replication_metric map[string]map[int32]int, replication_distribution_balanced_rate_metric map[string]float32,
	consumer_group_num_metric int, topic_partition_consumer_group_offsets map[string]map[string]int64, topic_partition_balance_rate_metric map[string]float32) {

	// 获取kafka的所有指标
	// kafka_config := Parse_kafka_config()
	// config := sarama.NewConfig()
	// config.Version = sarama.V2_7_0_0
	// var kafka_host string
	// if len(kafka_config.Cluster.Hosts) > 0 {
	// 	kafka_host = kafka_config.Cluster.Hosts[0]
	// }
	// fmt.Println("kafka_host: ", kafka_host)
	// client, err := sarama.NewClient([]string{kafka_host + ":9092"}, config)
	// if err != nil {
	// 	fmt.Printf("metadata_test try create client err :%s\n", err.Error())
	// 	return
	// }
	client, config := GetClient()
	defer client.Close()

	total_brokers, alive_brokers = getBrokerNums(client)
	fmt.Println("total: ", total_brokers, "alive: ", alive_brokers)
	topic_num_metric, topic_partition_metric, topic_partition_brokers, topic_partition_offsets_metric, topic_partition_replication_metric, replication_distribution_balanced_rate_metric, consumer_group_num_metric, topic_partition_consumer_group_offsets, topic_partition_balance_rate_metric = getTopicInfo(client, &config)

	diskStatus = getDiskUsage(client, &config)
	return
}

func Metadata_test() {
	fmt.Printf("metadata test\n")
	kafka_config := Parse_kafka_config()
	config := sarama.NewConfig()
	config.Version = sarama.V2_7_0_0
	// config.Version = sarama.V0_11_0_2
	var kafka_host string
	if len(kafka_config.Cluster.Hosts) > 0 {
		kafka_host = kafka_config.Cluster.Hosts[0]
	}
	fmt.Println("kafka_host: ", kafka_host)
	client, err := sarama.NewClient([]string{kafka_host + ":9092"}, config)
	if err != nil {
		fmt.Printf("metadata_test try create client err :%s\n", err.Error())
		return
	}
	defer client.Close()

	// get topic set
	topics, err := client.Topics()
	if err != nil {
		fmt.Printf("try get topics err %s\n", err.Error())
		return
	}

	fmt.Printf("topic个数: %d:\n", len(topics))
	topic_partitions := make(map[string][]int32)
	// topic_partitions["rh_ftp_041918_dev"]=0
	// topic_partitions["rh_opc_057493_dev"]=0
	// topic_partitions["Test"]=0
	// topic_partitions["Tttest02"]=0
	// topic_partitions["node1"]=0
	// topic_partitions["rhopc_057493-realtime"]=0
	for _, topic := range topics {
		fmt.Println("topic: ", topic)
		partitions, _ := client.Partitions(topic)
		fmt.Print("topic的分区数: ", len(partitions))
		for _, partition := range partitions {
			fmt.Print("partition: ", partition, "\t")
			replication_ids, _ := client.Replicas(topic, partition)
			fmt.Print("副本的ids: ", replication_ids, "\t")
			fmt.Println("副本的个数: ", len(replication_ids))
			offsets, err5 := client.GetOffset(topic, partition, sarama.OffsetNewest)
			if err5 != nil {
				fmt.Println("GetOffset error: ", err5.Error())
			} else {
				offset_info := fmt.Sprintf("topic:%s  partition:%d  offsets:%d", topic, partition, offsets)
				fmt.Println(offset_info)
			}
		}
		topic_partitions[topic] = make([]int32, 0)
		topic_partitions[topic] = partitions
	}

	fmt.Println("client.Config().Consumer.Offsets: ", client.Config().Consumer.Offsets)

	s := config.Admin
	// NewClusterAdminFromClient
	fmt.Println(s.Retry, s.Timeout)

}
