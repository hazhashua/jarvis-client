package utils

import (
	"fmt"
	"metric_exporter/config"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func checkError(err error) {
	// fmt.Println("error: ", err.Error())
	Logger.Printf("error: %s", err.Error())
}

type dbConfig struct {
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type dataGather struct {
	Id           int64
	ServiceName  string
	Ip           string
	Port         string
	ProtocolType string
	Remarks      string
	ServiceType  int64
	CreateName   string
	CreateTime   time.Time
	UpdateTime   time.Time
}

type DataStroe struct {
	Id         int
	DataName   string
	Ip         string
	Remarks    string
	DataType   string
	CreateName string
	CreateTime time.Time
	UpdateTime time.Time
	Path       string
}

type Data_store_configure_default DataStroe
type Data_store_configure_cype DataStroe
type Data_store_configure DataStroe

type gatherName struct {
	Id   int64
	Name string
}

// 创建数据库对象
func DbOpen(dbType string) (db *gorm.DB) {
	// 根据dbTye决定数据库类型
	// var err error
	//参数根据自己的数据库进行修改
	// db, err = sql.Open("postgres", "host=192.168.10.79 port=5432 user=postgres password=pwd@123 dbname=ahdb sslmode=disable")
	dc := ParseDbConfig()
	if dbType == "postgres" {
		// dsn := "host=192.168.10.68 user=postgres password=pwd@123 dbname=cluster port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", dc.Cluster.Postgres.Ip, dc.Cluster.Postgres.Username, dc.Cluster.Postgres.Password, dc.Cluster.Postgres.DatasourceInfo.Schema, dc.Cluster.Postgres.Port)
		var err error
		if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err == nil {
			// fmt.Println("*************************connect to db success")
			Logger.Printf("*************************connect to pg db success")
			return db
		} else {
			// fmt.Println("connet to db error!")
			Logger.Printf("connet to pg db error!")
			return nil
		}
	} else if dbType == "mysql" {
		// 创建mysql 连接对象
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", dc.Cluster.Mysql.Username, dc.Cluster.Mysql.Password, dc.Cluster.Mysql.Ip, dc.Cluster.Mysql.Port, dc.Cluster.Mysql.DefaultDB)
		fmt.Println("mysql 连接串: ", dsn)

		if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		}); err == nil {
			// DbMap[fmt.Sprintf("%s%d", config.Cluster.Postgres.Ip, config.Cluster.Postgres.Port)] = dbInstance
			Logger.Println("*************************connect to mysql db success")
			return db
		} else {
			Logger.Println("*************************connect to  mysql error")
			return nil
		}
	}

	return nil
}

// ID           int
// ServiceName  *string
// ChildService *string
// ClusterName  *string
// IP           *string
// Port         sql.NullInt64
// PortType     *string

// 基于sql查询数据库信息
func ServiceQuery(db *gorm.DB) (servicePort []ServicePort) {
	sps := make([]ServicePort, 0)
	sql := fmt.Sprintf(` SELECT dgc.id as id, gn.name as service_name, dgc.service_name as child_service, 
				-- case 
				-- 	when dgc.remarks='' then '大数据融合平台'
				-- else
				-- 	dgc.remarks
				-- end as cluster_name,
				'' as cluster_name,
				dgc.ip as ip, 
				case 
					when dgc.port!='' then cast(dgc.port as signed)
					else -1
				end as port, 
				dgc.protocol_type as port_type,
				dgc.remarks as comment,
				dgc.username as username,
				dgc.password as password
				FROM %s dgc 
				JOIN %s gn 
				ON dgc.service_type=gn.id `, DbConfig.Cluster.Postgres.GatherDetailTable, DbConfig.Cluster.Postgres.GatherTable)
	Logger.Println("ServiceQuery sql: ", sql)
	db.Raw(sql).Scan(&sps)
	return sps
}

// 查询exporter地址数据
//func PgDataStoreQuery(db *gorm.DB, table string) []Data_store_configure_default {
//	dss := make([]Data_store_configure_default, 0)

func DataStoreQuery(db *gorm.DB, table string) []Data_store_configure {
	dss := make([]Data_store_configure, 0)
	// sql := fmt.Sprintf(` SELECT dsc.id as id,
	// 		dsc.data_name as dataname,
	// 		dsc.ip as ip,
	// 		dsc.remarks as remarks,
	// 		dsc.data_type as datatype,
	// 		dsc.create_name as createname,
	// 		dsc.create_time as createtime,
	// 		dsc.update_time as updatetime,
	// 		dsc.path as path
	// 		FROM PUBLIC.%s dsc`, table) //data_store_configure
	// fmt.Println("sql: ", sql)
	// db.Raw(sql).Scan(&dss)
	// 使用gorm语法
	db.Find(&dss)
	return dss
}

// 数据写入data_store_cofigure_default表
func DataStoreInsert(db *gorm.DB, datas *Data_store_configure) {
	db = db.Create(datas)
}

// 删除本exporter关联的data_store信息
func DataStoreRemove(db *gorm.DB) {
	for data_name, ip := range config.MetricIpMap {
		db.Where("data_name=?", strings.ToLower(data_name)).Where("ip=?", ip).Delete(&Data_store_configure{})
		// db.Delete()
	}
}

// 查询基础服务存活数据的个数
func CountQuery(db *gorm.DB, sql string) int {
	if sql == "" {
		sql = fmt.Sprintf(` SELECT COUNT(*) 
		FROM 
		(
			SELECT DISTINCT service_name, ip, port, service_type 
			FROM  %s
		) tmp`, DbConfig.Cluster.Mysql.GatherDetailTable) //DbConfig.Cluster.Postgres.GatherDetailTable
	}
	var countValue int
	db.Raw(sql).Scan(&countValue)
	return countValue
}

// pg数据库中插入数据
func PgGatherNameInsert(db *gorm.DB, srcData gatherName) {
	Logger.Printf("插入表%s数据:  %v \n", DbConfig.Cluster.Postgres.GatherTable, srcData)
	db.Table(fmt.Sprintf("public.%s", DbConfig.Cluster.Postgres.GatherTable)).Create(&srcData)
}

func PgGatherNameConfigureInsert(db *gorm.DB, srcData dataGather) {

	tableName := fmt.Sprintf("public.%s", DbConfig.Cluster.Postgres.GatherDetailTable)
	Logger.Printf("插入表%s数据: %v ", tableName, srcData)
	db.Table(tableName).Create(&srcData)
}

// 临时迁移数据从mysql到pg
func Migirate() {
	// 创建数据库对象
	db := DbOpen("postgres")
	// tx := db.Begin()
	// tx.Exec("use public")
	// tx.Commit()
	_, servicePorts := Query("")
	gatherNames := make([]gatherName, 0)
	dataGathers := make([]dataGather, 0)
	lastName := ""
	newName := ""
	primaryId := 0
	// 将mysql中的数据源数据解析到结构体对象中
	for idx, serviceInfo := range servicePorts {
		newName = *serviceInfo.ServiceName
		// 新的服务名记录在列表中
		if lastName != newName {
			primaryId += 1
			gatherNames = append(gatherNames, gatherName{
				Id:   int64(primaryId),
				Name: *serviceInfo.ServiceName,
			})
		}
		lastName = newName
		Logger.Printf("serviceInfo: %v\n", serviceInfo)

		var port string
		if serviceInfo.Port.Valid {
			port = fmt.Sprintf("%d", serviceInfo.Port.Int64)
		} else {
			port = ""
		}
		dataGathers = append(dataGathers, dataGather{
			Id:           int64(idx) + 1,
			ServiceName:  *serviceInfo.ChildService,
			Ip:           *serviceInfo.IP,
			Port:         port,
			ProtocolType: *serviceInfo.PortType,
			Remarks:      "",
			ServiceType:  int64(primaryId),
			CreateName:   "",
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		})
	}

	// 分别写入gather_name 和 data_gather_configure中数据
	for _, nameInfo := range gatherNames {
		// 插入gather_name表数据
		// db.Table("gather_name")
		PgGatherNameInsert(db, nameInfo)
		Logger.Printf("插入%s数据: %v\n", DbConfig.Cluster.Postgres.GatherTable, nameInfo)
	}

	for _, gatherInfo := range dataGathers {
		// 插入data_gather_configure表数据
		PgGatherNameConfigureInsert(db, gatherInfo)
		Logger.Printf("插入%s数据: %v\n", DbConfig.Cluster.Postgres.GatherDetailTable, gatherInfo)
	}
}
