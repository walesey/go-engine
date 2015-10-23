package physics

import (
	"github.com/stretchr/testify/assert"
	vmath "github.com/walesey/go-engine/vectormath"
	"testing"
)

var pool *WorkerPool

func TestWorker(t *testing.T) {
	pool = NewWorkerPool()
	defer pool.Close()

	t1 := Triangle{vmath.Vector3{0, 0, 0}, vmath.Vector3{0, 1, 0}, vmath.Vector3{0, 0, 1}}
	t2 := Triangle{vmath.Vector3{0, 0, 0}, vmath.Vector3{0, 1, 0}, vmath.Vector3{1, 0, 0}}
	t3 := Triangle{vmath.Vector3{0, 1, 0}, vmath.Vector3{1, 0, 0}, vmath.Vector3{0, 0, 1}}
	triangles1 := []Triangle{t1, t2, t3}
	convexHull1 := NewConvexHull(triangles1)
	obj1 := NewPhysicsObject()
	convexHull1.AttachTo(&obj1)

	t4 := Triangle{vmath.Vector3{0, 0, 0}, vmath.Vector3{0, 1, 0}, vmath.Vector3{0, 0, -1}}
	t5 := Triangle{vmath.Vector3{0, 0, 0}, vmath.Vector3{0, 1, 0}, vmath.Vector3{-1, 0, 0}}
	t6 := Triangle{vmath.Vector3{0, 1, 0}, vmath.Vector3{-1, 0, 0}, vmath.Vector3{0, 0, -1}}
	triangles2 := []Triangle{t4, t5, t6}
	convexHull2 := NewConvexHull(triangles2)
	obj2 := NewPhysicsObject()
	convexHull2.AttachTo(&obj2)

	obj2.Position = vmath.Vector3{-0.3, 0, -0.3}
	workerA := pool.GetWorker()
	physicsPair := PhysicsPair{obj1, obj2}
	workerA.Write(physicsPair)

	obj2.Orientation = vmath.AngleAxis(3.14, vmath.Vector3{0, 1, 0})
	workerB := pool.GetWorker()
	physicsPair = PhysicsPair{obj1, obj2}
	workerB.Write(physicsPair)

	//assync collision tests can be processed concurrently
	assert.False(t, workerA.Read(), "Worker: ConvexHull Overlap should be false")
	assert.True(t, workerB.Read(), "Worker: ConvexHull Overlap should be true")

	pool.ReleaseWorker(workerA)
	pool.ReleaseWorker(workerB)
}
