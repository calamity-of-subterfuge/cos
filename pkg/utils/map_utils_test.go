package utils_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/calamity-of-subterfuge/cos/pkg/utils"
	"github.com/jakecoffman/cp"
)

func TestTriangle_Sample_simple(t *testing.T) {
	// increase number of iterations when debugging; can fix the
	// seed to try the relevant example
	var i int64
	for i = 0; i < 1_000; i++ {
		rand.Seed(i)

		tri := utils.NewTriangle(cp.Vector{X: 1, Y: 1}, cp.Vector{X: 2, Y: 0}, cp.Vector{X: 1, Y: 2})
		sample := tri.Sample()

		if !tri.Contains(sample) {
			t.Errorf("seed %d, sample %q not contained in triangle %q", i, sample, tri.Vertices)
			break
		}
	}
}

func TestTriangle_Sample_rand(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	tri := utils.MapHexInnerEdgeTriangularization.Triangles[rand.Intn(len(utils.MapHexEdgeTriangularization.Triangles))]
	for i := 0; i < 1000; i++ {
		sample := tri.Sample()
		if !tri.Contains(sample) {
			t.Errorf("For triangle verts = %q, sample = %q, not contained", tri.Vertices, sample)
		}
	}
}

func BenchmarkTriangle_Sample(b *testing.B) {
	tri := utils.NewTriangle(cp.Vector{X: 1, Y: 1}, cp.Vector{X: 2, Y: 0}, cp.Vector{X: 1, Y: 2})

	for i := 0; i < b.N; i++ {
		tri.Sample()
	}
}
