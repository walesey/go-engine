package physics

import "fmt"

type PhysicsPair struct {
	object1, object2 PhysicsObject
}

type PhysicsWorker struct {
	in  chan PhysicsPair
	out chan bool
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

func NewPhysicsWorker() *PhysicsWorker {
	in := make(chan PhysicsPair)
	out := make(chan bool)
	worker := &PhysicsWorker{in, out}
	go worker.run()
	return worker
}

// run worker calculates narrow phase overlaps
func (worker *PhysicsWorker) run() {
	for {
		state := <-worker.in
		worker.out <- state.object1.NarrowPhaseOverlap(state.object2)
	}
	fmt.Println("Physics Worker Died.")
}

// Write new data to be processed if worker is busy an error is returned
func (worker *PhysicsWorker) Write(pair PhysicsPair) error {
	select {
	case worker.in <- pair:
		return nil
	default:
		return fmt.Errorf("worker is busy")
	}
}
