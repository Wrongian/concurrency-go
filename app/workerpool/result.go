package workerpool

type Result struct {
	Payload interface{}
	Error   error
}

func ErrResult(e error) Result {
	return Result{
		Payload: nil,
		Error:   e,
	}
}

func SuccessResult(payload interface{}) Result {
	return Result{
		Payload: payload,
		Error:   nil,
	}
}
