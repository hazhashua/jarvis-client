package utils

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func checkError(err error) {
	fmt.Println("error: ", err.Error())
}

type dbConfig struct {
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type data_gather struct {
	Id           int64     `json:"id"`
	ServiceName  string    `json:"service_name"`
	Ip           string    `json:"ip"`
	Port         string    `json:"port"`
	ProtocolType string    `json:"protocol_type"`
	Remarks      string    `json:"remarks"`
	ServiceType  int64     `json:"service_type"`
	CreateName   string    `json:"create_name"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
}

type gather_name struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

var db *gorm.DB

// 创建数据库对象
func DbOpen(dbConfig *dbConfig) (db *gorm.DB) {
	// var err error
	//参数根据自己的数据库进行修改
	// db, err = sql.Open("postgres", "host=192.168.10.79 port=5432 user=postgres password=pwd@123 dbname=ahdb sslmode=disable")
	dsn := "host=192.168.10.68 user=postgres password=pwd@123 dbname=public port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err == nil {
		return db
	}
	return nil
}

// 基于sql查询数据库信息
func PgQuery(sql string) (servicePort []ServicePort) {
	sps := make([]ServicePort, 0)
	if sql == "" {
		sql = ` SELECT * 
					FROM public.data_gather_configure dgc 
					JOIN public.gather_name gn 
					ON dgc.service_type=gn.id `
		db.Raw(sql).Scan(&sps)
	}
	return sps
}

// pg数据库中插入数据
func PgInsert(srcData ...interface{}) {
	db.Create(&srcData)
}
