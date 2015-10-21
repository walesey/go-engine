package physics

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var pool *WorkerPool

func TestMain(m *testing.M) {
	//do setup
	pool = NewWorkerPool()
	//run tests
	os.Exit(m.Run())
}

func TestWorker(t *testing.T) {
	//	worker := pool.GetWorker()
	//	physicsPair := PhysicsPair{PhysicsObject1}
	assert.EqualValues(t, false, false, "")
}
