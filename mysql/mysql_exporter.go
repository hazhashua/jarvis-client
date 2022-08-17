package mysql

import (
	"fmt"
	"metric_exporter/utils"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

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
	dbNum := utils.ValueQuery("select count(schema_name) from information_schema.schemata")

	//查询当前table的个数
	tableNum := utils.ValueQuery("select count(table_name) from information_schema.tables")

	tableInfoDescriptions := make([]*prometheus.Desc, tableNum)
	dbInfoDescriptions := make([]*prometheus.Desc, dbNum)
	for idx := 0; idx < dbNum; idx++ {
		dbInfoDescriptions[idx] = prometheus.NewDesc(
			prometheus.BuildFQName("", "", "db_info"),
			"show the mysql db detail info",
			[]string{"db_name", "default_character_set_name", "cluster", "ip"},
			nil,
		)
	}

	for idx := 0; idx < tableNum; idx++ {
		tableInfoDescriptions[idx] = prometheus.NewDesc(
			prometheus.BuildFQName("", "", "table_info"),
			"show the mysql table detail info",
			[]string{"db_name", "table_name", "table_rows", "data_size", "index_size", "cluster", "ip"},
			nil,
		)
	}

	//返回exporter对象
	return &MysqlExporter{
		up:                 prometheus.NewDesc("up", "show whether the mysql instance is ok", []string{"cluster", "role", "ip"}, prometheus.Labels{}),
		maxConnections:     prometheus.NewDesc("max_connections", "show the max connections of the mysql instance", []string{"cluster", "ip"}, prometheus.Labels{}),
		maxUserConnections: prometheus.NewDesc("max_user_connections", "show the max user connections of the mysql instance", []string{"cluster", "ip"}, prometheus.Labels{}),

		currentConnections:  prometheus.NewDesc("current_connections", "show the current connections of the mysql instance", []string{"cluster", "ip"}, prometheus.Labels{}),
		executeQuerys:       prometheus.NewDesc("query_total", "show the total query by this mysql instance", []string{"cluster", "ip"}, prometheus.Labels{}),
		executeTransactions: prometheus.NewDesc("transaction_total", "show the total transaction by this mysql instance", []string{"cluser", "ip"}, prometheus.Labels{}),

		querySlowTotal: prometheus.NewDesc("query_slow_total", "show the total slow query by thie mysql instance", []string{"cluster", "ip"}, prometheus.Labels{}),

		dbInfos:    dbInfoDescriptions,
		tableInfos: tableInfoDescriptions,
	}

}

func (e *MysqlExporter) Describe(ch chan<- *prometheus.Desc) {
	// 实现exporter的describe方法
	ch <- e.up
	ch <- e.maxConnections
	ch <- e.maxUserConnections
	ch <- e.currentConnections
	ch <- e.executeQuerys
	ch <- e.executeTransactions

	ch <- e.querySlowTotal
	for _, description := range e.dbInfos {
		ch <- description
	}

	for _, description := range e.tableInfos {
		ch <- description
	}
}

func (e *MysqlExporter) Collect(ch chan<- prometheus.Metric) {
	// 实现exporter的collector方法
	e = NewMysqlExporter()
	mysqlConfig := Parse_mysql_config()
	fmt.Println(mysqlConfig.Cluster.Ips)
	fmt.Println(mysqlConfig.Cluster.Port)
	mysqlConnector := utils.MysqlConnect{
		Host:      mysqlConfig.Cluster.Ips[0],
		Port:      mysqlConfig.Cluster.Port,
		Username:  "root",
		Password:  "pwd@123",
		DefaultDB: "information_schema",
	}
	fmt.Println("mysqlConnector: ", mysqlConnector)
	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 1, mysqlConfig.Cluster.Name, mysqlConfig.Cluster.Role[0], mysqlConnector.Host)
	// 查询mysql连接信息
	variables := utils.ConnectionQuery(mysqlConnector)
	for _, variable := range variables {
		// variable.VariableName
		if variable.VariableName == "max_connections" {
			ch <- prometheus.MustNewConstMetric(e.maxConnections, prometheus.GaugeValue, float64(variable.Value), mysqlConfig.Cluster.Name, mysqlConnector.Host)
		}
		if variable.VariableName == "max_user_connections" {
			ch <- prometheus.MustNewConstMetric(e.maxUserConnections, prometheus.GaugeValue, float64(variable.Value), mysqlConfig.Cluster.Name, mysqlConnector.Host)
		}
		if variable.VariableName == "Threads_connected" {
			ch <- prometheus.MustNewConstMetric(e.currentConnections, prometheus.GaugeValue, float64(variable.Value), mysqlConfig.Cluster.Name, mysqlConnector.Host)
		}
	}

	// executeQuerys:       prometheus.NewDesc("query_total", "show the total query by this mysql instance", []string{"cluster", "ip"}, prometheus.Labels{}),
	// executeTransactions: prometheus.NewDesc("transaction_total", "show the total transaction by this mysql instance", []string{"cluser", "ip"}, prometheus.Labels{}),
	// querySlowTotal: prometheus.NewDesc("query_slow_total", "show the total slow query by thie mysql instance", []string{"cluster", "ip"}, prometheus.Labels{}),

	// ch <- prometheus.MustNewConstMetric(collector.kafkaMetrics.BrokerNum, collector.kafkaMetrics.BrokerNumValueType, float64(total_brokers), kafka_config.Cluster.Name)

	connections := make(map[string]int)
	variables = utils.QpsAndSlowSqlQuery(mysqlConnector)
	for _, variable := range variables {
		connections[variable.VariableName] = variable.Value
	}
	ch <- prometheus.MustNewConstMetric(e.executeQuerys, prometheus.CounterValue, float64(connections["Queries"]), mysqlConfig.Cluster.Name, mysqlConnector.Host)
	ch <- prometheus.MustNewConstMetric(e.querySlowTotal, prometheus.CounterValue, float64(connections["Slow_queries"]), mysqlConfig.Cluster.Name, mysqlConnector.Host)

	statuses := utils.TpsQuery(mysqlConnector)
	for _, status := range statuses {
		//
		fmt.Println("status.ExecutedGtidSet: ", status.ExecutedGtidSet)
		length := len(strings.Split(status.ExecutedGtidSet, "-"))
		tpsTotal, _ := strconv.Atoi(strings.Split(status.ExecutedGtidSet, "-")[length-1])
		ch <- prometheus.MustNewConstMetric(e.executeTransactions, prometheus.CounterValue, float64(tpsTotal), mysqlConfig.Cluster.Name, mysqlConnector.Host)
	}

	//"db_name", "default_character_set_name", "cluster", "ip"
	schemaTables := utils.SchemaQuery(mysqlConnector)
	for idx, schema := range schemaTables {
		fmt.Println("schema.SchemaName: ", schema.SchemaName)
		fmt.Println("schema.: ", schema.DefaultCharacterSetName)
		ch <- prometheus.MustNewConstMetric(e.dbInfos[idx], prometheus.GaugeValue, 1, schema.SchemaName, schema.DefaultCharacterSetName, mysqlConfig.Cluster.Name, mysqlConnector.Host)

	}

	//"db_name", "table_name", "table_rows", "data_size", "index_size", "cluster", "ip"
	tableTables := utils.TableQuery(mysqlConnector)
	for idx, table := range tableTables {
		fmt.Println(table.TableSchema, table.TableName, table.TableRows, table.DataSize, table.IndexSize)
		ch <- prometheus.MustNewConstMetric(e.tableInfos[idx], prometheus.GaugeValue, 1, table.TableSchema, table.TableName, fmt.Sprintf("%d", table.TableRows), fmt.Sprintf("%f", table.DataSize), fmt.Sprintf("%f", table.IndexSize), mysqlConfig.Cluster.Name, mysqlConnector.Host)
	}

}
