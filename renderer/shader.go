package renderer

type Shader struct {
	Program uint32
	Loaded  bool

	Uniforms          map[string]interface{}
	FragDataLocations []string
	InputBuffers      int

	textureUnitCounter int32
	TextureUnits       map[string]int32

	FragSrc, VertSrc, GeoSrc string
}

func NewShader() *Shader {
	return &Shader{
		FragDataLocations: []string{"outputColor"},
		Uniforms:          make(map[string]interface{}),
		TextureUnits:      make(map[string]int32),
		InputBuffers:      1,
	}
}

// AddTexture - allocates a texture unit to the shader
func (shader *Shader) AddTexture(textureName string) int32 {
	index, ok := shader.TextureUnits[textureName]
	if !ok {
		index = shader.textureUnitCounter
		shader.textureUnitCounter = shader.textureUnitCounter + 1
		shader.TextureUnits[textureName] = index
	}

	shader.Uniforms[textureName] = index
	return index
}

func (shader *Shader) Copy() (copy *Shader) {
	copy = NewShader()
	copy.FragDataLocations = shader.FragDataLocations
	copy.FragSrc = shader.FragSrc
	copy.VertSrc = shader.VertSrc
	copy.GeoSrc = shader.GeoSrc
	copy.InputBuffers = shader.InputBuffers
	return
}
