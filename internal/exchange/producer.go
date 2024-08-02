package exchange

type Producer interface {
	Publish(body []byte) error
}
