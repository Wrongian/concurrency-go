package workerpool

func worker(jobCh chan Job, resCh chan Result) {
	for job := range jobCh {
		res := job.Execute()

		resCh <- res
	}
}
