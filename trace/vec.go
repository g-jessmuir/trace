package trace

import "math"

// Vec implements a Vector type for vector math
type Vec struct {
	X, Y, Z float32
}

// Add o to v and return a new vector with the result.
func (v Vec) Add(o Vec) Vec {
	return Vec{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

// Sub o from v and return a new vector with the result.
func (v Vec) Sub(o Vec) Vec {
	return Vec{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

// Neg returns the negative v
func (v Vec) Neg() Vec {
	return Vec{-v.X, -v.Y, -v.Z}
}

// Mul multiplies v by f and returns the result
func (v Vec) Mul(f float32) Vec {
	return Vec{v.X * f, v.Y * f, v.Z * f}
}

// Div divides v by f and returns the result
func (v Vec) Div(f float32) Vec {
	return Vec{v.X / f, v.Y / f, v.Z / f}
}

// VMul multiplies v by o component-wise and returns the result
func (v Vec) VMul(o Vec) Vec {
	return Vec{v.X * o.X, v.Y * o.Y, v.Z * o.Z}
}

// Len returns the length of v
func (v Vec) Len() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// SqrLen returns the squared length of v
func (v Vec) SqrLen() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Unit returns a new unit vector of v
func (v Vec) Unit() Vec {
	l := v.Len()
	v = v.Div(l)
	return v
}

// Dot returns the dot product of v and  o
func (v Vec) Dot(o Vec) float32 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

// Cross returns the cross product of v and o
func (v Vec) Cross(o Vec) Vec {
	return Vec{
		v.Y*o.Z - v.Z*o.Y,
		v.Z*o.X - v.X*o.Z,
		v.X*o.Y - v.Y*o.X,
	}
}

// Reflect mirrors v across normal n
func (v Vec) Reflect(n Vec) Vec {
	return v.Sub(n.Mul(2 * v.Dot(n)))
}

// Refract refracts the ray according to Snell's law
func (v Vec) Refract(n Vec, niByNt float32) (Vec, bool) {
	uv := v.Unit()
	dt := uv.Dot(n)
	discriminant := 1 - niByNt*niByNt*(1-dt*dt)
	if discriminant > 0 {
		return uv.Sub(n.Mul(dt)).Mul(niByNt).Sub(n.Mul(Sqrt32(discriminant))), true
	}
	return Vec{-1, -1, -1}, false
}
