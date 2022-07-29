package hive

import (
	"database/sql"
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
//   scrapehost: bigdata-dev01
//   scrapeip: 192.168.10.220
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
		ScrapeHost string `yaml:"scrapehost"`
		ScrapeIp   string `yaml:"scrapeip"`
	}
}

type DBS struct {
	DbId         int            `json:"DB_ID"`
	Desc         sql.NullString `json:"DESC"`
	DbLocaionUri *string        `json:"DB_LOCATIONURI"`
	Name         *string        `json:"NAME"`
	OwnerName    *string        `json:"OWNERNAME"`
	OwnerType    *string        `json:"OWNER_TYPE"`
	CtlgName     *string        `json:"CTLG_NAME"`
}

type DBTables struct {
	Name          *string       `json:"NAME"`
	TblId         sql.NullInt16 `json:"TBLID"`
	DbId          *string       `json:"DBID"`
	Owner         *string       `json:"OWNER"`
	TblName       *string       `json:"TBLNAME"`
	TblType       *string       `json:"TBLTYPE"`
	IsPartitioned int           `json:"ISPARTITION"`
	NumFiles      int           `json:"NUMFILES"`
	TotalSize     int           `json:"TOTALSIZE"`
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

// 获取库个数
func GetDbs() []DBS {
	hive_config := Parse_hive_config()
	mysql_connection := utils.MysqlConnect{
		Host:     hive_config.Cluster.Mysql.Host,
		Port:     hive_config.Cluster.Mysql.Port,
		Username: hive_config.Cluster.Mysql.User,
		Password: hive_config.Cluster.Mysql.Password,
	}

	// db := utils.GetConnection(mysql_connection)
	// dbs := QueryDbs(mysql_connection)
	// fmt.Println("数据库个数: ", len(dbs))

	db := utils.GetConnection(mysql_connection)
	sqlstr := "SELECT * FROM DBS"
	stmt, _ := db.Prepare(sqlstr)
	defer stmt.Close()
	res, _ := stmt.Query()
	defer res.Close()
	dbs := make([]DBS, 0)
	for res.Next() {
		var db DBS
		db.DbLocaionUri = new(string)
		db.Name = new(string)
		db.OwnerName = new(string)
		db.OwnerType = new(string)
		db.CtlgName = new(string)
		err := res.Scan(&db.DbId, &db.Desc, db.DbLocaionUri, db.Name, db.OwnerName, db.OwnerType, db.CtlgName)
		if err != nil {
			fmt.Println("err: ", err.Error())
		}
		fmt.Println("数据库信息: ", *db.DbLocaionUri)
		dbs = append(dbs, db)
	}
	db.Close()
	return dbs
}

type output interface {
	output() (string, string, string, string, string)
}

type DbTables struct {
	Name     string `json:"name"`
	TableNum int    `json:"tablenum"`
}

func QueryTbls(mysql_connection utils.MysqlConnect) []DbTables {
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
	tables := make([]DbTables, 0)
	for res.Next() {
		// var name string
		// var num int
		var table DbTables
		err := res.Scan(&table.Name, &table.TableNum)
		if err != nil {
			fmt.Println("err: ", err.Error())
		}
		tables = append(tables, table)
		fmt.Println("数据库名: ", table.Name, " 表个数: ", table.TableNum)
	}
	db.Close()
	return tables
}

// 查询分区表的信息
func QueryPartitionTbls(mysql_connection utils.MysqlConnect) []DBTables {
	db := utils.GetConnection(mysql_connection)
	// 查询所有表及其是不是分区表
	sqlstr := "SELECT dbs.name, tbls.tbl_name, prts.tbl_id, tbls.tbl_type FROM tbls LEFT OUTER JOIN partitions prts ON tbls.tbl_id=prts.tbl_id LEFT OUTER JOIN dbs ON tbls.db_id=dbs.db_id ORDER BY dbs.name DESC"
	stmt, _ := db.Prepare(sqlstr)
	defer stmt.Close()
	res, _ := stmt.Query()
	defer res.Close()
	tables := make([]DBTables, 0)
	for res.Next() {
		// var name string
		// var tbl_name string
		// var tbl_id int
		table := new(DBTables)
		table.Name = new(string)
		table.DbId = new(string)
		table.Owner = new(string)
		table.TblName = new(string)
		table.TblType = new(string)
		var tbl_id sql.NullInt64
		err := res.Scan(table.Name, table.TblName, &tbl_id, table.TblType)
		if tbl_id.Valid {
			table.IsPartitioned = 1
		} else {
			table.IsPartitioned = 0
		}
		if err != nil {
			fmt.Println("err: ", err.Error())
		}
		tables = append(tables, *table)
	}
	db.Close()
	return tables
}

// type DBTables struct {
// 	Name          *string       `json:"NAME"`
// 	TblId         sql.NullInt16 `json:"TBLID"`
// 	DbId          *string       `json:"DBID"`
// 	Owner         *string       `json:OWNER`
// 	TblName       *string       `json:TBLNAME`
// 	TblType       *string       `json:"TBLTYPE"`
// 	IsPartitioned int           `json:"ISPARTITION"`
// }

// 查询表详细信息
func QueryDetailTbls(mysql_connection utils.MysqlConnect) []DBTables {
	db := utils.GetConnection(mysql_connection)
	// 查询所有表及其是不是分区表
	sqlstr := `select  name, tbl_name, tbl_type, tbl_id, 
	IF(param_key is not null, SUBSTRING_INDEX(param_key,',', 1),'') key1,  IF(param_value is not null, SUBSTRING_INDEX(param_value, ',', 1), '') value1 , 
	if(param_key is not null, SUBSTRING_INDEX(param_key,',', -1),'') key2, if( param_value is not null, SUBSTRING_INDEX(param_value, ',', -1), '') value2  
	FROM ( 
		select dbs.name name , tbls.tbl_name tbl_name, tbls.tbl_type tbl_type, prts.tbl_id tbl_id, GROUP_CONCAT(tp.PARAM_KEY) param_key, GROUP_CONCAT(tp.param_value) param_value  
		from tbls   
		LEFT OUTER JOIN partitions prts 
		ON tbls.tbl_id=prts.tbl_id   
		LEFT OUTER JOIN (
			SELECT * FROM TABLE_PARAMS WHERE PARAM_KEY in ('numFiles', 'totalSize') 
			) tp  
		on tp.tbl_id=tbls.tbl_id 
		LEFT OUTER JOIN dbs  
		ON dbs.db_id=tbls.db_id  
		GROUP BY dbs.name, tbls.tbl_name, tbls.tbl_type, prts.tbl_id 
	) tmp`
	stmt, _ := db.Prepare(sqlstr)
	defer stmt.Close()
	res, _ := stmt.Query()
	defer res.Close()
	var db_tables []DBTables
	for res.Next() {
		var tbl_id sql.NullInt64
		var table DBTables
		var k1, k2 string
		var v1, v2 int
		err := res.Scan(&table.Name, &table.TblName, &table.TblType, &tbl_id, &k1, &v1, &k2, &v2)
		if err != nil {
			fmt.Println("读取table表详细数据错误！")
		}
		if tbl_id.Valid {
			table.IsPartitioned = 1
		} else {
			table.IsPartitioned = 0
		}
		if k1 == "numFiles" {
			table.NumFiles = v1
			table.TotalSize = v2
		} else {
			table.NumFiles = v2
			table.TotalSize = v1
		}
		db_tables = append(db_tables, table)
	}
	db.Close()
	return db_tables
}

func QueryTableFileInfo(mysql_connection utils.MysqlConnect) {
	db := utils.GetConnection(mysql_connection)
	// 查询所有表及其是不是分区表
	sqlstr := `SELECT dbs.name name, tbls.tbl_name tbl_name, tp.param_key key, tp.param_value value
				FROM tbls
				LEFT OUTER JOIN
					(SELECT * 
						FROM TABLE_PARAMS 
						WHERE 
						PARAM_KEY IN ("numFiles","totalSize")
					) tp 
				ON tp.tbl_id=tbls.tbl_id 
				LEFT OUTER JOIN dbs 
				on dbs.db_id=tbls.db_id ORDER BY name, tbl_name asc`
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
