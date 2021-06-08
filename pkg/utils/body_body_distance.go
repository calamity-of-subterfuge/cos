package utils

import "github.com/jakecoffman/cp"

// BodyBodyDistanceSq approximates the squared distance between the two bodies.
// This is stable but may be an overestimate of the squared distance.
func BodyBodyDistanceSq(a *cp.Body, b *cp.Body) float64 {
	var resSq float64 = -1

	// I'm not really sure why there isn't a builtin function to do this?
	// I don't even think this is even necessarily always correct, but it
	// should be very close

	a.EachShape(func(s1 *cp.Shape) {
		b.EachShape(func(s2 *cp.Shape) {
			s1NearestPoint := s1.PointQuery(b.LocalToWorld(s2.CenterOfGravity()))
			s2NearestPoint := s2.PointQuery(a.LocalToWorld(s1.CenterOfGravity()))

			if s1NearestPoint.Distance <= 0 || s2NearestPoint.Distance <= 0 {
				resSq = 0
				return
			}

			distanceSq := s2NearestPoint.Point.Sub(s1NearestPoint.Point).LengthSq()
			if resSq == -1 || distanceSq < resSq {
				resSq = distanceSq
			}
		})
	})

	return resSq
}
