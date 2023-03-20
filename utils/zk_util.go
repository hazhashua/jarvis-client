package utils

import (
	"fmt"
	"metric_exporter/config"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var defaultPath string
var currentPath string
var ModelChan chan string

// 初始化zk连接相关信息
func init() {

	defaultPath = "/exporter/all_node"
	currentPath = "/exporter/current"
	ModelChan = make(chan string)
}

func getZkHost() []string {
	var zookeeper_config config.ZookeepeConfig
	var ok bool

	zookeeper_config, ok = (ConfigStruct.ConfigData[config.ZOOKEEPER]).(config.ZookeepeConfig)
	if ok == false {
		Logger.Println("获取zookeeper config配置失败!")
		return []string{}
	}
	fmt.Println("zookeeper_config.Cluster.Name: ", zookeeper_config.Cluster.Name)
	fmt.Println("zookeeper_config.Cluster.Hosts: ", zookeeper_config.Cluster.Hosts)
	fmt.Println("zookeeper_config.Cluster.ClientPort: ", zookeeper_config.Cluster.ClientPort)
	hosts := make([]string, 0)
	for _, host := range zookeeper_config.Cluster.Hosts {
		hosts = append(hosts, fmt.Sprintf("%s:%s", host, zookeeper_config.Cluster.ClientPort))
	}
	return hosts
}

// 获取到zookeeper的连接
func conn() *zk.Conn {
	hosts := getZkHost()
	if len(hosts) == 0 {
		return nil
	}

	// hosts_str := strings.Join(hosts, ",")
	conn, _, err := zk.Connect(hosts, time.Second*5)
	defer conn.Close()
	if err != nil {
		Logger.Println("zk connect fail  error: ", err.Error())
		return nil
	} else {
		return conn
	}
}

// 创建节点
func create(path string, temporary bool, data []byte) bool {
	if path == "" {
		path = "/test_zk_node"
	}
	var conn *zk.Conn = conn()
	if conn != nil {
		defer conn.Close()
	}
	if len(data) == 0 {
		data = []byte("test zk node 数据")
	}
	//flags有4种取值：
	//0:永久，除非手动删除
	//zk.FlagEphemeral = 1:短暂，session断开则改节点也被删除
	//zk.FlagSequence  = 2:会自动在节点后面添加序号
	//3:Ephemeral和Sequence，即，短暂且自动添加序号
	var flags int32 = 0
	if temporary {
		flags = zk.FlagEphemeral
	}
	var acls = zk.WorldACL(zk.PermAll) //控制访问权限模式
	p, err_create := conn.Create(path, data, flags, acls)
	if err_create != nil {
		fmt.Println(err_create)
		return false
	}
	fmt.Println("create node:", p)
	return true
}

// 设置zookeeper的节点数据
func set(path string, data string) {
	if path == "" {
		path = "/brokers"
	}
	connection := conn()
	defer connection.Close()
	// var path = "home"
	fmt.Println("connection.SessionID(): ", connection.SessionID())

	connection.Set(path, []byte(data), -1)

	bytes, stat, _ := connection.Get(path)
	fmt.Println(fmt.Printf("stat.Version: %d", stat.Version))
	fmt.Println("节点的数据为: ", string(bytes))

}

// 查询节点数据
func get(path string) {
	if path == "" {
		path = "/cluster"
	}
	connection := conn()
	defer connection.Close()
	bytes, states, err := connection.Get(path)
	if err != nil {
		fmt.Println("连接zookeeper失败: ", err.Error())
	}
	// states.EphemeralOwner 如果节点为临时节点，则这个值为这个节点拥有者的sessionid, 非临时节点这个值为 0
	node_data := fmt.Sprintf("path:%s  version: %d   data:%s   path stat: %d", path, states.Version, string(bytes), states.EphemeralOwner)
	fmt.Println("node data: ", node_data)

}

// 删除节点
func delete(path string) {
	if path == "" {
		path = "/spark"
	}
	connection := conn()
	defer connection.Close()
	err := connection.Delete(path, -1)
	if err != nil {
		fmt.Println("节点删除失败！")
	} else {
		fmt.Println(fmt.Sprintf("节点%s删除成功!", path))
	}
}

// 发布执行成功的采集模块
func Publish(model string, ch chan string) {
	ch <- model
	fmt.Printf("发布模块%s", model)
}

func register(model string) {
	// 这册模块到zookeeper
	connection := conn()
	var existPath bool = true
	if existFlag, zkState, _ := connection.Exists(currentPath); !existFlag {
		Logger.Printf("zkState: %v", zkState)
		// 不存在存储根节点，则创建存储根节点
		existPath = create(currentPath, false, []byte("存储当前连接状态的exporter"))
	}
	if !existPath {
		Logger.Printf("创建current节点失败")
		return
	}
	// 注册临时节点
	if createOk := create(fmt.Sprintf("%s/%s", currentPath, model), true, []byte(fmt.Sprintf("%s模块连接成功", model))); !createOk {
		Logger.Printf("%s模块连接创建失败！", model)
	} else {
		Logger.Printf("%s模块连接创建成功！", model)
	}

}

func RegisterDefaultAll() {
	// 注册所有的exporter
	// HADOOP HBASE HIVE KAFKA MICROSERVICE
	// MYSQL NODE REDIS SKYWALKING SPARK ZOOKEEPER ALIVE APISIX CONFIG

	modelAll := []string{config.HADOOP, config.HBASE, config.HIVE, config.KAFKA, config.MICROSERVICE,
		config.MYSQL, config.NODE, config.REDIS, config.SKYWALKING, config.SPARK, config.ZOOKEEPER,
		config.ZOOKEEPER, config.ALIVE, config.APISIX, config.CONFIG}

	data := strings.Join(modelAll, ",")
	if true == create(defaultPath, false, []byte("可存储的所有exporter信息")) {
		// 设置全部可以获取的exporter
		set(defaultPath, data)
		Logger.Printf("初始化zk全部exporter信息成功！")
	} else {
		Logger.Printf("创建zk节点失败, 尝试update全部exporter信息")
		set(defaultPath, data)
	}

}

func Consumer(ch chan string) {
	for i := 1; i < 2; i++ {
		model := <-ch
		Logger.Printf("获取已发布模块%s, 将写入zookeeper!", model)
		register(model)
		i--
	}

}

// 订阅zk的注册模块信息
// func subscribe(ch chan interface{}) {
// model := <-ch
// // 启动对应的exporter
// fmt.Printf("启动%s ...", model)
// 订阅zk目录，zk有新节点上线则
// }

func Callback(event zk.Event) {
	fmt.Println("*******************")
	fmt.Println("path:", event.Path)
	fmt.Println("type:", event.Type.String())
	fmt.Println("state:", event.State.String())
	fmt.Println("*******************")
}

// 添加watcher, 创建监听
func Watch(hosts []string) {
	//var hosts = []string{"192.168.10.220:2181", "192.168.10.221:2181", "192.168.10.222:2181"}
	option := zk.WithEventCallback(Callback)
	conn, _, err := zk.Connect(hosts, time.Second*5, option)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	var path = "/exporter/all_node"
	_, _, _, err = conn.ExistsW(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// // 创建节点
	// path = ""
	// create("test_watcher")
	// time.Sleep(time.Second * 2)
	// _, _, _, err = conn.ExistsW(path)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // 删除节点
	// delete(path)
}
