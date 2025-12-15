package workerpool

type Consumer interface {
	Consume(Result)
}
