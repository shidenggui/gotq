package brokers

type Broker interface {
	Request(string, int64) ([]byte, error)
	Delay([]byte, string) error
	QuickDelay([]byte, string) error
	Receive(string) ([]byte, error)
	Expire(string, int64) error
}
