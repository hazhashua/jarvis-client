# metric_exporter

## support metric collect for many bigdata services like hbase, hadoop, redis and so on

# docker部署步骤
* gitlab 拉取代码: http://192.168.1.33:8080/bigdata-fusion/fusion/metric_exporter.git

* 在项目目录下编译dockerfile文件：docker build -t exporter_all .
* 编译只启动特定模块的dockerfile文件：docker build -t exporter_all -f ./Dockerfile_xxx .

* 启动docker应用：docker run --rm -it -p 48080:38080 -v /etc:/etc --name exporterall   exporter_all



