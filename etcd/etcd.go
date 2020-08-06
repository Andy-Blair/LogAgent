package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var (
	cli *clientv3.Client
)

// LogEntryConf ...
type LogEntryConf struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

// Init 初始化Etcd连接
func Init(addr []string, timeout time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: timeout,
	})
	if err != nil {
		return
	}
	return
}

// GetConf 从Etcd获取数据
func GetConf(key string) (logEntryConf []*LogEntryConf, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	rep, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		return
	}
	for _, ev := range rep.Kvs {
		err = json.Unmarshal(ev.Value, &logEntryConf)
		if err != nil {
			return nil, err
		}
	}
	return logEntryConf, nil

}

// Watch 在Etcd创建监听
func Watch(key string, confChan chan<- []*LogEntryConf) {
	wh := cli.Watch(context.Background(), key)
	for wresp := range wh {
		for _, evt := range wresp.Events {
			fmt.Printf("Type:%v  Key:%v  Value:%v\n", evt.Type, string(evt.Kv.Key), string(evt.Kv.Value))
			var newEntry []*LogEntryConf
			_ = json.Unmarshal(evt.Kv.Value, &newEntry)
			confChan <- newEntry
		}
	}
}
