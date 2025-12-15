package main

import (
	"concurrency/app/workerpool"
	"fmt"
	"os"
	"strconv"
)

// extremely complicated way to using dp to get nth fibonacci number
func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("usage: <program name> <int>")
		os.Exit(1)
	}

	n, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("parameter not an integer")
		os.Exit(1)
	}

	if n <= 0 {
		fmt.Println("parameter cannot be negative")
		os.Exit(1)
	}

	num, err := fib(n)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("fibonacci number #%d is %d", n, num)
}

type FibPayload struct {
	isDone   bool
	nextArgs workerpool.Args
	resNum   int
}

func fibAction(args workerpool.Args) workerpool.Result {
	first, ok := args[0].(int)
	if !ok {
		return workerpool.ErrResult(fmt.Errorf("first argument invalid"))
	}

	second, ok := args[1].(int)
	if !ok {
		return workerpool.ErrResult(fmt.Errorf("second argument invalid"))
	}

	left, ok := args[2].(int)
	if !ok {
		return workerpool.ErrResult(fmt.Errorf("third argument invalid"))
	}

	sum := first + second

	left -= 1
	if left == 0 {
		return workerpool.SuccessResult(
			FibPayload{
				isDone:   true,
				nextArgs: workerpool.Args{},
				resNum:   sum,
			},
		)
	}

	return workerpool.SuccessResult(
		FibPayload{
			isDone:   false,
			nextArgs: workerpool.Args{second, sum, left},
			resNum:   -1,
		},
	)
}

func fibJob(fibNum int) workerpool.Job {
	job := workerpool.Job{
		Action: fibAction,
		Args:   workerpool.Args{0, 1, fibNum},
	}

	return job
}

func fib(num int) (int, error) {
	// base cases
	if num == 0 {
		return 0, nil
	}

	if num == 1 {
		return 1, nil
	}

	pool := workerpool.NewPool(2, 2)

	// num - 1 iterations
	initJob := fibJob(num - 1)
	pool.AddJob(initJob)
	pool.Start()
	var res int
	done := make(chan bool)

	go func() {
		for result := range pool.OutChannel {
			fibPayload := result.Payload.(FibPayload)
			if !fibPayload.isDone {
				pool.AddJob(workerpool.Job{
					Action: fibAction,
					Args:   fibPayload.nextArgs,
				})
			} else {
				res = fibPayload.resNum
				pool.Close()
				done <- true
				return
			}
		}
	}()
	// wait till done
	<-done
	return res, nil
}
