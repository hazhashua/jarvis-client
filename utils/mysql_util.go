package utils

// user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type ServicePort struct {
	ID           int
	ServiceName  *string
	ChildService *string
	ClusterName  *string
	IP           *string
	Port         sql.NullInt64
	PortType     *string
}

type MysqlConnect struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	DefaultDB string `json:"defaultdb"`
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
		sqlstr = "SELECT count(distinct service_name, child_service, cluster_name, ip , port_type) FROM test.service_port sp"
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

func ReflectNewByValue(target interface{}) *reflect.Value {
	if target == nil {
		fmt.Println("反射参数不能为空......")
		return nil
	}
	typet := reflect.TypeOf(target)
	if typet.Kind() == reflect.Ptr {
		typet = typet.Elem()
	}

	value := reflect.New(typet)
	return &value
}

func ReflectNewByString(typestr string) interface{} {
	// 根据typestr字符串，返回type类型的指针变量
	typestr = strings.ToLower(typestr)
	var typeValue interface{}
	if strings.Contains(typestr, "int") && !strings.Contains(typestr, "uint") {
		typeValue = new(int64)
	}
	if strings.Contains(typestr, "uint") {
		typeValue = new(uint64)
	}
	if strings.Contains(typestr, "float") {
		typeValue = new(float64)
	}
	if strings.Contains(typestr, "string") {
		typeValue = new(string)
	}
	return typeValue
}

// 返回参数列表
func argsList(columns []string, types []string) []interface{} {
	var anyVariable interface{}
	var anyVariables []interface{}
	for _, typev := range types {
		anyVariable = ReflectNewByString(typev)
		fmt.Println("*value: ", anyVariable)
		anyVariables = append(anyVariables, anyVariable)
	}
	return anyVariables
}

type SchemaTable struct {
	CatalogName             string `json:"catalog_name"`
	SchemaName              string `json:"schema_name"`
	DefaultCharacterSetName string `json:"default_character_set_name"`
	DefaultCollationName    string `json:"default_collation_name"`
	SqlPath                 string `json:"sql_path"`
}

// 查询数据库，返回json string类型
func ExecuteSchemaQuery(db *sql.DB, sqlStr string, columns []string, types []string) []SchemaTable {
	if &sqlStr == nil || sqlStr == "" {
		return nil
	}
	stmt, err2 := db.Prepare(sqlStr)
	defer stmt.Close()
	painc_err(err2)
	res, err := stmt.Query()
	defer res.Close()
	painc_err(err)
	schemaTables := make([]SchemaTable, 0)
	for res.Next() {
		st := new(SchemaTable)
		if len(columns) == 1 {
			fmt.Println("")
		} else if len(columns) == 2 {
			fmt.Println("查询两个字断...")
			res.Scan(&st.SchemaName, &st.DefaultCharacterSetName)
		} else if len(columns) == 3 {
			fmt.Println("查询三个字断...")
		} else if len(columns) == 4 {
			fmt.Println("查询4个字断...")
		} else if len(columns) == 5 {
			fmt.Println("查询5个字断...")
		}
		schemaTables = append(schemaTables, *st)
	}
	db.Close()
	return schemaTables
}

type TableTable struct {
	TableSchema string  `json:"table_schema"`
	TableName   string  `json:"table_name"`
	TableRows   int     `json:"table_rows"`
	DataSize    float32 `json:"data_size"`
	IndexSize   float32 `json:"index_size"`
}

func ExecuteTableQuery(db *sql.DB, sqlStr string, columns []string, types []string) []TableTable {
	if &sqlStr == nil || sqlStr == "" {
		return nil
	}
	stmt, err2 := db.Prepare(sqlStr)
	defer stmt.Close()
	painc_err(err2)
	res, err := stmt.Query()
	defer res.Close()
	painc_err(err)
	tableTables := make([]TableTable, 0)
	for res.Next() {
		tt := new(TableTable)
		if len(columns) == 1 {
			fmt.Println("")
		} else if len(columns) == 2 {
			fmt.Println("查询两个字断...")
		} else if len(columns) == 3 {
			fmt.Println("查询三个字断...")
		} else if len(columns) == 4 {
			fmt.Println("查询4个字断...")
		} else if len(columns) == 5 {
			fmt.Println("查询5个字断...")
			res.Scan(&tt.TableSchema, &tt.TableName, &tt.TableRows, &tt.DataSize, &tt.IndexSize)
		}
		tableTables = append(tableTables, *tt)
	}
	db.Close()
	return tableTables
}

type Variable struct {
	VariableName string `json:"variable_name"`
	Value        int    `json:"value"`
}

func ExecuteVariableQuery(db *sql.DB, sqlStr string, columns []string, types []string) []Variable {
	if &sqlStr == nil || sqlStr == "" {
		return nil
	}
	stmt, err2 := db.Prepare(sqlStr)
	defer stmt.Close()
	painc_err(err2)
	res, err := stmt.Query()
	defer res.Close()
	painc_err(err)
	variables := make([]Variable, 0)
	for res.Next() {
		variable := new(Variable)
		if len(columns) == 1 {
			fmt.Println("查询单个字段...")
		} else if len(columns) == 2 {
			fmt.Println("查询两个字段...")
			res.Scan(&variable.VariableName, &variable.Value)
		}
		variables = append(variables, *variable)
	}
	return variables
}

func ExecuteStatusQuery(db *sql.DB, sqlStr string, columns []string, types []string) []Status {
	if &sqlStr == nil || sqlStr == "" {
		return nil
	}
	stmt, err2 := db.Prepare(sqlStr)
	defer stmt.Close()
	painc_err(err2)
	res, err := stmt.Query()
	defer res.Close()
	painc_err(err)
	statuses := make([]Status, 0)
	for res.Next() {
		status := new(Status)
		if len(columns) == 1 {
			fmt.Println("查询单个字段...")
		} else if len(columns) == 5 {
			fmt.Println("查询5个字段...")
			res.Scan(&status.File, &status.Position, &status.BinlogDoDB, &status.BinlogIgnoreDB, &status.ExecutedGtidSet)
		}
		statuses = append(statuses, *status)
	}
	return statuses
}

func Query(sqlstr string) []ServicePort {
	/*
		查询服务和端口信息
	*/
	// 本机 test用户，默认密码为空
	// dsn := "test:@tcp(localhost:3306)/test?charset=utf8&parseTime=true"

	// 开发环境 root用户， pwd@123密码
	// "%s:%s@tcp(%s:%d)/hive?charset=utf8&parseTime=true"
	dsn := fmt.Sprintf("%s:%s@tcp(192.168.10.70:3306)/test?charset=utf8&parseTime=true", "root", "pwd@123")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// 执行sql查询
	if &sqlstr == nil || sqlstr == "" {
		// sqlstr = "SELECT * FROM test.datasource_alive da"
		sqlstr = "SELECT * FROM test.service_port sp order by service_name asc"
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
	fmt.Println("service_port table length: ", len(service_port_slice))
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
	// service_ports := Query("")
	db := DbOpen(nil)
	service_ports := PgServiceQuery(db, "")
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", mysql_connection.Username, mysql_connection.Password, mysql_connection.Host, mysql_connection.Port, mysql_connection.DefaultDB)
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

// 查询库数据
func SchemaQuery(mysqlConnector MysqlConnect) []SchemaTable {
	db := GetConnection(mysqlConnector)

	query := "SELECT SCHEMA_NAME, DEFAULT_CHARACTER_SET_NAME FROM schemata"
	columns := []string{"schema_name", "character_set"}
	types := []string{"string", "string"}

	st := ExecuteSchemaQuery(db, query, columns, types)
	// db.Close()
	return st
}

// 查询表数据
func TableQuery(mysqlConnector MysqlConnect) []TableTable {
	db := GetConnection(mysqlConnector)

	query := `SELECT 
			TABLE_SCHEMA, TABLE_NAME, TABLE_ROWS, 
			data_length/1024/1024 as data_size, 
			index_length/1024/1024 as index_size 
			FROM information_schema.tables ORDER BY data_length DESC;`
	columns := []string{"schema_name", "table_name", "table_rows", "data_size", "index_size"}
	types := []string{"string", "string", "int", "float", "float"}

	tt := ExecuteTableQuery(db, query, columns, types)
	// db.Close()
	return tt

}

func ConnectionQuery(mysqlConnector MysqlConnect) []Variable {

	maxconnectionQuery := "SHOW VARIABLES LIKE 'MAX_CONNECTIONS'"
	userconnectionQuery := "SHOW VARIABLES LIKE 'MAX_USER_CONNECTIONS'"
	currentConnectionQuery := "SHOW status LIKE '%Threads_connected'"

	columns := []string{"variable_name", "value"}
	types := []string{"string", "int"}

	db := GetConnection(mysqlConnector)
	variables := ExecuteVariableQuery(db, maxconnectionQuery, columns, types)

	variablesUser := ExecuteVariableQuery(db, userconnectionQuery, columns, types)

	variablesCurrent := ExecuteVariableQuery(db, currentConnectionQuery, columns, types)

	variables = append(variables, variablesUser[0])
	variables = append(variables, variablesCurrent[0])
	db.Close()
	return variables
}

func QpsAndSlowSqlQuery(mysqlConnector MysqlConnect) []Variable {
	qpsQuery := "SHOW GLOBAL STATUS LIKE 'Queries'"
	slowSqlQuery := "SHOW GLOBAL status like '%que%'"
	columns := []string{"variable_name", "value"}
	types := []string{"string", "int"}
	db := GetConnection(mysqlConnector)
	variables := ExecuteVariableQuery(db, qpsQuery, columns, types)
	variablesSlow := ExecuteVariableQuery(db, slowSqlQuery, columns, types)
	for _, variable := range variablesSlow {
		if variable.VariableName == "Slow_queries" {
			variables = append(variables, variable)
		}
	}
	db.Close()
	return variables
}

type Status struct {
	File            string `json:"file"`
	Position        string `json:"position"`
	BinlogDoDB      string `json:"binlog_do_db"`
	BinlogIgnoreDB  string `json:"binlog_ignore_db"`
	ExecutedGtidSet string `json:"executed_gtid_set"`
}

func TpsQuery(mysqlConnector MysqlConnect) []Status {
	tpsSql := "show master STATUS"

	columns := []string{"file", "position", "binlog_do_db", "binlog_ignore_db", "executed_gtid_set"}
	types := []string{"string", "int", "string", "string", "string"}

	db := GetConnection(mysqlConnector)
	statuses := ExecuteStatusQuery(db, tpsSql, columns, types)
	db.Close()

	return statuses

}
