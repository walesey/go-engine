package renderer

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"image"
	"image/draw"

	"goEngine/vectorMath"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const(
	MAX_LIGHTS int = 8
)

//Renderer API
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
	CreateMaterial( material *Material )
	DestroyMaterial( material *Material )
	DrawGeometry( geometry *Geometry )
	CreateLight( ar,ag,ab, dr,dg,db, sr,sg,sb float32, position vectorMath.Vector3, i int )
	DestroyLight( i int )
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
	lights []float32
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

	textureUniform := gl.GetUniformLocation(program, gl.Str("diffuse\x00"))
	gl.Uniform1i(textureUniform, 0)
	textureUniform = gl.GetUniformLocation(program, gl.Str("normal\x00"))
	gl.Uniform1i(textureUniform, 1)
	textureUniform = gl.GetUniformLocation(program, gl.Str("specular\x00"))
	gl.Uniform1i(textureUniform, 2)
	textureUniform = gl.GetUniformLocation(program, gl.Str("roughness\x00"))
	gl.Uniform1i(textureUniform, 3)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	//setup Lights
	glRenderer.lights = make([]float32, MAX_LIGHTS*16, MAX_LIGHTS*16)
	glRenderer.CreateLight( 100000,100000,100000,   1600000,1600000,1600000,   1200000,1200000,1200000,   vectorMath.Vector3{500, 500, 500}, 0 )

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

func convertVector( v vectorMath.Vector3 ) mgl32.Vec3{
	return mgl32.Vec3{(float32)(v.X), (float32)(v.Y), (float32)(v.Z)}
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

//
func (glRenderer *OpenglRenderer) CreateGeometry( geometry *Geometry ) {

	// Configure the vertex data
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(geometry.Verticies)*4, gl.Ptr(geometry.Verticies), gl.STATIC_DRAW)
	geometry.vboId = vbo

	var ibo uint32
	gl.GenBuffers(1, &ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(geometry.Indicies)*4, gl.Ptr(geometry.Indicies), gl.STATIC_DRAW)
	geometry.iboId = ibo
}

//
func (glRenderer *OpenglRenderer) DestroyGeometry( geometry *Geometry ) {

}

//setup Texture
func (glRenderer *OpenglRenderer) CreateMaterial( material *Material ) {
	if material.Diffuse != nil {
		material.diffuseId = glRenderer.newTexture( material.Diffuse, gl.TEXTURE0 )
	}
	if material.Normal != nil {
		material.normalId = glRenderer.newTexture( material.Normal, gl.TEXTURE1 )
	}
	if material.Specular != nil {
		material.specularId = glRenderer.newTexture( material.Specular, gl.TEXTURE2 )
	} 
	if material.Roughness != nil {
		material.roughnessId = glRenderer.newTexture( material.Roughness, gl.TEXTURE3 )
	}
}

//setup Texture
func (glRenderer *OpenglRenderer) newTexture( img image.Image, textureUnit uint32 ) uint32 {
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

//
func (glRenderer *OpenglRenderer) DestroyMaterial( material *Material ) {

}

//
func (glRenderer *OpenglRenderer) DrawGeometry( geometry *Geometry ) {

	gl.BindBuffer(gl.ARRAY_BUFFER, geometry.vboId)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, geometry.iboId)

	//set verticies attribute
	vertAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 14*4, gl.PtrOffset(0))

	//set normals/tangent attribute
	normAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("normal\x00")))
	gl.EnableVertexAttribArray(normAttrib)
	gl.VertexAttribPointer(normAttrib, 3, gl.FLOAT, false, 14*4, gl.PtrOffset(3*4))
	tangentAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("tangent\x00")))
	gl.EnableVertexAttribArray(tangentAttrib)
	gl.VertexAttribPointer(tangentAttrib, 3, gl.FLOAT, false, 14*4, gl.PtrOffset(6*4))
	bitangentAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("bitangent\x00")))
	gl.EnableVertexAttribArray(bitangentAttrib)
	gl.VertexAttribPointer(bitangentAttrib, 3, gl.FLOAT, false, 14*4, gl.PtrOffset(9*4))

	//set texture coord attribute
	texCoordAttrib := uint32(gl.GetAttribLocation(glRenderer.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 14*4, gl.PtrOffset(12*4))

	//setup textures
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, geometry.Material.diffuseId)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, geometry.Material.normalId)
	gl.ActiveTexture(gl.TEXTURE2)
	gl.BindTexture(gl.TEXTURE_2D, geometry.Material.specularId)
	gl.ActiveTexture(gl.TEXTURE3)
	gl.BindTexture(gl.TEXTURE_2D, geometry.Material.roughnessId)

	gl.DrawElements(gl.TRIANGLES, (int32)(len(geometry.Indicies)), gl.UNSIGNED_INT, gl.PtrOffset(0))
}

// ambient, diffuse and specular light values ( i is the light index )
func (glRenderer *OpenglRenderer) CreateLight( ar,ag,ab, dr,dg,db, sr,sg,sb float32, position vectorMath.Vector3, i int ){
	//position
	glRenderer.lights[(i*16)] = (float32)(position.X)
	glRenderer.lights[(i*16)+1] = (float32)(position.Y)
	glRenderer.lights[(i*16)+2] = (float32)(position.Z)
	glRenderer.lights[(i*16)+3] = 1
	//ambient
	glRenderer.lights[(i*16)+4] = ar
	glRenderer.lights[(i*16)+5] = ag
	glRenderer.lights[(i*16)+6] = ab
	glRenderer.lights[(i*16)+7] = 1
	//diffuse
	glRenderer.lights[(i*16)+8] = dr
	glRenderer.lights[(i*16)+9] = dg
	glRenderer.lights[(i*16)+10] = db
	glRenderer.lights[(i*16)+11] = 1
	//specular
	glRenderer.lights[(i*16)+12] = sr
	glRenderer.lights[(i*16)+13] = sg
	glRenderer.lights[(i*16)+14] = sb
	glRenderer.lights[(i*16)+15] = 1
	//set uniform array
	lightsUniform := gl.GetUniformLocation(glRenderer.program, gl.Str("lights\x00"))
	gl.Uniform4fv( lightsUniform, (int32)(MAX_LIGHTS*16), &glRenderer.lights[0] )
}

func (glRenderer *OpenglRenderer) DestroyLight( i int ){

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

var vertexShader string = `
#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec3 normal;
in vec3 tangent;
in vec3 bitangent;
in vec2 vertTexCoord;

out vec3 fragNormal;
out vec3 fragTangent;
out vec3 fragBitangent;
out vec2 fragTexCoord;
out vec4 vertexPosition;

void main() {
	fragTexCoord = vertTexCoord;
	fragNormal = normal;
	fragTangent = tangent;
	fragBitangent = bitangent;
	vertexPosition = projection * camera * model * vec4(vert, 1);
	gl_Position = vertexPosition;
}` + "\x00"

var fragmentShader = `
#version 330

#define MAX_LIGHTS 8

#define LIGHT_POSITION 0
#define LIGHT_AMBIENT 1
#define LIGHT_DIFFUSE 2
#define LIGHT_SPECULAR 3

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

uniform vec4 lights[ MAX_LIGHTS * 4 ];

uniform sampler2D diffuse;
uniform sampler2D normal;
uniform sampler2D specular;
uniform sampler2D roughness;

in vec3 fragNormal;
in vec3 fragTangent;
in vec3 fragBitangent;
in vec2 fragTexCoord;
in vec4 vertexPosition;

out vec4 outputColor;

void main() {
	vec4 diffuseValue = texture(diffuse, fragTexCoord);
	vec4 normalValue = texture(normal, fragTexCoord);
	vec4 specularValue = texture(specular, fragTexCoord);
	vec4 roughnessValue = texture(roughness, fragTexCoord);
	vec4 finalColor = vec4(0,0,0,1);

	//Normal calculations
 	vec3 normalMapValue = vec3( normalValue.rgb ) * 2 - 1;
 	if( abs(normalMapValue.x) < 0.1 && abs(normalMapValue.y) < 0.1 && abs(normalMapValue.z) < 0.1 ){
 		normalMapValue = vec3(0,0,1);
 	}

	//tangent space conversion
	mat3 TBNMatrix = mat3(fragTangent, fragBitangent, fragNormal);

	//lights
	for (int i=0;i<MAX_LIGHTS;i++){

		//light components
		vec4 LightPos = projection * camera * lights[(i*4)+LIGHT_POSITION];
		vec4 LightAmb = lights[(i*4)+LIGHT_AMBIENT];
		vec4 LightDiff = lights[(i*4)+LIGHT_DIFFUSE];
		vec4 LightSpec = lights[(i*4)+LIGHT_SPECULAR];

		//point light source 
		vec4 lightDirection = vec4( LightPos - vertexPosition );
		float lightDistanceSQ = lightDirection.x*lightDirection.x + lightDirection.y*lightDirection.y + lightDirection.z*lightDirection.z;
		float illuminanceMultiplier = 1 / lightDistanceSQ;
		vec4 eyeDirection = vec4( - vertexPosition ); // eyePos is zero

		//tangent space
		vec3 TBNlightDirection = normalize( lightDirection.xyz * TBNMatrix );
		vec3 TBNeyeDirection = normalize( eyeDirection.xyz * TBNMatrix );

		//ambient component
		vec4 ambientOut = diffuseValue * LightAmb;
		//diffuse component
	 	float diffuseMultiplier = max(0.0, dot(normalMapValue, TBNlightDirection));
		vec4 diffuseOut = diffuseValue * diffuseMultiplier * LightDiff;
		//specular component
		vec4 reflectedEye = vec4( reflect( TBNeyeDirection, normalMapValue ), 1);
	 	float specularMultiplier = max(0.0, pow( dot(reflectedEye.xyz, TBNlightDirection), 2.0));
		vec4 specularOut = vec4( specularValue.rgb, 1 ) * specularMultiplier * LightSpec;

		finalColor += (( ambientOut + diffuseOut + specularOut ) * illuminanceMultiplier);
   	}

	//final output
	outputColor = finalColor;
}` + "\x00"
