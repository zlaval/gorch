package main

import "Gorch/internal/worker"

func main() {
	a := worker.Cmd("worker-1")
	_ = a.Run()
}
