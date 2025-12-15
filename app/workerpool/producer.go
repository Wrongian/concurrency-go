package workerpool

type Producer interface {
	Produce() (Job, bool)
}
