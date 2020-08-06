# LogAgent 日志收集客户端
## V 1.0.0
```
-conf 指定基础配置文件，参考 conf/conf.ini。
-logdir 指定程序日志目录
```
1.读取基础配置

2.连接Kafka并启动日志发送线程

3.从ETCD请求日志收集配置，配置信息为json格式，包含日志文件绝对路径、Kafka推送Topic。

###### etcd日志收集配置样例：
`[{"path":"path/to/logfile","topic":"send_to_topic"}...]`

4.根据配置信息创建日志收集任务

5.监听ETCD配置变更事件，动态更新（新增、停止）日志收集任务
