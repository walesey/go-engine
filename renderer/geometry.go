package renderer

import (
	"image"
	"image/color"

	vmath "github.com/walesey/go-engine/vectormath"
)

const VertexStride = 12

const (
	MODE_UNLIT int32 = iota
	MODE_LIT
	MODE_EMIT
)

const (
	TRANSPARENCY_NON_EMISSIVE int = iota
	TRANSPARENCY_EMISSIVE
)

type Material struct {
	DiffuseId, NormalId, SpecularId, RoughnessId uint32
	loaded                                       bool
	LightingMode                                 int32
	Transparency                                 int
	DepthTest, DepthMask                         bool
	Diffuse, Normal, Specular, Roughness         image.Image
}

func CreateMaterial() *Material {
	return &Material{
		loaded:       false,
		LightingMode: MODE_LIT,
		Transparency: TRANSPARENCY_NON_EMISSIVE,
		DepthTest:    true,
		DepthMask:    true,
	}
}

//Geometry
type Geometry struct {
	VboId, IboId uint32
	loaded       bool
	VboDirty     bool
	Indicies     []uint32
	Verticies    []float32
	Material     *Material
	CullBackface bool
}

//vericies format : x,y,z,   nx,ny,nz,   u,v,  r,g,b,a
//indicies format : f1,f2,f3 (triangles)
func CreateGeometry(indicies []uint32, verticies []float32) *Geometry {
	mat := CreateMaterial()
	return &Geometry{
		Indicies:     indicies,
		Verticies:    verticies,
		Material:     mat,
		loaded:       false,
		CullBackface: true,
	}
}

func (geometry *Geometry) Copy() *Geometry {
	indicies := make([]uint32, len(geometry.Indicies))
	verticies := make([]float32, len(geometry.Verticies))
	copy(indicies, geometry.Indicies)
	copy(verticies, geometry.Verticies)
	return CreateGeometry(indicies, verticies)
}

func (geometry *Geometry) Draw(renderer Renderer) {
	geometry.load(renderer)
	renderer.DrawGeometry(geometry)
}

func (geometry *Geometry) load(renderer Renderer) {
	if !geometry.loaded && len(geometry.Indicies) != 0 && len(geometry.Verticies) != 0 {
		renderer.CreateGeometry(geometry)
		geometry.loaded = true
	}
	if !geometry.Material.loaded {
		renderer.CreateMaterial(geometry.Material)
		geometry.Material.loaded = true
	}
}

func (geometry *Geometry) Destroy(renderer Renderer) {
	renderer.DestroyGeometry(geometry)
	geometry.loaded = false
	if geometry.Material != nil {
		renderer.DestroyMaterial(geometry.Material)
		geometry.Material.loaded = false
	}
}

func (geometry *Geometry) Centre() vmath.Vector3 {
	return vmath.Vector3{0, 0, 0}
}

func (geometry *Geometry) ClearBuffers() {
	geometry.Indicies = geometry.Indicies[:0]
	geometry.Verticies = geometry.Verticies[:0]
}

func (geometry *Geometry) SetColor(color color.Color) {
	for i := 8; i < len(geometry.Verticies); i = i + VertexStride {
		r, g, b, a := color.RGBA()
		geometry.Verticies[i] = float32(r) / 65535.0
		geometry.Verticies[i+1] = float32(g) / 65535.0
		geometry.Verticies[i+2] = float32(b) / 65535.0
		geometry.Verticies[i+3] = float32(a) / 65535.0
	}
	geometry.VboDirty = true
}

func (geometry *Geometry) Transform(transform Transform) {
	geometry.transformRange(transform, 0)
}

func (geometry *Geometry) transformRange(transform Transform, from int) {
	for i := from; i < len(geometry.Verticies); i = i + VertexStride {
		v := transform.TransformCoordinate(vmath.Vector3{
			float64(geometry.Verticies[i]),
			float64(geometry.Verticies[i+1]),
			float64(geometry.Verticies[i+2]),
		})
		n := transform.TransformNormal(vmath.Vector3{
			float64(geometry.Verticies[i+3]),
			float64(geometry.Verticies[i+4]),
			float64(geometry.Verticies[i+5]),
		})
		geometry.Verticies[i] = float32(v.X)
		geometry.Verticies[i+1] = float32(v.Y)
		geometry.Verticies[i+2] = float32(v.Z)
		geometry.Verticies[i+3] = float32(n.X)
		geometry.Verticies[i+4] = float32(n.Y)
		geometry.Verticies[i+5] = float32(n.Z)
	}
	geometry.VboDirty = true
}

//load the verts/indicies of geometry into destination Geometry
func (geometry *Geometry) Optimize(destination *Geometry, transform Transform) {
	vertOffset := len(destination.Verticies)
	indexOffset := len(destination.Indicies)
	destination.Verticies = append(destination.Verticies, geometry.Verticies...)
	destination.Indicies = append(destination.Indicies, geometry.Indicies...)
	for i, _ := range geometry.Indicies {
		destination.Indicies[i+indexOffset] = geometry.Indicies[i] + uint32(vertOffset/VertexStride)
	}
	destination.transformRange(transform, vertOffset)
	geometry.VboDirty = true
}

func (geometry *Geometry) SetUVs(uvs ...float32) {
	for i := 0; i < len(uvs); i = i + 2 {
		geometry.Verticies[((i/2)*VertexStride)+6] = uvs[i]
		geometry.Verticies[((i/2)*VertexStride)+7] = uvs[i+1]
	}
	geometry.VboDirty = true
}

//Primitives
func CreateBox(width, height float32) *Geometry {
	return CreateBoxWithOffset(width, height, -width/2, -height/2)
}

func CreateBoxWithOffset(width, height, offsetX, offsetY float32) *Geometry {
	verticies := []float32{
		offsetX, height + offsetY, 0, 0, 1, 0, 0, 0, 1.0, 1.0, 1.0, 1.0,
		width + offsetX, height + offsetY, 0, 0, 1, 0, 1, 0, 1.0, 1.0, 1.0, 1.0,
		width + offsetX, offsetY, 0, 0, 1, 0, 1, 1, 1.0, 1.0, 1.0, 1.0,
		offsetX, offsetY, 0, 0, 1, 0, 0, 1, 1.0, 1.0, 1.0, 1.0,
	}
	indicies := []uint32{0, 1, 2, 2, 3, 0}
	return CreateGeometry(indicies, verticies)
}

func CreateSkyBox() *Geometry {
	return CreateGeometry(cubeIndicies, skyboxVerticies)
}

// CreateBeam - creates a square prism oriented along the vector
func CreateBeam(width float64, vector vmath.Vector3) *Geometry {
	direction := vector.Normalize()
	geo := CreateBoxWithOffset(float32(width), float32(width), float32(-width*0.5), float32(-width*0.5))
	geo2 := CreateBoxWithOffset(float32(width), float32(width), float32(-width*0.5), float32(-width*0.5))
	facingTx := CreateTransform()
	facingTx.From(vmath.Vector3{1, 1, 1}, vmath.Vector3{}, vmath.FacingOrientation(0, direction, vmath.Vector3{0, 0, 1}, vmath.Vector3{1, 0, 0}))
	geo.Transform(facingTx)
	facingTx.From(vmath.Vector3{1, 1, 1}, vector, vmath.FacingOrientation(0, direction, vmath.Vector3{0, 0, -1}, vmath.Vector3{1, 0, 0}))
	geo2.Optimize(geo, facingTx)
	geo.Indicies = append(geo.Indicies, 0, 1, 4, 4, 5, 0) //top
	geo.Indicies = append(geo.Indicies, 1, 2, 7, 7, 4, 1) //side
	geo.Indicies = append(geo.Indicies, 2, 3, 6, 6, 7, 2) //bottom
	geo.Indicies = append(geo.Indicies, 3, 0, 5, 5, 6, 3) //side
	geo.CullBackface = false
	return geo
}
