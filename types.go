package wipeout

import (
	"math"
)

type Vec2 struct {
	X, Y float32
}

type Vec2i struct {
	X, Y int32
}

type Vec3 struct {
	X, Y, Z float32
}

type Mat4 [16]float32

type Vertex struct {
	Pos   Vec3
	UV    Vec2
	Color RGBA
}

type RGBA struct {
	R, G, B, A float32
}

func NewRGBA(r, g, b, a float32) RGBA {
	return RGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

func NewVec2(x, y float32) Vec2 {
	return Vec2{
		X: x,
		Y: y,
	}
}

func NewVec3(x, y, z float32) Vec3 {
	return Vec3{
		X: x,
		Y: y,
		Z: z,
	}
}

func NewVec2i(x, y int32) Vec2i {
	return Vec2i{
		X: x,
		Y: y,
	}
}

func NewMat4(m [16]float32) Mat4 {
	return Mat4(m)
}

func NewMat4Identity() Mat4 {
	return NewMat4([16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})
}

func Vec2MulF(a Vec2, f float32) Vec2 {
	return NewVec2(
		a.X*f,
		a.Y*f,
	)
}

func Vec2iMulF(a Vec2i, f float32) Vec2i {
	return NewVec2i(
		a.X,
		a.Y,
	)
}

func Vec3Add(a, b Vec3) Vec3 {
	return NewVec3(
		a.X+b.X,
		a.Y+b.Y,
		a.Z+b.Z,
	)
}

func Vec3Sub(a, b Vec3) Vec3 {
	return NewVec3(
		a.X-b.X,
		a.Y-b.Y,
		a.Z-b.Z,
	)
}

func Vec3Mul(a, b Vec3) Vec3 {
	return NewVec3(
		a.X*b.X,
		a.Y*b.Y,
		a.Z*b.Z,
	)
}

func Vec3MulF(a Vec3, f float32) Vec3 {
	return NewVec3(
		a.X*f,
		a.Y*f,
		a.Z*f,
	)
}

func Vec3Inv(a Vec3) Vec3 {
	return NewVec3(
		-a.X,
		-a.Y,
		-a.Z,
	)
}

func Vec3DivF(a Vec3, f float32) Vec3 {
	return NewVec3(
		a.X/f,
		a.Y/f,
		a.Z/f,
	)
}

func Vec3Len(a Vec3) float32 {
	return float32(math.Sqrt(float64(a.X*a.X + a.Y*a.Y + a.Z*a.Z)))
}

func Vec3Cross(a, b Vec3) Vec3 {
	return NewVec3(
		a.Y*b.Z-a.Z*b.Y,
		a.Z*b.X-a.X*b.Z,
		a.X*b.Y-a.Y*b.X,
	)
}

func Vec3Dot(a, b Vec3) float32 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func Vec3Lerp(a, b Vec3, t float32) Vec3 {
	return NewVec3(
		a.X+t*(b.X-a.X),
		a.Y+t*(b.Y-a.Y),
		a.Z+t*(b.Z-a.Z),
	)
}

func Vec3Normalize(a Vec3) Vec3 {
	length := Vec3Len(a)
	return Vec3DivF(a, length)
}

func WrapAngle(a float32) float32 {
	a = float32(math.Mod(float64(a+math.Pi), math.Pi*2))
	if a < 0 {
		a += float32(math.Pi * 2)
	}
	return a - float32(math.Pi)
}
