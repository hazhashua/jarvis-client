package zookeeper

import (
	"fmt"
	"metric_exporter/utils"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

// 获取到zookeeper的连接
func conn() *zk.Conn {

	hosts := utils.GetZkHost()
	// hosts_str := strings.Join(hosts, ",")
	conn, _, err := zk.Connect(hosts, time.Second*5)
	defer conn.Close()
	if err != nil {
		utils.Logger.Println("zk connect fail  error: ", err.Error())
		return nil
	} else {
		return conn
	}
}

// 创建节点
func create(path string) bool {
	if path == "" {
		path = "/test_zk_node"
	}
	var conn *zk.Conn = conn()
	if conn != nil {
		defer conn.Close()
	}
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
