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
    loaded bool
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
        loaded : false, 
        CullBackface : true, 
        Flipbook: Flipbook{0, 0, 1.0, 1.0},
    }
}

func (geometry *Geometry) Draw( renderer Renderer ) {
    geometry.load( renderer )
    renderer.DrawGeometry( geometry )
}

func (geometry *Geometry) load( renderer Renderer ) {
    if !geometry.loaded {
        renderer.CreateGeometry( geometry )
        geometry.loaded = true
    }
    if !geometry.Material.loaded {
        renderer.CreateMaterial( geometry.Material )
        geometry.Material.loaded = true
    }
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

func (geometry *Geometry) Optimize( geom *Geometry, transform Transform ) {
    vertOffset := len(geom.Verticies)
    indexOffset := vertOffset / 18
    geom.Verticies = append(geom.Verticies, geometry.Verticies...)
    for i:=vertOffset ; i<len(geom.Verticies) ; i=i+18 {
        v := transform.TransformCoordinate( vectorMath.Vector3{
            float64(geom.Verticies[i]),
            float64(geom.Verticies[i+1]),
            float64(geom.Verticies[i+2]),
        })
        geom.Verticies[i] = float32(v.X)
        geom.Verticies[i+1] = float32(v.Y)
        geom.Verticies[i+2] = float32(v.Z)
    }
    for i,_ := range geometry.Indicies {
        geom.Indicies = append(geom.Indicies, geometry.Indicies[i] + uint32(indexOffset))
    }
    geom.Material = geometry.Material
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
