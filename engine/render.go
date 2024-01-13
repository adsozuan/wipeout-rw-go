package engine

import (
	"fmt"
	"math"
	"unsafe"

	gl "github.com/chsc/gogl/gl33"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	AtlasSize   = 64
	AtlasGrid   = 32
	AtlasBorder = 16

	RenderTrisBufferCapacity = 2048
	TextureMax               = 1024

	NearPlane                       = 16.0
	FarPlane                        = 64000.0
	RenderDepthBufferInternalFormat = gl.DEPTH_COMPONENT24

	RenderUseMipMaps  = true
	RenderFadeOutNear = 48000.0
	RenderFadeOutFar  = 64000.0
)

type RenderBlendMode byte

const (
	RenderBlendModeNormal RenderBlendMode = iota
	RenderBlendModeLighter
)

type RenderResolution byte

const (
	RenderResolutionNative RenderResolution = iota
	RenderResolution240p
	RenderResolution480p
)

type RenderPostEffect byte

const (
	RenderPostEffectNone RenderPostEffect = iota
	RenderPostEffectCRT
	NumRenderPostEffects
)

type RenderTexture struct {
	offset Vec2i
	size   Vec2i
}

var (
	// For pinning the array in memory otherwise can cause memory corruption
	trisBuffer [RenderTrisBufferCapacity]Tris

	RenderInstance *Render
)

type Render struct {
	vbo     gl.Uint
	trisLen int

	screenSize     Vec2i
	backBufferSize Vec2i

	atlasMap        [AtlasSize]uint32
	atlasTexture    gl.Uint
	renderBlendMode RenderBlendMode
	renderNoTexture uint16

	projectionMat2d Mat4
	projectionMatbb Mat4
	projectionMat3d Mat4
	spriteMat       Mat4
	viewMat         Mat4

	textures             [TextureMax]RenderTexture
	texturesLen          int
	textureMipMapIsDirty bool

	renderResolution      RenderResolution
	backBuffer            gl.Uint
	backBufferTexture     gl.Uint
	backBufferDepthBuffer gl.Uint

	programGame        *ProgramGame
	programPostEffect  *ProgramPostEffect
	programPostEffects [NumRenderPostEffects]*ProgramPostEffect
}

func NewRender() *Render {
	r := &Render{}
	r.trisLen = 0

	r.atlasTexture = 0
	r.renderBlendMode = RenderBlendModeNormal
	r.renderNoTexture = 0

	r.projectionMat2d = NewMat4Identity()
	r.projectionMatbb = NewMat4Identity()
	r.projectionMat3d = NewMat4Identity()
	r.spriteMat = NewMat4Identity()
	r.viewMat = NewMat4Identity()

	r.texturesLen = 0
	r.textureMipMapIsDirty = false

	r.backBuffer = 0
	r.backBufferTexture = 0
	r.backBufferDepthBuffer = 0

	RenderInstance = r

	return r
}

func (r *Render) Init(screenSize Vec2i) {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.GenTextures(1, &r.atlasTexture)
	gl.BindTexture(gl.TEXTURE_2D, r.atlasTexture)
	if RenderUseMipMaps {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	}
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	// var anisotropy float32 = 0
	// gl.GetFloatv(gl.MAX_TEXTURE_MAX_ANISTROPY_EXT, &anisotropy)
	// gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAX_ANISOTROPY_EXT, anisotropy)

	tw := gl.Sizei(AtlasSize * AtlasGrid)
	th := gl.Sizei(AtlasSize * AtlasGrid)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, tw, th, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	Logger.Printf("atlas texture %5d", r.atlasTexture)

	// Tris buffer
	gl.GenBuffers(1, &r.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)

	// Post effect shaders
	r.programPostEffects[RenderPostEffectNone] = ShaderPostEffectDefaultInit()
	r.programPostEffects[RenderPostEffectCRT] = ShaderPostEffectCRTInit()
	r.SetPostEffect(RenderPostEffectNone)

	// Game shader
	prgGame := ShaderGameInit()
	r.programGame = prgGame
	gl.UseProgram(prgGame.program)
	gl.BindVertexArray(prgGame.vao)

	r.SetView(Vec3{0, 0, 0}, Vec3{0, 0, 0})
	r.SetModelMat(&Mat4Id)

	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// Create default white texture
	white := []RGBA{
		RGBA{128, 128, 128, 255}, RGBA{128, 128, 128, 255},
		RGBA{128, 128, 128, 255}, RGBA{128, 128, 128, 255},
	}
	t, err := r.TextureCreate(2, 2, white)
	if err != nil {
		panic(err)
	}
	r.renderNoTexture = uint16(t)

	// Back buffer
	r.renderResolution = RenderResolutionNative
	r.SetScreenSize(screenSize)
}

func (r *Render) Cleanup() {
	// TODO see if this is needed
	// gl.DeleteTextures(1, &r.atlasTexture)
	// gl.DeleteBuffers(1, &r.vbo)
}

func (r *Render) SetScreenSize(size Vec2i) {
	r.screenSize = size
	r.projectionMatbb = r.Setup2dProjectionMat(size)

	r.SetResolution(r.renderResolution)
}

func (r *Render) Size() Vec2i {
	return r.backBufferSize
}

func (r *Render) Setup2dProjectionMat(size Vec2i) Mat4 {
	var near gl.Float = -1
	var far gl.Float = 1
	var left gl.Float = 0
	var right gl.Float = gl.Float(size.X)
	var top gl.Float = 0
	var bottom gl.Float = gl.Float(size.Y)

	var lr gl.Float = 1.0 / (left - right)
	var bt gl.Float = 1.0 / (bottom - top)
	var nf gl.Float = 1.0 / (near - far)

	return Mat4{
		-2 * lr, 0, 0, 0,
		0, -2 * bt, 0, 0,
		0, 0, 2 * nf, 0,
		lr * (left + right), bt * (top + bottom), nf * (far + near), 1,
	}
}

func (r *Render) Setup3dProjectionMat(size Vec2i) Mat4 {
	aspect := float32(size.X) / float32(size.Y)
	fov := (73.75 / 180.0) * math.Pi
	f := float32(1.0 / math.Tan(fov*0.5))
	nf := float32(1.0 / (NearPlane - FarPlane))
	return Mat4{
		gl.Float(f / aspect), 0, 0, 0,
		0, gl.Float(f), 0, 0,
		0, 0, gl.Float((FarPlane + NearPlane) * nf), -1,
		0, 0, gl.Float((2 * FarPlane * NearPlane) * nf), 0,
	}
}

func (r *Render) SetResolution(res RenderResolution) {
	r.renderResolution = res

	if res == RenderResolutionNative {
		r.backBufferSize = r.screenSize
	} else {
		aspect := float32(r.screenSize.X) / float32(r.screenSize.Y)
		if res == RenderResolution240p {
			r.backBufferSize = Vec2i{int32(aspect * 240), 240}
		} else if res == RenderResolution480p {
			r.backBufferSize = Vec2i{int32(aspect * 480), 480}
		} else {
			panic(fmt.Sprintf("invalid resolution %d", res))
		}
	}

	if r.backBuffer == 0 {
		gl.GenTextures(1, &r.backBufferTexture)
		gl.GenFramebuffers(1, &r.backBuffer)
		gl.GenRenderbuffers(1, &r.backBufferDepthBuffer)
	}

	gl.BindTexture(gl.TEXTURE_2D, r.backBufferTexture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, gl.Sizei(r.backBufferSize.X), gl.Sizei(r.backBufferSize.Y), 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.BindFramebuffer(gl.FRAMEBUFFER, r.backBuffer)
	gl.BindRenderbuffer(gl.RENDERBUFFER, r.backBufferDepthBuffer)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, r.backBufferDepthBuffer)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, r.backBufferTexture, 0)

	gl.BindRenderbuffer(gl.RENDERBUFFER, r.backBufferDepthBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, RenderDepthBufferInternalFormat, gl.Sizei(r.backBufferSize.X), gl.Sizei(r.backBufferSize.Y))

	r.projectionMat2d = r.Setup2dProjectionMat(r.backBufferSize)
	r.projectionMat3d = r.Setup3dProjectionMat(r.backBufferSize)

	// Use nearest filtering for 240p and 480p
	gl.BindTexture(gl.TEXTURE_2D, r.atlasTexture)
	if r.renderResolution == RenderResolutionNative {
		if RenderUseMipMaps {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		} else {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		}
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	}
	gl.Viewport(0, 0, gl.Sizei(r.backBufferSize.X), gl.Sizei(r.backBufferSize.Y))
}

func (r *Render) SetPostEffect(postEffect RenderPostEffect) error {
	if postEffect > NumRenderPostEffects {
		return fmt.Errorf("invalid post effect %d", postEffect)
	}
	r.programPostEffect = r.programPostEffects[postEffect]

	return nil
}

func (r *Render) FramePrepare() {
	gl.UseProgram(r.programGame.program)
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.backBuffer)
	gl.Viewport(0, 0, gl.Sizei(r.backBufferSize.X), gl.Sizei(r.backBufferSize.Y))

	gl.BindTexture(gl.TEXTURE_2D, r.atlasTexture)
	gl.Uniform2f(gl.Int(r.programGame.uniform.screen), 0, 0)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(gl.TRUE)
	gl.Disable(gl.POLYGON_OFFSET_FILL)
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)
}

func (r *Render) FrameEnd(cycleTime float64) {
	r.Flush()

	gl.UseProgram(r.programPostEffect.program)

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, gl.Sizei(r.screenSize.X), gl.Sizei(r.screenSize.Y))
	gl.BindTexture(gl.TEXTURE_2D, r.backBufferTexture)
	gl.UniformMatrix4fv(gl.Int(r.programPostEffect.uniform.projection), 1, gl.FALSE, (*gl.Float)(&r.projectionMatbb[0]))
	gl.Uniform1f(gl.Int(r.programPostEffect.uniform.time), gl.Float(cycleTime))
	gl.Uniform2f(gl.Int(r.programPostEffect.uniform.screenSize), gl.Float(r.screenSize.X), gl.Float(r.screenSize.Y))

	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	white := RGBA{128, 128, 128, 255}
	r.trisLen++
	trisBuffer[r.trisLen] = Tris{
		Vertices: [3]Vertex{
			Vertex{Pos: Vec3{0, gl.Float(r.screenSize.Y), 0}, UV: Vec2{0, 0}, Color: white},
			Vertex{Pos: Vec3{gl.Float(r.screenSize.X), 0, 0}, UV: Vec2{1, 1}, Color: white},
			Vertex{Pos: Vec3{0, 0, 0}, UV: Vec2{0, 1}, Color: white},
		},
	}
	r.trisLen++
	trisBuffer[r.trisLen] = Tris{
		Vertices: [3]Vertex{
			Vertex{Pos: Vec3{gl.Float(r.screenSize.X), gl.Float(r.screenSize.Y), 0}, UV: Vec2{1, 0}, Color: white},
			Vertex{Pos: Vec3{gl.Float(r.screenSize.X), 0, 0}, UV: Vec2{1, 1}, Color: white},
			Vertex{Pos: Vec3{0, gl.Float(r.screenSize.Y), 0}, UV: Vec2{0, 0}, Color: white},
		},
	}

	r.Flush()

}

func (r *Render) Flush() {
	if r.trisLen == 0 {
		return
	}

	if r.textureMipMapIsDirty {
		gl.GenerateMipmap(gl.TEXTURE_2D)
		r.textureMipMapIsDirty = false
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(unsafe.Sizeof(trisBuffer[0]) * uintptr(r.trisLen)), gl.Pointer(&trisBuffer[0]), gl.DYNAMIC_DRAW)
	gl.DrawArrays(gl.TRIANGLES, gl.Int(0), gl.Sizei(r.trisLen*3))
	r.trisLen = 0
}

func (r *Render) SetView(pos Vec3, angles Vec3) {
	r.Flush()
	r.SetDepthWrite(true)
	r.SetDepthTest(true)

	r.viewMat = NewMat4Identity()
	Mat4SetTranslation(&r.viewMat, Vec3{0, 0, 0})
	Mat4SetRollPitchYaw(&r.viewMat, Vec3{angles.X, -angles.Y + math.Pi, angles.Z + math.Pi})
	Mat4Translate(&r.viewMat, Vec3Inv(pos))
	Mat4SetYawPitchRoll(&r.spriteMat, Vec3{-angles.X, angles.Y - math.Pi, 0})

	r.SetModelMat(&Mat4Id)

	gl.UniformMatrix4fv(gl.Int(r.programGame.uniform.view), 1, gl.FALSE, (*gl.Float)(&r.viewMat[0]))
	gl.UniformMatrix4fv(gl.Int(r.programGame.uniform.projection), 1, gl.FALSE, (*gl.Float)(&r.projectionMat3d[0]))
	gl.Uniform3f(gl.Int(r.programGame.uniform.cameraPos), gl.Float(pos.X), gl.Float(pos.Y), gl.Float(pos.Z))
	gl.Uniform2f(gl.Int(r.programGame.uniform.fade), gl.Float(RenderFadeOutNear), gl.Float(RenderFadeOutFar))
}

func (r *Render) SetModelMat(m *Mat4) {
	r.Flush()

	gl.UniformMatrix4fv(gl.Int(r.programGame.uniform.model), 1, gl.FALSE, &m[0])
}

func (r *Render) SetView2d() {
	r.SetDepthWrite(true)
	r.SetDepthTest(false)

	r.SetModelMat(&Mat4Id)
	gl.Uniform3f((gl.Int(r.programGame.uniform.cameraPos)), gl.Float(0), gl.Float(0), gl.Float(0))
	gl.UniformMatrix4fv(gl.Int(r.programGame.uniform.view), gl.Sizei(1), gl.GLBool(false), &Mat4Id[0])
	gl.UniformMatrix4fv(gl.Int(r.programGame.uniform.projection), gl.Sizei(1), gl.GLBool(false), &r.projectionMat2d[0])
}

func (r *Render) SetDepthWrite(enable bool) {
	r.Flush()
	gl.DepthMask(gl.GLBool(enable))
}

func (r *Render) SetDepthTest(enable bool) {
	r.Flush()
	if enable {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}

func (r *Render) SetDepthOffset(offset float32) {
	r.Flush()
	if offset == 0 {
		gl.Disable(gl.POLYGON_OFFSET_FILL)
		return
	}

	gl.Enable(gl.POLYGON_OFFSET_FILL)
	gl.PolygonOffset(gl.Float(offset), 1.0)
}

func (r *Render) SetScreenPosition(pos Vec2i) {
	r.Flush()
	gl.Uniform2f(gl.Int(r.programGame.uniform.screen), gl.Float(pos.X), -gl.Float(pos.Y))
}

func (r *Render) SetBlendMode(newMode RenderBlendMode) {
	if newMode == r.renderBlendMode {
		return
	}
	r.Flush()

	r.renderBlendMode = newMode
	if r.renderBlendMode == RenderBlendModeNormal {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	} else if r.renderBlendMode == RenderBlendModeLighter {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	}
}

func (r *Render) SetCullBackface(enable bool) {
	r.Flush()
	if enable {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
}

func (r *Render) Transform(pos Vec3) Vec3 {
	return Vec3Transform(Vec3Transform(pos, &(r.viewMat)), &(r.projectionMat3d))
}

func (r *Render) PushTris(tris Tris, textureIndex int) error {
	if textureIndex >= r.texturesLen {
		return fmt.Errorf("invalid texture index %d", textureIndex)
	}

	if r.trisLen >= RenderTrisBufferCapacity {
		r.Flush()
	}

	t := &r.textures[textureIndex]

	for i := 0; i < 3; i++ {
		tris.Vertices[i].UV.X += gl.Float(t.offset.X)
		tris.Vertices[i].UV.Y += gl.Float(t.offset.Y)
	}
	r.trisLen++
	trisBuffer[r.trisLen] = tris

	return nil
}

func (r *Render) PushSprite(pos Vec3, size Vec2i, color RGBA, textureIndex int) error {
	if textureIndex >= r.texturesLen {
		return fmt.Errorf("invalid texture index %d", textureIndex)
	}

	p1 := Vec3Add(pos, Vec3Transform(Vec3{gl.Float(-size.X) * 0.5, gl.Float(-size.Y) * 0.5, 0}, &(r.spriteMat)))
	p2 := Vec3Add(pos, Vec3Transform(Vec3{gl.Float(size.X) * 0.5, gl.Float(-size.Y) * 0.5, 0}, &(r.spriteMat)))
	p3 := Vec3Add(pos, Vec3Transform(Vec3{gl.Float(-size.X) * 0.5, gl.Float(size.Y) * 0.5, 0}, &(r.spriteMat)))
	p4 := Vec3Add(pos, Vec3Transform(Vec3{gl.Float(size.X) * 0.5, gl.Float(size.Y) * 0.5, 0}, &(r.spriteMat)))

	t := &r.textures[textureIndex]

	r.PushTris(Tris{
		Vertices: [3]Vertex{
			Vertex{Pos: p1, UV: Vec2{0, 0}, Color: color},
			Vertex{Pos: p2, UV: Vec2{gl.Float(0 + t.size.X), 0}, Color: color},
			Vertex{Pos: p3, UV: Vec2{0, gl.Float(0 + t.size.Y)}, Color: color},
		},
	}, textureIndex)

	r.PushTris(Tris{
		Vertices: [3]Vertex{
			Vertex{Pos: p3, UV: Vec2{0, gl.Float(0 + t.size.Y)}, Color: color},
			Vertex{Pos: p2, UV: Vec2{gl.Float(0 + t.size.X), 0}, Color: color},
			Vertex{Pos: p4, UV: Vec2{gl.Float(0 + t.size.X), gl.Float(0 + t.size.Y)}, Color: color},
		},
	}, textureIndex)

	return nil
}

func (r *Render) Push2d(pos Vec2i, size Vec2i, color RGBA, textureIndex int) error {
	ts, err := r.TextureSize(textureIndex)
	if err != nil {
		return err
	}
	r.Push2dTitle(pos, Vec2i{0, 0}, ts, size, color, textureIndex)

	return nil
}

func (r *Render) Push2dTitle(pos Vec2i, uvOffset Vec2i, uvSize Vec2i, size Vec2i, color RGBA, textureIndex int) error {
	if textureIndex >= r.texturesLen {
		return fmt.Errorf("invalid texture index %d", textureIndex)
	}

	r.PushTris(Tris{
		Vertices: [3]Vertex{
			{Pos: Vec3{gl.Float(pos.X) + gl.Float(size.X), gl.Float(pos.Y) + gl.Float(size.Y), 0}, UV: Vec2{gl.Float(uvOffset.X), gl.Float(uvOffset.Y)}, Color: color},
			{Pos: Vec3{gl.Float(pos.X + uvSize.X), gl.Float(pos.Y), 0}, UV: Vec2{gl.Float(uvOffset.X + uvSize.X), gl.Float(uvOffset.Y)}, Color: color},
			{Pos: Vec3{gl.Float(pos.X), gl.Float(pos.Y), 0}, UV: Vec2{gl.Float(uvOffset.X), gl.Float(uvOffset.Y)}, Color: color},
		}}, textureIndex)

	r.PushTris(Tris{
		Vertices: [3]Vertex{
			{Pos: Vec3{gl.Float(pos.X) + gl.Float(size.X), gl.Float(pos.Y) + gl.Float(size.Y), 0}, UV: Vec2{gl.Float(uvOffset.X) + gl.Float(uvSize.X), gl.Float(uvOffset.Y) + gl.Float(uvSize.Y)}, Color: color},
			{Pos: Vec3{gl.Float(pos.X + uvSize.X), gl.Float(pos.Y), 0}, UV: Vec2{gl.Float(uvOffset.X + uvSize.X), gl.Float(uvOffset.Y)}, Color: color},
			{Pos: Vec3{gl.Float(pos.X), gl.Float(pos.Y), 0}, UV: Vec2{gl.Float(uvOffset.X), gl.Float(uvOffset.Y) + gl.Float(uvSize.Y)}, Color: color},
		}}, textureIndex)

	return nil
}

func (r *Render) TextureCreate(tw int, th int, pixels []RGBA) (int, error) {
	if r.texturesLen >= TextureMax {
		return 0, fmt.Errorf("texture max reached, wanted %d", r.texturesLen)
	}

	bw := tw + AtlasBorder*2
	bh := th + AtlasBorder*2

	// Find a spot in the atlas for this texture with added border
	gridWidth := (bw + AtlasGrid - 1) / AtlasGrid
	gridHeight := (bh + AtlasGrid - 1) / AtlasGrid
	gridX := 0
	gridY := AtlasSize - gridHeight + 1
	for cx := 0; cx < AtlasSize-gridWidth; cx++ {
		if r.atlasMap[cx] >= uint32(gridY) {
			continue
		}
		cy := r.atlasMap[cx]
		isBest := true

		for bx := cx; bx < cx+gridWidth; bx++ {
			if r.atlasMap[bx] >= uint32(gridY) {
				isBest = false
				cx = bx
				break
			}
			if r.atlasMap[bx] > cy {
				cy = r.atlasMap[bx]
			}
		}
		if isBest {
			gridX = cx
			gridY = int(cy)
		}
	}

	if gridY+gridHeight > AtlasSize {
		return 0, fmt.Errorf("render atlas full")
	}

	for cx := gridX; cx < gridX+gridWidth; cx++ {
		r.atlasMap[cx] = uint32(gridY + gridHeight)
	}

	// Add border pixels for this texture
	pb := make([]RGBA, bw*bh)

	if tw > 0 && th > 0 {
		// Top border
		for y := 0; y < AtlasBorder; y++ {
			for x := 0; x < bw; x++ {
				pb[x+y*bw+AtlasBorder] = pixels[0] // TODO something wrong here
			}
		}

		// Bottom border
		for y := 0; y < AtlasBorder; y++ {
			for x := 0; x < bw; x++ {
				pb[x+(bh-AtlasBorder+y)*bw] = pixels[(th-1)*tw]
			}
		}

		// Left border
		for y := 0; y < bh; y++ {
			for x := 0; x < AtlasBorder; x++ {
				pb[x+y*bw] = pixels[int(Clamp(y-AtlasBorder, 0, th-1))*tw]
			}
		}

		// Right border
		for y := 0; y < bh; y++ {
			for x := 0; x < AtlasBorder; x++ {
				pb[y*bw+(bw-AtlasBorder+x)] = pixels[Clamp(y-AtlasBorder, 0, th-1)*tw+tw-1]
			}
		}

		// Texture
		for y := 0; y < th; y++ {
			for x := 0; x < tw; x++ {
				pb[(x+AtlasBorder)+(y+AtlasBorder)*bw] = pixels[x+y*tw]
			}
		}

		for y := 0; y < th; y++ {
			pb[bw*(y+AtlasBorder)+AtlasBorder] = pixels[tw*y]
		}
	}

	x := gridX * AtlasGrid
	y := gridY * AtlasGrid
	gl.BindTexture(gl.TEXTURE_2D, r.atlasTexture)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, gl.Int(x), gl.Int(y), gl.Sizei(bw), gl.Sizei(bh), gl.RGBA, gl.UNSIGNED_BYTE, gl.Pointer(&pb[0]))

	r.textureMipMapIsDirty = RenderUseMipMaps
	textureIndex := r.texturesLen
	r.texturesLen++
	r.textures[textureIndex] = RenderTexture{
		Vec2i{int32(x + AtlasBorder), int32(y + AtlasBorder)},
		Vec2i{int32(tw), int32(th)},
	}

	return textureIndex, nil
}

func (r *Render) TextureSize(textureIndex int) (Vec2i, error) {
	if textureIndex >= r.texturesLen {
		return Vec2i{}, fmt.Errorf("invalid texture index %d", textureIndex)
	}

	return r.textures[textureIndex].size, nil
}

func (r *Render) TextureReplacePixels(textureIndex uint16, pixels []RGBA) error {
	if textureIndex >= uint16(r.texturesLen) {
		return fmt.Errorf("invalid texture index %d", textureIndex)
	}

	t := &r.textures[textureIndex]
	gl.BindTexture(gl.TEXTURE_2D, r.atlasTexture)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, gl.Int(t.offset.X), gl.Int(t.offset.Y), gl.Sizei(t.size.X), gl.Sizei(t.size.Y), gl.RGBA, gl.UNSIGNED_BYTE, gl.Pointer(&pixels[0]))

	return nil
}

func (r *Render) TexturesLen() int {
	return r.texturesLen
}

func (r *Render) TexturesReset(len uint16) error {
	if len > uint16(r.texturesLen) {
		return fmt.Errorf("invalid texture reset len %d >= %d", len, r.texturesLen)
	}

	r.texturesLen = int(len)
	for i, _ := range r.atlasMap {
		r.atlasMap[i] = 0
	}

	// Clear complete atlas and recreate the default white texture
	if len == 0 {
		whitePixels := []RGBA{
			RGBA{128, 128, 128, 255}, RGBA{128, 128, 128, 255},
			RGBA{128, 128, 128, 255}, RGBA{128, 128, 128, 255},
		}
		t, err := r.TextureCreate(2, 2, whitePixels)

		r.renderNoTexture = uint16(t)

		if err != nil {
			return err
		}
		return nil
	}

	// Replay all textures grid insertions up to the reset len
	for i := 0; i < r.texturesLen; i++ {
		gridX := (r.textures[i].offset.X - AtlasBorder) / AtlasGrid
		gridY := (r.textures[i].offset.Y - AtlasBorder) / AtlasGrid
		gridWidth := (r.textures[i].size.X + AtlasBorder*2 + AtlasGrid - 1) / AtlasGrid
		gridHeight := (r.textures[i].size.Y + AtlasBorder*2 + AtlasGrid - 1) / AtlasGrid
		for cx := gridX; cx < gridX+gridWidth; cx++ {
			r.atlasMap[cx] = uint32(gridY + gridHeight)
		}
	}

	return nil
}

func (r *Render) TexturesDump(path string) error {

	// Get current displayed image on screen via OpenGL


	// Get the image from the SDL surface


	width := AtlasSize * AtlasGrid
	height := AtlasSize * AtlasGrid
	pixels := make([]RGBA, width*height)
	gl.GetTexImage(gl.TEXTURE_2D, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Pointer(&pixels[0]))

	surface, err := sdl.CreateRGBSurfaceFrom(unsafe.Pointer(&pixels[0]), int32(width), int32(height), 32, int(width*4),
        0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000)
    if err != nil {
        Logger.Fatalf("Failed to create SDL surface: %s\n", err)
    }
    defer surface.Free()

	




    // Flip the image vertically (OpenGL's origin is bottom-left, SDL's is top-left)
    //flipSurface(surface)

    // Save the surface to an image file
    if err := img.SavePNG(surface, path); err != nil {
        Logger.Fatalf("Failed to save screenshot: %s\n", err)
    }


	// surface, err := sdl.CreateRGBSurfaceFrom(pixels, int32(width), int32(height), 32, int32(width*4),
    //     0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000)
    // if err != nil {
    //     Logger.Fatalf("Failed to create SDL surface: %s\n", err)
    // }
    // defer surface.Free()

	// TODO write into png
	return nil
}
