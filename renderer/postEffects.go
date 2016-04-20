package renderer

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/vectormath"
)

type Uniform struct {
	Name  string
	Value interface{}
}

type Shader struct {
	Name     string
	Uniforms []Uniform
}

type postEffect struct {
	program   uint32
	fboId     uint32
	dboId     uint32
	textureId uint32
	shader    Shader
}

//Set up the frame buffer for rendering each post effect filter pass
func (glRenderer *OpenglRenderer) initPostEffects() {
	//post effects quad
	verts := []float32{
		-1, -1, 0, 0,
		1, -1, 1, 0,
		-1, 1, 0, 1,
		1, 1, 1, 1,
	}
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	glRenderer.postEffectVbo = vbo
}

func (glRenderer *OpenglRenderer) CreatePostEffect(shader Shader) {

	//Create program
	shaderName := shader.Name
	program := programFromFile(shaderName+".vert", shaderName+".frag")
	gl.UseProgram(program)

	//Create Texture
	var fbo_texture uint32
	gl.ActiveTexture(gl.TEXTURE0)
	gl.GenTextures(1, &fbo_texture)
	gl.BindTexture(gl.TEXTURE_2D, fbo_texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(glRenderer.WindowWidth), int32(glRenderer.WindowHeight), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	//Create depth buffer
	var dbo uint32
	gl.GenRenderbuffers(1, &dbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, dbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT16, int32(glRenderer.WindowWidth), int32(glRenderer.WindowHeight))
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	//Create frame buffer
	var fbo uint32
	gl.GenFramebuffers(1, &fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, fbo_texture, 0)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, dbo)
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))
	gl.UseProgram(glRenderer.program)

	//add new postEffect to the queue
	newPe := postEffect{
		program:   program,
		textureId: fbo_texture,
		dboId:     dbo,
		fboId:     fbo,
		shader:    shader,
	}
	glRenderer.postEffects = append(glRenderer.postEffects, newPe)
}

func (glRenderer *OpenglRenderer) DestroyPostEffects(shader Shader) {
	shaderName := shader.Name
	for i, po := range glRenderer.postEffects {
		if po.shader.Name == shaderName {
			gl.DeleteRenderbuffers(1, &po.dboId)
			gl.DeleteTextures(1, &po.textureId)
			gl.DeleteFramebuffers(1, &po.fboId)
			glRenderer.postEffects = append(glRenderer.postEffects[:i], glRenderer.postEffects[i+1:]...)
			break
		}
	}
}

func (glRenderer *OpenglRenderer) renderPostEffect(pe postEffect) {
	gl.UseProgram(pe.program)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, pe.textureId)
	gl.Disable(gl.CULL_FACE)
	gl.BindBuffer(gl.ARRAY_BUFFER, glRenderer.postEffectVbo)

	//update uniforms
	for _, uniform := range pe.shader.Uniforms {
		uniformLocation := gl.GetUniformLocation(pe.program, gl.Str(uniform.Name+"\x00"))
		switch t := uniform.Value.(type) {
		default:
			fmt.Printf("unexpected type for shader uniform: %T\n", t)
		case float32:
			gl.Uniform1f(uniformLocation, uniform.Value.(float32))
		case float64:
			gl.Uniform1f(uniformLocation, float32(uniform.Value.(float64)))
		case int32:
			gl.Uniform1i(uniformLocation, uniform.Value.(int32))
		case int:
			gl.Uniform1i(uniformLocation, int32(uniform.Value.(int)))
		case mgl32.Vec2:
			vec := uniform.Value.(mgl32.Vec2)
			gl.Uniform2f(uniformLocation, vec[0], vec[1])
		case mgl32.Vec3:
			vec := uniform.Value.(mgl32.Vec3)
			gl.Uniform3f(uniformLocation, vec[0], vec[1], vec[2])
		case mgl32.Vec4:
			vec := uniform.Value.(mgl32.Vec4)
			gl.Uniform4f(uniformLocation, vec[0], vec[1], vec[2], vec[3])
		case vectormath.Vector3:
			vec := uniform.Value.(vectormath.Vector3)
			gl.Uniform3f(uniformLocation, float32(vec.X), float32(vec.Y), float32(vec.Z))
		}
	}

	vertAttrib := uint32(gl.GetAttribLocation(pe.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	texCoordAttrib := uint32(gl.GetAttribLocation(pe.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	gl.DisableVertexAttribArray(texCoordAttrib)
}
