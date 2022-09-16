package skywalking

// 对elasticsearch进行增删改查的基础类

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"metric_exporter/config"
	"metric_exporter/utils"
	"reflect"
	"regexp"
	"time"

	"github.com/olivere/elastic/v7"
)

var client *elastic.Client
var host string
var hosts []string

type Employee struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	About     string   `json:"about"`
	Interests []string `json:"interests"`
}

//初始化
func init() {
	//errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
	// skywalkingConfig := ParseSkyWalkingConfig()
	skywalkingConfig, _ := (utils.ConfigStruct.ConfigData[config.SKYWALKING]).(config.SkyWalkingConfig)

	for _, ip := range skywalkingConfig.Cluster.ElasticSearch.Ips {
		hostUrl := fmt.Sprintf("http://%s:%d", ip, skywalkingConfig.Cluster.ElasticSearch.Port)
		hosts = append(hosts, hostUrl)
	}
	if len(hosts) == 0 {
		utils.Logger.Printf("读取skywalking配置为空，退出init函数")
		return
	}
	host = hosts[0]
	var err error
	//这个地方有个小坑 不加上elastic.SetSniff(false) 会连接不上
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(host))
	if err != nil {
		utils.Logger.Printf("elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(host)) error: %s", err.Error())
		panic(err)
	}
	info, code, err1 := client.Ping(host).Do(context.Background())
	if err1 != nil {
		utils.Logger.Printf("client.Ping(host).Do(context.Background()) error: %s", err.Error())
		panic(err1)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	var esversion string
	esversion, err = client.ElasticsearchVersion(host)
	if err != nil {
		utils.Logger.Printf("client.ElasticsearchVersion(host) error: %s", err.Error())
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

}

/*下面是简单的CURD*/

//创建
func create() {

	//使用结构体
	e1 := Employee{"Jane", "Smith", 32, "I like to collect rock albums", []string{"music"}}
	put1, err := client.Index().
		Index("megacorp").
		Type("employee").
		Id("1").
		BodyJson(e1).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)

	//使用字符串
	e2 := `{"first_name":"John","last_name":"Smith","age":25,"about":"I love to go rock climbing","interests":["sports","music"]}`
	put2, err := client.Index().
		Index("megacorp").
		Type("employee").
		Id("2").
		BodyJson(e2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put2.Id, put2.Index, put2.Type)

	e3 := `{"first_name":"Douglas","last_name":"Fir","age":35,"about":"I like to build cabinets","interests":["forestry"]}`
	put3, err := client.Index().
		Index("megacorp").
		Type("employee").
		Id("3").
		BodyJson(e3).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put3.Id, put3.Index, put3.Type)

}

type SkwEvent struct {
	StartTime       uint64 `json:"start_time"`
	EndPoint        string `json:"endpoint"`
	Service         string `json:"service"`
	Name            string `json:"name"`
	EndTime         uint64 `json:"end_time"`
	TimeBucket      uint64 `json:"time_bucket"`
	ServiceInstance string `json:"service_instance"`
	Type            string `json:"type"`
	Message         string `json:"message"`
	Uuid            string `json:"uuid"`
	Parameters      string `json:"parameters"`
}

type SkwEndpointTraffic struct {
	ServiceId  string `json:"service_id"`
	Name       string `json:"name"`
	TimeBucket uint64 `json:"time_bucket"`
}

type SkwInstanceTraffic struct {
	LastPing   uint64 `json:"last_ping"`
	ServiceId  string `json:"service_id"`
	Name       string `json:"name"`
	TimeBucket uint64 `json:"time_bucket"`
	Properties string `json:"properties"`
}

type PrintInterface interface {
	printStruct(e interface{})
}

// 查找index, type下的全部数据
func GetAll(index string, types string, typ interface{}) (results []interface{}) {
	//通过id查找
	// get1, err := client.Get().Index("megacorp").Type("employee").Id("2").Do(context.Background())
	getRes, err := client.Search().Size(10000).Index(index).Type(types).Do(context.Background())
	if err != nil {
		// panic(err)
		fmt.Println("查询es错误: ", err.Error())
		utils.Logger.Printf("查询es错误: %s", err.Error())
		return nil
	}

	// var typ skwEvent
	printEvents(getRes, err, typ)

	typet := reflect.TypeOf(typ)
	results = getRes.Each(typet)
	return results
}

type MyCpmInfo struct {
	MetricTable string `json:"metric_table"`
	Value       int    `json:"value"`
	ServiceId   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	EntityId    string `json:"entity_id"`
	TimeBucket  int    `json:"time_bucket"`
	Entity      string `json:"entity"`
}

type CpmInfo struct {
	MetricTable string `json:"metric_table"`
	Total       int    `json:"total"`
	ServiceId   string `json:"service_id"`
	TimeBucket  int    `json:"time_bucket"`
	EntityId    string `json:"entity_id"`
	Value       int    `json:"value"`
}

// 抓取当天的service 和 endpoint cpm
func GetCpmInfo(metricTable string) (cpmInfo []MyCpmInfo) {
	if metricTable != "service_instance_cpm" && metricTable != "endpoint_cpm" {
		fmt.Println("参数错误...")
		return
	}

	// current := time.Now()
	beforeOneM := time.Now().Add(time.Duration(-1000000000 * 60 * 485))
	year, m, day := beforeOneM.Date()
	index := fmt.Sprintf("sw_metrics-cpm-%04d%02d%02d", year, m, day)
	fmt.Println("index: ", index)

	// index = "sw_metrics-cpm-20220817"
	// fmt.Println("index: ", index)

	cpms := make([]MyCpmInfo, 0)
	var cpm CpmInfo
	//"service_instance_cpm"
	termQuery := elastic.NewTermQuery("metric_table", metricTable)
	timeint := year*10000*10000 + int(m)*10000*100 + day*10000 + beforeOneM.Hour()*100 + beforeOneM.Minute()
	fmt.Println("timeint: ", timeint)
	rangeQuery := elastic.NewRangeQuery("time_bucket").Gte(timeint)

	boolQuery := elastic.NewBoolQuery().Must(termQuery, rangeQuery)

	// queryStr := elastic.NewQueryStringQuery("metric_table:service_instance_cpm")
	if searchRs, err := client.Search(index).Query(boolQuery).Size(10000).Do(context.Background()); err == nil {
		fmt.Println("搜索到数据...", searchRs.Hits)

		fmt.Println("hit数组长度为: ", searchRs.TotalHits())
		if searchRs.TotalHits() <= 0 {
			fmt.Println("hit数组长度<=0")
		}
		for _, rs := range searchRs.Each(reflect.TypeOf(cpm)) {
			cpmin := rs.(CpmInfo)
			fmt.Println("遍历搜索到的数据...", " serviceid: ", cpmin.ServiceId)
			// 转义service_id和entity_id信息
			s := analysisItem(cpmin.ServiceId)
			if len(s) == 2 || len(s) == 1 {
				fmt.Println("解析后的service_id数据: ", s[0])
			}
			se := analysisItem(cpmin.EntityId)
			if len(se) == 1 {
				fmt.Println("解析后的service_id数据: ", se[0])
			} else if len(se) == 2 {
				fmt.Println("解析后的service_id数据: ", se[0])
				fmt.Println("解析后的实例或endpoint数据: ", se[1])
			}

			cpms = append(cpms, MyCpmInfo{
				MetricTable: cpmin.MetricTable,
				Value:       cpmin.Value,
				ServiceId:   cpmin.ServiceId,
				ServiceName: se[0],
				EntityId:    cpmin.EntityId,
				Entity:      se[1],
				TimeBucket:  cpmin.TimeBucket,
			})
		}
	}

	return cpms

	// getRes, err := client.Search().Size(10000).Index(index).Type(types).Do(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	// // var typ skwEvent
	// printEvents(getRes, err, typ)

	// // var typ skwEvent
	// typet := reflect.TypeOf(typ)
	// results = getRes.Each(typet)
	// return results
}

// 根据base64密文解析成明文数据
func analysisItem(item string) []string {
	r := regexp.MustCompile("(.*)(\\.\\d+[_-]?)(.*)")
	bbs := r.FindSubmatch([]byte(item))
	fmt.Println("len(findsubmatch): ", len(bbs))
	if len(bbs) == 2 || len(bbs) == 3 {
		base641, _ := base64.StdEncoding.DecodeString(string(bbs[1]))
		return []string{string(base641)}
		// base642, _ := base64.StdEncoding.DecodeString(string(bbs[2]))
		// fmt.Println("match 2: ", string(bbs[2]), string(base642))
	} else if len(bbs) == 4 {
		base641, _ := base64.StdEncoding.DecodeString(string(bbs[1]))
		fmt.Println("match 1: ", string(bbs[1]), string(base641))
		// base642, _ := base64.StdEncoding.DecodeString(string(bbs[2]))
		// fmt.Println("match 2: ", string(bbs[2]), string(base642))
		base643, _ := base64.StdEncoding.DecodeString(string(bbs[3]))
		fmt.Println("match 3: ", string(bbs[3]), string(base643))
		return []string{string(base641), string(base643)}
	}

	return []string{}
}

// 获取es中的一条数据
func gets() {
	//通过id查找
	get1, err := client.Get().Index("megacorp").Type("employee").Id("2").Do(context.Background())
	if err != nil {
		panic(err)
	}

	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
		var se SkwEvent
		err := json.Unmarshal(get1.Source, &se)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println()
		fmt.Println(string(get1.Source))
	}

}

//
//删除
func delete() {
	res, err := client.Delete().Index("megacorp").
		Type("employee").
		Id("1").
		Do(context.Background())
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Printf("delete result %s\n", res.Result)
}

//
//修改
func update() {
	res, err := client.Update().
		Index("megacorp").
		Type("employee").
		Id("2").
		Doc(map[string]interface{}{"age": 88}).
		Do(context.Background())
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("update age %s\n", res.Result)

}

//
////搜索
func query() {
	var res *elastic.SearchResult
	var err error
	//取所有
	res, err = client.Search("megacorp").Type("employee").Do(context.Background())
	printEmployee(res, err)

	//字段相等
	q := elastic.NewQueryStringQuery("last_name:Smith")
	res, err = client.Search("megacorp").Type("employee").Query(q).Do(context.Background())
	if err != nil {
		println(err.Error())
	}
	printEmployee(res, err)

	//条件查询
	//年龄大于30岁的
	boolQ := elastic.NewBoolQuery()
	boolQ.Must(elastic.NewMatchQuery("last_name", "smith"))
	boolQ.Filter(elastic.NewRangeQuery("age").Gt(30))
	res, err = client.Search("megacorp").Type("employee").Query(q).Do(context.Background())
	printEmployee(res, err)

	//短语搜索 搜索about字段中有 rock climbing
	matchPhraseQuery := elastic.NewMatchPhraseQuery("about", "rock climbing")
	res, err = client.Search("megacorp").Type("employee").Query(matchPhraseQuery).Do(context.Background())
	printEmployee(res, err)

	//分析 interests
	aggs := elastic.NewTermsAggregation().Field("interests")
	res, err = client.Search("megacorp").Type("employee").Aggregation("all_interests", aggs).Do(context.Background())
	printEmployee(res, err)

}

//
////简单分页
func list(size, page int) {
	if size < 0 || page < 1 {
		fmt.Printf("param error")
		return
	}
	res, err := client.Search("megacorp").
		Type("employee").
		Size(size).
		From((page - 1) * size).
		Do(context.Background())
	printEmployee(res, err)

}

//
//打印查询到的Employee
func printEmployee(res *elastic.SearchResult, err error) {
	if err != nil {
		print(err.Error())
		return
	}
	var typ Employee
	for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
		t := item.(Employee)
		fmt.Printf("%#v\n", t)
	}
}

// 从结果中查询event数据
func printEvents(res *elastic.SearchResult, err error, typ interface{}) {
	if err != nil {
		print(err.Error())
		return
	}
	// var typ skwEvent
	typet := reflect.TypeOf(typ)
	for _, item := range res.Each(typet) { //从搜索结果中取数据的方法
		// t := item.(SkwEvent)
		fmt.Printf("%#v\n", item)
	}
}

func main() {
	// create()
	// delete()
	// update()
	gets()
	query()
	list(2, 1)
}
