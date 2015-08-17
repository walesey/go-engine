package Renderer

import (
    "errors"
    "fmt"
    "image"
    "image/draw"
    _ "image/png"
    "log"
    "os"
    "strings"

    "goEngine/VectorMath"

    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/glfw/v3.1/glfw"
    "github.com/go-gl/mathgl/mgl32"
)

type Renderer interface {
    Start()
    BackGroundColor( r,g,b,a float32 )
    Projection( angle, aspect, near, far float32 )
    Camera( location, lookat, up VectorMath.Vector3 )
    CreateGeometry( geometry *Geometry )
    DestroyGeometry( geometry *Geometry )
    DrawGeometry( geometry *Geometry )
}

type OpenglRenderer struct {
    programPtr *uint32
    Init, Update, Render func(renderer Renderer)
    WindowWidth, WindowHeight int
    WindowTitle string
}

func (glRenderer OpenglRenderer) Start() {
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
    glRenderer.programPtr = &program

    glRenderer.Projection( 45.0, float32(glRenderer.WindowWidth)/float32(glRenderer.WindowHeight), 0.1, 10000.0 )
    glRenderer.Camera( VectorMath.Vector3{3,3,3}, VectorMath.Vector3{0,0,0}, VectorMath.Vector3{0,1,0} )

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

    glRenderer.Init( glRenderer )

	for !window.ShouldClose() {

        glRenderer.Update(glRenderer)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render
		gl.UseProgram(program)

		model = mgl32.HomogRotate3D(float32(0), mgl32.Vec3{0, 1, 0})
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

        glRenderer.Render(glRenderer)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (glRenderer OpenglRenderer) ClearBuffers() {
      gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (glRenderer OpenglRenderer) BackGroundColor( r,g,b,a float32 ){
      gl.ClearColor( r,g,b,a )
}

func (glRenderer OpenglRenderer) Projection( angle, aspect, near, far float32 ) {
  	projection := mgl32.Perspective(mgl32.DegToRad(angle), aspect, near, far)
    projectionUniform := gl.GetUniformLocation((*glRenderer.programPtr), gl.Str("projection\x00"))
  	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
}

func (glRenderer OpenglRenderer) Camera( location, lookat, up VectorMath.Vector3 ) {
    camera := mgl32.LookAtV(convertVector(location), convertVector(lookat), convertVector(up))
	cameraUniform := gl.GetUniformLocation((*glRenderer.programPtr), gl.Str("camera\x00"))
    gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
}

func convertVector( v VectorMath.Vector3 ) mgl32.Vec3{
    return mgl32.Vec3{(float32)(v.X), (float32)(v.Y), (float32)(v.Z)}
}

func (glRenderer OpenglRenderer) CreateGeometry( geometry *Geometry ) {
    // Configure the vertex data
    var vao uint32
    gl.GenVertexArrays(1, &vao)
    gl.BindVertexArray(vao)
    (*geometry).vaoId = vao

    var vbo uint32
    gl.GenBuffers(1, &vbo)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, len((*geometry).Verticies)*4, gl.Ptr((*geometry).Verticies), gl.STATIC_DRAW)

    // Load the texture
	texture, err := newTexture("square.png")
	if err != nil {
		panic(err)
	}
    (*geometry).textureId = texture

    vertAttrib := uint32(gl.GetAttribLocation((*glRenderer.programPtr), gl.Str("vert\x00")))
    gl.EnableVertexAttribArray(vertAttrib)
    gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

    texCoordAttrib := uint32(gl.GetAttribLocation((*glRenderer.programPtr), gl.Str("vertTexCoord\x00")))
    gl.EnableVertexAttribArray(texCoordAttrib)
    gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
}

func (glRenderer OpenglRenderer) DestroyGeometry( geometry *Geometry ) {

}


func (glRenderer OpenglRenderer) DrawGeometry( geometry *Geometry ) {
    gl.BindVertexArray((*geometry).vaoId)

    gl.ActiveTexture(gl.TEXTURE0)
    gl.BindTexture(gl.TEXTURE_2D, (*geometry).textureId)

    gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
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
