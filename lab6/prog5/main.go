package main

import (
	"fmt"
)

type CalcRequest struct {
	A, B float64
	Op   rune
	Resp chan float64
}

func worker(reqs <-chan CalcRequest) {

	for req := range reqs {

		var res float64

		switch req.Op {
		case '+':
			res = req.A + req.B
		case '-':
			res = req.A - req.B
		case '*':
			res = req.A * req.B
		case '/':
			res = req.A / req.B
		}
		req.Resp <- res
	}
}

func main() {

	reqs := make(chan CalcRequest)

	for i := 0; i < 5; i++ {
		go worker(reqs)
	}

	jobs := []CalcRequest{
		{A: 1, B: 2, Op: '+', Resp: make(chan float64)},
		{A: 5, B: 3, Op: '-', Resp: make(chan float64)},
		{A: 5, B: 3, Op: '*', Resp: make(chan float64)},
		{A: 5, B: 2, Op: '/', Resp: make(chan float64)},
	}

	go func() {
		for _, j := range jobs {
			reqs <- j
		}
		close(reqs)
	}()

	for _, j := range jobs {
		fmt.Printf("%.2f\n", <-j.Resp)
	}

}
