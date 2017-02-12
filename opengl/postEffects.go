package opengl

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/walesey/go-engine/renderer"
)

type postEffect struct {
	program  uint32
	fboId    uint32
	dboId    uint32
	textures []*renderer.Texture
	shader   *renderer.Shader
}

func (pe postEffect) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, pe.fboId)
}

//Set up the frame buffer for rendering each post effect filter pass
func (glRenderer *OpenglRenderer) initPostEffects() {
	//post effects quad
	quadVertices := []float32{
		// Positions  // Texture Coords
		-1, -1, 0, 0,
		1, -1, 1, 0,
		-1, 1, 0, 1,
		1, 1, 1, 1,
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVertices)*4, gl.Ptr(quadVertices), gl.STATIC_DRAW)
	glRenderer.postEffectVbo = vbo
}

func (glRenderer *OpenglRenderer) CreatePostEffect(shader *renderer.Shader) {
	glRenderer.UseShader(shader)
	glRenderer.enableShader()

	//Get render buffer dimensions
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	var dims [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &dims[0])
	bufferWidth := dims[2]
	bufferHeight := dims[3]

	//Create depth buffer
	var dbo uint32
	gl.GenRenderbuffers(1, &dbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, dbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT16, bufferWidth, bufferHeight)

	//Create frame buffer
	var fbo uint32
	gl.GenFramebuffers(1, &fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, dbo)

	//Create Textures
	buffers := make([]uint32, shader.InputBuffers)
	fbo_textures := make([]*renderer.Texture, shader.InputBuffers)

	for i := 0; i < shader.InputBuffers; i++ {
		buffers[i] = gl.COLOR_ATTACHMENT0 + uint32(i)
		textureUnit := gl.TEXTURE0 + uint32(i)
		var fbo_texture uint32
		gl.GenTextures(1, &fbo_texture)
		gl.ActiveTexture(textureUnit)
		gl.BindTexture(gl.TEXTURE_2D, fbo_texture)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, bufferWidth, bufferHeight, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, buffers[i], gl.TEXTURE_2D, fbo_texture, 0)
		fbo_textures[i] = renderer.NewTexture(fmt.Sprintf("tex%v", i), nil, false)
		fbo_textures[i].TextureId = fbo_texture
		fbo_textures[i].Loaded = true
	}

	gl.DrawBuffers(int32(len(buffers)), &buffers[0])
	gl.BindRenderbuffer(gl.FRAMEBUFFER, 0)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	//add new postEffect to the queue
	newPe := postEffect{
		program:  shader.Program,
		textures: fbo_textures,
		dboId:    dbo,
		fboId:    fbo,
		shader:   shader,
	}
	glRenderer.postEffects = append(glRenderer.postEffects, newPe)
}

func (glRenderer *OpenglRenderer) DestroyPostEffects(shader *renderer.Shader) {
	for i, po := range glRenderer.postEffects {
		if po.shader == shader {
			gl.DeleteRenderbuffers(1, &po.dboId)
			gl.DeleteFramebuffers(1, &po.fboId)
			gl.DeleteTextures(int32(len(po.textures)), &po.textures[0].TextureId)
			glRenderer.postEffects = append(glRenderer.postEffects[:i], glRenderer.postEffects[i+1:]...)
			break
		}
	}
}

func (glRenderer *OpenglRenderer) renderPostEffect(pe postEffect) {
	glRenderer.UseShader(pe.shader)
	glRenderer.enableShader()
	gl.BindBuffer(gl.ARRAY_BUFFER, glRenderer.postEffectVbo)
	gl.Disable(gl.CULL_FACE)

	for _, texture := range pe.textures {
		textureUnit := pe.shader.AddTexture(texture) + gl.TEXTURE0
		gl.ActiveTexture(uint32(textureUnit))
		gl.BindTexture(gl.TEXTURE_2D, texture.TextureId)
	}

	setupUniforms(pe.shader)

	vertAttrib := uint32(gl.GetAttribLocation(pe.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(pe.program, gl.Str("texCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

}
