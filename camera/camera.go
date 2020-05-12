package camera

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

type Camera3D interface {
	GetProjectionMatrix(width, height int) mgl.Mat4
	Move(move mgl.Vec3)
	Rotate(a, b float32)
}

type Camera struct {
	LookAt mgl.Vec3
	Pos    mgl.Vec3
}

//GetProjectionMatrix Using the LookAt and Pos of the camera it returns the ProjectMatrix or MVPMatrix.
func (c *Camera) GetProjectionMatrix(width, height int) mgl.Mat4 {
	view := mgl.LookAtV(
		c.Pos,
		c.Pos.Add(c.LookAt),
		mgl.Vec3{0, 1, 0},
	)
	model := mgl.Ident4()
	projection := mgl.Perspective(mgl.DegToRad(80.0), float32(width)/float32(height), 0.1, 100.0)
	mvp := projection.Mul4(view.Mul4(model))
	return mvp
}

func (c *Camera) Move(move mgl.Vec3) {
	c.Pos = c.Pos.Add(move)
}
func (c *Camera) Rotate(a, b float32) {
	c.LookAt = mgl.Rotate3DY(a).Mul3x1(c.LookAt)
	c.LookAt = mgl.Rotate3DZ(b).Mul3x1(c.LookAt)
}

// Create camera3d
func CreateCamera3D(lookAt mgl.Vec3, pos mgl.Vec3) Camera3D {
	return &Camera{lookAt, pos}
}
