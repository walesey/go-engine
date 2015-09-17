package renderer

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/vectormath"
)

const (
	MAX_LIGHTS int = 8
)

//Renderer API
type Renderer interface {
	Start()
	BackGroundColor(r, g, b, a float32)
	Projection(angle, aspect, near, far float32)
	Camera(location, lookat, up vectormath.Vector3)
	CameraLocation() vectormath.Vector3
	PopTransform()
	PushTransform()
	EnableDepthTest(depthTest bool)
	ApplyTransform(transform Transform)
	CreateGeometry(geometry *Geometry)
	DestroyGeometry(geometry *Geometry)
	CreateMaterial(material *Material)
	DestroyMaterial(material *Material)
	DrawGeometry(geometry *Geometry)
	CreateLight(ar, ag, ab, dr, dg, db, sr, sg, sb float32, directional bool, position vectormath.Vector3, i int)
	DestroyLight(i int)
	ReflectionMap(cm CubeMap)
	CreatePostEffect(shader Shader)
	DestroyPostEffects(shader Shader)
}

//used to combine transformations
func (s *Stack) MultiplyAll() mgl32.Mat4 {
	result := mgl32.Ident4()
	for i := 0; i < s.size; i++ {
		tx := s.Get(i).(*GlTransform)
		result = result.Mul4(tx.Mat)
	}
	return result
}

///////////////////
//OPEN GL Renderer
type OpenglRenderer struct {
	Init, Update, Render      func()
	WindowWidth, WindowHeight int
	WindowTitle               string
	Window                    *glfw.Window
	matStack                  Stack
	program                   uint32
	envMapId                  uint32
	envMapLOD1Id              uint32
	envMapLOD2Id              uint32
	envMapLOD3Id              uint32
	illuminanceMapId          uint32
	modelUniform              int32
	lights                    []float32
	directionalLights         []float32
	cameraLocation            vectormath.Vector3
	postEffectVbo             uint32
	postEffects               []postEffect
}

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
	window, err := glfw.CreateWindow(glRenderer.WindowWidth, glRenderer.WindowHeight, glRenderer.WindowTitle, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glRenderer.Window = window

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program := programFromFile("shaders/main.vert", "shaders/main.frag")
	gl.UseProgram(program)
	glRenderer.program = program

	//set default camera
	glRenderer.Projection(45.0, float32(glRenderer.WindowWidth)/float32(glRenderer.WindowHeight), 0.1, 10000.0)
	glRenderer.Camera(vectormath.Vector3{3, 3, 3}, vectormath.Vector3{0, 0, 0}, vectormath.Vector3{0, 1, 0})

	//create mat stack for push pop stack
	matStack := CreateStack()
	glRenderer.matStack = matStack
	glRenderer.PushTransform()
	model := mgl32.Ident4()

	//set shader uniforms
	glRenderer.modelUniform = gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(glRenderer.modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("diffuse\x00"))
	gl.Uniform1i(textureUniform, 0)
	textureUniform = gl.GetUniformLocation(program, gl.Str("normal\x00"))
	gl.Uniform1i(textureUniform, 1)
	textureUniform = gl.GetUniformLocation(program, gl.Str("specular\x00"))
	gl.Uniform1i(textureUniform, 2)
	textureUniform = gl.GetUniformLocation(program, gl.Str("roughness\x00"))
	gl.Uniform1i(textureUniform, 3)
	textureUniform = gl.GetUniformLocation(program, gl.Str("environmentMap\x00"))
	gl.Uniform1i(textureUniform, 4)
	textureUniform = gl.GetUniformLocation(program, gl.Str("environmentMapLOD1\x00"))
	gl.Uniform1i(textureUniform, 5)
	textureUniform = gl.GetUniformLocation(program, gl.Str("environmentMapLOD2\x00"))
	gl.Uniform1i(textureUniform, 6)
	textureUniform = gl.GetUniformLocation(program, gl.Str("environmentMapLOD3\x00"))
	gl.Uniform1i(textureUniform, 7)
	textureUniform = gl.GetUniformLocation(program, gl.Str("illuminanceMap\x00"))
	gl.Uniform1i(textureUniform, 8)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Configure global settings
	gl.Enable(gl.TEXTURE_CUBE_MAP_SEAMLESS)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	//setup Lights
	glRenderer.lights = make([]float32, MAX_LIGHTS*16, MAX_LIGHTS*16)
	glRenderer.directionalLights = make([]float32, MAX_LIGHTS*16, MAX_LIGHTS*16)
	glRenderer.CreateLight(0.1, 0.1, 0.1, 1, 1, 1, 1, 1, 1, true, vectormath.Vector3{0, -1, 0}, 0)

	glRenderer.initPostEffects()

	glRenderer.Init()

	//Main loop
	for !window.ShouldClose() {

		glRenderer.Update()
		gl.UseProgram(program)
		if len(glRenderer.postEffects) == 0 {
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			glRenderer.Render()
		} else {

			//Render to the first post effect buffer
			gl.BindFramebuffer(gl.FRAMEBUFFER, glRenderer.postEffects[0].fboId)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
			glRenderer.Render()
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
		gl.UseProgram(program)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (glRenderer *OpenglRenderer) BackGroundColor(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

func (glRenderer *OpenglRenderer) Projection(angle, aspect, near, far float32) {
	projection := mgl32.Perspective(mgl32.DegToRad(angle), aspect, near, far)
	projectionUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
}

func (glRenderer *OpenglRenderer) Camera(location, lookat, up vectormath.Vector3) {
	camera := mgl32.LookAtV(convertVector(location), convertVector(lookat), convertVector(up))
	cameraUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
	//store camera location
	glRenderer.cameraLocation = location
}

func (glRenderer *OpenglRenderer) CameraLocation() vectormath.Vector3 {
	return glRenderer.cameraLocation
}

func (glRenderer *OpenglRenderer) PushTransform() {
	glRenderer.matStack.Push(&GlTransform{mgl32.Ident4()})
}

func (glRenderer *OpenglRenderer) PopTransform() {
	glRenderer.matStack.Pop()
	model := glRenderer.matStack.MultiplyAll()
	gl.UniformMatrix4fv(glRenderer.modelUniform, 1, false, &model[0])
}

func (glRenderer *OpenglRenderer) ApplyTransform(transform Transform) {
	tx := glRenderer.matStack.Pop().(*GlTransform)
	tx.ApplyTransform(transform)
	glRenderer.matStack.Push(tx)
	model := glRenderer.matStack.MultiplyAll()
	gl.UniformMatrix4fv(glRenderer.modelUniform, 1, false, &model[0])
}

func (glRenderer *OpenglRenderer) EnableDepthTest(depthTest bool) {
	if depthTest {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}

//
func (glRenderer *OpenglRenderer) CreateGeometry(geometry *Geometry) {

	// Configure the vertex data
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(geometry.Verticies)*4, gl.Ptr(geometry.Verticies), gl.DYNAMIC_DRAW)
	geometry.vboId = vbo

	var ibo uint32
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(geometry.Indicies)*4, gl.Ptr(geometry.Indicies), gl.DYNAMIC_DRAW)
	geometry.iboId = ibo
}

//
func (glRenderer *OpenglRenderer) DestroyGeometry(geometry *Geometry) {
	gl.DeleteBuffers(1, &geometry.vboId)
	gl.DeleteBuffers(1, &geometry.iboId)
}

//setup Texture
func (glRenderer *OpenglRenderer) CreateMaterial(material *Material) {
	if material.Diffuse != nil {
		material.diffuseId = glRenderer.newTexture(material.Diffuse, gl.TEXTURE0)
	}
	if material.Normal != nil {
		material.normalId = glRenderer.newTexture(material.Normal, gl.TEXTURE1)
	}
	if material.Specular != nil {
		material.specularId = glRenderer.newTexture(material.Specular, gl.TEXTURE2)
	}
	if material.Roughness != nil {
		material.roughnessId = glRenderer.newTexture(material.Roughness, gl.TEXTURE3)
	}
}

//setup Texture
func (glRenderer *OpenglRenderer) newTexture(img image.Image, textureUnit uint32) uint32 {
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.Fatal("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	var texId uint32
	gl.GenTextures(1, &texId)
	gl.ActiveTexture(textureUnit)
	gl.BindTexture(gl.TEXTURE_2D, texId)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	return texId
}

func (glRenderer *OpenglRenderer) ReflectionMap(cm CubeMap) {
	cm.Resize(512)
	glRenderer.envMapId = glRenderer.newCubeMap(cm.Right, cm.Left, cm.Top, cm.Bottom, cm.Back, cm.Front, gl.TEXTURE4)
	cm.Resize(64)
	glRenderer.envMapLOD1Id = glRenderer.newCubeMap(cm.Right, cm.Left, cm.Top, cm.Bottom, cm.Back, cm.Front, gl.TEXTURE5)
	cm.Resize(32)
	glRenderer.envMapLOD2Id = glRenderer.newCubeMap(cm.Right, cm.Left, cm.Top, cm.Bottom, cm.Back, cm.Front, gl.TEXTURE6)
	cm.Resize(16)
	glRenderer.envMapLOD3Id = glRenderer.newCubeMap(cm.Right, cm.Left, cm.Top, cm.Bottom, cm.Back, cm.Front, gl.TEXTURE7)
	cm.Resize(8)
	glRenderer.illuminanceMapId = glRenderer.newCubeMap(cm.Right, cm.Left, cm.Top, cm.Bottom, cm.Back, cm.Front, gl.TEXTURE8)
}

func (glRenderer *OpenglRenderer) newCubeMap(right, left, top, bottom, back, front image.Image, textureUnit uint32) uint32 {
	var texId uint32
	gl.GenTextures(1, &texId)
	gl.ActiveTexture(textureUnit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, texId)

	for i := 0; i < 6; i++ {
		img := right
		texIndex := (uint32)(gl.TEXTURE_CUBE_MAP_POSITIVE_X)
		if i == 1 {
			img = left
			texIndex = gl.TEXTURE_CUBE_MAP_NEGATIVE_X
		} else if i == 2 {
			img = top
			texIndex = gl.TEXTURE_CUBE_MAP_NEGATIVE_Y
		} else if i == 3 {
			img = bottom
			texIndex = gl.TEXTURE_CUBE_MAP_POSITIVE_Y
		} else if i == 4 {
			img = back
			texIndex = gl.TEXTURE_CUBE_MAP_NEGATIVE_Z
		} else if i == 5 {
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
			gl.Ptr(rgba.Pix))
	}
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	return texId
}

//
func (glRenderer *OpenglRenderer) DestroyMaterial(material *Material) {
	gl.DeleteTextures(1, &material.diffuseId)
	gl.DeleteTextures(1, &material.normalId)
	gl.DeleteTextures(1, &material.specularId)
	gl.DeleteTextures(1, &material.roughnessId)
}

//
func (glRenderer *OpenglRenderer) DrawGeometry(geometry *Geometry) {

	gl.BindBuffer(gl.ARRAY_BUFFER, geometry.vboId)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, geometry.iboId)

	//update buffers
	if geometry.VboDirty {
		gl.BufferData(gl.ARRAY_BUFFER, len(geometry.Verticies)*4, gl.Ptr(geometry.Verticies), gl.DYNAMIC_DRAW)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(geometry.Indicies)*4, gl.Ptr(geometry.Indicies), gl.DYNAMIC_DRAW)
		geometry.VboDirty = false
	}

	//set back face culling
	if geometry.CullBackface {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}

	//set lighting mode
	lightsUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("mode\x00"))
	gl.Uniform1i(lightsUniform, geometry.Material.LightingMode)

	//set verticies attribute
	vertAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 18*4, gl.PtrOffset(0))

	//set normals/tangent attribute
	normAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("normal\x00")))
	gl.EnableVertexAttribArray(normAttrib)
	gl.VertexAttribPointer(normAttrib, 3, gl.FLOAT, false, 18*4, gl.PtrOffset(3*4))
	tangentAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("tangent\x00")))
	gl.EnableVertexAttribArray(tangentAttrib)
	gl.VertexAttribPointer(tangentAttrib, 3, gl.FLOAT, false, 18*4, gl.PtrOffset(6*4))
	bitangentAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("bitangent\x00")))
	gl.EnableVertexAttribArray(bitangentAttrib)
	gl.VertexAttribPointer(bitangentAttrib, 3, gl.FLOAT, false, 18*4, gl.PtrOffset(9*4))

	//set texture coord attribute
	texCoordAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 18*4, gl.PtrOffset(12*4))
	//vertex color attribute
	colorAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("color\x00")))
	gl.EnableVertexAttribArray(colorAttrib)
	gl.VertexAttribPointer(colorAttrib, 4, gl.FLOAT, false, 18*4, gl.PtrOffset(14*4))

	//setup textures
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, geometry.Material.diffuseId)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, geometry.Material.normalId)
	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, geometry.Material.specularId)
	gl.ActiveTexture(gl.TEXTURE3)
	gl.BindTexture(gl.TEXTURE_2D, geometry.Material.roughnessId)
	gl.ActiveTexture(gl.TEXTURE4)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, glRenderer.envMapId)
	gl.ActiveTexture(gl.TEXTURE5)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, glRenderer.envMapLOD1Id)
	gl.ActiveTexture(gl.TEXTURE6)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, glRenderer.envMapLOD2Id)
	gl.ActiveTexture(gl.TEXTURE7)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, glRenderer.envMapLOD3Id)
	gl.ActiveTexture(gl.TEXTURE8)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, glRenderer.illuminanceMapId)

	gl.DrawElements(gl.TRIANGLES, (int32)(len(geometry.Indicies)), gl.UNSIGNED_INT, gl.PtrOffset(0))
}

// ambient, diffuse and specular light values ( i is the light index )
func (glRenderer *OpenglRenderer) CreateLight(ar, ag, ab, dr, dg, db, sr, sg, sb float32, directional bool, position vectormath.Vector3, i int) {
	lights := &glRenderer.lights
	if directional {
		lights = &glRenderer.directionalLights
	}

	//position
	(*lights)[(i * 16)] = (float32)(position.X)
	(*lights)[(i*16)+1] = (float32)(position.Y)
	(*lights)[(i*16)+2] = (float32)(position.Z)
	(*lights)[(i*16)+3] = 1
	//ambient
	(*lights)[(i*16)+4] = ar
	(*lights)[(i*16)+5] = ag
	(*lights)[(i*16)+6] = ab
	(*lights)[(i*16)+7] = 1
	//diffuse
	(*lights)[(i*16)+8] = dr
	(*lights)[(i*16)+9] = dg
	(*lights)[(i*16)+10] = db
	(*lights)[(i*16)+11] = 1
	//specular
	(*lights)[(i*16)+12] = sr
	(*lights)[(i*16)+13] = sg
	(*lights)[(i*16)+14] = sb
	(*lights)[(i*16)+15] = 1

	//set uniform array
	uniformName := "lights\x00"
	if directional {
		uniformName = "directionalLights\x00"
	}
	lightsUniform := gl.GetUniformLocation(glRenderer.program, gl.Str(uniformName))
	gl.Uniform4fv(lightsUniform, (int32)(MAX_LIGHTS*16), &(*lights)[0])
}

func (glRenderer *OpenglRenderer) DestroyLight(i int) {

}

func programFromFile(vertFilePath, fragFilePath string) uint32 {
	bufVert, err := ioutil.ReadFile(vertFilePath)
	if err != nil {
		panic(err)
	}
	vertexShader := string(bufVert) + "\x00"
	bufFrag, err := ioutil.ReadFile(fragFilePath)
	if err != nil {
		panic(err)
	}
	fragmentShader := string(bufFrag) + "\x00"
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	return program
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, errors.New(fmt.Sprintf("failed to link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
