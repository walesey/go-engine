package webgl

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl32/matstack"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"github.com/walesey/go-engine/renderer"
)

const (
	maxLights int = 8
)

type WebglRenderer struct {
	onInit, onUpdate, onRender func()
	WindowWidth, WindowHeight  int
	canvas                     *js.Object
	gl                         *webgl.Context
	shouldClose                bool
	defaultTexture             *js.Object

	nbLights            int32
	nbDirectionalLights int32
	lights              []float32
	directionalLights   []float32

	bufferIds uint32
	buffers   map[uint32]*js.Object

	modelUniform       *js.Object
	modelNormalUniform *js.Object
	transformStack     *matstack.TransformStack
	program            *js.Object

	envMap         *js.Object
	envMapLOD1     *js.Object
	envMapLOD2     *js.Object
	envMapLOD3     *js.Object
	illuminanceMap *js.Object

	cameraLocation mgl32.Vec3
	cameraUp       mgl32.Vec3
	cameraLookAt   mgl32.Vec3
	cameraAngle    float32
	cameraNear     float32
	cameraFar      float32
	cameraOrtho    bool
}

func NewWebRenderer(canvas *js.Object) *WebglRenderer {
	return &WebglRenderer{
		canvas:       canvas,
		WindowWidth:  800,
		WindowHeight: 600,
		buffers:      make(map[uint32]*js.Object),
	}
}

func (wr *WebglRenderer) Init(callback func()) {
	wr.onInit = callback
}

func (wr *WebglRenderer) Update(callback func()) {
	wr.onUpdate = callback
}

func (wr *WebglRenderer) Render(callback func()) {
	wr.onRender = callback
}

func (wr *WebglRenderer) Start() {

	wr.canvas.Call("setAttribute", "width", fmt.Sprint(wr.WindowWidth))
	wr.canvas.Call("setAttribute", "height", fmt.Sprint(wr.WindowHeight))

	attrs := webgl.DefaultAttributes()
	attrs.Alpha = false
	attrs.Depth = true

	gl, err := webgl.NewContext(wr.canvas, attrs)
	if err != nil {
		js.Global.Call("alert", "Error: "+err.Error())
	}
	wr.gl = gl

	// Default shader program
	vertShader := compileShader(gl, mainVert, gl.VERTEX_SHADER)
	fragShader := compileShader(gl, mainFrag, gl.FRAGMENT_SHADER)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertShader)
	gl.AttachShader(program, fragShader)
	gl.LinkProgram(program)
	status := gl.GetProgramParameterb(program, gl.LINK_STATUS)
	if !status {
		panic(fmt.Sprintf("failed to link program: %v", gl.GetProgramInfoLog(program)))
	}

	gl.DeleteShader(vertShader)
	gl.DeleteShader(fragShader)

	gl.UseProgram(program)
	wr.program = program

	//set shader uniforms
	wr.transformStack = matstack.NewTransformStack()
	wr.modelUniform = gl.GetUniformLocation(program, "model")
	wr.PushTransform(mgl32.Ident4())

	textureUniform := gl.GetUniformLocation(program, "diffuse")
	gl.Uniform1i(textureUniform, 0)
	textureUniform = gl.GetUniformLocation(program, "normal")
	gl.Uniform1i(textureUniform, 1)
	textureUniform = gl.GetUniformLocation(program, "specular")
	gl.Uniform1i(textureUniform, 2)
	textureUniform = gl.GetUniformLocation(program, "roughness")
	gl.Uniform1i(textureUniform, 3)
	textureUniform = gl.GetUniformLocation(program, "environmentMap")
	gl.Uniform1i(textureUniform, 4)
	textureUniform = gl.GetUniformLocation(program, "environmentMapLOD1")
	gl.Uniform1i(textureUniform, 5)
	textureUniform = gl.GetUniformLocation(program, "environmentMapLOD2")
	gl.Uniform1i(textureUniform, 6)
	textureUniform = gl.GetUniformLocation(program, "environmentMapLOD3")
	gl.Uniform1i(textureUniform, 7)
	textureUniform = gl.GetUniformLocation(program, "illuminanceMap")
	gl.Uniform1i(textureUniform, 8)

	defaultImage := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: 1, Y: 1},
	})
	defaultImage.SetRGBA(0, 0, color.RGBA{255, 255, 255, 255})
	wr.defaultTexture = wr.newTexture(defaultImage, 0)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	//setup Lights
	wr.lights = make([]float32, maxLights*16, maxLights*16)
	wr.directionalLights = make([]float32, maxLights*16, maxLights*16)

	wr.onInit()

	var fn *js.Object
	fn = js.MakeFunc(func(this *js.Object, args []*js.Object) interface{} {
		wr.onUpdate()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		wr.onRender()
		js.Global.Get("window").Call("requestAnimationFrame", fn)
		return 0
	})

	js.Global.Get("window").Call("requestAnimationFrame", fn)
}

func (wr *WebglRenderer) BackGroundColor(r, g, b, a float32) {
	wr.gl.ClearColor(r, g, b, a)
}

func (wr *WebglRenderer) WindowDimensions() mgl32.Vec2 {
	return mgl32.Vec2{float32(wr.WindowWidth), float32(wr.WindowHeight)}
}

// Ortho - set orthogonal rendering mode
func (wr *WebglRenderer) Ortho() {
	gl := wr.gl
	projection := mgl32.Ortho2D(0, float32(wr.WindowWidth), float32(wr.WindowHeight), 0)
	projectionUniform := gl.GetUniformLocation(wr.program, "projection")
	gl.UniformMatrix4fv(projectionUniform, false, projection[:])
	camera := mgl32.Ident4()
	cameraUniform := gl.GetUniformLocation(wr.program, "camera")
	gl.UniformMatrix4fv(cameraUniform, false, camera[:])
	wr.cameraOrtho = true
}

// Perspective - set Perspective camera mode
func (wr *WebglRenderer) Perspective(location, lookat, up mgl32.Vec3, angle, near, far float32) {
	gl := wr.gl
	projection := mgl32.Perspective(mgl32.DegToRad(angle), float32(wr.WindowWidth)/float32(wr.WindowHeight), near, far)
	projectionUniform := gl.GetUniformLocation(wr.program, "projection")
	gl.UniformMatrix4fv(projectionUniform, false, projection[:])
	camera := mgl32.LookAtV(location, lookat, up)
	cameraUniform := gl.GetUniformLocation(wr.program, "camera")
	gl.UniformMatrix4fv(cameraUniform, false, camera[:])
	wr.cameraAngle = angle
	wr.cameraNear = near
	wr.cameraFar = far
	wr.cameraLocation = location
	wr.cameraLookAt = lookat
	wr.cameraUp = up
	wr.cameraOrtho = false
	// wr.updateCameraVectors() // TODO: Frustrum culling
}

func (wr *WebglRenderer) CameraLocation() mgl32.Vec3 {
	return wr.cameraLocation
}

func (wr *WebglRenderer) PushTransform(transform mgl32.Mat4) {
	gl := wr.gl
	wr.transformStack.Push(transform)
	model := wr.transformStack.Peek()
	modelNormal := model.Inv().Transpose()
	gl.UniformMatrix4fv(wr.modelUniform, false, model[:])
	gl.UniformMatrix4fv(wr.modelNormalUniform, false, modelNormal[:])
}

func (wr *WebglRenderer) PopTransform() {
	gl := wr.gl
	wr.transformStack.Pop()
	model := wr.transformStack.Peek()
	modelNormal := model.Inv().Transpose()
	gl.UniformMatrix4fv(wr.modelUniform, false, model[:])
	gl.UniformMatrix4fv(wr.modelNormalUniform, false, modelNormal[:])
}

func (wr *WebglRenderer) FrustrumContainsSphere(radius float32) bool {
	return true // TODO: Frustrum culling
}

func (wr *WebglRenderer) EnableDepthTest(depthTest bool) {
	gl := wr.gl
	if depthTest {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}

func (wr *WebglRenderer) EnableDepthMask(depthMast bool) {
	gl := wr.gl
	gl.DepthMask(depthMast)
}

func (wr *WebglRenderer) CreateGeometry(geometry *renderer.Geometry) {
	gl := wr.gl

	// Configure the vertex data
	vbo := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, geometry.Verticies, gl.DYNAMIC_DRAW)
	geometry.VboId = wr.addBuffer(vbo)

	ibo := gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, geometry.Indicies, gl.DYNAMIC_DRAW)
	geometry.IboId = wr.addBuffer(ibo)
}

func (wr *WebglRenderer) addBuffer(bufferObject *js.Object) uint32 {
	wr.bufferIds++
	wr.buffers[wr.bufferIds] = bufferObject
	return wr.bufferIds
}

func (wr *WebglRenderer) getBuffer(id uint32) *js.Object {
	return wr.buffers[id]
}

func (wr *WebglRenderer) removeBuffer(id uint32) {
	delete(wr.buffers, id)
}

func (wr *WebglRenderer) DestroyGeometry(geometry *renderer.Geometry) {
	//TODO
}

func (wr *WebglRenderer) DrawGeometry(geometry *renderer.Geometry) {
	gl := wr.gl

	gl.BindBuffer(gl.ARRAY_BUFFER, wr.getBuffer(geometry.VboId))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, wr.getBuffer(geometry.IboId))

	//update buffers
	if geometry.VboDirty && len(geometry.Verticies) > 0 && len(geometry.Indicies) > 0 {
		gl.BufferData(gl.ARRAY_BUFFER, geometry.Verticies, gl.DYNAMIC_DRAW)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, geometry.Indicies, gl.DYNAMIC_DRAW)
		geometry.VboDirty = false
	}

	//set back face culling
	if geometry.CullBackface {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}

	//set depthbuffer modes
	wr.EnableDepthTest(geometry.Material.DepthTest)
	wr.EnableDepthMask(geometry.Material.DepthMask)

	//set lighting mode
	lightsUniform := gl.GetUniformLocation(wr.program, "mode")
	gl.Uniform1i(lightsUniform, int(geometry.Material.LightingMode))

	//world camera position
	cam := wr.CameraLocation()
	worldCamPosUniform := gl.GetUniformLocation(wr.program, "worldCamPos")
	gl.Uniform4f(worldCamPosUniform, cam[0], cam[1], cam[2], 0)

	//transparency mode
	if geometry.Material.Transparency == renderer.TRANSPARENCY_NON_EMISSIVE {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	} else if geometry.Material.Transparency == renderer.TRANSPARENCY_EMISSIVE {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	} else {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	}
	//set verticies attribute
	vertAttrib := gl.GetAttribLocation(wr.program, "vert")
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, renderer.VertexStride*4, 0)
	//set normals attribute
	normAttrib := gl.GetAttribLocation(wr.program, "vertNormal")
	gl.EnableVertexAttribArray(normAttrib)
	gl.VertexAttribPointer(normAttrib, 3, gl.FLOAT, false, renderer.VertexStride*4, 3*4)
	//set texture coord attribute
	texCoordAttrib := gl.GetAttribLocation(wr.program, "vertTexCoord")
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, renderer.VertexStride*4, 6*4)
	//vertex color attribute
	colorAttrib := gl.GetAttribLocation(wr.program, "color")
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(colorAttrib, 4, gl.FLOAT, false, renderer.VertexStride*4, 8*4)

	//setup textures
	diffuse := wr.getBuffer(geometry.Material.DiffuseId)
	normal := wr.getBuffer(geometry.Material.NormalId)
	specular := wr.getBuffer(geometry.Material.SpecularId)
	roughness := wr.getBuffer(geometry.Material.RoughnessId)

	gl.ActiveTexture(gl.TEXTURE0)
	if diffuse != nil {
		gl.BindTexture(gl.TEXTURE_2D, diffuse)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, wr.defaultTexture)
	}

	gl.ActiveTexture(gl.TEXTURE1)
	if normal != nil {
		gl.BindTexture(gl.TEXTURE_2D, normal)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, wr.defaultTexture)
	}

	gl.ActiveTexture(gl.TEXTURE2)
	if specular != nil {
		gl.BindTexture(gl.TEXTURE_2D, specular)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, wr.defaultTexture)
	}

	gl.ActiveTexture(gl.TEXTURE3)
	if roughness != nil {
		gl.BindTexture(gl.TEXTURE_2D, roughness)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, wr.defaultTexture)
	}

	gl.ActiveTexture(gl.TEXTURE4)
	if wr.envMap != nil {
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, wr.envMap)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, wr.defaultTexture)
	}

	gl.ActiveTexture(gl.TEXTURE5)
	if wr.envMapLOD1 != nil {
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, wr.envMapLOD1)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, wr.defaultTexture)
	}

	gl.ActiveTexture(gl.TEXTURE6)
	if wr.envMapLOD2 != nil {
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, wr.envMapLOD2)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, wr.defaultTexture)
	}

	gl.ActiveTexture(gl.TEXTURE7)
	if wr.envMapLOD3 != nil {
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, wr.envMapLOD3)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, wr.defaultTexture)
	}

	gl.ActiveTexture(gl.TEXTURE8)
	if wr.illuminanceMap != nil {
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, wr.illuminanceMap)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, wr.defaultTexture)
	}

	useVertexColorUniform := gl.GetUniformLocation(wr.program, "useVertexColor")
	if geometry.Material.Diffuse == nil {
		gl.Uniform1i(useVertexColorUniform, 1)
	} else {
		gl.Uniform1i(useVertexColorUniform, 0)
	}

	gl.DrawElements(gl.TRIANGLES, (len(geometry.Indicies)), gl.UNSIGNED_SHORT, 0)
}

func (wr *WebglRenderer) CreateMaterial(material *renderer.Material) {}

func (wr *WebglRenderer) DestroyMaterial(material *renderer.Material) {}

func (wr *WebglRenderer) CreateLight(ar, ag, ab, dr, dg, db, sr, sg, sb float32, directional bool, position mgl32.Vec3, i int) {
}

func (wr *WebglRenderer) DestroyLight(i int) {}

func (wr *WebglRenderer) ReflectionMap(cm *renderer.CubeMap)        {}
func (wr *WebglRenderer) CreatePostEffect(shader renderer.Shader)   {}
func (wr *WebglRenderer) DestroyPostEffects(shader renderer.Shader) {}
func (wr *WebglRenderer) LockCursor(lock bool)                      {}

//setup Texture
func (wr *WebglRenderer) newTexture(img image.Image, textureUnit int) *js.Object {
	gl := wr.gl
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	tex := gl.CreateTexture()
	gl.ActiveTexture(textureUnit)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, js.MakeWrapper(rgba))
	return tex
}
