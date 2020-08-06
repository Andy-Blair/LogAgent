package conf

// KafkaConf ...
type KafkaConf struct {
	Address     string `ini:"address"`
	Logchansize int    `ini:"log_channel_size"`
}

// EtcdConf ...
type EtcdConf struct {
	Address string `ini:"address"`
	TimeOut int    `ini:"timeout"`
	Key     string `ini:"key"`
}

// AgentConf ...
type AgentConf struct {
	KafkaConf   `ini:"kafka"`
	EtcdConf    `ini:"etcd"`
	LogFileName string `ini:"logfile"`
}
