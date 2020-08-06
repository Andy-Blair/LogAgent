package taillog

import (
	"context"
	"log"
	"time"

	"LogAgent/etcd"
	"LogAgent/kafka"

	"github.com/hpcloud/tail"
)

var (
	tasks       = make(map[string]*tailTask)
	newLogEntry = make(chan []*etcd.LogEntryConf)
)

type tailTask struct {
	path       string
	topic      string
	tail       *tail.Tail
	cancelFunc context.CancelFunc
}

//CreateTailTask 循环所有日志收集条目，创建日志收集任务
func CreateTailTask(conf []*etcd.LogEntryConf) (err error) {
	for _, logconf := range conf {
		err = newTailTask(logconf)
		if err != nil {
			return
		}
	}
	go UpdateTailTask()
	return nil
}

func newTailTask(logconf *etcd.LogEntryConf) (err error) {
	logPath := logconf.Path
	logTopic := logconf.Topic
	tk := tailTask{
		path:  logPath,
		topic: logTopic,
	}
	tk.init()
	ctx, cancel := context.WithCancel(context.Background())
	tk.cancelFunc = cancel
	tasks[logPath+"&"+logTopic] = &tk
	go tk.sendToLogChan(ctx, logTopic)
	return nil
}

// Init 初始化打开文件
func (t *tailTask) init() {
	config := tail.Config{
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件那个地方开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      true,
	}
	var err error
	t.tail, err = tail.TailFile(t.path, config)
	if err != nil {
		log.Printf("Init Task faild, err: %v", err)
	}
}

// sendToLogChan 向kafka日志通道发送日志信息
func (t *tailTask) sendToLogChan(ctx context.Context, topic string) {
	for {
		select {
		case line := <-t.tail.Lines:
			kafka.LogChan(topic, line.Text)
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second)
		}
	}
}

// NewEntryChan 接收更新收集日志条目通道
func NewEntryChan() chan<- []*etcd.LogEntryConf {
	return newLogEntry
}

//UpdateTailTask 更新后台日志收集任务
func UpdateTailTask() {
	for {
		logentrys := <-newLogEntry
		var newEntryMap = make(map[string]int)
		for _, entry := range logentrys {
			path := entry.Path
			topic := entry.Topic
			taskKey := path + "&" + topic
			newEntryMap[taskKey] = 1
			_, ok := tasks[taskKey]
			if !ok {
				_ = newTailTask(entry)
				log.Printf("新增任务%s\n", taskKey)
			}
		}
		for k, v := range tasks {
			_, ok := newEntryMap[k]
			if !ok {
				v.cancelFunc()
				log.Printf("%s任务已结束\n", k)
			}
		}
	}
}
