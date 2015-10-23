package physics

import "fmt"

type PhysicsPair struct {
	object1, object2 PhysicsObject
}

type PhysicsWorker struct {
	in   chan PhysicsPair
	out  chan bool
	kill bool
}

type WorkerPool struct {
	pool []*PhysicsWorker
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{make([]*PhysicsWorker, 0, 32)}
}

func (wp *WorkerPool) GetWorker() *PhysicsWorker {
	if len(wp.pool) > 0 {
		worker := wp.pool[len(wp.pool)-1]
		wp.pool = wp.pool[:len(wp.pool)-1]
		return worker
	}
	return NewPhysicsWorker()
}

func (wp *WorkerPool) ReleaseWorker(worker *PhysicsWorker) {
	wp.pool = append(wp.pool, worker)
}

func (wp *WorkerPool) Close() {
	for _, worker := range wp.pool {
		worker.Close()
	}
}

func NewPhysicsWorker() *PhysicsWorker {
	in := make(chan PhysicsPair)
	out := make(chan bool)
	worker := &PhysicsWorker{in: in, out: out}
	go worker.run()
	return worker
}

// run worker calculates narrow phase overlaps
func (worker *PhysicsWorker) run() {
	for {
		state := <-worker.in
		if worker.kill {
			break
		}
		worker.out <- state.object1.NarrowPhaseOverlap(state.object2)
	}
	fmt.Println("Physics Worker Died.")
}

// Write new data to be processed if worker is busy an error is returned
func (worker *PhysicsWorker) Write(pair PhysicsPair) {
	worker.in <- pair
}

func (worker *PhysicsWorker) Read() bool {
	return <-worker.out
}

func (worker *PhysicsWorker) Close() {
	worker.kill = true
	worker.Write(PhysicsPair{NewPhysicsObject(), NewPhysicsObject()})
}
