package config

const (
	HADOOP       string = "hadoop"
	HBASE        string = "hbase"
	HIVE         string = "hive"
	KAFKA        string = "kafka"
	MICROSERVICE string = "micro_service"
	MYSQL        string = "mysql"
	NODE         string = "node_exporter"
	REDIS        string = "redis"
	SKYWALKING   string = "skywalking"
	SPARK        string = "spark"
	ZOOKEEPER    string = "zookeeper"
	ALIVE        string = "alive"

	HADOOP_METRICPATH       string = "/hadoop/metrics"
	HBASE_METRICPATH        string = "/hbase/metrics"
	HIVE_METRICPATH         string = "/hive/metrics"
	KAFKA_METRICPATH        string = "/kafka/metrics"
	MICROSERVICE_METRICPATH string = "/micro_service/metrics"
	MYSQL_METRICPATH        string = "/mysql/metrics"
	NODE_METRICPATH         string = "/node/metrics"
	REDIS_METRICPATH        string = "/redis/metrics"
	SKYWALKING_METRICPATH   string = "/skywalking/metrics"
	SPARK_METRICPATH        string = "/spark/metrics"
	ZOOKEEPER_METRICPATH    string = "/zookeeper/metrics"
	ALIVE_METRICPATH        string = "/alive/metrics"
)

var MetricPathMap = map[string]string{
	HADOOP:       HADOOP_METRICPATH,
	HBASE:        HBASE_METRICPATH,
	HIVE:         HIVE_METRICPATH,
	KAFKA:        KAFKA_METRICPATH,
	MICROSERVICE: MICROSERVICE_METRICPATH,
	MYSQL:        MYSQL_METRICPATH,
	NODE:         NODE_METRICPATH,
	REDIS:        REDIS_METRICPATH,
	SKYWALKING:   SKYWALKING_METRICPATH,
	SPARK:        SPARK_METRICPATH,
	ZOOKEEPER:    ZOOKEEPER_METRICPATH,
	ALIVE:        ALIVE_METRICPATH,
}

// 保存注册的endpoint信息
var MetricIpMap = map[string]string{
	// HADOOP:       "",
	// HBASE:        "",
	// HIVE:         "",
	// KAFKA:        "",
	// MICROSERVICE: "",
	// MYSQL:        "",
	// NODE:         "",
	// REDIS:        "",
	// SKYWALKING:   "",
	// SPARK:        "",
	// ZOOKEEPER:    "",
	// ALIVE:        "",
}

type K8sConfig struct {
	ServiceURL      string `yaml:"serviceUrl,omitempty"`
	EndpointURL     string `yaml:"endpointUrl,omitempty"`
	NodeURL         string `yaml:"nodeUrl,omitempty"`
	PodURL          string `yaml:"podUrl,omitempty"`
	NodeResourceURL string `yaml:"nodeResourceUrl,omitempty"`
}

var k8sConfig K8sConfig = K8sConfig{
	// ServiceURL:  "http://124.65.131.14:38080/api/v1/services",
	// EndpointURL: "http://124.65.131.14:38080/api/v1/endpoints",
	ServiceURL:      "http://192.168.10.20:8080/api/v1/services",
	EndpointURL:     "http://192.168.10.20:8080/api/v1/endpoints",
	NodeURL:         "http://192.168.10.20:8080/api/v1/nodes",
	PodURL:          "http://192.168.10.20:8080/api/v1/pods",
	NodeResourceURL: "http://192.168.10.20:8080/apis/metrics.k8s.io/v1beta1/nodes",
}

type HbaseConfigure struct {
	Cluster struct {
		ClusterName         string   `yaml:"clustername"`
		MasterJmxPort       string   `yaml:"masterjmxport"`
		RegionserverJmxPort string   `yaml:"regionserverjmxport"`
		Hosts               []string `yaml:"hosts"`
		Names               []string `yaml:"names"`
		// names:
		//   - bigdata-dev01
		//   - bigdata-dev02
		//   - bigdata-dev03
	}
}

type HadoopConfigure struct {
	Cluster struct {
		Name                    string   `yaml:"name"`
		Services                []string `yaml:"services"`
		ServiceNum              int      `yaml:"servicenum"`
		Namenodes               []string `yaml:"namenodes"`
		NamenodeHosts           []string `yaml:"namenodehosts"`
		NamenodeHttpPort        int      `yaml:"namenodehttpport"`
		NamenodeRpcPort         int      `yaml:"namenoderpcport"`
		Datanodes               []string `yaml:"datanodes"`
		DatanodeHosts           []string `yaml:"datanodehosts"`
		DatanodeHttpPort        int      `yaml:"datanodehttpport"`
		DatanodeRpcPort         int      `yaml:"datanoderpcport"`
		ResourceManagers        []string `yaml:"resourcemanagers"`
		ResourceManagerHosts    []string `yaml:"resourcemanagerhosts"`
		ResourceManagerUrl      string   `yaml:"resourcemanagerurl"`
		ResourcemanagerHost     string   `yaml:"resourcemanagerhost"`
		ResourcemanagerHttpPort int      `yaml:"resourcemanagerhttpport"`
	}
}

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

type KafkaConfigure struct {
	Cluster struct {
		Name  string   `yaml:"name"`
		Hosts []string `yaml:"hosts"`
		Port  int      `yaml:"port"`
		Env   string   `yaml:"env"`
	}
}

type K8sYamlConfig struct {
	Cluster struct {
		Name          string   `yaml:"name"`
		Master        []string `yaml:"master"`
		Nodes         []string `yaml:"nodes"`
		ApiServerPort string   `yaml:"apiserverport"`
	}
}

// cluster:
//   name: bigdata-dev-cluster
//   ips:
//     - 192.168.10.70
//   port: 3306
//   username: root
//   password: pwd@123
//   defaultdb: information_schema
//   role:
//     - master

// 存储mysql配置文件信息
type MysqlConfig struct {
	Cluster struct {
		Name      string   `yaml:"name"`
		Ips       []string `yaml:"ips"`
		Port      int      `yaml:"port"`
		Username  string   `yaml:"username"`
		Password  string   `yaml:"password"`
		DefaultDB string   `yaml:"defaultdb"`
		Role      []string `yaml:"role"`
	}
}

type NodeConfig struct {
	// cluster:
	// 	name: bigdata-dev-cluster
	Cluster struct {
		Name string `name:"name"`
	}
}

type RedisConfig struct {
	Cluster struct {
		Name       string   `yaml:"name"`
		Ips        []string `yaml:"ips"`
		Hosts      []string `yaml:"hosts"`
		Ippwds     []string `yaml:"ippwds"`
		ScrapeHost string   `yaml:"scrapehost"`
		ScrapeIp   string   `yaml:"scrapeip"`
		RedisPort  int      `yaml:"redisport"`
	}
}

type SkyWalkingConfig struct {
	Cluster struct {
		Name          string `json:"name"`
		ElasticSearch struct {
			Ips  []string `json:"ips"`
			Port int      `json:"port"`
		}
	}
}

type SparkConfig struct {
	// cluster `yaml:"cluster"`
	// masterConf MasterConf `yaml:"masterHttp"`
	// applicationConf HttpConf   `yaml:"application_http"`
	Cluster    string `yaml:"cluster"`
	Masterhttp struct {
		Ips  []string `yaml:"ips"`
		Port int      `yaml:"port"`
		Path string   `yaml:"path"`
	}
	Workerhttp struct {
		Ips  []string `yaml:"ips"`
		Port int      `yaml:"port"`
		Path string   `yaml:"path"`
	}
	Applicationhttp struct {
		Ips          []string `yaml:"ips"`
		Ports        []int    `yaml:"ports"`
		MainPath     string   `yaml:"mainpath"`
		ExecutorPath string   `yaml:"executorpath"`
	}
}

// cluster:
//   name: test环境zookeeper
//   hosts:
//     - 192.168.10.220
//     - 192.168.10.221
//     - 192.168.10.222
//   clientport: 2181

type ZookeepeConfig struct {
	Cluster struct {
		Name       string   `yaml:"name"`
		Hosts      []string `yaml:"hosts"`
		ClientPort string   `yaml:"clientport"`
	}
}

// 存储utils工具包, 数据库源地址配置数据
type DbConfigure struct {
	Cluster struct {
		Name  string
		Mysql struct {
			Ip       string
			Port     int32
			Username string
			Password string
		}
		Postgres struct {
			Ip             string
			Port           int32
			Username       string
			Password       string
			DatasourceInfo struct {
				Schema string
			}

			GatherTable       string // gather_name
			GatherDetailTable string // data_gather_configure
			ExportTable       string // data_store_configure

		}
		HttpPort int
	}
}
