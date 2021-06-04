package utils

import (
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/jakecoffman/cp"
)

// Triangle describes a 2d simplex
type Triangle struct {
	// Vertices are the vertices of the 2d simplex; they must not
	// fall on a line
	Vertices [3]cp.Vector

	// Basis is a 2x2 matrix that makes up the basis for this triangle.
	// If V is a 2x1 vector such that both components are in the unit
	// interval, then BV is a point within the triangle. This can be used
	// to sample from the triangle quickly.
	//
	// the first row is Basis[0], Basis[1]. The second row is Basis[2] Basis[3]
	Basis [4]float64
}

// NewTriangle computes the basis of the given triangle then
// constructs it
func NewTriangle(a, b, c cp.Vector) *Triangle {
	basisA := b.X - a.X
	basisB := c.X - a.X
	basisC := b.Y - a.Y
	basisD := c.Y - a.Y

	return &Triangle{
		Vertices: [3]cp.Vector{a, b, c},
		Basis:    [4]float64{basisA, basisB, basisC, basisD},
	}
}

// Area computes the area of this triangle.
func (t *Triangle) Area() float64 {
	det := t.Basis[0]*t.Basis[3] - t.Basis[1]*t.Basis[2]
	return math.Abs(0.5 * det)
}

// Sample returns a random vector in this triangle
func (t *Triangle) Sample() cp.Vector {
	x := rand.Float64()
	y := rand.Float64()

	if x+y > 1 {
		x = 1 - x
		y = 1 - y
	}

	return cp.Vector{
		X: t.Vertices[0].X + t.Basis[0]*x + t.Basis[1]*y,
		Y: t.Vertices[0].Y + t.Basis[2]*x + t.Basis[3]*y,
	}
}

// Contains checks if the vector v is in this triangle. This is not
// the primary intended use of triangles within this package
// right now and is not optimized
func (t *Triangle) Contains(v cp.Vector) bool {
	det := t.Basis[0]*t.Basis[3] - t.Basis[1]*t.Basis[2]
	invDet := 1 / det

	invBasis := [4]float64{
		invDet * t.Basis[3], -invDet * t.Basis[1],
		-invDet * t.Basis[2], invDet * t.Basis[0],
	}

	relPt := v.Sub(t.Vertices[0])
	bx := invBasis[0]*relPt.X + invBasis[1]*relPt.Y
	if bx < 0 {
		return false
	}

	by := invBasis[2]*relPt.X + invBasis[3]*relPt.Y
	if by < 0 {
		return false
	}

	return (bx + by) <= 1
}

// Triangularization is a triangularization of some polygon which
// can rapidly sampled from
type Triangularization struct {
	// Triangles contains the triangles within this triangularization.
	// In general this should sorted according to descending area for
	// best performance, but this is not required
	Triangles []Triangle

	// AreaPartialSums contains the sum of the area of all triangles up
	// to and including the index.
	AreaPartialSums []float64
}

// NewTriangularization returns a new triangularization produced by
// the given triangles. The resulting triangularization is a copy of
// the given slice sorted according to descending area, with area
// partial sums filled in.
func NewTriangularization(triangles []Triangle) *Triangularization {
	newTriangles := append(make([]Triangle, 0, len(triangles)), triangles...)
	sort.Slice(newTriangles, func(i, j int) bool { return newTriangles[i].Area() > newTriangles[j].Area() })

	areaPartialSums := make([]float64, len(newTriangles))
	var currentSum float64 = 0
	for idx, tri := range newTriangles {
		currentSum += tri.Area()
		areaPartialSums[idx] = currentSum
	}

	return &Triangularization{Triangles: newTriangles, AreaPartialSums: areaPartialSums}
}

// NewHexTriangularization triangularizes a hexagon of the given radius
func NewHexTriangularization(radius float64) *Triangularization {
	triangles := make([]Triangle, 0, 4)

	vertexA := cp.Vector{X: radius, Y: 0}

	for vertexBInd := 1; vertexBInd < 6; vertexBInd += 2 {
		vertexBRads := float64(vertexBInd) * math.Pi / 3.0
		vertexCRads := vertexBRads + math.Pi/3.0

		vertexB := cp.Vector{X: radius * math.Cos(vertexBRads), Y: radius * math.Sin(vertexBRads)}
		vertexC := cp.Vector{X: radius * math.Cos(vertexCRads), Y: radius * math.Sin(vertexCRads)}

		triangles = append(triangles, *NewTriangle(vertexA, vertexB, vertexC))
	}

	return NewTriangularization(triangles)
}

// NewHexEdgeTriangularization triangularizes a hexagon of the given radius,
// but without the inner diamond
func NewHexEdgeTriangularization(radius float64) *Triangularization {
	triangles := make([]Triangle, 0, 4)

	for half := 0; half < 2; half++ {
		startRadians := float64(half) * math.Pi

		vecs := make([]cp.Vector, 4)
		for i := 0; i < 4; i++ {
			rads := startRadians + float64(i)*math.Pi/3.0
			vecs[i] = cp.Vector{X: radius * math.Cos(rads), Y: radius * math.Sin(rads)}
		}

		midVec := vecs[1].Add(vecs[2]).Mult(0.5)

		triangles = append(
			triangles,
			*NewTriangle(vecs[0], vecs[1], midVec),
			*NewTriangle(midVec, vecs[2], vecs[3]),
		)
	}

	return NewTriangularization(triangles)
}

// Sample a point uniformly from one of the triangles in this triangularization
func (t *Triangularization) Sample() cp.Vector {
	seed := rand.Float64() * t.AreaPartialSums[len(t.AreaPartialSums)-1]

	for idx, psum := range t.AreaPartialSums {
		if seed < psum {
			return t.Triangles[idx].Sample()
		}
	}

	log.Fatalf("Triangularization.Sample() invalid AreaPartialSums: %v (seed=%v)", t, seed)
	return cp.Vector{}
}

// MapHexTriangularization triangularizes one map hex. Note that this
// includes the wall within the hexagon! Use MapHexInnerTriangularization
// for just the non-wall part of a map hex
var MapHexTriangularization = NewHexTriangularization(MAP_HEX_RADIUS)

// MapHexEdgeTriangularization triangularizes the edges of one map hex. Note
// that this includes the walls within the hexagon! Use
// MapHexInnerEdgeTriangularization for just the non-wall part of a map-hex
var MapHexEdgeTriangularization = NewHexEdgeTriangularization(MAP_HEX_RADIUS)

// MapHexInnerTriangularization triangularizes the inner, non-wall part of a hex.
// Note the center hex has no walls.
var MapHexInnerTriangularization = NewHexTriangularization(MAP_HEX_RADIUS - MAP_HEX_WALL_THICKNESS*0.5)

// MapHexEdgeTriangularization triangularizes the edges of one map hex, excluding
// the walls
var MapHexInnerEdgeTriangularization = NewHexEdgeTriangularization(MAP_HEX_RADIUS - MAP_HEX_WALL_THICKNESS*0.5)

// MapHexCenters maps from the hex index to the center of the hex.
var MapHexCenters = []cp.Vector{
	{X: math.Sqrt(3) * MAP_HEX_RADIUS * math.Cos(math.Pi/6.0), Y: math.Sqrt(3) * MAP_HEX_RADIUS * math.Sin(math.Pi/6.0)},
	{X: math.Sqrt(3) * MAP_HEX_RADIUS * math.Cos(math.Pi/2.0), Y: math.Sqrt(3) * MAP_HEX_RADIUS * math.Sin(math.Pi/2.0)},
	{X: math.Sqrt(3) * MAP_HEX_RADIUS * math.Cos(5*math.Pi/6.0), Y: math.Sqrt(3) * MAP_HEX_RADIUS * math.Sin(5*math.Pi/6.0)},
	{X: math.Sqrt(3) * MAP_HEX_RADIUS * math.Cos(7*math.Pi/6.0), Y: math.Sqrt(3) * MAP_HEX_RADIUS * math.Sin(7*math.Pi/6.0)},
	{X: math.Sqrt(3) * MAP_HEX_RADIUS * math.Cos(3*math.Pi/2.0), Y: math.Sqrt(3) * MAP_HEX_RADIUS * math.Sin(3*math.Pi/2.0)},
	{X: math.Sqrt(3) * MAP_HEX_RADIUS * math.Cos(11*math.Pi/6.0), Y: math.Sqrt(3) * MAP_HEX_RADIUS * math.Sin(11*math.Pi/6.0)},
}
