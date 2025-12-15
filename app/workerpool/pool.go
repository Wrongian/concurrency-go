package workerpool

import "sync"

type Pool struct {
	InChannel  chan Job
	OutChannel chan Result
	numWorkers int
	wg         sync.WaitGroup
	isStarted  bool
	mu         sync.Mutex
}

func NewPool(numWorkers int, bufferSize int) *Pool {
	inChannel := make(chan Job, bufferSize)
	outChannel := make(chan Result, bufferSize)
	return &Pool{
		InChannel:  inChannel,
		OutChannel: outChannel,
		numWorkers: numWorkers,
	}
}

func (pool *Pool) Start() {
	// Prevent race condition on isStarted
	pool.mu.Lock()

	if pool.isStarted {
		return
	}

	pool.isStarted = true

	pool.mu.Unlock()

	for range pool.numWorkers {
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()
			worker(pool.InChannel, pool.OutChannel)
		}()
	}

	go func() {
		pool.Wait()
		close(pool.OutChannel)
	}()

}

// Adds a single job
func (pool *Pool) AddJob(job Job) {
	pool.InChannel <- job
}

// Add multiple jobs
func (pool *Pool) AddJobs(jobs []Job) {
	for _, job := range jobs {
		pool.InChannel <- job
	}
}

// Wait for all the jobs to be processed
// If hasn't been started yet, just returns
func (pool *Pool) Wait() {
	pool.wg.Wait()
}

// After done producing everything close the input channel to close the pool
// Required to be done for the pool the workers to return
func (pool *Pool) Close() {
	close(pool.InChannel)
}

// Streams jobs as they come in (closes after producer is done)
//
// Only one producer
func (pool *Pool) Stream(producer Producer) {
	go func() {
		for {
			job, isDone := producer.Produce()
			if isDone {
				break
			}
			pool.AddJob(job)
		}
		pool.Close()
	}()
}

// Streams a group of producers
func (pool *Pool) StreamMany(producers []Producer) {
	var wg sync.WaitGroup

	wg.Add(len(producers))

	for _, producer := range producers {
		go func(p Producer) {
			defer wg.Done()
			for {
				job, isDone := p.Produce()
				if isDone {
					break
				}
				pool.AddJob(job)
			}
		}(producer)
	}

	go func() {
		// wait for all producers to be done
		wg.Wait()
		// close the pool when done
		pool.Close()
	}()
}

// todo add pipelining where you add stages to a pipeline with a transform function that transforms a result returns a job and a bool indicating whether to go forwards
// hard because of backpressure coordination that overflows buffers

// todo add fan-out dispatcher to multiple pools with distribution strategies
// create tests for all these
