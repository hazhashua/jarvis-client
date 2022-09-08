package zookeeper

import (
	"fmt"
	"metric_exporter/config"
	"metric_exporter/utils"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

// 获取到zookeeper的连接
func conn() *zk.Conn {
	var zookeeper_config config.ZookeepeConfig
	var ok bool

	zookeeper_config, ok = (utils.ConfigStruct.ConfigData[config.ZOOKEEPER]).(config.ZookeepeConfig)
	if ok == false {
		fmt.Println("获取zookeeper config配置失败!")
		utils.Logger.Println("获取zookeeper config配置失败!")
		return nil
	}

	fmt.Println("zookeeper_config.Cluster.Name: ", zookeeper_config.Cluster.Name)
	fmt.Println("zookeeper_config.Cluster.Hosts: ", zookeeper_config.Cluster.Hosts)
	fmt.Println("zookeeper_config.Cluster.ClientPort: ", zookeeper_config.Cluster.ClientPort)
	hosts := make([]string, 0)
	for _, host := range zookeeper_config.Cluster.Hosts {
		hosts = append(hosts, fmt.Sprintf("%s:%s", host, zookeeper_config.Cluster.ClientPort))
	}
	// hosts_str := strings.Join(hosts, ",")
	conn, _, err := zk.Connect(hosts, time.Second*5)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		fmt.Println("连接成功! ")
		return conn
	}
}

// 创建节点
func create(path string) bool {
	if path == "" {
		path = "/test_zk_node"
	}
	var conn *zk.Conn = conn()
	defer conn.Close()
	var data = []byte("test zk node 数据")
	//flags有4种取值：
	//0:永久，除非手动删除
	//zk.FlagEphemeral = 1:短暂，session断开则改节点也被删除
	//zk.FlagSequence  = 2:会自动在节点后面添加序号
	//3:Ephemeral和Sequence，即，短暂且自动添加序号
	var flags int32 = 0
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
func set(path string) {
	if path == "" {
		path = "/brokers"
	}
	connection := conn()
	defer connection.Close()
	// var path = "home"
	fmt.Println("connection.SessionID(): ", connection.SessionID())
	data := []byte("测试节点数据")
	connection.Set(path, data, -1)

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

func Callback(event zk.Event) {
	fmt.Println("*******************")
	fmt.Println("path:", event.Path)
	fmt.Println("type:", event.Type.String())
	fmt.Println("state:", event.State.String())
	fmt.Println("*******************")
}

// 添加watcher, 创建监听
func Watch() {
	var hosts = []string{"192.168.10.220:2181", "192.168.10.221:2181", "192.168.10.222:2181"}
	option := zk.WithEventCallback(Callback)
	conn, _, err := zk.Connect(hosts, time.Second*5, option)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	var path = "/testzk"
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
