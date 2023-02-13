package kafka

import (
	"fmt"
	"metric_exporter/config"
	"metric_exporter/utils"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// disk_usage  broker_num  broker_alive_num  topic_num
// topic_lagging_rate  基于生产偏移量计算
// consumer_offset_lag 基于生产偏移量计算
// messages_per_sec    基于偏移量计算
// replication_balance_rate 副本的均衡率 (副本分布的节点/总节点)
// partition_balance_rate   分区的均衡率 (分区分布的节点/总节点)
// topic_consume_balance_rate 基于消费偏移量计算

type KafkaMetric struct {
	// 下面是master的指标
	DiskUsage                           []*prometheus.Desc
	DiskUsageValType                    prometheus.ValueType
	BrokerNum                           *prometheus.Desc
	BrokerNumValueType                  prometheus.ValueType
	BrokerAliveNum                      *prometheus.Desc
	BrokerAliveNumValType               prometheus.ValueType
	TopicNum                            *prometheus.Desc
	TopicNumValType                     prometheus.ValueType
	ConsumerGroupNum                    *prometheus.Desc
	ConsumerGroupNumType                prometheus.ValueType
	TopicPartitionReplicationNum        []*prometheus.Desc
	TopicPartitionReplicationNumValType prometheus.ValueType
	TopicPartitionNum                   []*prometheus.Desc
	TopicPartitionNumValType            prometheus.ValueType
	TopicProduceOffset                  []*prometheus.Desc
	TopicProduceOffsetValType           prometheus.ValueType
	TopicConsumergroupOffset            []*prometheus.Desc
	TopicConsumergroupOffsetValType     prometheus.ValueType
	ReplicationBalanceRate              []*prometheus.Desc
	ReplicationBalanceRateValType       prometheus.ValueType
	PartitionBalanceRate                []*prometheus.Desc
	PartitionBalanceRateValType         prometheus.ValueType
}

type kafkaCollector struct {
	kafkaMetrics KafkaMetric
}

type myDesc struct {
	prometheus.Desc
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func NewKafkaCollector() *kafkaCollector {
	var kafka_metrics KafkaMetric
	// 自动重载kafka配置
	utils.ReloadConfigFromDB(config.KAFKA)
	client, config := GetClient()
	if client == nil {
		return nil
	}
	defer client.Close()
	total_brokers, _, _ := getBrokerInfo(client)
	kafka_metrics.DiskUsage = make([]*prometheus.Desc, 0)
	for i := 0; i < total_brokers; i++ {
		kafka_metrics.DiskUsage = append(kafka_metrics.DiskUsage, prometheus.NewDesc("disk_usage", "show disk usage of the kafka cluster",
			[]string{"cluster", "host", "ip", "broker_id", "disk_path"},
			prometheus.Labels{}))
	}
	kafka_metrics.DiskUsageValType = prometheus.GaugeValue

	kafka_metrics.BrokerNum = prometheus.NewDesc("broker_num", "show broker num of the kafka cluster",
		[]string{"cluster"},
		prometheus.Labels{})
	kafka_metrics.BrokerNumValueType = prometheus.GaugeValue

	kafka_metrics.BrokerAliveNum = prometheus.NewDesc("broker_alive_num", "show the alive broker num of this cluster",
		[]string{"cluster"},
		prometheus.Labels{})
	kafka_metrics.BrokerAliveNumValType = prometheus.GaugeValue

	kafka_metrics.TopicNum = prometheus.NewDesc("topic_num", "the topic num of this kafka cluster",
		[]string{"cluster"},
		prometheus.Labels{})
	kafka_metrics.TopicNumValType = prometheus.GaugeValue

	kafka_metrics.ConsumerGroupNum = prometheus.NewDesc("consumer_group_num", "the consumer group num of this kafka cluster",
		[]string{"cluster"},
		prometheus.Labels{})
	kafka_metrics.ConsumerGroupNumType = prometheus.GaugeValue

	topic_num_metric, topic_partition_metric, topic_partition_brokers, topic_partition_offsets_metric, topic_partition_replication_metric, replication_distribution_balanced_rate_metric, consumer_group_num_metric, topic_partition_consumer_group_offsets, topic_partition_balance_rate, _ := getTopicInfo(client, &config)
	fmt.Println(topic_num_metric, topic_partition_brokers, topic_partition_offsets_metric, topic_partition_replication_metric, replication_distribution_balanced_rate_metric, consumer_group_num_metric, topic_partition_consumer_group_offsets, topic_partition_balance_rate)

	total_partitions := 0
	for _, partitions := range topic_partition_metric {
		total_partitions += partitions
	}
	kafka_metrics.TopicPartitionReplicationNum = make([]*prometheus.Desc, 0)
	for i := 0; i < total_partitions; i++ {
		kafka_metrics.TopicPartitionReplicationNum = append(kafka_metrics.TopicPartitionReplicationNum, prometheus.NewDesc("topic_partition_replication_num", "the replication_num of kafka's topic partition",
			[]string{"cluster", "topic", "partition"},
			prometheus.Labels{}))
	}
	kafka_metrics.TopicPartitionReplicationNumValType = prometheus.GaugeValue

	kafka_metrics.TopicPartitionNum = make([]*prometheus.Desc, 0)
	for i := 0; i < topic_num_metric; i++ {
		kafka_metrics.TopicPartitionNum = append(kafka_metrics.TopicPartitionNum, prometheus.NewDesc("topic_partition_num", "the partition num of kafka's topic partition",
			[]string{"cluster", "topic"},
			prometheus.Labels{}))
	}
	kafka_metrics.TopicPartitionNumValType = prometheus.GaugeValue

	kafka_metrics.TopicProduceOffset = make([]*prometheus.Desc, 0)
	for i := 0; i < total_partitions; i++ {
		kafka_metrics.TopicProduceOffset = append(kafka_metrics.TopicProduceOffset, prometheus.NewDesc("topic_partition_offset", "the partition offset of every kafka's topic partition",
			[]string{"cluster", "topic", "partition"},
			prometheus.Labels{}))
	}
	kafka_metrics.TopicProduceOffsetValType = prometheus.CounterValue

	kafka_metrics.TopicConsumergroupOffset = make([]*prometheus.Desc, 0)
	for i := 0; i < consumer_group_num_metric; i++ {
		for j := 0; j < total_partitions; j++ {
			kafka_metrics.TopicConsumergroupOffset = append(kafka_metrics.TopicConsumergroupOffset,
				prometheus.NewDesc("topic_consumer_group_offset", "the consumer group topic partition offset",
					[]string{"cluster", "consumer_group", "topic", "partition"},
					prometheus.Labels{}))
		}
	}
	kafka_metrics.TopicConsumergroupOffsetValType = prometheus.CounterValue

	// kafka_metrics.ReplicationBalanceRate = prometheus.NewDesc("topic_replication_balance_rate", "the topic replication_balance_rate of the topic",
	// 	[]string{"cluster", "topic"},
	// 	prometheus.Labels{})
	kafka_metrics.ReplicationBalanceRate = make([]*prometheus.Desc, 0)
	for i := 0; i < topic_num_metric; i++ {
		kafka_metrics.ReplicationBalanceRate = append(kafka_metrics.ReplicationBalanceRate,
			prometheus.NewDesc("topic_replication_balance_rate", "the topic's replication balance rate",
				[]string{"cluster", "topic"},
				prometheus.Labels{}))
	}
	kafka_metrics.ReplicationBalanceRateValType = prometheus.GaugeValue

	kafka_metrics.PartitionBalanceRate = make([]*prometheus.Desc, 0)
	for i := 0; i < topic_num_metric; i++ {
		kafka_metrics.PartitionBalanceRate = append(kafka_metrics.PartitionBalanceRate,
			prometheus.NewDesc("topic_partition_balance_rate", "the topic's partition balance rate",
				[]string{"cluster", "topic"},
				prometheus.Labels{}))
	}
	kafka_metrics.PartitionBalanceRateValType = prometheus.CounterValue

	return &kafkaCollector{kafkaMetrics: kafka_metrics}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *kafkaCollector) Describe(ch chan<- *prometheus.Desc) {
	if collector == nil {
		return
	}
	for _, disk_desc := range collector.kafkaMetrics.DiskUsage {
		ch <- disk_desc
	}
	ch <- collector.kafkaMetrics.BrokerNum
	ch <- collector.kafkaMetrics.BrokerAliveNum
	ch <- collector.kafkaMetrics.TopicNum
	ch <- collector.kafkaMetrics.ConsumerGroupNum
	for _, partition_replication_desc := range collector.kafkaMetrics.TopicPartitionReplicationNum {
		ch <- partition_replication_desc
	}

	for _, partition_desc := range collector.kafkaMetrics.TopicPartitionNum {
		ch <- partition_desc
	}
	for _, produce_desc := range collector.kafkaMetrics.TopicProduceOffset {
		ch <- produce_desc
	}
	// ch <- collector.kafkaMetrics.TopicProduceOffset
	for _, consumer_group_desc := range collector.kafkaMetrics.TopicConsumergroupOffset {
		ch <- consumer_group_desc
	}
	// ch <- collector.kafkaMetrics.TopicConsumergroupOffset
	for _, replication_balance_rate := range collector.kafkaMetrics.ReplicationBalanceRate {
		ch <- replication_balance_rate
	}
	// ch <- collector.kafkaMetrics.ReplicationBalanceRate
	for _, partition_balance_rate := range collector.kafkaMetrics.PartitionBalanceRate {
		ch <- partition_balance_rate
	}
	// ch <- collector.kafkaMetrics.PartitionBalanceRate
}

//Collect implements required collect function for all promehteus collectors
func (collector *kafkaCollector) Collect(ch chan<- prometheus.Metric) {

	// 获取kafka的配置信息
	// kafka_config := Parse_kafka_config()
	collector = NewKafkaCollector()
	if collector == nil {
		utils.Logger.Printf("创建kafka采集器失败!")
		return
	}
	kafka_config := (utils.ConfigStruct.ConfigData[config.KAFKA]).(config.KafkaConfigure)
	// 获取kafka的metric数据
	//    disk_status, total_brokers, alive_brokers, topic_num_metric, topic_partition_metric, topic_partition_brokers, topic_partition_offsets_metric, topic_partition_replication_metric, replication_distribution_balanced_rate_metric, consumer_group_num_metric, topic_partition_consumer_group_offsets, topic_partition_balance_rate_metric := GetKafkaMetrics()
	disk_status, total_brokers, alive_brokers, topic_num_metric, topic_partition_metric, _, topic_partition_offsets_metric, topic_partition_replication_metric, replication_distribution_balanced_rate_metric, consumer_group_num_metric, consumer_group_topic_partition_offsets, topic_partition_balance_rate_metric, _ := GetKafkaMetrics()
	// utils.Logger.Printf("topic_partition_brokers: %v\n", topic_partition_brokers)
	for idx, disk_desc := range collector.kafkaMetrics.DiskUsage {
		if idx >= len(disk_status) {
			fmt.Println("数组索引, 超过获取的kafka磁盘数......")
			break
		}
		ch <- prometheus.MustNewConstMetric(disk_desc, prometheus.GaugeValue,
			float64(float32(disk_status[idx].Used)/float32(disk_status[idx].All)),
			kafka_config.Cluster.Name, disk_status[idx].Host, "", fmt.Sprintf("%d", disk_status[idx].BrokerID), disk_status[idx].Path)
	}
	ch <- prometheus.MustNewConstMetric(collector.kafkaMetrics.BrokerNum, collector.kafkaMetrics.BrokerNumValueType, float64(total_brokers), kafka_config.Cluster.Name)
	ch <- prometheus.MustNewConstMetric(collector.kafkaMetrics.BrokerAliveNum, collector.kafkaMetrics.BrokerAliveNumValType, float64(alive_brokers), kafka_config.Cluster.Name)
	ch <- prometheus.MustNewConstMetric(collector.kafkaMetrics.TopicNum, collector.kafkaMetrics.TopicNumValType, float64(topic_num_metric), kafka_config.Cluster.Name)

	// consumer group num
	ch <- prometheus.MustNewConstMetric(collector.kafkaMetrics.ConsumerGroupNum, collector.kafkaMetrics.ConsumerGroupNumType, float64(consumer_group_num_metric), kafka_config.Cluster.Name)

	offsets_info := make([]int64, 0)
	topic_info := make([]string, 0)
	partition_info := make([]int, 0)
	consumer_group_info := make([]string, 0)

	// partition replication num
	topic_partition_replication_info := make([]int, 0)
	for topic, partition_replication := range topic_partition_replication_metric {
		for partition, replication_num := range partition_replication {
			topic_info = append(topic_info, topic)
			partition_info = append(partition_info, int(partition))
			topic_partition_replication_info = append(topic_partition_replication_info, replication_num)
		}
	}
	for idx, partition_replication_desc := range collector.kafkaMetrics.TopicPartitionReplicationNum {
		if idx >= len(topic_info) || idx >= len(partition_info) {
			utils.Logger.Printf("idx指向topic_info|partition_info数组越界!\n")
			break
		}
		ch <- prometheus.MustNewConstMetric(partition_replication_desc, collector.kafkaMetrics.TopicPartitionReplicationNumValType, float64(topic_partition_replication_info[idx]), kafka_config.Cluster.Name, topic_info[idx], fmt.Sprintf("%d", partition_info[idx]))
	}

	// topic partition num
	topic_info = make([]string, 0)
	partition_info = make([]int, 0)
	// topic partition infos
	for topic, partition := range topic_partition_metric {
		topic_info = append(topic_info, topic)
		partition_info = append(partition_info, partition)
	}
	for idx, partition_desc := range collector.kafkaMetrics.TopicPartitionNum {
		ch <- prometheus.MustNewConstMetric(partition_desc, collector.kafkaMetrics.TopicPartitionNumValType, float64(partition_info[idx]), kafka_config.Cluster.Name, topic_info[idx])
	}

	topic_info = make([]string, 0)
	partition_info = make([]int, 0)
	for topic, offsets := range topic_partition_offsets_metric {
		for partition, offset := range offsets {
			topic_info = append(topic_info, topic)
			partition_info = append(partition_info, partition)
			offsets_info = append(offsets_info, offset)
		}
	}
	for idx, produce_offset_desc := range collector.kafkaMetrics.TopicProduceOffset {
		if idx >= len(offsets_info) || idx >= len(topic_info) || idx >= len(partition_info) {
			utils.Logger.Printf("idx指向offsets_info|topic_info|partition_info数组越界!")
			break
		}
		ch <- prometheus.MustNewConstMetric(produce_offset_desc, collector.kafkaMetrics.TopicProduceOffsetValType, float64(offsets_info[idx]), kafka_config.Cluster.Name, topic_info[idx], fmt.Sprintf("%d", partition_info[idx]))
	}

	consumer_group_info = make([]string, 0)
	topic_info = make([]string, 0)
	partition_info = make([]int, 0)
	offsets_info = make([]int64, 0)
	// consumer_group_topic_partition_offsets
	for consumer_group, topic := range consumer_group_topic_partition_offsets {
		for topic_partition, offset := range topic {
			topic_partition_list := strings.Split(topic_partition, "_")
			// partition := topic_partition_list[len(topic_partition)-1]
			partition, _ := strconv.Atoi(topic_partition_list[len(topic_partition_list)-1])
			topic := strings.Join(topic_partition_list[:len(topic_partition_list)-1], "_")

			consumer_group_info = append(consumer_group_info, consumer_group)
			topic_info = append(topic_info, topic)
			partition_info = append(partition_info, partition)
			offsets_info = append(offsets_info, offset)
		}
	}

	for idx, consumer_topic_partition_offsets_desc := range collector.kafkaMetrics.TopicConsumergroupOffset {
		if len(topic_info) > idx {
			ch <- prometheus.MustNewConstMetric(consumer_topic_partition_offsets_desc, collector.kafkaMetrics.TopicConsumergroupOffsetValType, float64(offsets_info[idx]), kafka_config.Cluster.Name, consumer_group_info[idx], topic_info[idx], fmt.Sprintf("%d", partition_info[idx]))

		}
	}

	// 处理分区副本均衡率
	topic_info = make([]string, 0)
	replication_balance_rate_info := make([]float32, 0)
	// collector.kafkaMetrics.ReplicationBalanceRate
	for topic, replication_rate := range replication_distribution_balanced_rate_metric {

		topic_info = append(topic_info, topic)
		replication_balance_rate_info = append(replication_balance_rate_info, replication_rate)
	}
	for idx, replication_balance_rate_desc := range collector.kafkaMetrics.ReplicationBalanceRate {
		if idx >= len(replication_balance_rate_info) || idx >= len(topic_info) {
			utils.Logger.Printf("idx指向topic_info|replication_balance_rate_info数组越界!\n")
			break
		}

		ch <- prometheus.MustNewConstMetric(replication_balance_rate_desc, collector.kafkaMetrics.ReplicationBalanceRateValType, float64(replication_balance_rate_info[idx]), kafka_config.Cluster.Name, topic_info[idx])
	}

	// 处理分区的均衡率
	topic_info = make([]string, 0)
	partition_balance_rate_info := make([]float32, 0)
	for topic, rate := range topic_partition_balance_rate_metric {
		topic_info = append(topic_info, topic)
		partition_balance_rate_info = append(partition_balance_rate_info, rate)
	}

	for idx, partiton_balance_rate_desc := range collector.kafkaMetrics.PartitionBalanceRate {
		if idx >= len(partition_balance_rate_info) || idx >= len(topic_info) {
			utils.Logger.Printf("idx指向topic_info|partition_balance_rate_info数组越界!\n")
			break
		}
		ch <- prometheus.MustNewConstMetric(partiton_balance_rate_desc, collector.kafkaMetrics.PartitionBalanceRateValType, float64(partition_balance_rate_info[idx]), kafka_config.Cluster.Name, topic_info[idx])

	}

}
