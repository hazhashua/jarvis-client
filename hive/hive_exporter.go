package hive

import (
	"fmt"
	"metric_exporter/utils"

	"github.com/prometheus/client_golang/prometheus"
)

type HiveExporter struct {
	// 保存当前处理的redis实例的索引
	// redisAddr         string
	hiveConfig  HiveConfig
	hiveCluster string

	// 存储hive集群模式
	clusterMode               prometheus.Gauge
	totalScrapes              prometheus.Counter
	scrapeDuration            prometheus.Summary
	targetScrapeRequestErrors prometheus.Counter

	metricDescriptions map[string]*prometheus.Desc

	// 存储所有表的基础信息
	dbInfoDescriptions []*prometheus.Desc

	// 存储所有表的基础信息
	tableInfoDescriptions []*prometheus.Desc

	// metricMapCounters map[string]string
	// metricMapGauges   map[string]string
}

func NewHiveExporter() *HiveExporter {
	hive_exporter := HiveExporter{
		hiveConfig:  *Parse_hive_config(),
		hiveCluster: Parse_hive_config().Cluster.Name,
	}

	// 注册常量指标
	for metric, desc := range map[string]struct {
		txt  string
		lbls []string
	}{
		"partition_table_num":    {txt: `the num of partition table`, lbls: []string{"cluster", "exporter_host", "exporter_ip"}},
		"nonpartition_table_num": {txt: `the num of non partition table`, lbls: []string{"cluster", "exporter_host", "exporter_ip"}},
		"external_table_num":     {txt: "the num of external table", lbls: []string{"cluster", "exporter_host", "exporter_ip"}},
		"nonexternal_table_num":  {txt: "the number of non external table", lbls: []string{"cluster", "exporter_host", "exporter_ip"}},
		// "table_info":             {txt: "the detail info of table", lbls: []string{"cluster", "exporter_host", "exporter_ip", "db", "name", "file_num", "external", "partitioned", "capacity"}},
	} {
		hive_exporter.metricDescriptions[metric] = prometheus.NewDesc(prometheus.BuildFQName("", "", metric), desc.txt, desc.lbls, nil)
	}

	db_num := len(GetDbs())
	hive_exporter.dbInfoDescriptions = make([]*prometheus.Desc, 0)

	// DbId         int     `json:"DB_ID"`
	// Desc         *string `json:"DESC"`
	// DbLocaionUri *string `json:"DB_LOCATIONURI"`
	// Name         *string `json:"NAME"`
	// OwnerName    *string `json:"OWNERNAME"`
	// OwnerType    *string `json:"OWNER_TYPE"`
	// CtlgName     *string `json:"CTLG_NAME"`

	for idx := 0; idx < db_num; idx++ {
		hive_exporter.dbInfoDescriptions[idx] = prometheus.NewDesc(
			prometheus.BuildFQName("", "", "db_info"),
			"show the db detail info",
			[]string{"db_desc", "db_location_uri", "name", "owner_name", "cluster", "exporter_host", "exporter_ip"},
			nil,
		)
	}
	hive_config := Parse_hive_config()
	mysql_connection := utils.MysqlConnect{
		Host:     hive_config.Cluster.Mysql.Host,
		Port:     hive_config.Cluster.Mysql.Port,
		Username: hive_config.Cluster.Mysql.User,
		Password: hive_config.Cluster.Mysql.Password,
	}
	hive_exporter.tableInfoDescriptions = make([]*prometheus.Desc, 0)
	db_tables := QueryDetailTbls(mysql_connection)
	for idx := 0; idx < len(db_tables); idx++ {
		hive_exporter.dbInfoDescriptions[idx] = prometheus.NewDesc(
			prometheus.BuildFQName("", "", "table_info"),
			"show the table detail info",
			[]string{"db_name", "table_name", "is_external", "is_partitioned", "cluster", "exporter_host", "exporter_ip"},
			nil,
		)
	}

	// hive_exporter.dbInfoDescriptions = make([]*prometheus.Desc, 0)
	// hive_exporter.tableInfoDescriptions = make([]*prometheus.Desc, 0)
	return &hive_exporter
}

func (e *HiveExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range e.metricDescriptions {
		ch <- desc
	}

	db_num := len(GetDbs())
	e.dbInfoDescriptions = make([]*prometheus.Desc, db_num)
	for _, db_desc := range e.dbInfoDescriptions {
		// 创建counter desc 并写入 ch
		fmt.Println("db_desc: ", db_desc)
		ch <- db_desc
	}

	// 获取table表数据
	hive_config := Parse_hive_config()
	mysql_connection := utils.MysqlConnect{
		Host:     hive_config.Cluster.Mysql.Host,
		Port:     hive_config.Cluster.Mysql.Port,
		Username: hive_config.Cluster.Mysql.User,
		Password: hive_config.Cluster.Mysql.Password,
	}
	db_tables := QueryTbls(mysql_connection)
	table_num := 0
	for _, table := range db_tables {
		table_num += table.TableNum
	}
	e.tableInfoDescriptions = make([]*prometheus.Desc, table_num)

	for _, table_desc := range e.tableInfoDescriptions {
		// 创建counter desc 并写入 ch
		fmt.Println("table_desc: ", table_desc)
		ch <- table_desc
	}

	ch <- e.clusterMode.Desc()

}

func (e *HiveExporter) Collect(ch chan<- prometheus.Metric) {
	// hive 指标采集
	// db_num := len(dbs)

	hive_config := Parse_hive_config()
	mysql_connection := utils.MysqlConnect{
		Host:     hive_config.Cluster.Mysql.Host,
		Port:     hive_config.Cluster.Mysql.Port,
		Username: hive_config.Cluster.Mysql.User,
		Password: hive_config.Cluster.Mysql.Password,
	}
	// QueryTbls(mysql_connection)
	db_tables := QueryPartitionTbls(mysql_connection)
	// QueryTableFileInfo(mysql_connection)

	// 顺序写入metric数据
	// "partition_table_num"
	// "nonpartition_table_num"
	// "external_table_num"
	// "nonexternal_table_num"
	// for metric_name, desc := range e.metricDescriptions {
	// 	if metric_name == "partition_table_num" || metric_name == "nonpartition_table_num" {

	// 获取分区相关的指标信息
	partition_value := 0
	for _, table := range db_tables {
		partition_value += table.IsPartitioned
	}
	// 写分区的相关指标
	ch <- prometheus.MustNewConstMetric(e.metricDescriptions["partition_table_num"], prometheus.GaugeValue, float64(partition_value), hive_config.Cluster.Name, hive_config.Cluster.ScrapeHost, hive_config.Cluster.ScrapeIp)
	ch <- prometheus.MustNewConstMetric(e.metricDescriptions["nonpartition_table_num"], prometheus.GaugeValue, float64(len(db_tables)-partition_value), hive_config.Cluster.Name, hive_config.Cluster.ScrapeHost, hive_config.Cluster.ScrapeIp)

	// 获取内外部表相关指标信息
	external_value := 0
	for _, table := range db_tables {
		if *table.TblType == "EXTERNAL_TABLE" {
			external_value += 1
		}
	}
	// 写内外部表的相关指标
	ch <- prometheus.MustNewConstMetric(e.metricDescriptions["external_table_num"], prometheus.GaugeValue, float64(external_value), hive_config.Cluster.Name, hive_config.Cluster.ScrapeHost, hive_config.Cluster.ScrapeIp)
	ch <- prometheus.MustNewConstMetric(e.metricDescriptions["nonexternal_table_num"], prometheus.GaugeValue, float64(len(db_tables)-external_value), hive_config.Cluster.Name, hive_config.Cluster.ScrapeHost, hive_config.Cluster.ScrapeIp)

	// for _, desc := range e.metricDescriptions {
	// 	prometheus.MustNewConstMetric()
	// }

	// 写数据库的详细指标数据
	dbs := GetDbs()
	// []string{"db_desc", "db_location_uri", "name", "owner_name", "cluster", "exporter_host", "exporter_ip"},
	for idx, desc := range e.dbInfoDescriptions {
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(1), *dbs[idx].Desc, *dbs[idx].DbLocaionUri, *dbs[idx].Name, *dbs[idx].OwnerName, hive_config.Cluster.Name, hive_config.Cluster.ScrapeHost, hive_config.Cluster.ScrapeIp)
	}

	// 写表的详细指标数据
	db_tables = QueryDetailTbls(mysql_connection)
	// []string{"db_name", "table_name", "is_external", "is_partitioned", "cluster", "exporter_host", "exporter_ip"},

	for idx, desc := range e.dbInfoDescriptions {
		is_external := "0"
		if *db_tables[idx].TblType == "EXTERNAL_TABLE" {
			is_external = "1"
		}
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(1), *db_tables[idx].Name, *db_tables[idx].TblName, is_external, string(db_tables[idx].IsPartitioned), hive_config.Cluster.Name, hive_config.Cluster.ScrapeHost, hive_config.Cluster.ScrapeIp)
	}
	ch <- e.clusterMode
}
