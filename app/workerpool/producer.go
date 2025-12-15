package workerpool

import "fmt"

type Producer interface {
	Produce() (Job, bool)
}

type TestProducer struct{}

func (p *TestProducer) Produce() (Job, bool) {
	return Job{
		Action: func(a Args) Result {

			num, exists := a[0].(int)

			if !exists {
				return Result{
					Payload: nil,
					Error:   fmt.Errorf("obj in channel not a number"),
				}
			}

			return Result{
				Payload: num + 1,
				Error:   nil,
			}
		},
		Args: []interface{}{1},
	}, true
}
