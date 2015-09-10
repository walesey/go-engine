package renderer

import(
    "image"
    "image/color"
    
    "github.com/Walesey/goEngine/vectorMath"
)

type Material struct {
    diffuseId, normalId, specularId, glossId, roughnessId uint32
    loaded bool
    LightingMode int32
    Diffuse, Normal, Specular, Roughness image.Image
}

func CreateMaterial() Material {
    return Material{ loaded : false, LightingMode : MODE_LIT }
}

type Flipbook struct {
    IndexX, IndexY int
    FrameSizeX, FrameSizeY float32
}

//Geometry
type Geometry struct {
    vboId, iboId uint32
    loadedLen int
    VboDirty bool
    Indicies []uint32
    Verticies []float32
    Material *Material
    CullBackface bool
    Flipbook Flipbook
}

//vericies format : x,y,z,   nx,ny,nz,tx,ty,tz,btx,bty,btz,   u,v
//indicies format : f1,f2,f3 (triangles)
func CreateGeometry( indicies []uint32, verticies []float32 ) Geometry {
    mat := CreateMaterial()
    return Geometry{ 
        Indicies : indicies, 
        Verticies : verticies, 
        Material: &mat, 
        loadedLen: 0,
        CullBackface : true, 
        Flipbook: Flipbook{0, 0, 1.0, 1.0},
    }
}

func (geometry *Geometry) Draw( renderer Renderer ) {
    geometry.load( renderer )
    renderer.DrawGeometry( geometry )
}

func (geometry *Geometry) load( renderer Renderer ) {
    if geometry.loadedLen < len(geometry.Verticies) && len(geometry.Indicies) != 0 && len(geometry.Verticies) != 0 {
        renderer.CreateGeometry( geometry )
        geometry.loadedLen = len(geometry.Verticies)
    }
    if !geometry.Material.loaded {
        renderer.CreateMaterial( geometry.Material )
        geometry.Material.loaded = true
    }
}

func (geometry *Geometry) ClearBuffers() {
    geometry.Indicies = geometry.Indicies[:0]
    geometry.Verticies = geometry.Verticies[:0]
}

func (geometry *Geometry) SetColor(color color.NRGBA) {
    for i:=14 ;i<len(geometry.Verticies) ; i=i+18 {
        geometry.Verticies[i] = float32(color.R) / 255.0
        geometry.Verticies[i+1] = float32(color.G) / 255.0
        geometry.Verticies[i+2] = float32(color.B) / 255.0
        geometry.Verticies[i+3] = float32(color.A) / 255.0
    }
    geometry.VboDirty = true
}

func (geometry *Geometry) Transform(transform Transform) {
    geometry.transformRange(transform, 0)
}

func (geometry *Geometry) transformRange(transform Transform, from int) {
    for i:=from ; i<len(geometry.Verticies) ; i=i+18 {
        v := transform.TransformCoordinate( vectorMath.Vector3{
            float64(geometry.Verticies[i]),
            float64(geometry.Verticies[i+1]),
            float64(geometry.Verticies[i+2]),
        })
        n := transform.TransformNormal( vectorMath.Vector3{
            float64(geometry.Verticies[i+3]),
            float64(geometry.Verticies[i+4]),
            float64(geometry.Verticies[i+5]),
        })
        t := transform.TransformNormal( vectorMath.Vector3{
            float64(geometry.Verticies[i+6]),
            float64(geometry.Verticies[i+7]),
            float64(geometry.Verticies[i+8]),
        })
        bt := transform.TransformNormal( vectorMath.Vector3{
            float64(geometry.Verticies[i+9]),
            float64(geometry.Verticies[i+10]),
            float64(geometry.Verticies[i+11]),
        })
        geometry.Verticies[i] = float32(v.X)
        geometry.Verticies[i+1] = float32(v.Y)
        geometry.Verticies[i+2] = float32(v.Z)
        geometry.Verticies[i+3] = float32(n.X)
        geometry.Verticies[i+4] = float32(n.Y)
        geometry.Verticies[i+5] = float32(n.Z)
        geometry.Verticies[i+6] = float32(t.X)
        geometry.Verticies[i+7] = float32(t.Y)
        geometry.Verticies[i+8] = float32(t.Z)
        geometry.Verticies[i+9] = float32(bt.X)
        geometry.Verticies[i+10] = float32(bt.Y)
        geometry.Verticies[i+11] = float32(bt.Z)
    }
    geometry.VboDirty = true
}

//load the verts/indicies of geometry into geom
func (geometry *Geometry) Optimize( geom *Geometry, transform Transform ) {
    vertOffset := len(geom.Verticies)
    indexOffset := len(geom.Indicies)
    geom.Verticies = append(geom.Verticies, geometry.Verticies...)
    geom.Indicies = append(geom.Indicies, geometry.Indicies...)
    for i,_ := range geometry.Indicies {
        geom.Indicies[i+indexOffset] = geometry.Indicies[i] + uint32( vertOffset/18 )
    }
    geom.transformRange(transform, vertOffset)
    geometry.VboDirty = true
}

//Primitives
func CreateBox( height, width float32 ) Geometry {
    verticies := []float32{
        -width/2,0,height/2,  0,1,0, 1,0,-1, -1,0,-1, 0,0, 1.0,1.0,1.0,1.0,
        width/2,0,height/2,   0,1,0, 1,0,-1, -1,0,-1, 1,0, 1.0,1.0,1.0,1.0,
        width/2,0,-height/2,  0,1,0, 1,0,-1, -1,0,-1, 1,1, 1.0,1.0,1.0,1.0,
        width/2,0,-height/2,  0,1,0, 1,0,-1, -1,0,-1, 1,1, 1.0,1.0,1.0,1.0,
        -width/2,0,-height/2, 0,1,0, 1,0,-1, -1,0,-1, 0,1, 1.0,1.0,1.0,1.0,
        -width/2,0,height/2,  0,1,0, 1,0,-1, -1,0,-1, 0,0, 1.0,1.0,1.0,1.0,
    }
    indicies := []uint32{0,1,2,3,4,5}
    return CreateGeometry(indicies, verticies)
}
