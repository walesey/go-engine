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

func (shader *Shader) AddTexture(tex *Texture) int32 {
	index, ok := shader.TextureUnits[tex.TextureName]
	if !ok {
		index = shader.textureUnitCounter
		shader.textureUnitCounter = shader.textureUnitCounter + 1
		shader.TextureUnits[tex.TextureName] = index
	}

	shader.Uniforms[tex.TextureName] = index
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
