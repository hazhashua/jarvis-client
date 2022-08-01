package mysql

import "github.com/prometheus/client_golang/prometheus"

type MysqlExporter struct {
	up                  *prometheus.Desc
	maxConnections      *prometheus.Desc
	maxUserConnections  *prometheus.Desc
	currentConnections  *prometheus.Desc
	executeQuerys       *prometheus.Desc
	executeTransactions *prometheus.Desc
	querySlowTotal      *prometheus.Desc

	// 存储数据库相关信息
	dbInfos []*prometheus.Desc
	// 存储表相关信息
	tableInfos []*prometheus.Desc
}

func NewMysqlExporter() *MysqlExporter {
	// 查询当前db的个数

	//查询当前table的个数

	//返回exporter对象
	return &MysqlExporter{
		up:                 prometheus.NewDesc("up", "show whether the mysql instance is ok", []string{"cluster", "role", "host", "ip"}, prometheus.Labels{}),
		maxConnections:     prometheus.NewDesc("max_connections", "show the max connections of the mysql instance", []string{"cluster", "host", "ip"}, prometheus.Labels{}),
		maxUserConnections: prometheus.NewDesc("", "", []string{}, prometheus.Labels{}),

		currentConnections:  prometheus.NewDesc("", "", []string{}, prometheus.Labels{}),
		executeQuerys:       prometheus.NewDesc("", "", []string{}, prometheus.Labels{}),
		executeTransactions: prometheus.NewDesc("", "", []string{}, prometheus.Labels{}),

		querySlowTotal: prometheus.NewDesc("", "", []string{}, prometheus.Labels{}),
	}

}

func (*MysqlExporter) Describe(ch chan<- *prometheus.Desc) {
	// 实现exporter的describe方法
}

func (*MysqlExporter) Collect(ch chan<- prometheus.Metric) {
	// 实现exporter的collector方法

}
