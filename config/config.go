package config

type Config struct {
	Broker *BrokerCfg
}

type BrokerCfg struct {
	Host     string
	Port     int64
	Password string
	DB       int64
}
