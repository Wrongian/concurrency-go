package workerpool

type Job struct {
	Action func(Args) Result
	Args   Args
}

func (job Job) Execute() Result {
	return job.Action(job.Args)
}
