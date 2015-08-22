package renderer

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	_ "image/jpeg"
	"log"
	"os"
	"strings"

	"goEngine/vectorMath"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Renderer interface {
	Start()
	BackGroundColor( r,g,b,a float32 )
	Projection( angle, aspect, near, far float32 )
	Camera( location, lookat, up vectorMath.Vector3 )
	PopTransform()
	PushTransform()
	ApplyTransform( transform Transform )
	CreateGeometry( geometry *Geometry )
	DestroyGeometry( geometry *Geometry )
	DrawGeometry( geometry *Geometry )
}

type Transform interface {
	ApplyTransform( transform Transform ) 
}

type GlTransform struct {
	Mat mgl32.Mat4
}

func (glTx *GlTransform) ApplyTransform( transform Transform ) {
	switch v := transform.(type) {
    default:
        fmt.Printf("unexpected type for ApplyTransform GlTransform: %T", v)
    case *GlTransform:
		glTx.Mat = glTx.Mat.Mul4( transform.(*GlTransform).Mat )
    } 
}

//used to combine transformations
func (s *Stack) MultiplyAll() mgl32.Mat4 {
	result := mgl32.Ident4()
	for i:=0 ; i<s.size ; i++ {
		tx := s.Get(i).(*GlTransform)
		result = result.Mul4(tx.Mat)
	}
	return result
}

///////////////////
//OPEN GL Renderer
type OpenglRenderer struct {
	Init, Update, Render func()
	WindowWidth, WindowHeight int
	WindowTitle string
	matStack *Stack
	program uint32
	modelUniform int32
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

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)
	glRenderer.program = program

	glRenderer.Projection( 45.0, float32(glRenderer.WindowWidth)/float32(glRenderer.WindowHeight), 0.1, 10000.0 )
	glRenderer.Camera( vectorMath.Vector3{3,3,3}, vectorMath.Vector3{0,0,0}, vectorMath.Vector3{0,1,0} )

	matStack := CreateStack()
	glRenderer.matStack = matStack
	model := mgl32.Ident4()
	glRenderer.modelUniform = gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(glRenderer.modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.2, 0.0, 0.0, 1.0)

	glRenderer.Init()

	for !window.ShouldClose() {

		glRenderer.Update()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render
		gl.UseProgram(program)

		glRenderer.Render()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (glRenderer *OpenglRenderer) BackGroundColor( r,g,b,a float32 ){
	gl.ClearColor( r,g,b,a )
}

func (glRenderer *OpenglRenderer) Projection( angle, aspect, near, far float32 ) {
	projection := mgl32.Perspective(mgl32.DegToRad(angle), aspect, near, far)
	projectionUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
}

func (glRenderer *OpenglRenderer) Camera( location, lookat, up vectorMath.Vector3 ) {
	camera := mgl32.LookAtV(convertVector(location), convertVector(lookat), convertVector(up))
	cameraUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
}

func (glRenderer *OpenglRenderer) PushTransform(){
	glRenderer.matStack.Push( &GlTransform{ mgl32.Ident4() } )
}

func (glRenderer *OpenglRenderer) PopTransform(){
	glRenderer.matStack.Pop()
	model := glRenderer.matStack.MultiplyAll()
	gl.UniformMatrix4fv(glRenderer.modelUniform, 1, false, &model[0])
}

func (glRenderer *OpenglRenderer) ApplyTransform( transform Transform ){
	tx := glRenderer.matStack.Pop().(*GlTransform)
	tx.ApplyTransform( transform )
	glRenderer.matStack.Push(tx)
	model := glRenderer.matStack.MultiplyAll()
	gl.UniformMatrix4fv(glRenderer.modelUniform, 1, false, &model[0])
}

func convertVector( v vectorMath.Vector3 ) mgl32.Vec3{
	return mgl32.Vec3{(float32)(v.X), (float32)(v.Y), (float32)(v.Z)}
}

func (glRenderer *OpenglRenderer) CreateGeometry( geometry *Geometry ) {
	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(geometry.Verticies)*4, gl.Ptr(geometry.Verticies), gl.STATIC_DRAW)

	var elementbuffer uint32
	gl.GenBuffers(1, &elementbuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, elementbuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(geometry.Indicies)*4, gl.Ptr(geometry.Indicies), gl.STATIC_DRAW)
	geometry.vaoId = elementbuffer

	// Load the texture
	texture, err := newTexture("square.jpg")
	if err != nil {
		panic(err)
	}
	(*geometry).textureId = texture
}

func (glRenderer *OpenglRenderer) DestroyGeometry( geometry *Geometry ) {

}

func (glRenderer *OpenglRenderer) DrawGeometry( geometry *Geometry ) {

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, geometry.vaoId)

	vertAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(len(geometry.Verticies)*4))
	
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, geometry.textureId)

	gl.DrawElements(gl.TRIANGLES, (int32)(len(geometry.Indicies)), gl.UNSIGNED_INT, gl.PtrOffset(0))
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

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
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

	return texture, nil
}

////////////////
//TEST data

var vertexShader string = `
#version 330
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;
void main() {
	fragTexCoord = vertTexCoord;
	gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"

var fragmentShader = `
#version 330
uniform sampler2D tex;
in vec2 fragTexCoord;
out vec4 outputColor;
void main() {
	outputColor = texture(tex, fragTexCoord);
}
` + "\x00"
