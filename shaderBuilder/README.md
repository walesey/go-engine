# ShaderBuilder

A Shader preprocessor that handles the following # macros
* #include - works like c include for local shader files only
* #vert / #endvert - code in this section is only outputted to the vertex shader.
* #frag / #endfrag
* #geo / #geo

The output of the preprocessor is a single shader file to be compiled by opengl etc.

### Example Usage

```
	shaderBuilder path/to/file.glsl vert > out.vert
	shaderBuilder path/to/file.glsl frag > out.frag
```
