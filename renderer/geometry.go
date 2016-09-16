package renderer

import (
	"image"
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/util"
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
	VboId, IboId    uint32
	loaded          bool
	VboDirty        bool
	Indicies        []uint32
	Verticies       []float32
	Material        *Material
	CullBackface    bool
	FrustrumCulling bool
	boundingRadius  float32
}

//vericies format : x,y,z,   nx,ny,nz,   u,v,  r,g,b,a
//indicies format : f1,f2,f3 (triangles)
func CreateGeometry(indicies []uint32, verticies []float32) *Geometry {
	mat := CreateMaterial()
	return &Geometry{
		Indicies:        indicies,
		Verticies:       verticies,
		Material:        mat,
		loaded:          false,
		CullBackface:    true,
		FrustrumCulling: true,
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
	if !geometry.FrustrumCulling || renderer.FrustrumContainsSphere(geometry.boundingRadius) {
		renderer.DrawGeometry(geometry)
	}
}

func (geometry *Geometry) load(renderer Renderer) {
	if !geometry.loaded && len(geometry.Indicies) != 0 && len(geometry.Verticies) != 0 {
		geometry.boundingRadius = geometry.MaximalPointFromGeometry().Len()
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

func (geometry *Geometry) Centre() mgl32.Vec3 {
	return mgl32.Vec3{0, 0, 0}
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
	geometry.updateGeometry()
}

func (geometry *Geometry) Transform(transform mgl32.Mat4) {
	geometry.transformRange(transform, 0)
}

func (geometry *Geometry) transformRange(transform mgl32.Mat4, from int) {
	for i := from; i < len(geometry.Verticies); i = i + VertexStride {
		v := mgl32.TransformCoordinate(mgl32.Vec3{
			geometry.Verticies[i],
			geometry.Verticies[i+1],
			geometry.Verticies[i+2],
		}, transform)
		n := mgl32.TransformNormal(mgl32.Vec3{
			geometry.Verticies[i+3],
			geometry.Verticies[i+4],
			geometry.Verticies[i+5],
		}, transform)
		geometry.Verticies[i] = v.X()
		geometry.Verticies[i+1] = v.Y()
		geometry.Verticies[i+2] = v.Z()
		geometry.Verticies[i+3] = n.X()
		geometry.Verticies[i+4] = n.Y()
		geometry.Verticies[i+5] = n.Z()
	}
	geometry.updateGeometry()
}

//load the verts/indicies of geometry into destination Geometry
func (geometry *Geometry) Optimize(destination *Geometry, transform mgl32.Mat4) {
	vertOffset := len(destination.Verticies)
	indexOffset := len(destination.Indicies)
	destination.Verticies = append(destination.Verticies, geometry.Verticies...)
	destination.Indicies = append(destination.Indicies, geometry.Indicies...)
	for i, _ := range geometry.Indicies {
		destination.Indicies[i+indexOffset] = geometry.Indicies[i] + uint32(vertOffset/VertexStride)
	}
	destination.transformRange(transform, vertOffset)
	geometry.updateGeometry()
}

func (geometry *Geometry) SetUVs(uvs ...float32) {
	for i := 0; i < len(uvs); i = i + 2 {
		geometry.Verticies[((i/2)*VertexStride)+6] = uvs[i]
		geometry.Verticies[((i/2)*VertexStride)+7] = uvs[i+1]
	}
	geometry.updateGeometry()
}

func (geometry *Geometry) updateGeometry() {
	geometry.boundingRadius = geometry.MaximalPointFromGeometry().Len()
	geometry.VboDirty = true
}

// Gets the furthest point from the origin in this geometry
func (geometry *Geometry) MaximalPointFromGeometry() mgl32.Vec3 {
	var longestSq float32
	var longest *mgl32.Vec3
	for i := 0; i < len(geometry.Indicies); i = i + 1 {
		index := geometry.Indicies[i]
		v := mgl32.Vec3{
			geometry.Verticies[index*VertexStride],
			geometry.Verticies[index*VertexStride+1],
			geometry.Verticies[index*VertexStride+2],
		}
		if longest == nil || util.Vec3LenSq(v) > longestSq {
			longestSq = util.Vec3LenSq(v)
			longest = &v
		}
	}
	if longest != nil {
		return *longest
	}
	return mgl32.Vec3{}
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
func CreateBeam(width float32, vector mgl32.Vec3) *Geometry {
	direction := vector.Normalize()
	geo := CreateBoxWithOffset(width, width, -width*0.5, -width*0.5)
	geo2 := CreateBoxWithOffset(width, width, -width*0.5, -width*0.5)
	facingTx := util.Mat4From(mgl32.Vec3{1, 1, 1}, mgl32.Vec3{}, util.FacingOrientation(0, direction, mgl32.Vec3{0, 0, 1}, mgl32.Vec3{1, 0, 0}))
	geo.Transform(facingTx)
	facingTx = util.Mat4From(mgl32.Vec3{1, 1, 1}, vector, util.FacingOrientation(0, direction, mgl32.Vec3{0, 0, -1}, mgl32.Vec3{1, 0, 0}))
	geo2.Optimize(geo, facingTx)
	geo.Indicies = append(geo.Indicies, 0, 1, 4, 4, 5, 0) //top
	geo.Indicies = append(geo.Indicies, 1, 2, 7, 7, 4, 1) //side
	geo.Indicies = append(geo.Indicies, 2, 3, 6, 6, 7, 2) //bottom
	geo.Indicies = append(geo.Indicies, 3, 0, 5, 5, 6, 3) //side
	geo.CullBackface = false
	return geo
}
