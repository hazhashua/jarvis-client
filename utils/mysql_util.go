package utils

// user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type ServicePort struct {
	ID           int
	ServiceName  *string
	ChildService *string
	ClusterName  *string
	IP           *string
	Port         int
	PortType     *string
}

type MysqlConnect struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// type Sortable interface {
// 	Len() int
// 	Less(int, int) bool
// 	Swap(int, int)
// }

func ValueQuery(sqlstr string) int {
	// dsn := "test:@tcp(localhost:3306)/test?charset=utf8&parseTime=true"
	//      root:pwd@123@tcp(127.0.0.1:3306)/test
	dsn := "root:pwd@123@tcp(192.168.10.70:3306)/test?charset=utf8&parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if &sqlstr == nil || sqlstr == "" {
		// sqlstr = "SELECT * FROM test.datasource_alive da"
		sqlstr = "SELECT count(distinct service_name, child_service, cluster_name, ip , port, port_type) FROM test.service_port sp"
	}
	stmt, err2 := db.Prepare(sqlstr)
	defer stmt.Close()
	painc_err(err2)
	res, err := stmt.Query()
	defer res.Close()
	painc_err(err)
	var count_value int
	for res.Next() {
		// var datasource_alive DatsourceAlive
		// err = res.Scan(&datasource_alive.ID, &datasource_alive.Cluster, &datasource_alive.Name, &datasource_alive.DatasourceType, &datasource_alive.MetricTimestamp, &datasource_alive.MetricValue)
		err = res.Scan(&count_value)
		painc_err(err)
		fmt.Println("count is ", count_value)
	}
	db.Close()
	return count_value
}

func Query(sqlstr string) []ServicePort {
	/*
		查询服务和端口信息
	*/
	// 本机 test用户，默认密码为空
	// dsn := "test:@tcp(localhost:3306)/test?charset=utf8&parseTime=true"

	// 开发环境 root用户， pwd@123密码
	dsn := "root:pwd\\@123@tcp(192.168.10.70:3306)/test?charset=utf8&parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if &sqlstr == nil || sqlstr == "" {
		// sqlstr = "SELECT * FROM test.datasource_alive da"
		sqlstr = "SELECT * FROM test.service_port sp"
	}
	stmt, err2 := db.Prepare(sqlstr)
	defer stmt.Close()
	painc_err(err2)
	res, err := stmt.Query()
	defer res.Close()
	painc_err(err)
	service_port_slice := make([]ServicePort, 0)
	for res.Next() {
		// var datasource_alive DatsourceAlive
		// err = res.Scan(&datasource_alive.ID, &datasource_alive.Cluster, &datasource_alive.Name, &datasource_alive.DatasourceType, &datasource_alive.MetricTimestamp, &datasource_alive.MetricValue)
		var service_port ServicePort
		err = res.Scan(&service_port.ID, &service_port.ServiceName, &service_port.ChildService, &service_port.ClusterName, &service_port.IP, &service_port.Port, &service_port.PortType)
		service_port_slice = append(service_port_slice, service_port)
		painc_err(err)
		fmt.Println(service_port.ID, *service_port.ServiceName, *service_port.ChildService, *service_port.ClusterName, *service_port.IP, service_port.Port, *service_port.PortType)
	}
	db.Close()
	fmt.Println("len: ", len(service_port_slice))
	return service_port_slice
}

func Insert(servicePort ServicePort) bool {
	/*
		往数据库中插入数据
	*/
	return true
}

func Serilize() bool {
	//对结构体数据进行序列化操作
	// var serilize_data []byte
	serilize_data := make([]byte, 0)
	service_ports := Query("")
	for _, service_port := range service_ports {
		service_port_seriaize, _ := json.Marshal(service_port)
		fmt.Println("&&&&&&&", string(service_port_seriaize))
		for _, ele := range service_port_seriaize {
			serilize_data = append(serilize_data, ele)
		}
		serilize_data = append(serilize_data, '\n')
	}
	// fmt.Println("************", string(serilize_data), "************")

	err2 := ioutil.WriteFile("./port_info.txt", serilize_data[:len(serilize_data)-1], 0666) //写入文件(字节数组)
	painc_err(err2)
	return true
}

func ReSerialize() []ServicePort {
	/*
		对序列化后的数据进行反序列化
	*/
	// port_infos := make([]byte, 0)
	port_infos, err := ioutil.ReadFile("./port_info.txt")
	painc_err(err)
	// fmt.Println("port_infos: ", string(port_infos))
	servicePorts := make([]ServicePort, 0)
	for _, ele := range strings.Split(string(port_infos), "\n") {
		if ele == "" {
			fmt.Println("line data empty, continue...")
			continue
		}
		// fmt.Println("ele: ", ele)
		var servicePort = new(ServicePort)
		json.Unmarshal([]byte(ele), servicePort)
		fmt.Println(servicePort.ID, servicePort.ServiceName, servicePort.ChildService, servicePort.ClusterName, servicePort.IP, servicePort.Port, servicePort.PortType)

		servicePorts = append(servicePorts, *servicePort)
	}
	return servicePorts

}

func painc_err(err error) {
	if err != nil {
		panic(err)
	}
}

func GetConnection(mysql_connection MysqlConnect) *sql.DB {
	// 获取初始化的mysql db 结构体
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/test?charset=utf8&parseTime=true", mysql_connection.Username, mysql_connection.Password, mysql_connection.Host, mysql_connection.Port)
	fmt.Println("mysql 连接串: ", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("mysql打开连接失败......")
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db

}

type Dbs struct {
	Db_id           int    `json:"DB_ID"`
	Desc            string `json:"DESC"`
	Db_location_uri string `json:"DB_LOCATION_URI"`
	Name            string `json:"NAME"`
	Owner_name      string `json:"OWNER"`
	Owner_type      string `json:"OWNER_TYPE"`
	Ctlg_name       string `json:"CTLG_NAME"`
}

func (*Dbs) limitf() (num int) {
	return 100
}

func QueryDbs(mysql_connection MysqlConnect) {
	var mysql_db *gorm.DB
	dsn_str := fmt.Sprintf("%s:%s@tcp(%s:%d)/test?charset=utf8", mysql_connection.Username, mysql_connection.Password, mysql_connection.Host, mysql_connection.Port)
	mysql_db, err := gorm.Open("mysql", dsn_str)
	if err != nil {
		fmt.Println("failed to connect database:", err)
		return
	}
	fmt.Println("connect database success")
	mysql_db.SingularTable(true)
	var dbs []Dbs
	mysql_db.Limit(100).Find(&dbs)
	fmt.Println("end QueryDbs......")
	fmt.Println(dbs[0])

	// mysql_db.

}
