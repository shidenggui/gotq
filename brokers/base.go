package brokers

type Broker interface {
	Dispatch([]byte) error
	Delay([]byte, string) error
	Receive(string) ([]byte, error)
}
