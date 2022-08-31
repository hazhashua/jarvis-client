package config

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
