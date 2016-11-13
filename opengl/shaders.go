package opengl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
)

func newProgram(shaders ...uint32) (uint32, error) {
	program := gl.CreateProgram()

	for _, shader := range shaders {
		gl.AttachShader(program, shader)
	}
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

	for _, shader := range shaders {
		gl.DeleteShader(shader)
	}

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
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

func setupUniforms(shader *renderer.Shader) {
	for name, uniform := range shader.Uniforms {
		uniformLocation := gl.GetUniformLocation(shader.Program, gl.Str(name+"\x00"))
		switch t := uniform.(type) {
		case float32:
			gl.Uniform1f(uniformLocation, t)
		case float64:
			gl.Uniform1f(uniformLocation, float32(t))
		case int32:
			gl.Uniform1i(uniformLocation, t)
		case int:
			gl.Uniform1i(uniformLocation, int32(t))
		case mgl32.Vec2:
			gl.Uniform2f(uniformLocation, t[0], t[1])
		case mgl32.Vec3:
			gl.Uniform3f(uniformLocation, t[0], t[1], t[2])
		case mgl32.Vec4:
			gl.Uniform4f(uniformLocation, t[0], t[1], t[2], t[3])
		case mgl32.Mat4:
			gl.UniformMatrix4fv(uniformLocation, 1, false, &t[0])
		case []float32:
			gl.Uniform4fv(uniformLocation, (int32)(len(t)), &t[0])
		default:
			fmt.Printf("unexpected type for shader uniform: %T\n", t)
		}
	}
}
