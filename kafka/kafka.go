package kafka

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
)

type logData struct {
	topic string
	data  string
}

var (
	client sarama.SyncProducer
	logMsg chan *logData
)

// Init 初始化连接
func Init(addr []string, logchansize int) (err error) {
	// 配置生产者
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	//连接kafka
	client, err = sarama.NewSyncProducer(addr, config)
	if err != nil {
		// fmt.Println("connection kafka faile")
		return
	}
	logMsg = make(chan *logData, logchansize)
	go sendMsg()
	return
}

// SendMsg 发送消息到Kafka
func sendMsg() {
	for {
		select {
		case logmsg := <-logMsg:
			//构造消息
			message := &sarama.ProducerMessage{}
			message.Topic = logmsg.topic
			message.Value = sarama.StringEncoder(logmsg.data)
			//发送消息
			_, _, err := client.SendMessage(message)
			if err != nil {
				log.Println("send msg failed, err:", err)
			}
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}

// LogChan 向logMsg存储即将发送的日志数据结构体指针
func LogChan(topic, data string) {
	ld := logData{
		topic: topic,
		data:  data,
	}
	logMsg <- &ld
}
