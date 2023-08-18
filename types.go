package wipeout

import (
	"math"
)

type Vec2 struct {
	X, Y float64
}

type Vec2i struct {
	X, Y int32
}

type Vec3 struct {
	X, Y, Z float64
}

type Mat4 [16]float64

type Vertex struct {
	Pos   Vec3
	UV    Vec2
	Color RGBA
}

type RGBA struct {
	R, G, B, A float64
}

func NewRGBA(r, g, b, a float64) RGBA {
	return RGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

func NewVec2(x, y float64) Vec2 {
	return Vec2{
		X: x,
		Y: y,
	}
}

func NewVec3(x, y, z float64) Vec3 {
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

func NewMat4(m [16]float64) Mat4 {
	return Mat4(m)
}

func NewMat4Identity() Mat4 {
	return NewMat4([16]float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})
}

func Vec2MulF(a Vec2, f float64) Vec2 {
	return NewVec2(
		a.X*f,
		a.Y*f,
	)
}

func Vec2iMulF(a Vec2i, f float64) Vec2i {
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

func Vec3MulF(a Vec3, f float64) Vec3 {
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

func Vec3DivF(a Vec3, f float64) Vec3 {
	return NewVec3(
		a.X/f,
		a.Y/f,
		a.Z/f,
	)
}

func Vec3Len(a Vec3) float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

func Vec3Cross(a, b Vec3) Vec3 {
	return NewVec3(
		a.Y*b.Z-a.Z*b.Y,
		a.Z*b.X-a.X*b.Z,
		a.X*b.Y-a.Y*b.X,
	)
}

func Vec3Dot(a, b Vec3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func Vec3Lerp(a, b Vec3, t float64) Vec3 {
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

func WrapAngle(a float64) float64 {
	a = math.Mod(float64(a+math.Pi), math.Pi*2)
	if a < 0 {
		a += math.Pi * 2
	}
	return a - math.Pi
}

func Vec3WrapAngle(a Vec3) Vec3 {
	return Vec3{
		X: WrapAngle(a.X),
		Y: WrapAngle(a.Y),
		Z: WrapAngle(a.Z),
	}
}

func Vec3Angle(a, b Vec3) float64 {
	magnitude := Vec3Len(a) * Vec3Len(b)
	cosine := Vec3Dot(a, b) / magnitude
	return math.Acos(Clamp(cosine, -1, 1))
}

func Vec3Transform(a Vec3, mat *Mat4) Vec3 {
	w := mat[3]*a.X + mat[7]*a.Y + mat[11]*a.Z + mat[15]
	if w == 0 {
		w = 1
	}
	return Vec3{
		X: (mat[0]*a.X + mat[4]*a.Y + mat[8]*a.Z + mat[12]) / w,
		Y: (mat[1]*a.X + mat[5]*a.Y + mat[9]*a.Z + mat[13]) / w,
		Z: (mat[2]*a.X + mat[6]*a.Y + mat[10]*a.Z + mat[14]) / w,
	}
}

func Vec3ProjectToRay(p, r0, r1 Vec3) Vec3 {
	ray := Vec3Normalize(Vec3Sub(r1, r0))
	dp := Vec3Dot(Vec3Sub(p, r0), ray)
	return Vec3Add(r0, Vec3MulF(ray, dp))
}

func Vec3DistanceToPlane(p, planePos, planeNormal Vec3) float64 {
	dotProduct := Vec3Dot(Vec3Sub(planePos, p), planeNormal)
	normDotProduct := Vec3Dot(Vec3MulF(planeNormal, -1), planeNormal)
	return dotProduct / normDotProduct
}

func Vec3Reflect(incidence, normal Vec3, f float64) Vec3 {
	return Vec3Add(incidence, Vec3MulF(normal, Vec3Dot(normal, Vec3MulF(incidence, -1))*f))
}

func Clamp(value, min, max float64) float64 {
	if value < min {
		return float64(min)
	}
	if value > max {
		return float64(max)
	}
	return float64(value)
}

func Mat4SetTranslation(mat *Mat4, pos Vec3) {
	mat[12] = pos.X
	mat[13] = pos.Y
	mat[14] = pos.Z
}

func Mat4SetYawPitchRoll(mat *Mat4, rot Vec3) {
	sx := math.Sin(rot.X)
	sy := math.Sin(-rot.Y)
	sz := math.Sin(-rot.Z)
	cx := math.Cos(rot.X)
	cy := math.Cos(-rot.Y)
	cz := math.Cos(-rot.Z)

	mat[0] = cy*cz + sx*sy*sz
	mat[1] = cz*sx*sy - cy*sz
	mat[2] = cx * sy
	mat[4] = cx * sz
	mat[5] = cx * cz
	mat[6] = -sx
	mat[8] = -cz*sy + cy*sx*sz
	mat[9] = cy*cz*sx + sy*sz
	mat[10] = cx * cy
}

func Mat4SetRollPitchYaw(mat *Mat4, rot Vec3) {
	sx := math.Sin(rot.X)
	sy := math.Sin(-rot.Y)
	sz := math.Sin(-rot.Z)
	cx := math.Cos(rot.X)
	cy := math.Cos(-rot.Y)
	cz := math.Cos(-rot.Z)

	mat[0] = cy*cz - sx*sy*sz
	mat[1] = -cx * sz
	mat[2] = cz*sy + cy*sx*sz
	mat[4] = cz*sx*sy + cy*sz
	mat[5] = cx * cz
	mat[6] = -cy*cz*sx + sy*sz
	mat[8] = -cx * sy
	mat[9] = sx
	mat[10] = cx * cy
}

func Mat4Translate(mat *Mat4, translation Vec3) {
	mat[12] = mat[0]*translation.X + mat[4]*translation.Y + mat[8]*translation.Z + mat[12]
	mat[13] = mat[1]*translation.X + mat[5]*translation.Y + mat[9]*translation.Z + mat[13]
	mat[14] = mat[2]*translation.X + mat[6]*translation.Y + mat[10]*translation.Z + mat[14]
	mat[15] = mat[3]*translation.X + mat[7]*translation.Y + mat[11]*translation.Z + mat[15]
}

func Mat4Mul(res, a, b *Mat4) {
	res[0] = b[0]*a[0] + b[1]*a[4] + b[2]*a[8] + b[3]*a[12]
	res[1] = b[0]*a[1] + b[1]*a[5] + b[2]*a[9] + b[3]*a[13]
	res[2] = b[0]*a[2] + b[1]*a[6] + b[2]*a[10] + b[3]*a[14]
	res[3] = b[0]*a[3] + b[1]*a[7] + b[2]*a[11] + b[3]*a[15]
	res[4] = b[4]*a[0] + b[5]*a[4] + b[6]*a[8] + b[7]*a[12]
	res[5] = b[4]*a[1] + b[5]*a[5] + b[6]*a[9] + b[7]*a[13]
	res[6] = b[4]*a[2] + b[5]*a[6] + b[6]*a[10] + b[7]*a[14]
	res[7] = b[4]*a[3] + b[5]*a[7] + b[6]*a[11] + b[7]*a[15]
	res[8] = b[8]*a[0] + b[9]*a[4] + b[10]*a[8] + b[11]*a[12]
	res[9] = b[8]*a[1] + b[9]*a[5] + b[10]*a[9] + b[11]*a[13]
	res[10] = b[8]*a[2] + b[9]*a[6] + b[10]*a[10] + b[11]*a[14]
	res[11] = b[8]*a[3] + b[9]*a[7] + b[10]*a[11] + b[11]*a[15]
	res[12] = b[12]*a[0] + b[13]*a[4] + b[14]*a[8] + b[15]*a[12]
	res[13] = b[12]*a[1] + b[13]*a[5] + b[14]*a[9] + b[15]*a[13]
	res[14] = b[12]*a[2] + b[13]*a[6] + b[14]*a[10] + b[15]*a[14]
	res[15] = b[12]*a[3] + b[13]*a[7] + b[14]*a[11] + b[15]*a[15]
}
