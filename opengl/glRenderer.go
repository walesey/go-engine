package opengl

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"runtime"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
)

const MAX_POINT_LIGHTS = 4
const MAX_DIRECTIONAL_LIGHTS = 4

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// OpenglRenderer - opengl implementation
type OpenglRenderer struct {
	onInit, onUpdate, onRender func()
	WindowWidth, WindowHeight  int
	FullScreen                 bool
	WindowTitle                string
	Window                     *glfw.Window
	camera                     *renderer.Camera

	postEffectVbo uint32
	postEffects   []postEffect

	material, activeMaterial *renderer.Material
	shader, activeShader     *renderer.Shader
	cubeMap, activeCubeMap   *renderer.CubeMap

	transparency   renderer.Transparency
	rendererParams renderer.RendererParams

	depthTest, depthMast, cullFace, unlit, useTextures bool

	lights                  []*renderer.Light
	ambientLightValue       mgl32.Vec3
	nbPointLights           int32
	pointLightValues        []float32
	pointLightPositions     []float32
	nbDirectionalLights     int32
	directionalLightValues  []float32
	directionalLightVectors []float32
}

// NewOpenglRenderer - create new renderer
func NewOpenglRenderer(WindowTitle string, WindowWidth, WindowHeight int, FullScreen bool) *OpenglRenderer {
	return &OpenglRenderer{
		WindowTitle:             WindowTitle,
		WindowWidth:             WindowWidth,
		WindowHeight:            WindowHeight,
		FullScreen:              FullScreen,
		pointLightValues:        make([]float32, MAX_POINT_LIGHTS*4),
		pointLightPositions:     make([]float32, MAX_POINT_LIGHTS*4),
		directionalLightValues:  make([]float32, MAX_DIRECTIONAL_LIGHTS*4),
		directionalLightVectors: make([]float32, MAX_DIRECTIONAL_LIGHTS*4),
	}
}

func (glRenderer *OpenglRenderer) SetInit(callback func()) {
	glRenderer.onInit = callback
}

func (glRenderer *OpenglRenderer) SetUpdate(callback func()) {
	glRenderer.onUpdate = callback
}

func (glRenderer *OpenglRenderer) SetRender(callback func()) {
	glRenderer.onRender = callback
}

func (glRenderer *OpenglRenderer) SetCamera(camera *renderer.Camera) {
	glRenderer.camera = camera
}

//Start -
func (glRenderer *OpenglRenderer) Start() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	if glRenderer.FullScreen {
		glRenderer.WindowWidth = glfw.GetPrimaryMonitor().GetVideoMode().Width
	} else if glRenderer.WindowWidth == 0 {
		glRenderer.WindowWidth = glfw.GetPrimaryMonitor().GetVideoMode().Width * 95 / 100
	}
	if glRenderer.FullScreen {
		glRenderer.WindowHeight = glfw.GetPrimaryMonitor().GetVideoMode().Height
	} else if glRenderer.WindowHeight == 0 {
		glRenderer.WindowHeight = glfw.GetPrimaryMonitor().GetVideoMode().Height * 95 / 100
	}

	var monitor *glfw.Monitor
	if glRenderer.FullScreen {
		monitor = glfw.GetPrimaryMonitor()
	}
	window, err := glfw.CreateWindow(glRenderer.WindowWidth, glRenderer.WindowHeight, glRenderer.WindowTitle, monitor, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glRenderer.Window = window

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.Enable(gl.TEXTURE_CUBE_MAP_SEAMLESS)
	gl.Enable(gl.BLEND)
	gl.DepthFunc(gl.LEQUAL)
	gl.CullFace(gl.BACK)

	glRenderer.initPostEffects()

	glRenderer.onInit()

	window.SetRefreshCallback(func(w *glfw.Window) {
		glRenderer.mainLoop()
		window.SwapBuffers()
	})

	//Main loop
	for !window.ShouldClose() {
		glRenderer.mainLoop()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (glRenderer *OpenglRenderer) mainLoop() {
	glRenderer.onUpdate()
	glRenderer.updateLights()

	//set defaults
	glRenderer.UseRendererParams(renderer.DefaultRendererParams())
	glRenderer.UseMaterial(nil)

	if len(glRenderer.postEffects) == 0 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glRenderer.onRender()
	} else {
		//Render to the first post effect buffer
		gl.BindFramebuffer(gl.FRAMEBUFFER, glRenderer.postEffects[0].fboId)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glRenderer.onRender()
		//Render Post effects
		for i := 0; i < len(glRenderer.postEffects)-1; i = i + 1 {
			gl.BindFramebuffer(gl.FRAMEBUFFER, glRenderer.postEffects[i+1].fboId)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			glRenderer.renderPostEffect(glRenderer.postEffects[i])
		}
		//Render final post effect to the frame buffer
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glRenderer.renderPostEffect(glRenderer.postEffects[len(glRenderer.postEffects)-1])
	}
}

// BackGroundColor - set background color for the scene
func (glRenderer *OpenglRenderer) BackGroundColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (glRenderer *OpenglRenderer) WindowDimensions() mgl32.Vec2 {
	return mgl32.Vec2{float32(glRenderer.WindowWidth), float32(glRenderer.WindowHeight)}
}

func (glRenderer *OpenglRenderer) setTransparency(transparency renderer.Transparency) {
	if transparency == glRenderer.transparency {
		return
	}
	switch transparency {
	case renderer.NON_EMISSIVE:
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	case renderer.EMISSIVE:
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	default:
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	}
	glRenderer.transparency = transparency
}

func (glRenderer *OpenglRenderer) enableCullFace(cullFace bool) {
	if cullFace == glRenderer.cullFace {
		return
	}
	if cullFace {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
	glRenderer.cullFace = cullFace
}

func (glRenderer *OpenglRenderer) enableDepthTest(depthTest bool) {
	if depthTest == glRenderer.depthTest {
		return
	}
	if depthTest {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
	glRenderer.depthTest = depthTest
}

func (glRenderer *OpenglRenderer) enableDepthMask(depthMask bool) {
	if depthMask == glRenderer.depthMast {
		return
	}
	gl.DepthMask(depthMask)
	glRenderer.depthMast = depthMask
}

func (glRenderer *OpenglRenderer) enableUnlit(unlit bool) {
	glRenderer.unlit = unlit
}

func (glRenderer *OpenglRenderer) enableMaterial() {
	if glRenderer.material == glRenderer.activeMaterial {
		return
	}

	glRenderer.activeMaterial = glRenderer.material
	glRenderer.useTextures = (glRenderer.activeMaterial != nil && len(glRenderer.activeMaterial.Textures) > 0)

	// setup material
	if glRenderer.activeShader != nil && glRenderer.activeMaterial != nil {
		glRenderer.createMaterial(glRenderer.activeMaterial)
		textures := glRenderer.activeMaterial.Textures
		for i, tex := range textures {
			textureUnit := gl.TEXTURE0 + uint32(i)
			glRenderer.activeShader.Uniforms[tex.TextureName] = int32(i)
			gl.ActiveTexture(textureUnit)
			gl.BindTexture(gl.TEXTURE_2D, tex.TextureId)

		}
	}
}

func (glRenderer *OpenglRenderer) enableCubeMap() {
	if glRenderer.cubeMap == glRenderer.activeCubeMap {
		return
	}

	glRenderer.activeCubeMap = glRenderer.cubeMap
	cubeMap := glRenderer.activeCubeMap

	// setup cubeMap
	if glRenderer.activeShader != nil && cubeMap != nil {
		glRenderer.createCubeMap(glRenderer.activeCubeMap)
		gl.ActiveTexture(gl.TEXTURE10)
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubeMap.Id)
		glRenderer.activeShader.Uniforms[cubeMap.Name] = int32(10)
	}
}

func (glRenderer *OpenglRenderer) enableShader() {
	if glRenderer.shader == glRenderer.activeShader {
		return
	}

	glRenderer.createShader(glRenderer.shader)
	glRenderer.activeShader = glRenderer.shader
	gl.UseProgram(glRenderer.activeShader.Program)
}

// CreateGeometry - add geometry to the renderer
func (glRenderer *OpenglRenderer) createGeometry(geometry *renderer.Geometry) {
	if geometry.Loaded || len(geometry.Verticies) == 0 || len(geometry.Indicies) == 0 {
		return
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(geometry.Verticies)*4, gl.Ptr(geometry.Verticies), gl.DYNAMIC_DRAW)
	geometry.VboId = vbo

	var ibo uint32
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(geometry.Indicies)*4, gl.Ptr(geometry.Indicies), gl.DYNAMIC_DRAW)
	geometry.IboId = ibo
	geometry.Loaded = true
}

//
func (glRenderer *OpenglRenderer) DestroyGeometry(geometry *renderer.Geometry) {
	gl.DeleteBuffers(1, &geometry.VboId)
	gl.DeleteBuffers(1, &geometry.IboId)
}

// CreateMaterial load material
func (glRenderer *OpenglRenderer) createMaterial(material *renderer.Material) {
	for i, tex := range material.Textures {
		if !tex.Loaded {
			textureUnit := gl.TEXTURE0 + uint32(i)
			tex.TextureId = glRenderer.loadTexture(tex.Img, textureUnit, tex.Lod)
			tex.Loaded = true
		}
	}
}

//
func (glRenderer *OpenglRenderer) DestroyMaterial(material *renderer.Material) {
	for _, tex := range material.Textures {
		if tex.Loaded {
			gl.DeleteTextures(1, &tex.TextureId)
			tex.Loaded = false
		}
	}
}

func (glRenderer *OpenglRenderer) createCubeMap(cubeMap *renderer.CubeMap) {
	if cubeMap.Loaded {
		return
	}

	cm := cubeMap
	cubeMap.Id = glRenderer.loadCubeMap(cm.Right, cm.Left, cm.Top, cm.Bottom, cm.Back, cm.Front, uint32(10), cm.Lod)
	cubeMap.Loaded = true
}

func (glRenderer *OpenglRenderer) DestroyCubeMap(cubeMap *renderer.CubeMap) {
	if !cubeMap.Loaded {
		return
	}
	gl.DeleteTextures(1, &cubeMap.Id)
	cubeMap.Loaded = false
}

func (glRenderer *OpenglRenderer) createShader(shader *renderer.Shader) {
	if shader.Loaded {
		return
	}

	var shaders []uint32

	if len(shader.VertSrc) > 0 {
		if s, err := compileShader(shader.VertSrc+"\x00", gl.VERTEX_SHADER); err == nil {
			shaders = append(shaders, s)
		} else {
			fmt.Println("Error Compiling Vert Shader: ", err)
		}
	}

	if len(shader.FragSrc) > 0 {
		if s, err := compileShader(shader.FragSrc+"\x00", gl.FRAGMENT_SHADER); err == nil {
			shaders = append(shaders, s)
		} else {
			fmt.Println("Error Compiling Frag Shader: ", err)
		}
	}

	if len(shader.GeoSrc) > 0 {
		if s, err := compileShader(shader.GeoSrc+"\x00", gl.GEOMETRY_SHADER); err == nil {
			shaders = append(shaders, s)
		} else {
			fmt.Println("Error Compiling Geo Shader: ", err)
		}
	}

	program, err := newProgram(shaders...)
	if err != nil {
		fmt.Println("Error Creating Shader Program: ", err)
	}
	shader.Program = program
	shader.Loaded = true

	gl.UseProgram(program)
	gl.BindFragDataLocation(program, 0, gl.Str(fmt.Sprintf("%v\x00", shader.FragDataLocation)))
}

func (glRenderer *OpenglRenderer) loadTexture(img image.Image, textureUnit uint32, lod bool) uint32 {
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.Fatal("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	var texId uint32
	gl.GenTextures(1, &texId)
	gl.ActiveTexture(textureUnit)
	gl.BindTexture(gl.TEXTURE_2D, texId)
	if rgba.Rect.Size().X == 0 || rgba.Rect.Size().Y == 0 {
		return texId
	}
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	if lod {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.GenerateMipmap(gl.TEXTURE_2D)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	}
	return texId
}

func (glRenderer *OpenglRenderer) loadCubeMap(right, left, top, bottom, back, front image.Image, textureUnit uint32, lod bool) uint32 {
	var texId uint32
	gl.GenTextures(1, &texId)
	gl.ActiveTexture(textureUnit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, texId)

	for i := 0; i < 6; i++ {
		img := right
		var texIndex uint32 = gl.TEXTURE_CUBE_MAP_POSITIVE_X
		switch i {
		case 1:
			img = left
			texIndex = gl.TEXTURE_CUBE_MAP_NEGATIVE_X
		case 2:
			img = top
			texIndex = gl.TEXTURE_CUBE_MAP_NEGATIVE_Y
		case 3:
			img = bottom
			texIndex = gl.TEXTURE_CUBE_MAP_POSITIVE_Y
		case 4:
			img = back
			texIndex = gl.TEXTURE_CUBE_MAP_NEGATIVE_Z
		case 5:
			img = front
			texIndex = gl.TEXTURE_CUBE_MAP_POSITIVE_Z
		}
		img = imaging.FlipV(img)
		rgba := image.NewRGBA(img.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			log.Fatal("unsupported stride")
		}
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
		gl.TexImage2D(
			texIndex,
			0,
			gl.RGBA,
			int32(rgba.Rect.Size().X),
			int32(rgba.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix),
		)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
		if lod {
			gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
			gl.GenerateMipmap(gl.TEXTURE_CUBE_MAP)
		} else {
			gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		}
	}
	return texId
}

func (glRenderer *OpenglRenderer) UseRendererParams(params renderer.RendererParams) {
	glRenderer.rendererParams = params
}

func (glRenderer *OpenglRenderer) UseShader(shader *renderer.Shader) {
	glRenderer.shader = shader
}

func (glRenderer *OpenglRenderer) UseMaterial(material *renderer.Material) {
	glRenderer.material = material
}

func (glRenderer *OpenglRenderer) UseCubeMap(cubeMap *renderer.CubeMap) {
	glRenderer.cubeMap = cubeMap
}

func (glRenderer *OpenglRenderer) DrawGeometry(geometry *renderer.Geometry, transform mgl32.Mat4) {
	glRenderer.enableShader()
	glRenderer.enableMaterial()
	glRenderer.enableCubeMap()

	if glRenderer.activeShader == nil {
		panic("ERROR: No shader is configured.")
	}

	glRenderer.createGeometry(geometry)

	shader := glRenderer.activeShader
	program := shader.Program
	params := glRenderer.rendererParams

	glRenderer.enableDepthTest(params.DepthTest)
	glRenderer.enableDepthMask(params.DepthMask)
	glRenderer.enableCullFace(params.CullBackface)
	glRenderer.enableUnlit(params.Unlit)
	glRenderer.setTransparency(params.Transparency)

	// set buffers
	gl.BindBuffer(gl.ARRAY_BUFFER, geometry.VboId)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, geometry.IboId)

	// update buffers
	if geometry.VboDirty && len(geometry.Verticies) > 0 && len(geometry.Indicies) > 0 {
		gl.BufferData(gl.ARRAY_BUFFER, len(geometry.Verticies)*4, gl.Ptr(geometry.Verticies), gl.DYNAMIC_DRAW)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(geometry.Indicies)*4, gl.Ptr(geometry.Indicies), gl.DYNAMIC_DRAW)
		geometry.VboDirty = false
	}

	// set uniforms
	modelNormal := transform.Inv().Transpose()
	shader.Uniforms["model"] = transform
	shader.Uniforms["modelNormal"] = modelNormal

	// set camera uniforms
	cam := glRenderer.camera
	win := glRenderer.WindowDimensions()
	shader.Uniforms["cameraTranslation"] = cam.Translation
	if cam.Ortho {
		shader.Uniforms["projection"] = mgl32.Ortho2D(0, win.X(), win.Y(), 0)
		shader.Uniforms["camera"] = mgl32.Ident4()
	} else {
		shader.Uniforms["projection"] = mgl32.Perspective(mgl32.DegToRad(cam.Angle), win.X()/win.Y(), cam.Near, cam.Far)
		shader.Uniforms["camera"] = mgl32.LookAtV(cam.Translation, cam.Lookat, cam.Up)
	}

	shader.Uniforms["unlit"] = glRenderer.unlit
	shader.Uniforms["useTextures"] = glRenderer.useTextures

	shader.Uniforms["ambientLightValue"] = glRenderer.ambientLightValue

	shader.Uniforms["nbPointLights"] = glRenderer.nbPointLights
	shader.Uniforms["pointLightValues"] = glRenderer.pointLightValues
	shader.Uniforms["pointLightPositions"] = glRenderer.pointLightPositions

	shader.Uniforms["nbDirectionalLights"] = glRenderer.nbDirectionalLights
	shader.Uniforms["directionalLightValues"] = glRenderer.directionalLightValues
	shader.Uniforms["directionalLightVectors"] = glRenderer.directionalLightVectors

	// set custom uniforms
	setupUniforms(shader)

	// set verticies attribute
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, renderer.VertexStride*4, gl.PtrOffset(0))
	// set normals attribute
	normAttrib := uint32(gl.GetAttribLocation(program, gl.Str("normal\x00")))
	gl.EnableVertexAttribArray(normAttrib)
	gl.VertexAttribPointer(normAttrib, 3, gl.FLOAT, false, renderer.VertexStride*4, gl.PtrOffset(3*4))
	// set texture coord attribute
	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("texCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, renderer.VertexStride*4, gl.PtrOffset(6*4))
	// vertex color attribute
	colorAttrib := uint32(gl.GetAttribLocation(program, gl.Str("color\x00")))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(colorAttrib, 4, gl.FLOAT, false, renderer.VertexStride*4, gl.PtrOffset(8*4))

	gl.DrawElements(gl.TRIANGLES, (int32)(len(geometry.Indicies)), gl.UNSIGNED_INT, gl.PtrOffset(0))
}

func (glRenderer *OpenglRenderer) LockCursor(lock bool) {
	glRenderer.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}

func (glRenderer *OpenglRenderer) updateLights() {
	glRenderer.nbPointLights = 0
	glRenderer.nbDirectionalLights = 0
	for _, light := range glRenderer.lights {
		c := light.Color
		p := light.Position
		d := light.Direction
		switch light.LightType {
		case renderer.AMBIENT:
			glRenderer.ambientLightValue = mgl32.Vec3{c[0], c[1], c[2]}
		case renderer.POINT:
			i := glRenderer.nbPointLights
			glRenderer.pointLightValues[i*4], glRenderer.pointLightValues[i*4+1], glRenderer.pointLightValues[i*4+2] = c[0], c[1], c[2]
			glRenderer.pointLightPositions[i*4], glRenderer.pointLightPositions[i*4+1], glRenderer.pointLightPositions[i*4+2] = p[0], p[1], p[2]
			glRenderer.nbPointLights++
		case renderer.DIRECTIONAL:
			i := glRenderer.nbDirectionalLights
			glRenderer.directionalLightValues[i*4], glRenderer.directionalLightValues[i*4+1], glRenderer.directionalLightValues[i*4+2] = c[0], c[1], c[2]
			glRenderer.directionalLightVectors[i*4], glRenderer.directionalLightVectors[i*4+1], glRenderer.directionalLightVectors[i*4+2] = d[0], d[1], d[2]
			glRenderer.nbDirectionalLights++
		}
	}
}

func (glRenderer *OpenglRenderer) AddLight(light *renderer.Light) {
	glRenderer.lights = append(glRenderer.lights, light)
}

func (glRenderer *OpenglRenderer) RemoveLight(light *renderer.Light) {
	for i, l := range glRenderer.lights {
		if l == light {
			glRenderer.lights[i] = glRenderer.lights[len(glRenderer.lights)-1]
			glRenderer.lights[len(glRenderer.lights)-1] = nil
			glRenderer.lights = glRenderer.lights[:len(glRenderer.lights)-1]
			break
		}
	}
}
