package renderer

type Shader struct {
	Program uint32
	loaded  bool

	Uniforms         map[string]interface{}
	FragDataLocation string

	FragSrc, VertSrc, GeoSrc string
}

func NewShader() *Shader {
	return &Shader{
		FragDataLocation: "outputColor",
		Uniforms:         make(map[string]interface{}),
	}
}

func (shader *Shader) Copy() (copy *Shader) {
	copy = NewShader()
	copy.FragDataLocation = shader.FragDataLocation
	copy.FragSrc = shader.FragSrc
	copy.VertSrc = shader.VertSrc
	copy.GeoSrc = shader.GeoSrc
	return
}
