package main

import (
	"LogAgent/conf"
	"LogAgent/etcd"
	"LogAgent/kafka"
	"LogAgent/taillog"

	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/ini.v1"
)

var (
	cfg = new(conf.AgentConf)
)

func logInit(logdir string) {
	logfile := "LogAgent.log"
	if len(cfg.LogFileName) != 0 {
		logfile = cfg.LogFileName
	}
	logFile, err := os.OpenFile(filepath.Join(logdir, logfile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lmicroseconds | log.Ldate)
}

func main() {
	// 控制台参数
	ConfFile := flag.String("conf", "conf/conf.ini", "specify a config file")
	logDir := flag.String("logdir", "./", "specify Agent log dir path")
	flag.Parse()
	// 加载基础配置文件
	err := ini.MapTo(cfg, *ConfFile)
	if err != nil {
		fmt.Println("Load ini conf failed, err: ", err)
		return
	}

	// 初始化日志文件
	logInit(*logDir)

	//初始化Kafka
	if len(cfg.KafkaConf.Address) == 0 {
		fmt.Println("no Kafka Address !!!")
		return
	}
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.Logchansize)
	if err != nil {
		fmt.Println("init kafka failed, err: ", err)
		return
	}

	//初始化Etcd
	if len(cfg.EtcdConf.Address) == 0 {
		fmt.Println("no Etcd Address !!!")
		return
	}
	err = etcd.Init([]string{cfg.EtcdConf.Address}, time.Duration(cfg.EtcdConf.TimeOut)*time.Second)
	if err != nil {
		fmt.Println("init etcd failed,err: ", err)
		return
	}
	// 请求日志收集配置
	logentry, err := etcd.GetConf(cfg.EtcdConf.Key)
	if err != nil {
		fmt.Println("get key failed, err: ", err)
		return
	}

	//创建日志收集后台任务
	err = taillog.CreateTailTask(logentry)
	if err != nil {
		fmt.Println("Create Task faild, err: %v", err)
		return
	}

	var wg sync.WaitGroup
	newConfChan := taillog.NewEntryChan()
	wg.Add(1)
	go etcd.Watch(cfg.EtcdConf.Key, newConfChan)
	wg.Wait()

}
