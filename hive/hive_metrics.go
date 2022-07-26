package hive

import (
	"fmt"
	"io/ioutil"
	"metric_exporter/utils"

	"gopkg.in/yaml.v2"
)

//   cluster:
//   name: 测试环境hive
//   hosts:
//     - 192.168.10.220
//     - 192.168.10.221
//     - 192.168.10.222
//   rpcport: 10000
//   mysql:
//       host: 192.168.10.223
//       port: 3306
//       user: root
//       password: pwd@123
type HiveConfig struct {
	Cluster struct {
		Name    string   `yaml:"name"`
		Hosts   []string `yaml:"hosts"`
		Rpcport string   `yaml:"rpcport"`
		Mysql   struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
		}
	}
}

func Parse_hive_config() *HiveConfig {
	// var bytes []byte
	hive_config := new(HiveConfig)
	if bytes, err := ioutil.ReadFile("./hive/config.yaml"); err != nil {
		fmt.Println("读配置文件 hive/config.yaml出错！")
		return nil
	} else {
		yaml.Unmarshal(bytes, hive_config)
	}
	return hive_config
}

// service_alive	hive各个服务的存活状态
// partition_table_num	分区表个数
// nonpartition_table_num	未分区大表个数
// table_small_file_num	每个表下的小文件数
// hive_closed_ops	Hive 操作数量: 关闭的操作数量
// hive_finished_ops	Hive 操作数量: 完成的操作数量
// hive_canceled_ops	Hive 操作数量: 取消的操作数量
// hive_error_ops	Hive 操作数量: 出错的操作数量

// 库，表个数
func GetDbs() {
	hive_config := Parse_hive_config()
	// host := hive_config.Cluster.Mysql.Host
	// port := hive_config.Cluster.Mysql.Port
	// username := hive_config.Cluster.Mysql.User
	// password := hive_config.Cluster.Mysql.Password
	mysql_connection := utils.MysqlConnect{
		Host:     hive_config.Cluster.Mysql.Host,
		Port:     hive_config.Cluster.Mysql.Port,
		Username: hive_config.Cluster.Mysql.User,
		Password: hive_config.Cluster.Mysql.Password,
	}

	// db := utils.GetConnection(mysql_connection)
	dbs := utils.QueryDbs(mysql_connection)

	for _, db := range dbs {
		fmt.Print(db.Db_id, db.Desc, db.Db_location_uri, db.Name, db.Owner_name, db.Owner_type, db.Ctlg_name, "\n")
	}
	fmt.Println("数据库个数: ", len(dbs))

}

type output interface {
	output() (string, string, string, string, string)
}

func QueryTbls(mysql_connection utils.MysqlConnect) {
	// dsn := "root:pwd@123@tcp(192.168.10.70:3306)/test?charset=utf8&parseTime=true"
	// db, err := sql.Open("mysql", dsn)
	// if err != nil {
	// 	panic(err)
	// }
	// // See "Important settings" section.
	// db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)
	db := utils.GetConnection(mysql_connection)
	sqlstr := "SELECT dbs.name, COUNT(tbls.tbl_id) as tables FROM tbls JOIN dbs ON tbls.db_id=dbs.db_id GROUP BY dbs.name"
	stmt, _ := db.Prepare(sqlstr)
	defer stmt.Close()
	res, _ := stmt.Query()
	defer res.Close()
	for res.Next() {
		var name string
		var num int
		err := res.Scan(&name, &num)
		if err != nil {
			fmt.Println("err: ", err.Error())
		}
		fmt.Println("数据库名: ", name, " 表个数: ", num)
	}
	db.Close()
}

// 查询分区表的信息
func QueryPartitionTbls(mysql_connection utils.MysqlConnect) {
	db := utils.GetConnection(mysql_connection)
	// 查询所有表及其是不是分区表
	sqlstr := "SELECT dbs.name, tbls.tbl_name, prts.tbl_id FROM tbls LEFT OUTER JOIN partitions prts ON tbls.tbl_id=prts.tbl_id LEFT OUTER JOIN dbs ON tbls.db_id=dbs.db_id ORDER BY dbs.name DESC"
	stmt, _ := db.Prepare(sqlstr)
	defer stmt.Close()
	res, _ := stmt.Query()
	defer res.Close()
	for res.Next() {
		var name string
		var tbl_name string
		var tbl_id int
		err := res.Scan(&name, &tbl_name, &tbl_id)
		if err != nil {
			fmt.Println("数据库名: ", name, " 表名: ", tbl_name, " 非分区表")
		} else {
			fmt.Println("数据库名: ", name, " 表名: ", tbl_name, "分区表id: ", tbl_id)
		}
	}
	db.Close()
}

// 查询内部表外部表信息
func QueryExternalTbls(mysql_connection utils.MysqlConnect) {
	db := utils.GetConnection(mysql_connection)
	// 查询所有表及其是不是分区表
	sqlstr := "SELECT dbs.name, tbls.tbl_name, tbls.tbl_type FROM  tbls join dbs on tbls.db_id=dbs.db_id"
	stmt, _ := db.Prepare(sqlstr)
	defer stmt.Close()
	res, _ := stmt.Query()
	defer res.Close()
	for res.Next() {
		var name string
		var tbl_name string
		var tbl_type string
		err := res.Scan(&name, &tbl_name, &tbl_type)
		if err != nil {
			fmt.Println("读取内外部表数据错误！")
		}
		if tbl_type == "MANAGED_TABLE" {
			fmt.Println("数据库名: ", name, " 表名: ", tbl_name, " 表类型: 内部表")
		} else if tbl_type == "EXTERNAL_TABLE" {
			fmt.Println("数据库名: ", name, " 表名: ", tbl_name, "表类型: 外部表")
		}
	}
	db.Close()
}

func QueryTableFileInfo(mysql_connection utils.MysqlConnect) {
	db := utils.GetConnection(mysql_connection)
	// 查询所有表及其是不是分区表
	sqlstr := `SELECT dbs.name, tbls.tbl_name, tp.param_key, tp.param_value
				FROM 
					(SELECT * 
						FROM TABLE_PARAMS 
						WHERE 
						PARAM_KEY IN ("numFiles","totalSize")
					) tp 
				JOIN tbls 
				ON tp.tbl_id=tbls.tbl_id 
				JOIN dbs 
				on dbs.db_id=tbls.db_id`
	stmt, _ := db.Prepare(sqlstr)
	defer stmt.Close()
	res, _ := stmt.Query()
	defer res.Close()
	for res.Next() {
		var name string
		var tbl_name string
		var tbl_type string
		var type_value string
		err := res.Scan(&name, &tbl_name, &tbl_type, &type_value)
		if err != nil {
			fmt.Println("读取表存储相关数据错误！")
		}

		fmt.Println("数据库名: ", name, " 表名: ", tbl_name, " 表指标类型: ", tbl_type, " 表指标值: ", type_value)

	}
	db.Close()

}
